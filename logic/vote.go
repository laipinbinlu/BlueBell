package logic

import (
	"blue_bell/dao/mysql"
	"blue_bell/dao/redis"
	"blue_bell/models"
	"strconv"

	"go.uber.org/zap"
)

// 实际的帖子投票需求
func VoteForPost(uid int64, p *models.ParamVoteData) error {
	zap.L().Debug("VoteForPost", zap.Int64("uid", uid),
		zap.String("PostID", p.PostID), zap.Int8("Direction", p.Direction))

	// 将字符串格式的postID转换为int64
	postID, err := strconv.ParseInt(p.PostID, 10, 64)
	if err != nil {
		return err
	}

	// 1. 判断投票限制并更新Redis    ----> 先将redis中的投票记录更新
	if err := redis.VoteForPost(p.PostID, strconv.Itoa(int(uid)), float64(p.Direction)); err != nil {
		return err
	}

	// 2. 保存投票记录到MySQL  --->在将投票记录存储到数据库中
	if err := mysql.SavePostVote(postID, uid, p.Direction); err != nil {
		// 如果MySQL保存失败，记录错误日志
		zap.L().Error("mysql.SavePostVote failed",
			zap.Int64("post_id", postID),
			zap.Int64("user_id", uid),
			zap.Error(err))
		// 注意：这里我们选择继续执行而不是立即返回错误，因为Redis的投票已经成功
		// 我们可以通过后台任务来重试MySQL的保存操作
	}

	return nil
}
