// Package handlers 路由和并处理任何请求
package handlers

import (
	"context"
	"html/template"
	"io/fs"
	"log/slog"
	"net/http"
	"path/filepath"
	"time"

	"github.com/RiverMint78/pone-quest/internal/embed"
	"github.com/RiverMint78/pone-quest/internal/search"
	"github.com/RiverMint78/pone-quest/ui"
)

type Handler struct {
	Engine        *search.Engine
	Embed         *embed.Client
	TemplateCache map[string]*template.Template
}

// NewTemplateCache 生成渲染cache
func NewTemplateCache() (map[string]*template.Template, error) {
	cache := make(map[string]*template.Template)

	pages, err := fs.Glob(ui.Files, "html/pages/*.tmpl")
	if err != nil {
		return nil, err
	}

	for _, page := range pages {
		name := filepath.Base(page)

		patterns := []string{
			"html/base.tmpl",
			"html/partials/*.tmpl",
			page,
		}

		ts, err := template.New(name).ParseFS(ui.Files, patterns...)
		if err != nil {
			return nil, err
		}

		cache[name] = ts
	}

	return cache, nil
}

func (h *Handler) RegisterRoutes(mux *http.ServeMux) {
	fileServer := http.FileServer(http.FS(ui.Files))
	mux.Handle("GET /static/", fileServer)

	mux.HandleFunc("GET /", h.handleIndex)
	mux.HandleFunc("GET /search", h.handleSearch)
}

func (h *Handler) handleIndex(w http.ResponseWriter, r *http.Request) {
	ts, ok := h.TemplateCache["index.tmpl"]
	if !ok {
		http.Error(w, "模板不存在", 500)
		return
	}
	ts.ExecuteTemplate(w, "base.tmpl", nil)
}

func (h *Handler) handleSearch(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query().Get("q")
	if query == "" {
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), 10*time.Second)
	defer cancel()

	vec, err := h.Embed.GetVector(ctx, query)
	if err != nil {
		slog.Error("向量化失败", "err", err)
		return
	}

	results := h.Engine.Search(vec, 12)

	// 渲染局部组件
	ts, ok := h.TemplateCache["index.tmpl"]
	if !ok {
		return
	}

	ts.ExecuteTemplate(w, "results.tmpl", results)
}
