// Package embed 用于请求文本嵌入
package embed

import (
	"fmt"

	lru "github.com/hashicorp/golang-lru/v2"
	ksearch "github.com/kelindar/search"
)

// EmbeddingRequest 定义标准的 API 请求格式
type EmbeddingRequest struct {
	Model string `json:"model"`
	Input string `json:"input"`
}

// EmbeddingResponse 定义所需响应字段
type EmbeddingResponse struct {
	Data []struct {
		Embedding []float32 `json:"embedding"`
	} `json:"data"`
}

// Client 封装本地 Vectorizer
type Client struct {
	vectorizer *ksearch.Vectorizer
	prefix     string
	cache      *lru.Cache[string, []float32]
}

// NewClient 加载本地 GGUF 模型并预设指令
func NewClient(modelPath string, instruction string) (*Client, error) {
	// 不考虑 GPU 推理
	vectorizer, err := ksearch.NewVectorizer(modelPath, 0)
	if err != nil {
		return nil, fmt.Errorf("init model error: %w", err)
	}
	cache, err := lru.New[string, []float32](4096)
	if err != nil {
		return nil, fmt.Errorf("init cache error: %w", err)
	}

	return &Client{
		vectorizer: vectorizer,
		prefix:     instruction,
		cache:      cache,
	}, nil
}

// Close 释放 Vectorizer 资源
func (c *Client) Close() {
	if c.vectorizer != nil {
		c.vectorizer.Close()
	}
}

// GetVector 获取文本描述的向量
func (c *Client) GetVector(text string, isQuery bool) ([]float32, error) {
	input := text
	if isQuery && c.prefix != "" {
		input = c.prefix + text
	}

	// cache read
	if vec, ok := c.cache.Get(input); ok {
		return vec, nil
	}

	// cache miss
	vec, err := c.vectorizer.EmbedText(input)
	if err != nil {
		return nil, err
	}

	// cache add
	c.cache.Add(input, vec)

	return vec, nil
}
