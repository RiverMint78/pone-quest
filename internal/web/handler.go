// Package web 负责网络层
package web

import (
	"fmt"
	"html/template"
	"log/slog"

	"github.com/RiverMint78/pone-quest/internal/embed"
	"github.com/RiverMint78/pone-quest/internal/search"
)

// Handler 封装了 Web 层所需的所有依赖
type Handler struct {
	Engine        *search.Engine
	Embed         *embed.Client
	TemplateCache map[string]*template.Template
	Logger        *slog.Logger
}

// New 创建并初始化 Web 处理程序
func New(engine *search.Engine, embed *embed.Client, logger *slog.Logger) (*Handler, error) {
	switch {
	case engine == nil:
		return nil, fmt.Errorf("missing search engine")
	case embed == nil:
		return nil, fmt.Errorf("missing embed client")
	case logger == nil:
		return nil, fmt.Errorf("missing logger")
	}
	cache, err := NewTemplateCache()
	if err != nil {
		return nil, err
	}

	return &Handler{
		Engine:        engine,
		Embed:         embed,
		TemplateCache: cache,
		Logger:        logger,
	}, nil
}
