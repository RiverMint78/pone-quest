package main

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"os"

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

	// init local embed client
	modelPath := os.Getenv("PQ_EMBEDDING_MODEL")
	if modelPath == "" {
		logger.Error("missing PQ_EMBEDDING_MODEL")
		os.Exit(1)
	}
	embClient, err := embed.NewClient(modelPath, 0)
	if err != nil {
		logger.Error("加载本地模型失败", "err", err)
		os.Exit(1)
	}
	defer embClient.Close()

	// read image descriptions
	raw, err := os.ReadFile(os.Getenv("PQ_IMAGEITEM_FILE"))
	if err != nil {
		logger.Error("读取数据源失败", "err", err)
		os.Exit(1)
	}

	var items []pone.ImageItem
	if err := json.Unmarshal(raw, &items); err != nil {
		logger.Error("数据解析失败", "err", err)
		os.Exit(1)
	}

	idx := ksearch.NewIndex[string]()
	logger.Info("开始执行索引任务", "count", len(items))

	for i, item := range items {
		l := logger.With("id", item.ID, "at", fmt.Sprintf("%d/%d", i+1, len(items)))

		l.Info("请求向量")
		vec, err := embClient.GetVector(item.Description)
		if err != nil {
			l.Warn("向量获取失败", "err", err)
			continue
		}

		idx.Add(vec, item.ID)
	}

	logger.Info("开始写入索引", "path", indexPath)
	if err := idx.WriteFile(indexPath); err != nil {
		logger.Error("序列化文件失败", "err", err)
		os.Exit(1)
	}

	logger.Info("任务处理完毕")
}
