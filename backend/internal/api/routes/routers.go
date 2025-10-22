package routes

import (
	"time"

	"github.com/ai-companion/backend/internal/api/handlers"
	"github.com/ai-companion/backend/internal/service/chat"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func SetupRouters(router *gin.Engine) {

	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:5173", "http://127.0.0.1:5173"},
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization", "X-Requested-With"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,           // 若前端带 cookie/凭证需要
		MaxAge:           12 * time.Hour, // 预检缓存
	}))

	// 创建API v1组
	api := router.Group("/api")
	{
		// 创建聊天服务和处理器
		chatService := chat.NewService()
		chatHandler := handlers.NewChatHandler(chatService)

		// 聊天相关路由
		api.POST("/chat", chatHandler.Chat)
		api.GET("/chatStream", chatHandler.ChatStream)
	}

	// 根路径
	router.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "Hello Go",
		})
	})
}
