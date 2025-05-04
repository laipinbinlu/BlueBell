package logic

import (
	"blue_bell/dao/redis"
	"blue_bell/models"
	"blue_bell/pkg/rabbitmq"
	"strconv"

	"go.uber.org/zap"
)

// VoteForPost 实现投票的功能
func VoteForPost(uid int64, p *models.ParamVoteData) error {
	zap.L().Debug("VoteForPost",
		zap.Int64("uid", uid),
		zap.String("PostID", p.PostID),
		zap.Int8("Direction", p.Direction))

	// 将字符串格式的postID转换为int64
	postID, err := strconv.ParseInt(p.PostID, 10, 64)
	if err != nil {
		return err
	}

	// 1. 判断投票限制并更新Redis    ----> 先将redis中的投票记录更新
	if err := redis.VoteForPost(p.PostID, strconv.Itoa(int(uid)), float64(p.Direction)); err != nil {
		return err
	}

	// 2. 异步保存投票记录到MySQL  --->通过消息队列异步处理
	go func() {
		if err := rabbitmq.RMQ.PublishVoteMessage(postID, uid, p.Direction); err != nil {
			zap.L().Error("发送投票消息到RabbitMQ失败",
				zap.Error(err),
				zap.Int64("post_id", postID),
				zap.Int64("user_id", uid),
				zap.Int8("direction", p.Direction))
			// 注意：即使消息发送失败，我们也不中断流程，因为Redis中的数据已更新成功
			// 这里可以考虑将失败的消息写入本地日志或重试队列
		}
	}()

	return nil
}
