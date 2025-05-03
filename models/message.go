package models

import "time"
// 用户实现私聊功能存入数据库的结构体
// Message 私信消息模型  
type Message struct {
	MsgID       int64     `json:"msg_id,string" db:"msg_id"`                    // 消息ID
	SenderUID   int64     `json:"sender_uid,string" db:"sender_uid"`            // 发送者ID
	ReceiverUID int64     `json:"receiver_uid,string" db:"receiver_uid"`        // 接收者ID
	Content     string    `json:"content" db:"content"`                         // 消息内容
	CreateTime  time.Time `json:"create_time" db:"create_time"`                 // 创建时间
	IsRead      bool      `json:"is_read" db:"is_read"`                        // 是否已读
}

// MessageRequest 发送消息请求
type MessageRequest struct {
	ReceiverUID int64  `json:"receiver_uid,string" binding:"required"`
	Content     string `json:"content" binding:"required"`
}

// MessageResponse 消息推送给接受者
type MessageResponse struct {
	Message
	SenderName   string `json:"sender_name"`   // 发送者用户名
	ReceiverName string `json:"receiver_name"` // 接收者用户名
} 