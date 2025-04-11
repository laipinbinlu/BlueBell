package logic

import (
	"blue_bell/dao/redis"
	"blue_bell/models"
	"go.uber.org/zap"
	"strconv"
)

// 实际的帖子投票需求
func VoteForPost(uid int64, p *models.ParamVoteData) error {
	zap.L().Debug("VoteForPost", zap.Int64("uid", uid),
		zap.String("PostID", p.PostID), zap.Int8("Direction", p.Direction))
	//以下操作都是需要调用dao层
	//1.判断帖子投票限制
	//2.更新帖子分数
	//3.将用户投票信息写入对应的redis中
	return redis.VoteForPost(p.PostID, strconv.Itoa(int(uid)), float64(p.Direction))
}
