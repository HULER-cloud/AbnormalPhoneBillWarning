package core

import (
	"AbnormalPhoneBillWarning/global"
	"context"
	"fmt"
	"github.com/go-redis/redis"
	"log"
	"time"
)

// 初始化redis连接
func InitRedis() {
	redisCfg := global.Config.Redis
	fmt.Println(redisCfg.Addr())
	rdb := redis.NewClient(&redis.Options{
		Addr:     redisCfg.Addr(),
		Password: redisCfg.Password,
		DB:       0,
		PoolSize: redisCfg.PoolSize,
	})
	_, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	res, err := rdb.Ping().Result()
	fmt.Println(res)
	if err != nil {
		log.Fatalf("redis[%s]连接失败！", redisCfg.Addr())
		return
	}
	global.Redis = rdb
}
