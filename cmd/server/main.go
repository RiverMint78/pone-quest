package main

import (
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
	// init slog and env
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
	_ = godotenv.Load()

	// db init
	db, err := store.New("pone.db")
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
			os.Getenv("EMBEDDING_API_URL"),
			os.Getenv("EMBEDDING_API_KEY"),
			os.Getenv("EMBEDDING_MODEL"),
		),
		TemplateCache: templateCache,
	}

	// router reg
	mux := http.NewServeMux()
	h.RegisterRoutes(mux)

	// server launch
	srv := &http.Server{
		Addr:         ":8080",
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
