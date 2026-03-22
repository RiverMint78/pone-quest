// Package search 用于进行向量搜索
package search

import (
	"encoding/binary"
	"math"
	"sort"

	ksearch "github.com/kelindar/search"
)

// Engine 包装泛型向量搜索索引
type Engine struct {
	index *ksearch.Index[[]byte]
}

// SearchResult 包装搜索结果
type SearchResult struct {
	LineID int32
	Score  float32
}

// NewEngine 创建并初始化引擎
func NewEngine(idx *ksearch.Index[[]byte]) *Engine {
	return &Engine{index: idx}
}

// Search 执行向量相似度搜索，返回 topK 个结果，跳过前 offset 个结果
func (e *Engine) Search(queryVec []float32, topK int, offset int) []SearchResult {
	results := e.index.Search(queryVec, topK+offset)
	if offset > len(results) {
		return nil
	}
	results = results[offset:]

	// 结果提取
	out := make([]SearchResult, 0, len(results))
	for i := range results {
		lineID, ok := decodeLineID(results[i].Value)
		if !ok {
			continue
		}

		out = append(out, SearchResult{
			LineID: lineID,
			Score:  float32(results[i].Relevance),
		})
	}

	// 二级排序：相似度 和 LineID
	sort.SliceStable(out, func(i, j int) bool {
		si := score4(out[i].Score)
		sj := score4(out[j].Score)
		if si == sj {
			return out[i].LineID < out[j].LineID
		}
		return si > sj
	})

	return out
}

func score4(v float32) int32 {
	return int32(math.Round(float64(v) * 10000))
}

// decodeLineID 从原始字节中提取 line_id 整数
func decodeLineID(raw []byte) (int32, bool) {
	if len(raw) != 4 {
		return 0, false
	}

	v := int32(binary.LittleEndian.Uint32(raw))
	if v < 0 {
		return 0, false
	}

	return v, true
}
