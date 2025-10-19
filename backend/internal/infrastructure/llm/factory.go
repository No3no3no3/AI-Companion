package llm

import (
	"slices"

	"github.com/ai-companion/backend/internal/pkg/logger"
)

var openAIMap = []string{"openai_compatible_llm",
	"openai_llm", "gemini_llm", "zhipu_llm", "deepseek_llm",
	"groq_llm", "mistral_llm", "lmstudio_llm"}

func CreateLLM(llmProvider string, cfg map[string]string) Handle {
	logger.Info("initialize llm")
	if slices.Contains(openAIMap, llmProvider) {
		return NewOpenAILLM(cfg)
	}
	// 通过Http发送请求
	if llmProvider == "stateless_llm_with_template" {
		return nil
	}
	if llmProvider == "ollama_llm" {
		return NewOllamaLLM(cfg)
	}
	if llmProvider == "claude_llm" {
		return NewClaudeLLM(cfg)
	}
	logger.Errorf("unsupported llm provider:%s", llmProvider)
	return nil

}
