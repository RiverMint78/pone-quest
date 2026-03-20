package main

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"os"
	"sort"
	"strconv"
	"strings"

	"github.com/RiverMint78/pone-quest/internal/embed"
	"github.com/RiverMint78/pone-quest/internal/pone"
	ksearch "github.com/kelindar/search"

	"github.com/joho/godotenv"
)

func main() {
	// init slog
	opts := &slog.HandlerOptions{
		ReplaceAttr: func(groups []string, a slog.Attr) slog.Attr {
			if a.Key == slog.TimeKey {
				return slog.String(a.Key, a.Value.Time().Format("2006-01-02 15:04:05"))
			}
			return a
		},
	}
	logger := slog.New(slog.NewTextHandler(os.Stdout, opts))

	// load .env
	if err := godotenv.Load(); err != nil {
		logger.Warn("未找到 .env 配置文件")
	}

	indexPath := os.Getenv("PQ_INDEX_FILE")
	if indexPath == "" {
		logger.Error("missing PQ_INDEX_FILE")
		os.Exit(1)
	}

	episodeFile := os.Getenv("PQ_EPISODEITEM_FILE")
	if episodeFile == "" {
		logger.Error("missing PQ_EPISODEITEM_FILE")
		os.Exit(1)
	}

	// init local embed client
	modelPath := os.Getenv("PQ_EMBEDDING_MODEL")
	if modelPath == "" {
		logger.Error("missing PQ_EMBEDDING_MODEL")
		os.Exit(1)
	}
	embClient, err := embed.NewClient(modelPath, "")
	if err != nil {
		logger.Error("加载本地模型失败", "err", err)
		os.Exit(1)
	}
	defer embClient.Close()

	raw, err := os.ReadFile(episodeFile)
	if err != nil {
		logger.Error("读取数据源失败", "err", err)
		os.Exit(1)
	}

	var episodes map[int]pone.EpisodeItem
	if err := json.Unmarshal(raw, &episodes); err != nil {
		logger.Error("数据解析失败", "err", err)
		os.Exit(1)
	}

	idx := ksearch.NewIndex[string]()

	ids := make([]int, 0, len(episodes))
	for id := range episodes {
		ids = append(ids, id)
	}
	sort.Ints(ids)

	totalLines := 0
	for _, id := range ids {
		for _, line := range episodes[id].Transcript {
			if strings.TrimSpace(line.Line) != "" {
				totalLines++
			}
		}
	}
	logger.Info("开始执行索引任务", "episodes", len(ids), "lines", totalLines)

	processed := 0
	for _, episodeID := range ids {
		ep := episodes[episodeID]
		for i, line := range ep.Transcript {
			text := strings.TrimSpace(line.Line)
			if text == "" {
				continue
			}

			processed++
			lineID := pone.MakeLineID(episodeID, i+1)
			l := logger.With("line_id", lineID, "at", fmt.Sprintf("%d/%d", processed, totalLines))

			l.Info("请求向量")
			vec, err := embClient.GetVector(text, false)
			if err != nil {
				l.Warn("向量获取失败", "err", err)
				continue
			}

			idx.Add(vec, strconv.Itoa(lineID))
		}
	}

	logger.Info("开始写入索引", "path", indexPath)
	if err := idx.WriteFile(indexPath); err != nil {
		logger.Error("序列化文件失败", "err", err)
		os.Exit(1)
	}

	logger.Info("任务处理完毕")
}
