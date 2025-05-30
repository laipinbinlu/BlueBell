package redis

import (
	"blue_bell/models"
	"context"
	"strconv"
	"time"

	"github.com/redis/go-redis/v9"
)

func getIDsFormKey(key string, Page, Size int64) ([]string, error) {
	// 根据设置，不管是最新还是最热都是分数从高到低，所有直接从redis中查询数据即可
	start := (Page - 1) * Size
	end := start + Size - 1
	// 开始redis查询
	return client.ZRevRange(context.Background(), key, start, end).Result()
}

func GetPostIDsInOrder(p *models.ParamPostList) ([]string, error) {
	// 查询redis中的时间排序或者score排序的zset
	key := getRedisKey(KeyPostTimeZSet) // 默认是时间
	// 传入参数
	if p.Order == models.OrderScore { // 如果是分数，那么修改key
		key = getRedisKey(KeyPostScoreZSet)
	}
	return getIDsFormKey(key, p.Page, p.Size)
}

// 按照社区查询出排好序的ids
func GetCommunityPostIDsInOrder(p *models.ParamPostList) ([]string, error) {
	//TODO:
	//cKey := getRedisKey(KeyCommunitySetPF + strconv.Itoa(int(communityID)))
	orderKey := getRedisKey(KeyPostTimeZSet) // 默认是时间
	// 传入参数
	if p.Order == models.OrderScore { // 如果是分数，那么修改key
		orderKey = getRedisKey(KeyPostScoreZSet)
	}

	key := orderKey + strconv.Itoa(int(p.CommunityID))
	if client.Exists(context.Background(), key).Val() < 1 {
		pipeline := client.Pipeline()
		pipeline.ZInterStore(context.Background(), key, &redis.ZStore{
			Aggregate: "MAX",
		})
		pipeline.Expire(context.Background(), key, time.Second*60)
		_, err := pipeline.Exec(context.Background())
		if err != nil {
			return nil, err
		}
	}

	return getIDsFormKey(key, p.Page, p.Size)
}
