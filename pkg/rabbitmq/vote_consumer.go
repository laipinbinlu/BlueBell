package rabbitmq

import (
	"blue_bell/dao/mysql"
	"encoding/json"

	"go.uber.org/zap"
)

// StartVoteConsumer 启动投票消息消费者    --->启动后端协程来不断地将投票消息消费给mysql记录
func StartVoteConsumer() {
	msgs, err := RMQ.channel.Consume(
		RMQ.VoteQueueName,
		"",    // consumer
		false, // auto-ack
		false, // exclusive
		false, // no-local
		false, // no-wait
		nil,   // args
	)
	if err != nil {
		zap.L().Fatal("Failed to register a consumer", zap.Error(err))
		return
	}

	forever := make(chan bool)

	go func() {
		for d := range msgs {
			var msg VoteMessage
			if err := json.Unmarshal(d.Body, &msg); err != nil {
				zap.L().Error("解析投票消息失败", zap.Error(err))
				d.Nack(false, true) // 重新入队
				continue
			}

			// 保存投票记录到MySQL
			if err := mysql.SavePostVote(msg.PostID, msg.UserID, msg.Direction); err != nil {
				zap.L().Error("保存投票记录到MySQL失败",
					zap.Error(err),
					zap.Int64("post_id", msg.PostID),
					zap.Int64("user_id", msg.UserID))
				d.Nack(false, true) // 重新入队
				continue
			}

			// 确认消息已处理
			d.Ack(false)
		}
	}()

	zap.L().Info("Vote consumer started")
	<-forever
}
