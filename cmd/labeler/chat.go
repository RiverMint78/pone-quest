package main

type ContentPart struct {
	Type     string          `json:"type"`
	Text     string          `json:"text,omitempty"`
	ImageURL *ImageURLConfig `json:"image_url,omitempty"`
}

type ImageURLConfig struct {
	URL    string `json:"url"`
	Detail string `json:"detail,omitempty"`
}

type Message struct {
	Role    string        `json:"role"`
	Content []ContentPart `json:"content"`
}

type ChatRequest struct {
	Model    string    `json:"model"`
	Messages []Message `json:"messages"`
}

type ChatResponse struct {
	Choices []struct {
		Message struct {
			Content string `json:"content"`
		} `json:"message"`
	} `json:"choices"`
}
