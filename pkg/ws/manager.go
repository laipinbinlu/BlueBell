package ws

import (
	"github.com/gorilla/websocket"
	"go.uber.org/zap"
	"sync"
)

// Manager WebSocket连接管理器
type Manager struct {
	Clients    map[int64]*Client // 保存所有在线用户的WebSocket连接，key是用户ID
	ClientLock sync.RWMutex      // 读写锁，用于并发安全
}

// Client WebSocket客户端连接
type Client struct {
	ID     int64           // 用户ID
	Socket *websocket.Conn // WebSocket连接
	Send   chan []byte     // 发送消息的通道
}

var (
	// WebsocketManager 全局WebSocket连接管理器
	WebsocketManager = &Manager{
		Clients:    make(map[int64]*Client),
		ClientLock: sync.RWMutex{},
	}
)

// NewClient 创建新的WebSocket客户端
func (manager *Manager) NewClient(userID int64, conn *websocket.Conn) *Client {
	client := &Client{
		ID:     userID,
		Socket: conn,
		Send:   make(chan []byte, 100),
	}

	// 将新客户端添加到管理器中
	manager.ClientLock.Lock()
	manager.Clients[userID] = client
	manager.ClientLock.Unlock()

	return client
}

// RemoveClient 移除WebSocket客户端   --- 用户下线了，需要进行移除
func (manager *Manager) RemoveClient(client *Client) {
	manager.ClientLock.Lock()
	defer manager.ClientLock.Unlock()

	// 关闭发送通道
	if _, ok := manager.Clients[client.ID]; ok {
		close(client.Send)
		delete(manager.Clients, client.ID)
	}
}

// GetClient 获取指定用户的WebSocket客户端
func (manager *Manager) GetClient(userID int64) *Client {
	manager.ClientLock.RLock()
	defer manager.ClientLock.RUnlock()

	if client, ok := manager.Clients[userID]; ok {    
		return client
	}
	return nil
}

// SendMessage 发送消息到指定用户
func (manager *Manager) SendMessage(userID int64, message []byte) bool {
	manager.ClientLock.RLock()
	defer manager.ClientLock.RUnlock()

	if client, ok := manager.Clients[userID]; ok {  // 如果当前用户存在wb客户端的话，发送消息
		select {
		case client.Send <- message:
			return true
		default:
			zap.L().Error("send message failed, channel full")
			return false
		}
	}
	return false
} 