package myredis

import (
	"context"
	"time"

	"github.com/go-redis/redis/v8"
)

type Redis struct {
	// redis client
	rc *redis.Client
}

func NewRedis(uri string) (myRedis *Redis, err error) {
	rc := redis.NewClient(&redis.Options{
		Addr: uri,
	})
	if _, err = rc.Ping(context.TODO()).Result(); err != nil {
		return
	}
	myRedis = &Redis{
		rc: rc,
	}
	return
}

// key列表長度
func (r *Redis) LLen(context context.Context, key string) int64 {
	return r.rc.LLen(context, key).Val()
}

// key是否存在
func (r *Redis) Exists(context context.Context, key string) int64 {
	return r.rc.Exists(context, key).Val()
}

// 返回支持事務的Pipeline
func (r *Redis) TxPipeline() redis.Pipeliner {
	return r.rc.TxPipeline()
}

// 依據ip在時間區間內訪問的頻率判斷是否可在訪問
// 如可訪問返回含當次訪問次數 不能返回-1
func (r *Redis) AllowRequest(ip string, limitTimes int64, limitDuration time.Duration) (isAllow bool, currTimes int64) {
	ctx := context.TODO()
	preTimes := r.LLen(ctx, ip)
	if preTimes >= limitTimes {
		return false, -1
	}

	// 如果第一次訪問，則加到redis
	if v := r.Exists(ctx, ip); v == 0 {
		// 支持事務的pipeline
		pipe := r.TxPipeline()
		pipe.RPush(ctx, ip, ip)
		// 設置過期
		pipe.Expire(ctx, ip, limitDuration)
		// 命令批次送出
		_, _ = pipe.Exec(ctx)
	} else {
		r.rc.RPushX(ctx, ip, ip)
	}

	return true, preTimes + 1
}
