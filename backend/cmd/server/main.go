package main

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/ai-companion/backend/global"
	"github.com/ai-companion/backend/internal/api/routes"
	"github.com/gin-gonic/gin"
)

func main() {
	// 初始化路由器
	router := gin.Default()

	// 设置路由
	routes.SetupRouters(router)
	// 打印启动信息
	fmt.Printf("🚀 AI Companion Server starting on port %s\n", global.Cfg.Server.Port)
	// 创建HTTP服务器
	srv := &http.Server{
		Addr:    ":" + global.Cfg.Server.Port,
		Handler: router,
	}

	// 启动服务器
	go func() {
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			fmt.Printf("listen: %s\n", err)
		}
	}()

	// 等待中断信号
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	_ = srv.Shutdown(context.Background())
}
