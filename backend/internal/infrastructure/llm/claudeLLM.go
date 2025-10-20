package llm

import (
	"context"

	"github.com/ai-companion/backend/internal/pkg/config"
	"github.com/ai-companion/backend/internal/pkg/logger"
	"github.com/tmc/langchaingo/llms/anthropic"
)

type ClaudeLLM struct {
	llm *anthropic.LLM
}

func NewClaudeLLM(cfg *config.LLMConfig) *ClaudeLLM {
	llm, err := anthropic.New(
		anthropic.WithModel(cfg.Model),
		anthropic.WithToken(cfg.Token),
	)
	if err != nil {
		logger.Errorf("new anthropic claude llm error:%s", err.Error())
		return nil
	}
	return &ClaudeLLM{
		llm: llm,
	}
}

func (o *ClaudeLLM) GenerateChat(ctx context.Context, req *ChatRequest) (*ChatResponse, error) {
	return nil, nil
}

func (o *ClaudeLLM) GenerateStream(ctx context.Context, req *ChatRequest) (<-chan *StreamChunk, error) {
	return nil, nil
}

func (o *ClaudeLLM) ValidateConfig() error {
	return nil
}
