package redis

import (
	"blue_bell/dao/mysql"
	"context"
	"errors"
	"math"
	"strconv"
	"time"

	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
)

/*
投票分为四种情况：1.投赞成票 2.投反对票 3.取消投票 4.反转投票

记录文章参与投票的人
更新文章分数：赞成票要加分；反对票减分

v=1时，有两种情况
1.之前没投过票，现在要投赞成票
2.之前投过反对票，现在要改为赞成票
v=0时，有两种情况
1.之前投过赞成票，现在要取消
2.之前投过反对票，现在要取消
v=-1时，有两种情况
1.之前没投过票，现在要投反对票
2.之前投过赞成票，现在要改为反对票
*/

const (
	OneWeekInSeconds         = 7 * 24 * 3600
	VoteScore        float64 = 432
	PostPerAge               = 20
)

var (
	ErrorVoteTimeExpire = errors.New("已过投票时间")
	ErrorVoted          = errors.New("已经投过票了")
)

func VoteForPost(postID, userID string, v float64) (err error) {
	// 1.判断帖子投票限制
	// 先取出当前帖子的发布时间
	postTime := client.ZScore(context.Background(), getRedisKey(KeyPostTimeZSet), postID).Val()
	if float64(time.Now().Unix())-postTime > OneWeekInSeconds {
		// 超过了时间，不允许投票了
		return ErrorVoteTimeExpire
	}

	// 2.更新帖子分数
	// 需要先查看该用户为该帖子的投票记录
	key := getRedisKey(KeyPostVotedHashPF + postID)
	ov, err := client.HGet(context.Background(), key, userID).Float64()
	if err != nil && err != redis.Nil {
		return err
	}
	// 如果投票的值相同，则表示已经投过了票了，没有必要再投票了
	if v == ov {
		return ErrorVoted
	}

	var op float64
	if v > ov { //确定分数的结果方向
		op = 1
	} else {
		op = -1
	}
	diffAbs := math.Abs(ov - v) // 差值绝对值

	// 使用redis事务
	pipeline := client.TxPipeline()

	// 更新帖子分数
	pipeline.ZIncrBy(context.Background(), getRedisKey(KeyPostScoreZSet), op*diffAbs*VoteScore, postID)

	// 3.记录用户为该帖子投票的数据
	if v == 0 {
		// 取消投票，删除投票记录
		pipeline.HDel(context.Background(), key, userID)
	} else {
		// 更新投票记录
		pipeline.HSet(context.Background(), key, userID, v)
	}

	_, err = pipeline.Exec(context.Background())
	return err
}

// 创建新的表到redis数据库中
func CreatePost(postID, communityID int64) error {
	// 使用redis事务操作
	pipeline := client.TxPipeline()

	// 帖子的创建时间
	pipeline.ZAdd(context.Background(), getRedisKey(KeyPostTimeZSet), redis.Z{
		Score:  float64(time.Now().Unix()),
		Member: postID,
	})

	// 初始化 帖子的分数(根据时间，保证最新最热的帖子被人发现)
	pipeline.ZAdd(context.Background(), getRedisKey(KeyPostScoreZSet), redis.Z{
		Score:  float64(time.Now().Unix()),
		Member: postID,
	})
	// 补充： 初始化社区id对应的帖子id的表  set
	cKey := getRedisKey(KeyCommunitySetPF + strconv.Itoa(int(communityID)))
	pipeline.SAdd(context.Background(), cKey, postID)

	_, err := pipeline.Exec(context.Background())

	return err
}

// GetPostVoteData 根据ids查询每篇帖子的投赞成票的数据
func GetPostVoteData(ids []string) (data []int64, err error) {
	// 查询每个帖子的赞成票数量
	data = make([]int64, 0, len(ids))

	for _, id := range ids {
		// 1. 先从Redis获取投票数据
		key := getRedisKey(KeyPostVotedHashPF + id)
		vals, err := client.HVals(context.Background(), key).Result()
		if err != nil {
			return nil, err
		}

		// 2. 如果Redis中有数据，直接统计
		var count int64
		if len(vals) > 0 {
			for _, val := range vals {
				if val == "1" {
					count++
				}
			}
		} else {
			// 3. 如果Redis中没有数据，从MySQL中获取
			postID, _ := strconv.ParseInt(id, 10, 64)
			upCount, _, err := mysql.GetPostVoteCount(postID)
			if err != nil {
				// 如果发生错误，记录日志但继续处理其他帖子
				zap.L().Error("mysql.GetPostVoteCount failed",
					zap.String("post_id", id),
					zap.Error(err))
				count = 0
			} else {
				count = int64(upCount)
			}
		}

		data = append(data, count)
	}

	return data, nil
}
