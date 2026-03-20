package main

import (
	"context"
	"encoding/base64"
	"fmt"
	"mime"
	neturl "net/url"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/openai/openai-go"
	"github.com/openai/openai-go/option"
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

func normalizeOpenAIBaseURL(rawURL string) string {
	trimmed := strings.TrimSpace(rawURL)
	if trimmed == "" {
		return ""
	}

	u, err := neturl.Parse(trimmed)
	if err != nil {
		return strings.TrimRight(trimmed, "/")
	}

	cleanPath := strings.TrimRight(u.Path, "/")
	for _, suffix := range []string{"/chat/completions", "/completions"} {
		if strings.HasSuffix(cleanPath, suffix) {
			cleanPath = strings.TrimSuffix(cleanPath, suffix)
			break
		}
	}
	u.Path = strings.TrimRight(cleanPath, "/")
	u.RawQuery = ""
	u.Fragment = ""

	return strings.TrimRight(u.String(), "/")
}

func requestAPILabel(apiURL, key, path string) (string, error) {
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

	clientOpts := []option.RequestOption{
		option.WithAPIKey(key),
		option.WithRequestTimeout(30 * time.Second),
	}
	if baseURL := normalizeOpenAIBaseURL(apiURL); baseURL != "" {
		clientOpts = append(clientOpts, option.WithBaseURL(baseURL))
	}

	client := openai.NewClient(clientOpts...)
	resp, err := client.Chat.Completions.New(context.Background(), openai.ChatCompletionNewParams{
		Model: openai.ChatModel(model),
		Messages: []openai.ChatCompletionMessageParamUnion{
			openai.UserMessage([]openai.ChatCompletionContentPartUnionParam{
				openai.ImageContentPart(openai.ChatCompletionContentPartImageImageURLParam{
					URL:    dataURL,
					Detail: "low",
				}),
				openai.TextContentPart(prompt),
			}),
		},
	})
	if err != nil {
		return "", fmt.Errorf("API 请求失败: %w", err)
	}

	if len(resp.Choices) == 0 {
		return "", fmt.Errorf("API 未返回有效内容")
	}

	description := strings.TrimSpace(resp.Choices[0].Message.Content)
	if description == "" {
		return "", fmt.Errorf("API 返回为空")
	}

	// 清理描述
	description = strings.ReplaceAll(description, "\n", " ")

	return description, nil
}
