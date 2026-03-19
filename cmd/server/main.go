package main

import (
	"flag"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/RiverMint78/pone-quest/internal/embed"
	"github.com/RiverMint78/pone-quest/internal/search"
	"github.com/RiverMint78/pone-quest/internal/web"

	"github.com/joho/godotenv"
	ksearch "github.com/kelindar/search"
)

func main() {
	// CLI parse
	port := flag.Int("port", 8080, "HTTP service port")
	debug := flag.Bool("debug", false, "Enable debug logging")
	flag.Parse()

	// init slog
	level := slog.LevelInfo
	addSource := false
	if *debug {
		level = slog.LevelDebug
		addSource = true
	}
	opts := &slog.HandlerOptions{
		Level:     level,
		AddSource: addSource,
		ReplaceAttr: func(groups []string, a slog.Attr) slog.Attr {
			if a.Key == slog.TimeKey {
				return slog.String(a.Key, a.Value.Time().Format("2006-01-02 15:04:05"))
			}
			if a.Key == slog.SourceKey {
				source := a.Value.Any().(*slog.Source)
				return slog.String("src", fmt.Sprintf("%s:%d", filepath.Base(source.File), source.Line))
			}
			if a.Value.Kind() == slog.KindDuration {
				d := a.Value.Duration()
				ms := float64(d.Nanoseconds()) / 1e6
				return slog.String(a.Key, fmt.Sprintf("%.4fms", ms))
			}
			return a
		},
	}
	logger := slog.New(slog.NewTextHandler(os.Stdout, opts))

	// load .env
	if err := godotenv.Load(); err != nil {
		logger.Warn("未找到 .env 配置文件")
	}

	// cache tmpl
	templateCache, err := web.NewTemplateCache()
	if err != nil {
		logger.Error("模板缓存初始化失败", "err", err)
		os.Exit(1)
	}

	// init local embed
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

	// local embed health check
	testVec, err := embClient.GetVector("thumbs")
	if err != nil {
		logger.Error("无法产生测试向量", "err", err)
		os.Exit(1)
	}

	// init idx
	indexPath := os.Getenv("PQ_INDEX_FILE")
	if indexPath == "" {
		logger.Error("missing PQ_INDEX_FILE")
		os.Exit(1)
	}
	idx := ksearch.NewIndex[string]()
	if err := idx.ReadFile(indexPath); err != nil {
		logger.Error("索引加载失败", "err", err)
		os.Exit(1)
	}
	logger.Info("索引加载完毕", "size", idx.Len())

	// init search
	searchEngine := search.NewEngine(idx)

	// search health check
	testRes := searchEngine.Search(testVec, 5)
	logger.Debug("搜索健康检查", "count", len(testRes), "top_result", testRes)

	// init handler
	h := &web.Handler{
		Engine:        searchEngine,
		Embed:         embClient,
		TemplateCache: templateCache,
		Logger:        logger,
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

	logger.Info("PoneQuest 启动成功", "addr", srv.Addr)
	if err := srv.ListenAndServe(); err != nil {
		logger.Error("服务器关闭", "err", err)
	}
}
