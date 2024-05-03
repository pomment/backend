package dao

import (
	"context"
	"errors"
	"github.com/redis/go-redis/v9"
	"pomment-go/common"
	"time"
)

var RedisClient *redis.Client
var ExpireTime = 86400
var isEnabled = false

var ctx = context.Background()

func ConnectToRedisServer(enabled bool, addr string, password string, db int) {
	isEnabled = enabled
	if !isEnabled {
		return
	}
	RedisClient = redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password,
		DB:       db,
	})
}

func SetCache(key string, value string) (e error) {
	if !isEnabled {
		return nil
	}
	err := RedisClient.Set(ctx, key, value, time.Duration(ExpireTime)*time.Second).Err()
	if err != nil {
		return err
	}
	return nil
}

func GetCache(key string) (v string, e error) {
	if !isEnabled {
		return "", nil
	}
	val, err := RedisClient.Get(ctx, key).Result()
	if errors.Is(err, redis.Nil) {
		return "", nil
	}
	if err != nil {
		return "", err
	}
	return val, nil
}

func DeleteCache(key string) (e error) {
	if !isEnabled {
		return nil
	}
	_, err := RedisClient.Del(ctx, key).Result()
	if errors.Is(err, redis.Nil) {
		return nil
	}
	if err != nil {
		return err
	}
	return nil
}

func DeleteAllCache() (e error) {
	if !isEnabled {
		return nil
	}
	_, err := RedisClient.FlushAll(ctx).Result()
	if errors.Is(err, redis.Nil) {
		return nil
	}
	if err != nil {
		return err
	}
	return nil
}

func DeleteCacheForThread(item *common.Thread) (err error) {
	if !isEnabled {
		return nil
	}
	err = DeleteCache(common.CachePostIDKeyPrefix + item.ID)
	if err != nil {
		return err
	}
	err = DeleteCache(common.CachePostURLKeyPrefix + item.URL)
	return err
}
