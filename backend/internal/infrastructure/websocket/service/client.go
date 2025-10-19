package service

import (
	"context"
	"fmt"
	"runtime/debug"
	"time"

	"github.com/ai-companion/backend/internal/pkg/logger"
	"github.com/gorilla/websocket"
)

const (
	// 用户连接超时时间
	heartbeatExpirationTime = 6 * 60
)

// Client 用户连接
type Client struct {
	Addr          string          // 客户端地址
	Socket        *websocket.Conn // 用户连接
	Send          chan []byte     `json:"-"` // 待发送的数据
	AppID         uint32          // 登录的平台ID app/web/ios
	UserID        string          // 用户ID，用户登录以后才有
	FirstTime     uint64          // 首次连接事件
	HeartbeatTime uint64          // 用户上次心跳时间
	LoginTime     uint64          // 登录时间 登录以后才有
	Device        string
	Action        string
}

// NewClient 初始化 每次请求都将更新心跳时间
func NewClient(addr string, socket *websocket.Conn, firstTime uint64, deviceHeader string, action string) (client *Client) {

	//topicStruct := models2.DeviceParseWebsocket(context.TODO(), deviceHeader)

	client = &Client{
		Addr: addr,
		//UserID:        topicStruct.DeviceID,
		Socket:        socket,
		Send:          make(chan []byte, 100),
		FirstTime:     firstTime,
		HeartbeatTime: firstTime,
		Device:        deviceHeader,
		Action:        action,
	}
	return
}

// 读取客户端数据
func (c *Client) read(ctx context.Context) {
	defer func() {
		if r := recover(); r != nil {
			logger.Info(ctx, "read client message stop", string(debug.Stack()), r)
		}
	}()
	defer func() {
		logger.Info(ctx, "read client message closed", map[string]string{"RemoteAddr": c.Addr})
		close(c.Send)
	}()
	for {
		_, message, err := c.Socket.ReadMessage()

		if err != nil {
			logger.Info(ctx, "read client message error", fmt.Sprintf(`{"client":%s, "error": %s}`, c.Addr, err))
			return
		}
		c.Heartbeat(uint64(time.Now().Unix()))
		// 处理程序
		ProcessData(c, message)
	}
}

// 向客户端写数据
func (c *Client) write(ctx context.Context) {
	defer func() {
		if r := recover(); r != nil {
			logger.Info(ctx, "send data to client stop", string(debug.Stack()), r)
		}
	}()
	defer func() {
		logger.Info(ctx, "send data to client defer", map[string]string{"RemoteAddr": c.Addr})
		WebsocketClientManager.Unregister <- c
		_ = c.Socket.Close()
	}()
	for {
		select {
		case message, ok := <-c.Send:

			if !ok {
				// 发送数据错误 关闭连接
				logger.Info(ctx, "send data to client closed", c.Addr, "ok", ok)
				return
			}
			err := c.Socket.WriteMessage(websocket.TextMessage, message)
			if err != nil {
				logger.Info(ctx, "send data to client error", fmt.Sprintf(`{"client":%s, "error": %s}`, c.Addr, err))
				return
			}
			c.Heartbeat(uint64(time.Now().Unix()))

		}
	}
}

// SendMsg 发送数据
func (c *Client) SendMsg(ctx context.Context, msg []byte) {
	if c == nil {
		return
	}
	defer func() {
		if r := recover(); r != nil {
			logger.Info(ctx, "SendMsg stop", r, string(debug.Stack()))
		}
	}()

	if c.Send == nil {
		logger.Error(ctx, "SendMsg send is nil")
		return
	}

	c.Send <- msg
}

// ResponseJson 向客户端发送JSON格式的响应。
// 参数:
//
//	c *Client: 客户端连接对象，用于发送消息。
//	requestID string: 请求ID，用于跟踪请求和响应。
//	action string: 操作类型，描述响应的操作。
//	code uint32: 响应码，表示响应的状态。
//	msg string: 附加消息，提供更详细的响应信息。
//	data interface{}: 响应数据，可以是任意类型的数据。
func (c *Client) ResponseJson(ctx context.Context, code uint32, msg string, data interface{}) {
	if msg == "" {
		//	msg = object.GetErrorMessage(code, msg)
		msg = "message is empty"
	}
	//requestId := logger.GetRequestId(ctx)
	// 创建响应头对象，封装响应的所有信息。
	//responseHead := models2.NewResponse(requestId, code, msg, data)

	// 将响应头对象序列化为JSON格式的字节切片。
	//headByte, err := json.Marshal(responseHead)
	//if err != nil {
	//	// 如果序列化失败，则记录错误日志，无需返回错误给调用者。
	//	logger.Error(ctx, "Websocket SendJson Marshal error", err)
	//	return
	//}

	// 使用客户端连接对象发送序列化后的响应数据。
	//c.SendMsg(ctx, headByte)
	c.SendMsg(ctx, []byte(msg))
}

func (c *Client) ResponseManage(ctx context.Context, data []byte) {
	// 使用客户端连接对象发送序列化后的响应数据。
	c.SendMsg(ctx, data)
}

func Response(deviceId string, ctx context.Context, data []byte) {
	conn := WebsocketClientManager.Users[deviceId]
	if conn == nil {
		logger.Error(ctx, "wsPage Response conn is nil")
		return
	}
	logger.Info(ctx, "wsPage ResponseSlice success", map[string]string{
		"RemoteAddr": conn.Addr,
	})

	conn.SendMsg(ctx, data)
}

// close 关闭客户端连接
func (c *Client) close() {
	close(c.Send)
}

// Heartbeat 用户心跳
func (c *Client) Heartbeat(currentTime uint64) {
	c.HeartbeatTime = currentTime
	return
}

// IsHeartbeatTimeout 心跳超时
func (c *Client) IsHeartbeatTimeout(currentTime uint64) (timeout bool) {
	if c.HeartbeatTime+heartbeatExpirationTime <= currentTime {
		timeout = true
	}
	return
}
