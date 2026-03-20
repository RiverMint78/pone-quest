package main

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"mime"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"
)

func isImage(name string) bool {
	ext := strings.ToLower(filepath.Ext(name))
	return ext == ".jpg" || ext == ".jpeg" || ext == ".png" || ext == ".gif" || ext == ".webp"
}

// getMimeType 根据扩展名返回 MIME 类型
func getMimeType(ext string) string {
	t := mime.TypeByExtension(ext)
	if t == "" {
		return mime.TypeByExtension(".jpeg")
	}
	return t
}

func requestAPILabel(url, key, path string) (string, error) {
	// image to BASE64
	data, err := os.ReadFile(path)
	if err != nil {
		return "", fmt.Errorf("读取文件失败: %w", err)
	}

	mimeType := getMimeType(filepath.Ext(path))

	base64Data := base64.StdEncoding.EncodeToString(data)
	dataURL := fmt.Sprintf("data:%s;base64,%s", mimeType, base64Data)

	// 准备请求体
	prompt := os.Getenv("PQ_LABELER_PROMPT")
	model := os.Getenv("PQ_LABELER_MODEL")

	reqBody := ChatRequest{
		Model: model,
		Messages: []Message{
			{
				Role: "user",
				Content: []ContentPart{
					{
						Type: "image_url",
						ImageURL: &ImageURLConfig{
							URL:    dataURL,
							Detail: "low",
						},
					},
					{
						Type: "text",
						Text: prompt,
					},
				},
			},
		},
	}

	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return "", fmt.Errorf("序列化 JSON 失败: %w", err)
	}

	// 发送 HTTP 请求
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return "", fmt.Errorf("创建请求失败: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+key)

	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("网络请求失败: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("API 错误 (状态码 %d): %s", resp.StatusCode, string(body))
	}

	// 解析结果
	var chatResp ChatResponse
	if err := json.NewDecoder(resp.Body).Decode(&chatResp); err != nil {
		return "", fmt.Errorf("解码响应失败: %w", err)
	}

	if len(chatResp.Choices) == 0 {
		return "", fmt.Errorf("API 未返回有效内容")
	}

	// 清理描述
	description := strings.TrimSpace(chatResp.Choices[0].Message.Content)
	description = strings.ReplaceAll(description, "\n", " ")

	return description, nil
}
