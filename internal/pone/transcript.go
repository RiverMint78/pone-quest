// Package pone 定义了 pone-quest 项目的核心数据结构
package pone

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

// MakeLineID 生成台词主键，eg. 230012 表示第23集的第12行。
func MakeLineID(episodeID int, seqInEpisode int) int {
	return episodeID*10000 + seqInEpisode
}
