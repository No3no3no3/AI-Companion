package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/ai-companion/backend/global"
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

func (h *ChatHandler) ChatStream(c *gin.Context) {
	var req chat_domain.Request
	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusBadRequest, common.NewRequestError())
		return
	}
	// 设置SSE响应头
	c.Header("Content-Type", "text/event-stream")
	c.Header("Cache-Control", "no-cache")
	c.Header("Connection", "keep-alive")

	// 获取客户端通知通道
	notify := c.Request.Context().Done()

	var chatRes chat_domain.Response
	chatRes.MessageID = global.UUID.String()
	chatRes.Timestamp = time.Now().Unix()

	stream, err := h.chatService.ProcessStreamMessage(c, &req)
	if err != nil {
		//c.JSON(http.StatusInternalServerError, common.NewInternalError())
		fmt.Printf("Failed to process stream message: %s", err.Error())
		chatRes.Reply = "Failed to process stream message: " + err.Error()
		c.SSEvent("error", chatRes)
		return
	}

	sendSSEEvent(c, "star", chatRes)

	for {
		select {
		case chunk, ok := <-stream:
			if !ok {
				// 流关闭了
				sendSSEEvent(c, "end", nil)
				return
			}
			chatRes.Reply = chunk.Message
			sendSSEEvent(c, "message", chatRes)
		case <-notify:
			return
		}
	}
}

func sendSSEEvent(c *gin.Context, eventType string, data interface{}) {
	var dataStr string
	if data == nil {
		c.SSEvent(eventType, nil)
	}
	marshal, err := json.Marshal(data)
	if err != nil {
		dataStr = fmt.Sprintf(`{"error": "Failed to marshal data: %s"}`, err.Error())
		return
	}
	dataStr = string(marshal)
	c.SSEvent(eventType, dataStr)

}
