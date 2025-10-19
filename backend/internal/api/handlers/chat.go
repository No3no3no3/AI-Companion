package handlers

import (
	"github.com/ai-companion/backend/internal/service/chat"
	"github.com/gin-gonic/gin"
)

type ChatHandler struct {
	chatService *chat.Service
}

func NewChatHandler(chatService *chat.Service) *ChatHandler {
	return &ChatHandler{chatService: chatService}
}

func (h *ChatHandler) SendMsg(c *gin.Context) {

}
