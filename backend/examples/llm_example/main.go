package main

import (
	"context"
	"fmt"

	"github.com/ai-companion/backend/internal/infrastructure/llm"
	"github.com/ai-companion/backend/internal/pkg/config"
)

func main() {
	cfg := config.Load()
	fmt.Println("config: ", cfg.LLM.Model)

	handle := llm.CreateLLM(cfg.LLM.Provider, map[string]string{"model": cfg.LLM.Model, "baseUrl": cfg.LLM.BaseUrl})
	ctx := context.Background()
	req := &llm.ChatRequest{Message: "你好啊 ，今天星期几"}
	res, err := handle.GenerateChat(ctx, req)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(res.Object)

}
