package rabbitmq

import (
	"blue_bell/dao/mysql"
	"encoding/json"
	"time"

	"github.com/rabbitmq/amqp091-go"
	"go.uber.org/zap"
)

// StartVoteConsumer 启动投票消息消费者
func StartVoteConsumer() {
	const (
		batchSize = 100             // 批量处理的消息数量
		maxWait   = 3 * time.Second // 最大等待时间
	)

	// 设置 QoS，一次预取100条消息
	err := RMQ.channel.Qos(
		batchSize, // prefetch count
		0,         // prefetch size
		false,     // global
	)
	if err != nil {
		zap.L().Error("设置QoS失败", zap.Error(err))
		return
	}

	msgs, err := RMQ.channel.Consume(
		RMQ.VoteQueueName,
		"",    // consumer
		false, // auto-ack，设为false以便手动确认
		false, // exclusive
		false, // no-local
		false, // no-wait
		nil,   // args
	)
	if err != nil {
		zap.L().Fatal("注册消费者失败", zap.Error(err))
		return
	}

	forever := make(chan bool)

	go func() {
		var msgBatch []amqp091.Delivery
		ticker := time.NewTicker(maxWait)
		defer ticker.Stop()

		for {
			select {
			case msg := <-msgs:
				msgBatch = append(msgBatch, msg)

				// 当积累的消息达到批处理大小时，进行处理
				if len(msgBatch) >= batchSize {
					processMsgBatch(msgBatch)
					msgBatch = nil // 清空批次
				}

			case <-ticker.C:
				// 定时处理已积累的消息，即使未达到批处理大小
				if len(msgBatch) > 0 {
					processMsgBatch(msgBatch)
					msgBatch = nil
				}
			}
		}
	}()

	zap.L().Info("Vote consumer started")
	<-forever
}

// processMsgBatch 批量处理消息
func processMsgBatch(msgs []amqp091.Delivery) {
	for _, msg := range msgs {
		var voteData struct {
			PostID    int64 `json:"post_id"`
			UserID    int64 `json:"user_id"`
			Direction int8  `json:"direction"`
		}

		if err := json.Unmarshal(msg.Body, &voteData); err != nil {
			zap.L().Error("解析投票消息失败", zap.Error(err))
			msg.Nack(false, true) // 解析失败，重新入队
			continue
		}

		// 保存投票记录
		if err := mysql.SavePostVote(voteData.PostID, voteData.UserID, voteData.Direction); err != nil {
			zap.L().Error("保存投票记录失败",
				zap.Error(err),
				zap.Int64("post_id", voteData.PostID),
				zap.Int64("user_id", voteData.UserID))
			msg.Nack(false, true) // 处理失败，重新入队
			continue
		}

		msg.Ack(false) // 处理成功，确认消息
	}
}
