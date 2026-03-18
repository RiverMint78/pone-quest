package web

import (
	"context"
	"log/slog"
	"net/http"
	"time"
)

func (h *Handler) handleIndex(w http.ResponseWriter, r *http.Request) {
	h.render(w, r, http.StatusOK, "index.tmpl", "base.tmpl", nil)
}

func (h *Handler) handleSearch(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query().Get("q")
	if query == "" {
		h.render(w, r, http.StatusOK, "index.tmpl", "base.tmpl", nil)
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), 10*time.Second)
	defer cancel()

	// 获取向量
	vec, err := h.Embed.GetVector(ctx, query)
	if err != nil {
		slog.Error("向量化失败", "query", query, "err", err)
		http.Error(w, "搜索暂时不可用", http.StatusInternalServerError)
		return
	}

	// 搜索
	results := h.Engine.Search(vec, 12)

	// 渲染
	targetTemplate := "base.tmpl"
	if r.Header.Get("HX-Request") != "" {
		targetTemplate = "results.tmpl"
	}

	h.render(w, r, http.StatusOK, "index.tmpl", targetTemplate, results)
}
