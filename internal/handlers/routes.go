// Package handlers 处理请求
package handlers

import (
	"context"
	"html/template"
	"log/slog"
	"net/http"
	"time"

	"github.com/RiverMint78/pone-quest/internal/embed"
	"github.com/RiverMint78/pone-quest/internal/search"
)

type Handler struct {
	Engine *search.Engine
	Embed  *embed.Client
	Tmpl   *template.Template
}

// RegisterRoutes 将路由绑定到传入的 ServeMux 上
func (h *Handler) RegisterRoutes(mux *http.ServeMux) {
	// 静态资源
	fs := http.FileServer(http.Dir("ui/static"))
	mux.Handle("GET /static/", http.StripPrefix("/static/", fs))

	// 页面路由
	mux.HandleFunc("GET /", h.handleIndex)
	mux.HandleFunc("GET /search", h.handleSearch)
}

// handleIndex 渲染主页
func (h *Handler) handleIndex(w http.ResponseWriter, r *http.Request) {
	if err := h.Tmpl.ExecuteTemplate(w, "base.tmpl", nil); err != nil {
		slog.Error("渲染模板失败", "err", err)
		http.Error(w, "Internal Server Error", 500)
	}
}

// handleSearch 处理搜索请求
func (h *Handler) handleSearch(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query().Get("q")
	if query == "" {
		return
	}

	// 1. 调用嵌入模型 API
	ctx, cancel := context.WithTimeout(r.Context(), 10*time.Second)
	defer cancel()

	vec, err := h.Embed.GetVector(ctx, query)
	if err != nil {
		slog.Error("向量化查询失败", "query", query, "err", err)
		return
	}

	// 在内存中进行向量检索
	results := h.Engine.Search(vec, 12)

	// 渲染局部 HTML 片段
	if err := h.Tmpl.ExecuteTemplate(w, "results.tmpl", results); err != nil {
		slog.Error("渲染搜索结果失败", "err", err)
	}
}
