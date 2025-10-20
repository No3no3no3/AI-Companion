package chat_domain

// Request 聊天请求结构
type Request struct {
	Message string `json:"message" binding:"required"`
	UserID  string `json:"userId,omitempty"`
}

// Response 聊天响应结构
type Response struct {
	Reply     string `json:"reply"`
	MessageID string `json:"messageId"`
	Timestamp int64  `json:"timestamp"`
}
