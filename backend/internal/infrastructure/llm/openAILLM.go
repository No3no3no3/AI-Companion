package llm

import (
	"context"
	"time"

	"github.com/ai-companion/backend/internal/pkg/config"
	"github.com/ai-companion/backend/internal/pkg/logger"
	"github.com/tmc/langchaingo/llms"
	"github.com/tmc/langchaingo/llms/openai"
)

type OpenAILLM struct {
	llm *openai.LLM
}

var openAILLM *OpenAILLM

func NewOpenAILLM(cfg *config.LLMConfig) *OpenAILLM {
	opts := []openai.Option{
		openai.WithToken(cfg.Token),
		openai.WithModel(cfg.Model),
	}
	if cfg.BaseUrl != "" {
		opts = append(opts, openai.WithBaseURL(cfg.BaseUrl))
	}
	llm, err := openai.New(
		opts...,
	//openai.WithOrganization("org-id"), // Organization ID
	//openai.WithAPIVersion("2023-12-01") ,// API version)
	)
	if err != nil {
		logger.Errorf("connection openAI error : %s", err.Error())
		return nil
	}
	return &OpenAILLM{llm: llm}
}

func (o *OpenAILLM) GenerateChat(ctx context.Context, req *ChatRequest) (*ChatResponse, error) {
	res, err := llms.GenerateFromSinglePrompt(ctx, o.llm, req.Message)
	if err != nil {
		return nil, err
	}
	return &ChatResponse{Object: res}, nil
}

func (o *OpenAILLM) GenerateStream(ctx context.Context, req *ChatRequest) (<-chan *StreamChunk, error) {
	// 创建带缓冲的 channel，避免阻塞
	resChan := make(chan *StreamChunk, 10)

	// 在 goroutine 中处理流式响应
	go func() {
		defer close(resChan) // 确保 channel 在结束时关闭

		streamCtx, streamCancel := context.WithTimeout(ctx, 30*time.Second)
		defer streamCancel()

		_, err := o.llm.GenerateContent(
			streamCtx,
			[]llms.MessageContent{
				llms.TextParts(llms.ChatMessageTypeSystem, "你是一个非常有用的助理"),
				llms.TextParts(llms.ChatMessageTypeHuman, req.Message),
			},
			// 尝试启用流式输出
			llms.WithStreamingFunc(func(streamCtx context.Context, chunk []byte) error {
				chunkStr := string(chunk)
				if chunkStr != "" {
					select {
					case resChan <- &StreamChunk{
						Message: chunkStr,
					}:
					case <-streamCtx.Done():
						return streamCtx.Err()
					}
				}
				return nil
			}),
		)
		if err != nil {
			// 错误处理
			select {
			case resChan <- &StreamChunk{
				Error: err,
				Done:  true,
			}:
				logger.Error("Failed to create stream: " + err.Error())
			}
		}
	}()

	return resChan, nil
}

func (o *OpenAILLM) ValidateConfig() error {
	return nil
}
