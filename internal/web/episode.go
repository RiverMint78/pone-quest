package web

import (
	"net/http"
	"strconv"
	"time"
)

type EpisodeTranscriptLineView struct {
	LineID       int
	SeqInEpisode int
	Character    string
	CharacterCSS string
	Text         string
	Highlighted  bool
}

type EpisodeTranscriptView struct {
	EpisodeID     int
	EpisodeCode   string
	EpisodeTitle  string
	HighlightLine int
	Lines         []EpisodeTranscriptLineView
}

func (h *Handler) handleEpisode(w http.ResponseWriter, r *http.Request) {
	start := time.Now()

	episodeID, err := strconv.Atoi(r.URL.Query().Get("episodeID"))
	if err != nil || episodeID <= 0 {
		http.Error(w, "invalid episodeID", http.StatusBadRequest)
		return
	}

	highlightLineID, _ := strconv.Atoi(r.URL.Query().Get("lineID"))

	ep, ok := h.Store.GetEpisode(episodeID)
	if !ok {
		http.Error(w, "episode not found", http.StatusNotFound)
		return
	}

	lines := h.Store.GetEpisodeLines(episodeID)
	items := make([]EpisodeTranscriptLineView, 0, len(lines))
	for _, line := range lines {
		items = append(items, EpisodeTranscriptLineView{
			LineID:       line.LineID,
			SeqInEpisode: line.SeqInEpisode,
			Character:    line.Character,
			CharacterCSS: characterColorClass(line.Character),
			Text:         line.Text,
			Highlighted:  line.LineID == highlightLineID,
		})
	}

	view := EpisodeTranscriptView{
		EpisodeID:     episodeID,
		EpisodeCode:   episodeIDToSxxEyy(episodeID),
		EpisodeTitle:  ep.Title,
		HighlightLine: highlightLineID,
		Lines:         items,
	}

	h.Logger.Debug("剧集全文查看",
		"episode_id", episodeID,
		"line_count", len(items),
		"highlight_line", highlightLineID,
		"elapsed", time.Since(start),
	)

	h.render(w, r, http.StatusOK, "index.tmpl", "episode-viewer.tmpl", view)
}
