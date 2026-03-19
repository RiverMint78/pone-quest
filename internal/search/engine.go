// Package search 用于进行向量搜索
package search

import (
	ksearch "github.com/kelindar/search"
)

// Engine 包装泛型向量搜索索引
type Engine struct {
	index *ksearch.Index[string]
}

// NewEngine 创建并初始化引擎
func NewEngine(idx *ksearch.Index[string]) *Engine {
	return &Engine{index: idx}
}

// Search 执行向量相似度搜索并返回前 K 个结果
func (e *Engine) Search(queryVec []float32, topK int) []string {
	results := e.index.Search(queryVec, topK)

	// 结果提取
	out := make([]string, len(results))
	for i, r := range results {
		out[i] = r.Value
	}

	return out
}
