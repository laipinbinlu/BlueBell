package ws

import (
	"encoding/json"
	"github.com/gorilla/websocket"
	"go.uber.org/zap"
	"time"
)

const (
	// 写入超时
	writeWait = 10 * time.Second

	// 读取超时
	pongWait = 60 * time.Second

	// 发送ping的频率
	pingPeriod = (pongWait * 9) / 10

	// 最大消息大小
	maxMessageSize = 512
)

// StartClient 启动WebSocket客户端的读写协程
func (c *Client) StartClient() {
	// 启动读协程
	go c.readPump()
	// 启动写协程
	go c.writePump()
}

// readPump 处理WebSocket读取
func (c *Client) readPump() {
	defer func() {
		WebsocketManager.RemoveClient(c)
		c.Socket.Close()
	}()

	c.Socket.SetReadLimit(maxMessageSize)   
	c.Socket.SetReadDeadline(time.Now().Add(pongWait))
	c.Socket.SetPongHandler(func(string) error {
		c.Socket.SetReadDeadline(time.Now().Add(pongWait))
		return nil
	})

	for {
		_, message, err := c.Socket.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				zap.L().Error("unexpected close error", zap.Error(err))
			}
			break
		}
		// 处理接收到的消息
		var msg map[string]interface{}
		if err := json.Unmarshal(message, &msg); err != nil {
			zap.L().Error("unmarshal message failed", zap.Error(err))
			continue
		}
		// 这里可以处理接收到的消息，比如保存到数据库等
		zap.L().Info("received message", zap.Any("message", msg))
	}
}

// writePump 处理WebSocket写入
func (c *Client) writePump() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		c.Socket.Close()
	}()

	for {
		select {
		case message, ok := <-c.Send:
			c.Socket.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				c.Socket.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			w, err := c.Socket.NextWriter(websocket.TextMessage)
			if err != nil {
				return
			}
			w.Write(message)

			if err := w.Close(); err != nil {
				return
			}
		case <-ticker.C:
			c.Socket.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.Socket.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
} 