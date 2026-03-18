// Package embed 用于请求文本嵌入
package embed

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
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

// Client 封装大模型 API 配置
type Client struct {
	APIURL     string
	APIKey     string
	Model      string
	HTTPClient *http.Client
}

// NewClient 创建一个新的嵌入模型客户端
func NewClient(url, key, model string) *Client {
	return &Client{
		APIURL: url,
		APIKey: key,
		Model:  model,
		HTTPClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

// GetVector 获取图片描述的向量
func (c *Client) GetVector(ctx context.Context, text string) ([]float32, error) {
	if c.APIURL == "" || c.APIKey == "" {
		return nil, fmt.Errorf("client 未配置完全: APIURL 或 APIKey 为空")
	}

	reqBody := EmbeddingRequest{
		Model: c.Model,
		Input: text,
	}

	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return nil, fmt.Errorf("序列化请求体失败: %w", err)
	}

	// 绑定上下文
	req, err := http.NewRequestWithContext(ctx, "POST", c.APIURL, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("创建请求对象失败: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+c.APIKey)

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("执行 HTTP 请求失败: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API 错误: 状态码 %d", resp.StatusCode)
	}

	var apiResp EmbeddingResponse
	if err := json.NewDecoder(resp.Body).Decode(&apiResp); err != nil {
		return nil, fmt.Errorf("解码响应体失败: %w", err)
	}

	if len(apiResp.Data) == 0 || len(apiResp.Data[0].Embedding) == 0 {
		return nil, fmt.Errorf("API 响应数据为空")
	}

	return apiResp.Data[0].Embedding, nil
}
