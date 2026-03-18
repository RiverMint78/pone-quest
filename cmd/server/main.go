package main

import (
	"flag"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/RiverMint78/pone-quest/internal/embed"
	"github.com/RiverMint78/pone-quest/internal/handlers"
	"github.com/RiverMint78/pone-quest/internal/search"
	"github.com/RiverMint78/pone-quest/internal/store"

	"github.com/joho/godotenv"
)

func main() {
	// CLI parse
	port := flag.Int("port", 8080, "HTTP service port")
	flag.Parse()

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

	// db init
	db, err := store.New(os.Getenv("PQ_DATABASE"))
	if err != nil {
		slog.Error("打开数据库失败", "err", err)
		os.Exit(1)
	}
	defer db.Close()

	items, err := db.LoadAll()
	if err != nil {
		slog.Error("加载数据失败", "err", err)
		os.Exit(1)
	}

	// cache tmpl
	templateCache, err := handlers.NewTemplateCache()
	if err != nil {
		slog.Error("模板缓存初始化失败", "err", err)
		os.Exit(1)
	}

	// init handler
	h := &handlers.Handler{
		Engine: search.NewEngine(items),
		Embed: embed.NewClient(
			os.Getenv("PQ_EMBEDDING_API_URL"),
			os.Getenv("PQ_EMBEDDING_API_KEY"),
			os.Getenv("PQ_EMBEDDING_MODEL"),
		),
		TemplateCache: templateCache,
	}

	// router reg
	mux := http.NewServeMux()
	h.RegisterRoutes(mux)

	// server launch
	addr := fmt.Sprintf(":%d", *port)
	srv := &http.Server{
		Addr:         addr,
		Handler:      mux,
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	slog.Info("PoneQuest 启动成功", "addr", srv.Addr)
	if err := srv.ListenAndServe(); err != nil {
		slog.Error("服务器关闭", "err", err)
	}
}
