package llm

import (
	"slices"

	"github.com/ai-companion/backend/internal/pkg/config"
	"github.com/ai-companion/backend/internal/pkg/logger"
)

var openAIMap = []string{"openai_compatible_llm",
	"openai_llm", "gemini_llm", "zhipu_llm", "deepseek_llm",
	"groq_llm", "mistral_llm", "lmstudio_llm"}

func CreateLLM(cfg *config.LLMConfig) Handle {
	logger.Info("initialize llm")
	if slices.Contains(openAIMap, cfg.Provider) {
		return NewOpenAILLM(cfg)
	}
	// 通过Http发送请求
	if cfg.Provider == "stateless_llm_with_template" {
		return nil
	}
	if cfg.Provider == "ollama_llm" {
		return NewOllamaLLM(cfg)
	}
	if cfg.Provider == "claude_llm" {
		return NewClaudeLLM(cfg)
	}
	logger.Errorf("unsupported llm provider:%s", cfg.Provider)
	return nil

}
