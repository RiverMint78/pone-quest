package web

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"
)

const (
	pageBaseSize = 100
	pageMaxSize  = 400
)

// SearchResultsView 是用于结果局部渲染的分页视图模型。
type SearchResultsView struct {
	Items    []SearchViewItem
	Query    string
	NextPage int
	HasMore  bool
}

func pageLimit(page int) int {
	if page < 0 {
		page = 0
	}

	limit := pageBaseSize << page
	if limit > pageMaxSize {
		return pageMaxSize
	}

	return limit
}

func pageOffset(page int) int {
	if page <= 0 {
		return 0
	}

	offset := 0
	for i := 0; i < page; i++ {
		offset += pageLimit(i)
	}

	return offset
}

// SearchViewItem 是用于模板渲染的搜索结果。
type SearchViewItem struct {
	LineID       int
	Score        float32
	Character    string
	CharacterCSS string
	Text         string
	EpisodeID    int
	EpisodeCode  string
	EpisodeTitle string
	SeqInEpisode int
}

// episodeIDToSxxEyy 将全局 episode_id 转为 sXXeYY，不是存在的剧集则返回空字符串。
func episodeIDToSxxEyy(episodeID int) string {
	if episodeID < 1 {
		return ""
	}

	var season int
	var numberInSeason int

	switch {
	case episodeID <= 26:
		season = 1
		numberInSeason = episodeID
	case episodeID <= 52:
		season = 2
		numberInSeason = episodeID - 26
	case episodeID <= 65:
		season = 3
		numberInSeason = episodeID - 52
	default:
		offset := episodeID - 66
		season = 4 + (offset / 26)
		numberInSeason = (offset % 26) + 1
	}

	if season >= 10 {
		return ""
	}

	return fmt.Sprintf("s%02de%02d", season, numberInSeason)
}

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

	page := 0
	if parsed, err := strconv.Atoi(r.URL.Query().Get("page")); err == nil && parsed >= 0 {
		page = parsed
	}

	limit := pageLimit(page)
	offset := pageOffset(page)

	h.Logger.Debug("收到搜索请求",
		"query", rawQuery,
		"page", page,
		"limit", limit,
		"offset", offset,
		"is_htmx", isHTMX,
		"UA", r.UserAgent(),
	)

	if rawQuery == "" {
		h.render(w, r, http.StatusOK, "index.tmpl", targetTemplate, SearchResultsView{})
		return
	}

	const maxRunes = 500
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
	rawResults := h.Engine.Search(vec, limit, offset)
	h.Logger.Debug("索引检索完成",
		"count", len(rawResults),
		"page", page,
		"elapsed", time.Since(searchStart),
	)

	// 按 line_id 查询台词数据。
	results := make([]SearchViewItem, 0, len(rawResults))
	for _, hit := range rawResults {
		line, ok := h.Store.GetLine(int(hit.LineID))
		if !ok {
			h.Logger.Warn("命中结果缺少台词记录", "line_id", hit.LineID)
			continue
		}

		ep, ok := h.Store.GetEpisode(line.EpisodeID)
		title := ""
		if ok {
			title = ep.Title
		}

		results = append(results, SearchViewItem{
			LineID:       line.LineID,
			Score:        hit.Score,
			Character:    line.Character,
			CharacterCSS: characterColorClass(line.Character),
			Text:         line.Text,
			EpisodeID:    line.EpisodeID,
			EpisodeCode:  episodeIDToSxxEyy(line.EpisodeID),
			EpisodeTitle: title,
			SeqInEpisode: line.SeqInEpisode,
		})
	}

	view := SearchResultsView{
		Items:    results,
		Query:    query,
		NextPage: page + 1,
		HasMore:  len(rawResults) == limit,
	}

	// 渲染结果
	h.Logger.Debug("完成搜索和渲染",
		"target", targetTemplate,
		"is_htmx", isHTMX,
		"has_more", view.HasMore,
		"next_page", view.NextPage,
		"total_latency", time.Since(start),
	)

	h.render(w, r, http.StatusOK, "index.tmpl", targetTemplate, view)
}
