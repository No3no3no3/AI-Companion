package handlers

import (
	"net/http"

	"github.com/ai-companion/backend/internal/common"
	"github.com/ai-companion/backend/internal/domain/chat_domain"
	"github.com/ai-companion/backend/internal/service/chat"
	"github.com/gin-gonic/gin"
)

type ChatHandler struct {
	chatService *chat.Service
}

func NewChatHandler(chatService *chat.Service) *ChatHandler {
	return &ChatHandler{chatService: chatService}
}

// Chat 处理聊天消息发送请求
func (h *ChatHandler) Chat(c *gin.Context) {
	var req chat_domain.Request
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, common.NewRequestError())
		return
	}
	reply, err := h.chatService.ProcessMessage(c, &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, common.NewInternalError())
		return
	}
	c.JSON(http.StatusOK, common.NewSuccess(reply))
}
