// Package pone 定义了 pone-quest 项目的核心数据结构
package pone

// ImageItem 是图片元数据记录
type ImageItem struct {
	ID          string    `json:"id"`          // 如 "smile.jpg"
	Description string    `json:"description"` // 图片文本描述，比如 "正在微笑"
	ImagePath   string    `json:"image_path"`  // 图片的 URL 路径，比如 "/static/images/smile.jpg"
	Vector      []float32 `json:"vector"`      // 文本描述向量
}
