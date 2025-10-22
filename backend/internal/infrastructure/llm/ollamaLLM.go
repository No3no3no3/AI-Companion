package llm

import (
	"context"

	"github.com/ai-companion/backend/internal/pkg/config"
	"github.com/ai-companion/backend/internal/pkg/logger"
	"github.com/tmc/langchaingo/llms"
	"github.com/tmc/langchaingo/llms/ollama"
)

type OllamaLLM struct {
	llm *ollama.LLM
}

func NewOllamaLLM(cfg *config.LLMConfig) *OllamaLLM {
	llm, err := ollama.New(
		ollama.WithModel(cfg.Model),
	)
	if err != nil {
		logger.Errorf("new ollama llm error:%s", err.Error())
		return nil
	}
	return &OllamaLLM{
		llm: llm,
	}
}

func (o *OllamaLLM) GenerateChat(ctx context.Context, req *ChatRequest) (*ChatResponse, error) {
	res, err := llms.GenerateFromSinglePrompt(ctx, o.llm, req.Message)
	if err != nil {
		logger.Errorf("generate chat_domain error : %s", err.Error())
		return nil, err
	}
	return &ChatResponse{Object: res}, nil
}

func (o *OllamaLLM) GenerateStream(ctx context.Context, req *ChatRequest) (<-chan *StreamChunk, error) {
	resChan := make(chan *StreamChunk, 10)
	go func() {
		defer close(resChan) // 确保 channel 在结束时关闭

		_, err := o.llm.GenerateContent(
			ctx,
			[]llms.MessageContent{
				llms.TextParts(llms.ChatMessageTypeSystem, "你是一个非常有用的助理"),
				llms.TextParts(llms.ChatMessageTypeHuman, req.Message), // 使用请求中的消息
			},
			// 开启流式输出
			llms.WithStreamingFunc(func(_ context.Context, chunk []byte) error {
				select {
				case resChan <- &StreamChunk{
					Message: string(chunk),
				}:
				case <-ctx.Done():
					// 如果上下文被取消，停止发送
					return ctx.Err()
				}
				return nil
			}),
		)

		if err != nil {
			// 发送错误信息到 channel
			select {
			case resChan <- &StreamChunk{
				Error: err,
				Done:  true,
			}:
			case <-ctx.Done():
			}
		}
	}()
	return resChan, nil
}

func (o *OllamaLLM) ValidateConfig() error {
	return nil
}
