package common

// Response 统一响应结构
type Response struct {
	Code int         `json:"code"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data,omitempty"`
}

type SSEResponse struct {
	EventType string      `json:"eventType"`
	Data      interface{} `json:"data"`
}

// PageResponse 分页响应结构
type PageResponse struct {
	Code    int         `json:"code"`
	Msg     string      `json:"msg"`
	Data    interface{} `json:"data,omitempty"`
	Page    int         `json:"page"`
	Size    int         `json:"size"`
	Total   int64       `json:"total"`
	HasMore bool        `json:"hasMore"`
}

// 错误码常量
const (
	CodeSuccess       = 0   // 成功
	CodeBadRequest    = 400 // 请求参数错误
	CodeUnauthorized  = 401 // 未授权
	CodeForbidden     = 403 // 禁止访问
	CodeNotFound      = 404 // 资源不存在
	CodeInternalError = 500 // 服务器内部错误
	CodeServiceError  = 502 // 服务错误
	CodeRateLimit     = 429 // 请求过于频繁
)

// 响应消息常量
const (
	MsgSuccess       = "success"
	MsgBadRequest    = "请求参数错误"
	MsgUnauthorized  = "未授权访问"
	MsgForbidden     = "禁止访问"
	MsgNotFound      = "资源不存在"
	MsgInternalError = "服务器内部错误"
	MsgServiceError  = "服务暂时不可用"
	MsgRateLimit     = "请求过于频繁，请稍后重试"
)

// NewSuccess 成功响应
func NewSuccess(data interface{}) *Response {
	return &Response{
		Code: CodeSuccess,
		Msg:  MsgSuccess,
		Data: data,
	}
}

// NewSuccessWithMessage 带自定义成功消息的响应
func NewSuccessWithMessage(msg string, data interface{}) *Response {
	return &Response{
		Code: CodeSuccess,
		Msg:  msg,
		Data: data,
	}
}

// NewError 错误响应
func NewError(code int, msg string) *Response {
	return &Response{
		Code: code,
		Msg:  msg,
	}
}

func NewRequestError() *Response {
	return &Response{
		Code: CodeBadRequest,
		Msg:  MsgBadRequest,
	}
}

func NewInternalError() *Response {
	return &Response{
		Code: CodeInternalError,
		Msg:  MsgInternalError,
	}
}

// NewPageResponse 分页响应
func NewPageResponse(data interface{}, page, size int, total int64) *PageResponse {
	return &PageResponse{
		Code:    CodeSuccess,
		Msg:     MsgSuccess,
		Data:    data,
		Page:    page,
		Size:    size,
		Total:   total,
		HasMore: int64(page*size) < total,
	}
}
