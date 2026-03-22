package web

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/RiverMint78/pone-quest/internal/pone"
)

const (
	pageBaseSize = 100
	pageMaxParam = 999
)

// SearchResultsView 是用于结果局部渲染的分页视图模型。
type SearchResultsView struct {
	Items    []SearchViewItem
	Query    string
	NextPage int
	HasMore  bool
	Error    string
}

func pageOffset(page int) int {
	return max(page*pageBaseSize, 0)
}

func normalizePage(raw string) int {
	parsed, err := strconv.Atoi(raw)
	if err != nil || parsed <= 0 {
		return 0
	}
	return min(parsed, pageMaxParam)
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
	h.render(w, r, http.StatusOK, "index.tmpl", "base.tmpl", SearchResultsView{})
}

func (h *Handler) handleSearch(w http.ResponseWriter, r *http.Request) {
	start := time.Now()

	isHTMX := r.Header.Get("HX-Request") != ""
	targetTemplate := "base.tmpl"
	if isHTMX {
		targetTemplate = "results.tmpl"
	}

	rawQuery := strings.TrimSpace(r.URL.Query().Get("q"))

	page := normalizePage(r.URL.Query().Get("page"))
	offset := pageOffset(page)

	h.Logger.Debug("收到搜索请求",
		"query", rawQuery,
		"page", page,
		"limit", pageBaseSize,
		"offset", offset,
		"is_htmx", isHTMX,
		"UA", r.UserAgent(),
	)

	if rawQuery == "" {
		h.render(w, r, http.StatusOK, "index.tmpl", targetTemplate, SearchResultsView{Query: ""})
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

	query = pone.NormalizeSearchText(query)

	// 本地向量化
	embedStart := time.Now()
	vec, err := h.Embed.GetVector(query, true)
	if err != nil {
		h.Logger.Error("向量化失败",
			"query_normed", query,
			"err", err,
			"elapsed", time.Since(embedStart),
		)
		h.render(w, r, http.StatusOK, "index.tmpl", targetTemplate, SearchResultsView{
			Query: query,
			Error: "搜索服务暂时不可用，请稍后重试。",
		})
		return
	}
	h.Logger.Debug("向量化成功", "elapsed", time.Since(embedStart))

	// 索引检索
	searchStart := time.Now()
	rawResults := h.Engine.Search(vec, pageBaseSize, offset)
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
		HasMore:  len(rawResults) == pageBaseSize,
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
