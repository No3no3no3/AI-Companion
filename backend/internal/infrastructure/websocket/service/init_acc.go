package service

import (
	"fmt"
	"net/http"
	"time"

	"github.com/ai-companion/backend/internal/pkg/config"
	"github.com/ai-companion/backend/internal/pkg/logger"
	"github.com/ai-companion/backend/internal/pkg/utils"
	"github.com/gorilla/websocket"
)

const WsReadBufferSize = 1024 * 4
const WsWriteBufferSize = 1024 * 6
const WsMaxSize = 1024 * 480

/*

  这个文件是整个WebSocket模块的入口点，负责：

  1. 服务接入: 提供WebSocket连接的接入点
  2. 连接管理: 处理客户端连接的建立和初始化
  3. 协议处理: 实现HTTP到WebSocket的协议升级
  4. 并发基础: 为后续的语音处理和设备管理提供连接基础


*/

// WebsocketClientManager 管理者
// serverIp 服务器IP
var (
	WebsocketClientManager = NewClientManager() // 管理者
	serverIp               string
)

// StartWebSocket 启动WebSocket服务
// - 获取服务器IP地址和端口配置
// - 设置HTTP路由处理器
// - 启动客户端管理器
// - 监听指定端口启动WebSocket服务
func StartWebSocket() {
	serverIp = utils.GetServerIp()
	fmt.Println("serve  external ip :", serverIp)
	webSocketPort := config.GetString("app.webSocketPort")
	http.HandleFunc("/", wsPage)

	// 添加处理程序
	go WebsocketClientManager.start()
	fmt.Println("StartWebSocket success.")
	fmt.Println(fmt.Sprintf(`{"serverIp":%s,  "webSocketPort":%s, "rpcPort":%s}`, serverIp, webSocketPort))
	_ = http.ListenAndServe(":"+webSocketPort, nil)
}

// wsPage 处理WebSocket连接请求
// - 协议升级: 将HTTP连接升级为WebSocket连接
// - 连接配置: 设置读写缓冲区大小限制
// - 读缓冲区：4KB
// - 写缓冲区：6KB
// - 最大消息大小：480KB
// - 客户端创建: 根据连接信息创建客户端实例
// - 并发处理: 为每个客户端启动读写goroutine
func wsPage(w http.ResponseWriter, req *http.Request) {
	var upgrade = websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			logger.Info("wsPage upgrade")
			return true
		},
		// 设置读写缓冲区大小（单位：字节）
		ReadBufferSize:  WsReadBufferSize,
		WriteBufferSize: WsWriteBufferSize,
	}

	// 在处理函数中：
	conn, err := upgrade.Upgrade(w, req, nil)
	if err != nil {
		logger.Warn("wsPage upgrade error", map[string]string{"Error": err.Error()})
		http.NotFound(w, req)
		return
	}

	conn.SetReadLimit(WsMaxSize)
	logger.Info(req.Context(), "wsPage connect success", map[string]string{
		"RemoteAddr": conn.RemoteAddr().String(),
	})
	currentTime := uint64(time.Now().Unix())

	userID := req.Header.Get("userID")

	client := NewClient(conn.RemoteAddr().String(), conn, currentTime, userID)
	go client.read(req.Context())
	go client.write(req.Context())

	// 用户连接事件
	WebsocketClientManager.Register <- client
}
