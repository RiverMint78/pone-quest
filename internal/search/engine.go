// Package search 用于进行向量搜索
package search

import (
	ksearch "github.com/kelindar/search"
)

// Engine 包装泛型向量搜索索引
type Engine struct {
	index *ksearch.Index[string]
}

// SearchResult 包装搜索结果
type SearchResult struct {
	ID    string
	Score float64
}

// NewEngine 创建并初始化引擎
func NewEngine(idx *ksearch.Index[string]) *Engine {
	return &Engine{index: idx}
}

// Search 执行向量相似度搜索并返回前 K 个结果
func (e *Engine) Search(queryVec []float32, topK int) []SearchResult {
	results := e.index.Search(queryVec, topK)

	// 结果提取
	out := make([]SearchResult, len(results))
	for i := range results {
		out[i].ID = results[i].Value
		out[i].Score = results[i].Relevance
	}
	// API 保证是高到低排序

	return out
}
