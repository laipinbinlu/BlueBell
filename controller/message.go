package controller

import (
	"blue_bell/dao/mysql"
	"blue_bell/models"
	"blue_bell/pkg/ws"
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"go.uber.org/zap"
)

// upgrader 是一个 websocket.Upgrader 实例,用于将 HTTP 连接升级为 WebSocket 连接
// 配置包括:
// - CheckOrigin: 跨域检查函数,这里允许所有跨域请求
// - ReadBufferSize: 读取缓冲区大小为 1024 字节
// - WriteBufferSize: 写入缓冲区大小为 1024 字节
// - Subprotocols: 支持的子协议,这里指定使用 JSON
var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true // 允许所有跨域请求
	},
	// 添加握手相关配置
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	// 添加子协议支持
	Subprotocols: []string{"json"},
}

// WebSocketHandler WebSocket连接处理
func WebSocketHandler(c *gin.Context) {
	// 从上下文中获取用户ID
	userID, err := getCurrentUser(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"code": CodeNeedLogin, "msg": "需要登录"})
		return
	}

	// 升级HTTP连接为WebSocket连接
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		zap.L().Error("upgrade websocket failed", zap.Error(err))
		return
	}

	// 创建新的WebSocket客户端
	client := ws.WebsocketManager.NewClient(userID, conn)

	// 获取离线消息并发送
	go sendOfflineMessages(client)

	// 启动客户端
	client.StartClient()
}

// SendMessageHandler 发送消息处理
func SendMessageHandler(c *gin.Context) {
	// 获取当前用户ID
	userID, err := getCurrentUser(c)
	if err != nil {
		ResponseError(c, CodeNeedLogin)
		return
	}

	// 解析请求参数   ---> 发送消息的内容和对象
	var req models.MessageRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		ResponseError(c, CodeInvalidParam)
		return
	}

	// 创建消息
	message := &models.Message{
		SenderUID:   userID,
		ReceiverUID: req.ReceiverUID,
		Content:     req.Content,
		CreateTime:  time.Now(),
		IsRead:      false,
	}

	// 保存消息到数据库
	if err := mysql.CreateMessage(message); err != nil {
		zap.L().Error("create message failed", zap.Error(err))
		ResponseError(c, CodeServerBusy)
		return
	}

	// 如果接收者在线，直接发送消息
	messageJSON, _ := json.Marshal(message)
	if ws.WebsocketManager.SendMessage(req.ReceiverUID, messageJSON) {
		message.IsRead = true
		// 更新消息状态为已读
		_ = mysql.MarkMessageAsRead(message.MsgID)
	}

	ResponseSuccess(c, nil)
}

// GetMessageListHandler 获取消息列表
func GetMessageListHandler(c *gin.Context) {
	// 获取当前用户ID
	userID, err := getCurrentUser(c)
	if err != nil {
		ResponseError(c, CodeNeedLogin)
		return
	}

	// 获取分页参数
	page, size := getPageInfo(c)

	// 获取消息列表
	messages, err := mysql.GetMessagesByUserID(userID, (page-1)*size, size)
	if err != nil {
		zap.L().Error("get messages failed", zap.Error(err))
		ResponseError(c, CodeServerBusy)
		return
	}

	ResponseSuccess(c, messages)
}

// GetUnreadCountHandler 获取未读消息数量
func GetUnreadCountHandler(c *gin.Context) {
	// 获取当前用户ID
	userID, err := getCurrentUser(c)
	if err != nil {
		ResponseError(c, CodeNeedLogin)
		return
	}

	// 获取未读消息数量
	count, err := mysql.GetUnreadMessageCount(userID)
	if err != nil {
		zap.L().Error("get unread count failed", zap.Error(err))
		ResponseError(c, CodeServerBusy)
		return
	}

	ResponseSuccess(c, gin.H{"count": count})
}

// GetChatHistoryHandler 获取聊天历史
func GetChatHistoryHandler(c *gin.Context) {
	// 获取当前用户ID
	userID, err := getCurrentUser(c)
	if err != nil {
		ResponseError(c, CodeNeedLogin)
		return
	}

	// 获取对方用户ID   -->  根据用户传入的url参数决定目标用户是谁。
	otherUID := c.Query("other_uid")

	if otherUID == "" {
		ResponseError(c, CodeInvalidParam)
		return
	}

	// 获取分页参数
	page, size := getPageInfo(c)
	otherUIDInt, _ := strconv.ParseInt(otherUID, 10, 64) // 先转化为最终的id
	// 获取聊天历史
	messages, err := mysql.GetChatHistory(userID, otherUIDInt, (page-1)*size, size)
	if err != nil {
		zap.L().Error("get chat history failed", zap.Error(err))
		ResponseError(c, CodeServerBusy)
		return
	}

	ResponseSuccess(c, messages)
}

// 发送离线消息  --- 获取数据库中的离线消息
func sendOfflineMessages(client *ws.Client) {
	messages, err := mysql.GetMessagesByUserID(client.ID, 0, 100)
	if err != nil {
		zap.L().Error("get offline messages failed", zap.Error(err))
		return
	}
	// 将未读取的消息推送给客户端
	for _, message := range messages {
		if !message.IsRead {
			messageJSON, _ := json.Marshal(message)
			if ws.WebsocketManager.SendMessage(client.ID, messageJSON) { // 推送给前端再将消息标记为已读消息。
				_ = mysql.MarkMessageAsRead(message.MsgID)
			}
		}
	}
}
