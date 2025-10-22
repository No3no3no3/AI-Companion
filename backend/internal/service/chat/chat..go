package chat

import (
	"context"
	"time"

	"github.com/ai-companion/backend/global"
	"github.com/ai-companion/backend/internal/domain/chat_domain"
	"github.com/ai-companion/backend/internal/infrastructure/llm"
	"github.com/ai-companion/backend/internal/pkg/logger"
)

type Service struct {
}

var llmHandle = llm.CreateLLM(&global.Cfg.LLM)

// NewService 创建新的聊天服务实例
func NewService() *Service {
	return &Service{}
}

// ProcessMessage 处理用户消息并生成AI回复
func (s *Service) ProcessMessage(c context.Context, req *chat_domain.Request) (*chat_domain.Response, error) {
	ctx, cancel := context.WithTimeout(c, 15*time.Second)
	defer cancel()

	result, err := llmHandle.GenerateChat(ctx, &llm.ChatRequest{
		Message: req.Message,
	})
	if err != nil {
		logger.Error("AI GenerateChat error: %s", err.Error())
		return nil, err
	}

	reply := &chat_domain.Response{
		Reply:     result.Object,
		MessageID: global.UUID.String(),
		Timestamp: time.Now().Unix(),
	}
	return reply, nil
}

// ProcessStreamMessage 流式处理用户消息并生成AI回复
func (s *Service) ProcessStreamMessage(c context.Context, req *chat_domain.Request) (<-chan *llm.StreamChunk, error) {
	return llmHandle.GenerateStream(c, &llm.ChatRequest{
		Message: req.Message,
	})
}
