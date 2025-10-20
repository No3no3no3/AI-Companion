package service

import (
	"context"
	"sync"
	"time"

	"github.com/ai-companion/backend/internal/pkg/logger"
)

var RegisterDeviceMap = make(map[string]int)

// ClientManager 连接管理
type ClientManager struct {
	Clients     map[*Client]bool   // 全部的连接
	ClientsLock sync.RWMutex       // 读写锁
	Users       map[string]*Client // 登录的用户 // appID+uuid
	UserLock    sync.RWMutex       // 读写锁
	Register    chan *Client       // 连接连接处理
	Unregister  chan *Client       // 断开连接处理程序
	Broadcast   chan []byte        // 广播 向全部成员发送数据
}

// NewClientManager 创建连接管理
func NewClientManager() (clientManager *ClientManager) {
	clientManager = &ClientManager{
		Clients:    make(map[*Client]bool),
		Users:      make(map[string]*Client),
		Register:   make(chan *Client, 1000),
		Unregister: make(chan *Client, 1000),
		Broadcast:  make(chan []byte, 1000),
	}
	return
}

// GetClients 获取所有客户端
func (manager *ClientManager) GetClients() (clients map[*Client]bool) {
	clients = make(map[*Client]bool)
	manager.ClientsRange(func(client *Client, value bool) (result bool) {
		clients[client] = value
		return true
	})
	return
}

// ClientsRange 遍历
func (manager *ClientManager) ClientsRange(f func(client *Client, value bool) (result bool)) {
	manager.ClientsLock.RLock()
	defer manager.ClientsLock.RUnlock()
	for key, value := range manager.Clients {
		result := f(key, value)
		if result == false {
			return
		}
	}
	return
}

// GetClientsLen GetClientsLen
func (manager *ClientManager) GetClientsLen() (clientsLen int) {
	clientsLen = len(manager.Clients)
	return
}

// EventUnRegister 用户关闭连接事件
func (manager *ClientManager) EventUnRegister(client *Client) {
	//	client2.DeleteDeviceOnline(context.TODO(), client.UserID)
	logger.Info(context.Background(), "EventUnRegister", map[string]string{"UserId": client.UserID, "Addr": client.Addr})
	delete(RegisterDeviceMap, client.UserID)
}

// EventRegister 用户建立连接事件
func (manager *ClientManager) EventRegister(client *Client) {
	//client2.SetDeviceOnline(context.TODO(), client.UserID)
	manager.ClientsLock.Lock()
	defer manager.ClientsLock.Unlock()
	logger.Info(context.Background(), "EventRegister", map[string]string{"UserId": client.UserID, "Addr": client.Addr})
	manager.Clients[client] = true
	manager.Users[client.UserID] = client
}

// 管道处理程序
func (manager *ClientManager) start() {
	for {
		select {
		case conn := <-manager.Unregister:
			// 断开连接事件
			manager.EventUnRegister(conn)
		case conn := <-manager.Register:
			// 建立连接事件
			manager.EventRegister(conn)
		case message := <-manager.Broadcast:
			// 广播事件
			clients := manager.GetClients()
			for conn := range clients {
				select {
				case conn.Send <- message:
				default:
					close(conn.Send)
				}
			}
		}
	}
}

// ClearTimeoutConnections 定时清理超时连接
func ClearTimeoutConnections() {
	currentTime := uint64(time.Now().Unix())
	clients := WebsocketClientManager.GetClients()
	for client := range clients {
		if client.IsHeartbeatTimeout(currentTime) {
			logger.Info(context.Background(), "ClearTimeoutConnections", map[string]string{"UserId": client.UserID, "Addr": client.Addr})
			client.Socket.Close()
			delete(WebsocketClientManager.Clients, client)
			if WebsocketClientManager.Users[client.UserID] == client {
				delete(WebsocketClientManager.Users, client.UserID)
				logger.Info(context.Background(), "Delete WebsocketClientManager", map[string]string{"UserId": client.UserID, "Addr": client.Addr})

			}
		}
	}
}
