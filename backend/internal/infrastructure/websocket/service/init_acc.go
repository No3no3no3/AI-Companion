package service

import (
	"fmt"
	"net/http"
	"time"

	"github.com/ai-companion/backend/internal/infrastructure/websocket/models"
	"github.com/ai-companion/backend/internal/pkg/config"
	"github.com/ai-companion/backend/internal/pkg/logger"
	"github.com/ai-companion/backend/internal/pkg/utils"
	"github.com/gorilla/websocket"
)

const WsReadBufferSize = 1024 * 48
const WsWriteBufferSize = 1024 * 600
const WsMaxSize = 1024 * 480

/*

  这个文件是整个WebSocket模块的入口点，负责：

  1. 服务接入: 提供WebSocket连接的接入点
  2. 连接管理: 处理客户端连接的建立和初始化
  3. 协议处理: 实现HTTP到WebSocket的协议升级
  4. 并发基础: 为后续的语音处理和设备管理提供连接基础


*/

// clientManager 管理者
// appIDs 全部的平台
// serverIp 服务器IP
// serverPort 服务器端口
var (
	WebsocketClientManager = NewClientManager() // 管理者
	serverIp               string
	serverPort             string
)

// GetServer 创建并返回服务器实例
func GetServer() (server *models.Server) {
	server = models.NewServer(serverIp, serverPort)
	return
}

// logMiddleware 是一个中间件函数，用于初始化日志 context
//func logMiddleware(next http.HandlerFunc) http.HandlerFunc {
//	return func(w http.ResponseWriter, r *http.Request) {
//		// 创建请求数据实例
//		data := logger.GlobalRequestDataStruct{
//			// 从context中获取或生成一个新的请求ID
//			RequestID: uuid.NewString(),
//			// 获取客户端IP地址
//			ClientIP: r.Host,
//			RemoteIP: r.RemoteAddr,
//			// 获取请求的URL路径
//			Url: r.RequestURI,
//			// 获取Content-Type请求头的值
//			ContentType: r.Header.Get("Content-Type"),
//			// 获取User-Agent请求头的值
//			UserAgent: r.UserAgent(),
//			// 获取请求内容的长度
//			ContentLength: r.ContentLength,
//			// 获取Accept请求头的值
//			Accept: r.Header.Get("Accept"),
//			// 记录当前时间的纳秒级时间戳
//			RequestTime: time.Now().UnixNano(),
//		}
//
//		// 将请求数据实例存储到context中，以便后续的处理器可以访问
//		ctx := context.WithValue(r.Context(), logger.RequestDataKey, data)
//		r.WithContext(ctx)
//		// 调用下一个处理函数
//		next(w, r)
//	}
//}

// StartWebSocket 启动WebSocket服务
//- 获取服务器IP地址和端口配置
//- 设置HTTP路由处理器
//- 启动客户端管理器
//- 监听指定端口启动WebSocket服务

func StartWebSocket() {
	serverIp = utils.GetServerIp()
	webSocketPort := config.GetString("app.webSocketPort")
	rpcPort := config.GetString("app.rpcPort")
	serverPort = rpcPort
	http.HandleFunc("/", wsPage)

	// 添加处理程序
	go WebsocketClientManager.start()
	fmt.Println("StartWebSocket success.")
	fmt.Println(fmt.Sprintf(`{"serverIp":%s, "serverPort":%s, "webSocketPort":%s, "rpcPort":%s}`, serverIp, serverPort, webSocketPort, rpcPort))
	_ = http.ListenAndServe(":"+webSocketPort, nil)
}

// wsPage 处理WebSocket连接请求
// - 协议升级: 将HTTP连接升级为WebSocket连接
// - 连接配置: 设置读写缓冲区大小限制
// - 读缓冲区：48KB
// - 写缓冲区：600KB
// - 最大消息大小：480KB
// - 客户端创建: 根据连接信息创建客户端实例
// - 并发处理: 为每个客户端启动读写goroutine
func wsPage(w http.ResponseWriter, req *http.Request) {
	var upgrade = websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			logger.Info(r.Context(), "wsPage upgrade")
			return true
		},
		// 设置读写缓冲区大小（单位：字节）
		ReadBufferSize:  WsReadBufferSize,
		WriteBufferSize: WsWriteBufferSize,
	}

	// 在处理函数中：
	conn, err := upgrade.Upgrade(w, req, nil)
	if err != nil {
		logger.Warn(req.Context(), "wsPage upgrade error", map[string]string{"Error": err.Error()})
		http.NotFound(w, req)
		return
	}

	conn.SetReadLimit(WsMaxSize)
	logger.Info(req.Context(), "wsPage connect success", map[string]string{
		"RemoteAddr": conn.RemoteAddr().String(),
	})
	currentTime := uint64(time.Now().Unix())

	deviceHeader := req.Header.Get("Device")
	action := req.Header.Get("Action")

	client := NewClient(conn.RemoteAddr().String(), conn, currentTime, deviceHeader, action)
	go client.read(req.Context())
	go client.write(req.Context())

	// 用户连接事件
	WebsocketClientManager.Register <- client
}
