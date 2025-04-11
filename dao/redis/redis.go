package redis

import (
	"blue_bell/settings"
	"context"
	"fmt"

	"github.com/redis/go-redis/v9"
)

var client *redis.Client

// 更具配置文件连接redis
func Init(redconfig *settings.RedisConfig) (err error) {
	client = redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", redconfig.Host, redconfig.Port),
		Password: redconfig.Password, // 密码
		DB:       redconfig.DB,       // 数据库
		PoolSize: redconfig.PoolSize, // 连接池大小
	})

	_, err = client.Ping(context.Background()).Result()
	if err != nil {
		return err
	}
	return
}

func Close() {
	_ = client.Close()
}
