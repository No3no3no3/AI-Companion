package service

import (
	"context"
	"fmt"
	"runtime/debug"
	"sync"

	"github.com/ai-companion/backend/internal/pkg/logger"
)

// DisposeFunc 处理函数
type DisposeFunc func(client *Client, ctx context.Context, message []byte)

var (
	handlers        = make(map[string]DisposeFunc)
	handlersRWMutex sync.RWMutex
)

// Register 注册
func Register(key string, value DisposeFunc) {
	handlersRWMutex.Lock()
	defer handlersRWMutex.Unlock()
	handlers[key] = value
	fmt.Println("websocket register success", handlers)
	return
}

// InitWSCtx 初始化WebSocket请求上下文。
// 该函数根据客户端信息和请求信息创建一个全局请求数据对象，并将其放入上下文中。
//
// 参数:
//
//	client *Client: 客户端连接信息，用于获取客户端IP。
//	req models2.WSRequestStruct: WebSocket请求结构体，包含请求ID和动作。
//
// 返回值:
//
//	ctx context.Context: 包含请求数据的上下文，用于在不同层级的处理中传递请求相关信息。
//func InitWSCtx(ctx context.Context, client *Client) context.Context {
//	// 创建全局请求数据对象，填充请求ID、客户端IP、请求动作和时间戳。
//	data := map[string]string{
//		"RequestID":   uuid.NewString(),
//		"ClientIP":    client.Socket.RemoteAddr().String(),
//		"Url":         client.Device,
//		"RequestTime": time.Now().String(),
//	}
//	// 使用RequestDataKey将全局请求数据对象放入上下文中。
//	return context.WithValue(ctx, "requestData", data)
//}

func getHandlers(key string) (value DisposeFunc, ok bool) {
	handlersRWMutex.RLock()
	defer handlersRWMutex.RUnlock()
	value, ok = handlers[key]
	return
}

// ProcessData 处理数据
func ProcessData(client *Client, message []byte) {

	ctx := context.Background()

	//ctx = InitWSCtx(ctx, client)

	defer func() {
		if r := recover(); r != nil {
			logger.Info(ctx, "Process client data stop", r, string(debug.Stack()))
		}
	}()

	// requestParams
	//logger.Info(ctx, "Process client data params", message)

	// 采用 map 注册的方式
	if value, ok := getHandlers(""); ok {
		value(client, ctx, message)
	} else {
		logger.Error(ctx, "getHandlers error router not exist", "")
	}
	return
}
