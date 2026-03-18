package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"os"
	"time"

	"github.com/RiverMint78/pone-quest/internal/embed"
	"github.com/RiverMint78/pone-quest/internal/pone"
	"github.com/RiverMint78/pone-quest/internal/store"

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
	slog.SetDefault(logger)

	// load .env
	if err := godotenv.Load(); err != nil {
		slog.Warn("未找到 .env 配置文件")
	}

	// init embed client
	embClient := embed.NewClient(
		os.Getenv("PG_EMBEDDING_API_URL"),
		os.Getenv("PG_EMBEDDING_API_KEY"),
		os.Getenv("PG_EMBEDDING_MODEL"),
	)

	// open db
	db, err := store.New(os.Getenv("PQ_DATABASE"))
	if err != nil {
		slog.Error("数据库连接失败", "err", err)
		os.Exit(1)
	}
	defer db.Close()

	// read data
	raw, err := os.ReadFile(os.Getenv("PG_IMAGEITEM"))
	if err != nil {
		slog.Error("读取数据源失败", "err", err)
		os.Exit(1)
	}

	var items []pone.ImageItem
	if err := json.Unmarshal(raw, &items); err != nil {
		slog.Error("数据解析失败", "err", err)
		os.Exit(1)
	}

	slog.Info("开始执行索引任务", "count", len(items))

	for i, item := range items {
		l := slog.With("id", item.ID, "at", fmt.Sprintf("%d/%d", i+1, len(items)))

		func() {
			ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
			defer cancel()

			item.ImagePath = "/static/images/" + item.ID

			l.Info("请求向量")
			vec, err := embClient.GetVector(ctx, item.Description)
			if err != nil {
				l.Warn("向量获取失败", "err", err)
				return
			}
			item.Vector = vec

			l.Info("写入DB")
			if err := db.SaveImageItem(item); err != nil {
				l.Error("写入失败", "err", err)
				return
			}
		}()
	}

	slog.Info("🎉 任务处理完毕")
}
