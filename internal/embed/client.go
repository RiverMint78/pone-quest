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
}

// NewClient 加载本地 GGUF 模型
// modelPath: 模型文件路径
// gpuLayers: GPU 加速层数 (0 表示纯 CPU)
func NewClient(modelPath string, gpuLayers int) (*Client, error) {
	m, err := ksearch.NewVectorizer(modelPath, gpuLayers)
	if err != nil {
		return nil, fmt.Errorf("初始化模型 %s 失败: %w", modelPath, err)
	}

	return &Client{
		vectorizer: m,
	}, nil
}

// Close 释放 Vectorizer 资源
func (c *Client) Close() {
	if c.vectorizer != nil {
		c.vectorizer.Close()
	}
}

// GetVector 获取文本描述的向量
func (c *Client) GetVector(text string) ([]float32, error) {
	vec, err := c.vectorizer.EmbedText(text)
	if err != nil {
		return nil, fmt.Errorf("生成文本向量失败: %w", err)
	}
	return vec, nil
}
