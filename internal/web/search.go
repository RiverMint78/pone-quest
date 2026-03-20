package web

import (
	"net/http"
	"strconv"
	"strings"
	"time"
)

func (h *Handler) handleIndex(w http.ResponseWriter, r *http.Request) {
	h.render(w, r, http.StatusOK, "index.tmpl", "base.tmpl", nil)
}

func (h *Handler) handleSearch(w http.ResponseWriter, r *http.Request) {
	start := time.Now()

	isHTMX := r.Header.Get("HX-Request") != ""
	targetTemplate := "base.tmpl"
	if isHTMX {
		targetTemplate = "results.tmpl"
	}

	rawQuery := strings.TrimSpace(r.URL.Query().Get("q"))

	// 解析topK
	topKStr := r.URL.Query().Get("topK")
	topK, err := strconv.Atoi(topKStr)
	if err != nil {
		topK = 10
	}
	topK = max(5, min(50, topK))

	h.Logger.Debug("收到搜索请求",
		"query", rawQuery,
		"topK", topK,
		"is_htmx", isHTMX,
		"UA", r.UserAgent(),
	)

	if rawQuery == "" {
		h.render(w, r, http.StatusOK, "index.tmpl", targetTemplate, nil)
		return
	}

	const maxRunes = 200
	query := rawQuery
	runeCount := 0

	for i := range rawQuery {
		// loop for utf-8
		if runeCount == maxRunes {
			query = rawQuery[:i]
			h.Logger.Warn("截断超长输入", "original_bytes", len(rawQuery))
			break
		}
		runeCount++
	}

	// 本地向量化
	embedStart := time.Now()
	vec, err := h.Embed.GetVector(query, true)
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
	results := h.Engine.Search(vec, topK)
	h.Logger.Debug("索引检索完成",
		"count", len(results),
		"elapsed", time.Since(searchStart),
	)

	// 渲染结果
	h.Logger.Debug("完成搜索和渲染",
		"target", targetTemplate,
		"is_htmx", isHTMX,
		"total_latency", time.Since(start),
	)

	h.render(w, r, http.StatusOK, "index.tmpl", targetTemplate, results)
}
