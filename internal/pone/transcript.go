// Package pone 定义了 pone-quest 项目的核心数据结构
package pone

import (
	"encoding/json"
	"fmt"
	"os"
	"sort"
	"strings"
)

// LineItem 是原始台词单元，包含角色和台词文本。
type LineItem struct {
	Character string `json:"character"`
	Line      string `json:"line"`
}

// EpisodeItem 是原始单集数据结构，包含元信息和台词列表。
type EpisodeItem struct {
	Title          string     `json:"title"`
	NumberInSeason *int       `json:"number_in_season,omitempty"`
	NumberOverall  *int       `json:"number_overall,omitempty"`
	Writers        string     `json:"writers"`
	Airdate        string     `json:"airdate"`
	URL            string     `json:"url"`
	TranscriptURL  string     `json:"transcript_url"`
	GalleryURL     string     `json:"gallery_url"`
	Transcript     []LineItem `json:"transcript"`
	Season         *int       `json:"season,omitempty"`
	WrittenBy      string     `json:"written_by"`
	StoryBoard     string     `json:"storyboard"`
	Synopsis       string     `json:"synopsis"`
}

// EpisodeMeta 是规范化后的单集元信息。
type EpisodeMeta struct {
	EpisodeID      int    `json:"episode_id"` // e.g. 1, 2, 3...
	Title          string `json:"title"`
	NumberInSeason *int   `json:"number_in_season,omitempty"`
	NumberOverall  *int   `json:"number_overall,omitempty"`
	Writers        string `json:"writers"`
	Airdate        string `json:"airdate"`
	URL            string `json:"url"`
	TranscriptURL  string `json:"transcript_url"`
	GalleryURL     string `json:"gallery_url"`
	Season         *int   `json:"season,omitempty"`
	WrittenBy      string `json:"written_by"`
	StoryBoard     string `json:"storyboard"`
	Synopsis       string `json:"synopsis"`
}

// LineRecord 是规范化后的台词记录。
type LineRecord struct {
	LineID       int    `json:"line_id"`    // e.g. 10012, 2330100
	EpisodeID    int    `json:"episode_id"` // 外键，指向所属剧集
	SeqInEpisode int    `json:"seq_in_episode"`
	Character    string `json:"character"`
	Text         string `json:"text"`
}

// TranscriptStore 是规范化容器。
type TranscriptStore struct {
	Episodes map[int]EpisodeMeta `json:"episodes"`
	Lines    map[int]LineRecord  `json:"lines"`
}

// LoadTranscriptStore 从原始 EpisodeItem JSON 文件构建规范化存储。
func LoadTranscriptStore(path string) (*TranscriptStore, error) {
	raw, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("read transcript source: %w", err)
	}

	var episodes map[int]EpisodeItem
	if err := json.Unmarshal(raw, &episodes); err != nil {
		return nil, fmt.Errorf("parse transcript source: %w", err)
	}

	store := &TranscriptStore{
		Episodes: make(map[int]EpisodeMeta, len(episodes)),
		Lines:    make(map[int]LineRecord),
	}

	for episodeID, ep := range episodes {
		store.Episodes[episodeID] = EpisodeMeta{
			EpisodeID:      episodeID,
			Title:          ep.Title,
			NumberInSeason: ep.NumberInSeason,
			NumberOverall:  ep.NumberOverall,
			Writers:        ep.Writers,
			Airdate:        ep.Airdate,
			URL:            ep.URL,
			TranscriptURL:  ep.TranscriptURL,
			GalleryURL:     ep.GalleryURL,
			Season:         ep.Season,
			WrittenBy:      ep.WrittenBy,
			StoryBoard:     ep.StoryBoard,
			Synopsis:       ep.Synopsis,
		}

		for i, line := range ep.Transcript {
			text := strings.TrimSpace(line.Line)
			if text == "" {
				continue
			}

			seq := i + 1
			lineID := MakeLineID(episodeID, seq)
			store.Lines[lineID] = LineRecord{
				LineID:       lineID,
				EpisodeID:    episodeID,
				SeqInEpisode: seq,
				Character:    line.Character,
				Text:         text,
			}
		}
	}

	return store, nil
}

// GetLine 按主键查找台词记录。
func (s *TranscriptStore) GetLine(lineID int) (LineRecord, bool) {
	if s == nil {
		return LineRecord{}, false
	}
	line, ok := s.Lines[lineID]
	return line, ok
}

// GetEpisode 按主键查找剧集元信息。
func (s *TranscriptStore) GetEpisode(episodeID int) (EpisodeMeta, bool) {
	if s == nil {
		return EpisodeMeta{}, false
	}
	ep, ok := s.Episodes[episodeID]
	return ep, ok
}

// GetEpisodeLines 返回某一集所有台词，按集内顺序升序排列。
func (s *TranscriptStore) GetEpisodeLines(episodeID int) []LineRecord {
	if s == nil {
		return nil
	}

	lines := make([]LineRecord, 0)
	for _, line := range s.Lines {
		if line.EpisodeID == episodeID {
			lines = append(lines, line)
		}
	}

	sort.Slice(lines, func(i, j int) bool {
		return lines[i].SeqInEpisode < lines[j].SeqInEpisode
	})

	return lines
}

// MakeLineID 生成台词主键，eg. 230012 表示第23集的第12行。
func MakeLineID(episodeID int, seqInEpisode int) int {
	return episodeID*10000 + seqInEpisode
}
