package main

import (
	"cmp"
	"encoding/json"
	"flag"
	"log/slog"
	"os"
	"path/filepath"
	"slices"
	"sync"
	"time"

	"github.com/RiverMint78/pone-quest/internal/pone"
	"github.com/joho/godotenv"
)

type task struct {
	id   string
	path string
}

type result struct {
	id   string
	desc string
	err  error
}

func main() {
	// cli arguments
	fullRedo := flag.Bool("full", false, "是否全部重做标注")
	concurrency := flag.Int("c", 8, "最大并发处理数量")
	rpm := flag.Int("rpm", 1000, "每分钟最大请求数")
	maxRetry := flag.Int("retry", 5, "失败后的最大重试次数")
	flag.Parse()

	_ = godotenv.Load()
	logger := slog.Default()

	imgDir := os.Getenv("PQ_IMAGE_DIR")
	jsonPath := os.Getenv("PQ_IMAGEITEM_FILE")
	apiKey := os.Getenv("PQ_LABELER_APIKEY")
	apiURL := os.Getenv("PQ_LABELER_APIURL")

	// 读 JSON 索引
	var existingItems []pone.ImageItem
	existingMap := make(map[string]string)
	if raw, err := os.ReadFile(jsonPath); err == nil {
		_ = json.Unmarshal(raw, &existingItems)
		for _, it := range existingItems {
			existingMap[it.ID] = it.Description
		}
	}

	// 产生任务
	absPath, _ := filepath.Abs(imgDir)
	logger.Info("正在扫描目录", "path", absPath)
	files, _ := os.ReadDir(imgDir)
	var tasks []task
	finalResults := make(map[string]string)

	for _, f := range files {
		if f.IsDir() || !isImage(f.Name()) {
			continue
		}
		id := f.Name()
		if !*fullRedo {
			// 不重复产生
			if desc, ok := existingMap[id]; ok {
				finalResults[id] = desc
				continue
			}
		}
		tasks = append(tasks, task{id: id, path: filepath.Join(imgDir, id)})
	}

	if len(tasks) == 0 {
		logger.Info("没有发现需要标注的新图片")
		return
	}

	// concurrency and RPM
	taskChan := make(chan task, len(tasks))
	resChan := make(chan result, len(tasks))
	var wg sync.WaitGroup

	// RPM rate limiter
	interval := time.Duration(60000/(*rpm)) * time.Millisecond
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	// worker
	logger.Info("开始标注任务", "total", len(tasks), "concurrency", *concurrency, "rpm", *rpm)
	for i := 0; i < *concurrency; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for t := range taskChan {
				<-ticker.C // rate limit

				var desc string
				var err error
				// retry
				for r := 0; r <= *maxRetry; r++ {
					if r > 0 {
						logger.Warn("正在重试", "file", t.id, "attempt", r)
						time.Sleep(time.Second * time.Duration(r*2))
					}
					desc, err = requestAPILabel(apiURL, apiKey, t.path)
					if err == nil {
						break
					}
					logger.Warn("单次请求失败", "file", t.id, "attempt", r, "err", err)
				}
				resChan <- result{id: t.id, desc: desc, err: err}
			}
		}()
	}

	// launch tasks
	go func() {
		for _, t := range tasks {
			taskChan <- t
		}
		close(taskChan)
	}()

	go func() {
		wg.Wait()
		close(resChan)
	}()

	for r := range resChan {
		if r.err != nil {
			logger.Error("标注失败", "file", r.id, "err", r.err)
			continue
		}
		finalResults[r.id] = r.desc
	}

	// 写索引 JSON
	var finalItems []pone.ImageItem
	for id, desc := range finalResults {
		finalItems = append(finalItems, pone.ImageItem{ID: id, Description: desc})
	}
	slices.SortFunc(finalItems, func(a, b pone.ImageItem) int {
		return cmp.Compare(a.ID, b.ID)
	})

	finalRaw, _ := json.MarshalIndent(finalItems, "", "  ")
	_ = os.WriteFile(jsonPath, finalRaw, 0644)
	logger.Info("任务全部完成", "total_files", len(finalItems))
}
