package web

import (
	"net/http"
	"time"
)

func (h *Handler) handleIndex(w http.ResponseWriter, r *http.Request) {
	h.render(w, r, http.StatusOK, "index.tmpl", "base.tmpl", nil)
}

func (h *Handler) handleSearch(w http.ResponseWriter, r *http.Request) {
	start := time.Now()

	query := r.URL.Query().Get("q")
	h.Logger.Debug("收到搜索请求",
		"query", query,
		"UA", r.UserAgent(),
	)

	if query == "" {
		h.render(w, r, http.StatusOK, "index.tmpl", "base.tmpl", nil)
		return
	}

	// 本地向量化
	embedStart := time.Now()
	vec, err := h.Embed.GetVector(query)
	if err != nil {
		h.Logger.Error("向量化失败",
			"query", query,
			"err", err,
			"elapsed", time.Since(embedStart),
		)
		http.Error(w, "搜索暂时不可用", http.StatusInternalServerError)
		return
	}
	h.Logger.Debug("向量化成功", "elapsed", time.Since(embedStart))

	// 索引检索
	searchStart := time.Now()
	ids := h.Engine.Search(vec, 12)
	h.Logger.Debug("索引检索完成",
		"count", len(ids),
		"elapsed", time.Since(searchStart),
	)

	results := ids

	// // ID -> 渲染对象?
	// // TODO: 添加描述等
	// results := make([]pone.ImageItem, len(ids))
	// for i, id := range ids {
	// 	results[i] = pone.ImageItem{
	// 		ID: id,
	// 	}
	// }

	// 渲染
	targetTemplate := "base.tmpl"
	isHTMX := r.Header.Get("HX-Request") != ""
	if isHTMX {
		targetTemplate = "results.tmpl"
	}

	h.Logger.Debug("完成搜索和渲染",
		"target", targetTemplate,
		"is_htmx", isHTMX,
		"total_latency", time.Since(start),
	)

	h.render(w, r, http.StatusOK, "index.tmpl", targetTemplate, results)
}
