package rabbitmq

import (
	"context"
	"encoding/json"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

var RMQ *RabbitMQ

// 使用rabbitmq来实现用户投票的解耦操作   --- > 将用户投票的消息发送到rabbitmq中，然后由消费者来处理
type RabbitMQ struct {
	conn    *amqp.Connection
	channel *amqp.Channel
	// 投票相关的交换机和队列
	VoteExchangeName string
	VoteQueueName    string
	VoteRoutingKey   string
}

type VoteMessage struct {
	PostID    int64 `json:"post_id"`
	UserID    int64 `json:"user_id"`
	Direction int8  `json:"direction"`
	Timestamp int64 `json:"timestamp"`
}

// Init 初始化 RabbitMQ 连接
func Init() error {
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	if err != nil {
		return err
	}

	ch, err := conn.Channel()
	if err != nil {
		return err
	}

	RMQ = &RabbitMQ{
		conn:             conn,
		channel:          ch,
		VoteExchangeName: "post_vote_exchange",
		VoteQueueName:    "post_vote_queue",
		VoteRoutingKey:   "post.vote",
	}

	// 声明交换机
	err = ch.ExchangeDeclare(
		RMQ.VoteExchangeName,
		"direct",
		true,  // durable
		false, // auto-deleted
		false, // internal
		false, // no-wait
		nil,   // arguments
	)
	if err != nil {
		return err
	}

	// 声明队列
	_, err = ch.QueueDeclare(
		RMQ.VoteQueueName,
		true,  // durable
		false, // auto-deleted
		false, // exclusive
		false, // no-wait
		nil,   // arguments
	)
	if err != nil {
		return err
	}

	// 绑定队列到交换机
	err = ch.QueueBind(
		RMQ.VoteQueueName,
		RMQ.VoteRoutingKey,
		RMQ.VoteExchangeName,
		false, // no-wait
		nil,   // arguments
	)
	if err != nil {
		return err
	}

	return nil
}

// Close 关闭连接
func (r *RabbitMQ) Close() {
	if r.channel != nil {
		r.channel.Close()
	}
	if r.conn != nil {
		r.conn.Close()
	}
}

// PublishVoteMessage 发布投票消息
func (r *RabbitMQ) PublishVoteMessage(postID, userID int64, direction int8) error {
	msg := VoteMessage{
		PostID:    postID,
		UserID:    userID,
		Direction: direction,
		Timestamp: time.Now().Unix(),
	}

	body, err := json.Marshal(msg)
	if err != nil {
		return err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	return r.channel.PublishWithContext(
		ctx,
		r.VoteExchangeName,
		r.VoteRoutingKey,
		false, // mandatory
		false, // immediate
		amqp.Publishing{
			ContentType:  "application/json",
			Body:         body,
			DeliveryMode: amqp.Persistent,
		},
	)
}
