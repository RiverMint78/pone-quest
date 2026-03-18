// Package search 用于进行向量搜索
package search

import (
	"math"
	"sort"

	"github.com/RiverMint78/pone-quest/internal/models"
)

// Engine 是内存搜索引擎
type Engine struct {
	Items []models.ImageItem
}

// NewEngine 创建并初始化引擎
func NewEngine(items []models.ImageItem) *Engine {
	return &Engine{Items: items}
}

// Search 执行向量相似度搜索并返回前 K 个结果
func (e *Engine) Search(queryVec []float32, topK int) []models.ImageItem {
	type scoreItem struct {
		item  models.ImageItem
		score float32
	}

	scores := make([]scoreItem, len(e.Items))

	// 直接计算查询向量与库中每个向量的余弦相似度
	for i, item := range e.Items {
		scores[i] = scoreItem{
			item:  item,
			score: cosineSimilarity(queryVec, item.Vector),
		}
	}

	// 相似度高到低排序
	sort.Slice(scores, func(i, j int) bool {
		return scores[i].score > scores[j].score
	})

	// 提取前 K 个结果
	resultLen := topK
	if len(scores) < topK {
		resultLen = len(scores)
	}

	results := make([]models.ImageItem, resultLen)
	for i := 0; i < resultLen; i++ {
		results[i] = scores[i].item
	}

	return results
}

// cosineSimilarity 计算两个向量的余弦相似度
func cosineSimilarity(a, b []float32) float32 {
	if len(a) != len(b) {
		return 0
	}
	var dotProduct, normA, normB float64
	for i := range a {
		dotProduct += float64(a[i]) * float64(b[i])
		normA += float64(a[i]) * float64(a[i])
		normB += float64(b[i]) * float64(b[i])
	}
	if normA == 0 || normB == 0 {
		return 0
	}
	return float32(dotProduct / (math.Sqrt(normA) * math.Sqrt(normB)))
}
