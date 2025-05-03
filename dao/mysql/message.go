package mysql

import (
	"blue_bell/models"
	"database/sql"
)

// CreateMessage 创建新消息--- 将当前用户私聊的消息存入数据库中
func CreateMessage(message *models.Message) error {
	sqlStr := `INSERT INTO message (
		sender_uid, receiver_uid, content, create_time, is_read
	) VALUES (?, ?, ?, ?, ?)`

	_, err := db.Exec(sqlStr,
		message.SenderUID,
		message.ReceiverUID,
		message.Content,
		message.CreateTime,
		message.IsRead,
	)
	return err
}

// GetMessagesByUserID 获取用户的所有的消息列表   -->  包括接收者和发送者的消息
func GetMessagesByUserID(userID int64, offset, limit int64) ([]*models.Message, error) {
	sqlStr := `SELECT 
		msg_id, sender_uid, receiver_uid, content, create_time, is_read
	FROM message 
	WHERE receiver_uid = ? OR sender_uid = ?
	ORDER BY create_time DESC
	LIMIT ?, ?`

	messages := make([]*models.Message, 0)
	err := db.Select(&messages, sqlStr, userID, userID, offset, limit)
	return messages, err
}

// GetUnreadMessageCount 获取用户未读消息数量
func GetUnreadMessageCount(userID int64) (int64, error) {
	sqlStr := `SELECT COUNT(*) FROM message WHERE receiver_uid = ? AND is_read = 0`

	var count int64
	err := db.Get(&count, sqlStr, userID)
	return count, err
}

// MarkMessageAsRead 将消息标记为已读
func MarkMessageAsRead(msgID int64) error {
	sqlStr := `UPDATE message SET is_read = 1 WHERE msg_id = ?`

	result, err := db.Exec(sqlStr, msgID)
	if err != nil {
		return err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return sql.ErrNoRows
	}
	return nil
}

// GetChatHistory 获取两个用户之间的聊天历史
func GetChatHistory(userID1, userID2 int64, offset, limit int64) ([]*models.Message, error) {
	sqlStr := `SELECT 
		msg_id, sender_uid, receiver_uid, content, create_time, is_read
	FROM message 
	WHERE (sender_uid = ? AND receiver_uid = ?) OR (sender_uid = ? AND receiver_uid = ?)
	ORDER BY create_time DESC
	LIMIT ?, ?`

	messages := make([]*models.Message, 0)
	err := db.Select(&messages, sqlStr, userID1, userID2, userID2, userID1, offset, limit)
	return messages, err
}
