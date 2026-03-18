package web

import (
	"context"
	"log/slog"
	"net/http"
	"time"
)

func (h *Handler) handleIndex(w http.ResponseWriter, r *http.Request) {
	start := time.Now()
	defer func() {
		h.Logger.DebugContext(r.Context(), "访问首页",
			slog.Duration("latency", time.Since(start)),
			slog.String("ip", r.RemoteAddr),
			slog.String("method", r.Method),
		)
	}()
	h.render(w, r, http.StatusOK, "index.tmpl", "base.tmpl", nil)
}

func (h *Handler) handleSearch(w http.ResponseWriter, r *http.Request) {
	start := time.Now()
	query := r.URL.Query().Get("q")
	h.Logger.DebugContext(r.Context(), "收到搜索请求",
		slog.String("query", query),
		slog.String("user_agent", r.UserAgent()),
	)
	if query == "" {
		h.render(w, r, http.StatusOK, "index.tmpl", "base.tmpl", nil)
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), 10*time.Second)
	defer cancel()

	// 获取向量
	embedStart := time.Now()
	vec, err := h.Embed.GetVector(ctx, query)
	if err != nil {
		h.Logger.ErrorContext(ctx, "向量化失败",
			"query", query,
			"err", err,
			"elapsed", time.Since(embedStart),
		)
		http.Error(w, "搜索暂时不可用", http.StatusInternalServerError)
		return
	}
	h.Logger.DebugContext(ctx, "向量化成功", "elapsed", time.Since(embedStart))

	// 搜索
	searchStart := time.Now()
	results := h.Engine.Search(vec, 12)
	h.Logger.DebugContext(ctx, "索引检索完成",
		"count", len(results),
		"elapsed", time.Since(searchStart),
	)

	// 渲染
	targetTemplate := "base.tmpl"
	isHTMX := r.Header.Get("HX-Request") != ""
	if isHTMX {
		targetTemplate = "results.tmpl"
	}
	h.Logger.DebugContext(ctx, "准备渲染模板",
		"target", targetTemplate,
		"is_htmx", isHTMX,
		"total_latency", time.Since(start),
	)

	h.render(w, r, http.StatusOK, "index.tmpl", targetTemplate, results)
}
