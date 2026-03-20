// Package embed 用于请求文本嵌入
package embed

import (
	"fmt"

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
}

// NewClient 加载本地 GGUF 模型并预设指令
func NewClient(modelPath string, instruction string) (*Client, error) {
	// 不考虑 GPU 推理
	m, err := ksearch.NewVectorizer(modelPath, 0)
	if err != nil {
		return nil, fmt.Errorf("init model error: %w", err)
	}
	return &Client{
		vectorizer: m,
		prefix:     instruction,
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
	return c.vectorizer.EmbedText(input)
}
