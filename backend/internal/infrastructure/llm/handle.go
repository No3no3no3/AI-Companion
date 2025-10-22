package llm

import (
	"context"
	"time"
)

type Handle interface {
	//GenerateChat 生成聊天回复内容
	GenerateChat(ctx context.Context, req *ChatRequest) (*ChatResponse, error)

	//GenerateStream 流式生成聊天恢复内容
	GenerateStream(ctx context.Context, req *ChatRequest) (<-chan *StreamChunk, error)

	//ValidateConfig 验证配置信息
	ValidateConfig() error
}

type ChatRequest struct {
	Message string
}

type ChatResponse struct {
	ID      string
	Object  string
	Created time.Time
	Model   string
	Choices []Choice
	Usage   Usage
}
type Choice struct {
	Index        int
	Message      string
	FinishReason string
	Delta        string
}

type Usage struct {
	PromptTokens     int
	CompletionTokens int
	TotalTokens      int
}

type StreamChunk struct {
	ID      string
	Message string
	Done    bool
	Error   error
}
