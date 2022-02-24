package redis

import (
	"context"
	"errors"
	"time"

	"github.com/go-redis/redis/v8"
)

var rdb *redis.Client

var rdbNilErr = errors.New("redis client is not initialized")

func Setup(db int, addr string, password string) error {
	rdb = redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password,
		DB:       db,
	})
	_, err := rdb.Ping(context.Background()).Result()
	return err
}

func Set(ctx context.Context, key, value string, second uint) error {
	if rdb == nil {
		return rdbNilErr
	}
	return rdb.Set(ctx, key, value, time.Duration(second)*time.Second).Err()
}

func Get(ctx context.Context, key string) (string, error) {
	if rdb == nil {
		return "", rdbNilErr
	}
	return rdb.Get(ctx, key).Result()
}

func Del(ctx context.Context, key string) error {
	if rdb == nil {
		return rdbNilErr
	}
	return rdb.Del(ctx, key).Err()
}

func MSet(ctx context.Context, data map[string]string, second uint) error {
	if rdb == nil {
		return rdbNilErr
	}
	_, err := rdb.Pipelined(ctx, func(pipe redis.Pipeliner) error {
		for k, v := range data {
			err := rdb.Set(ctx, k, v, time.Duration(second)*time.Second).Err()
			if err != nil {
				return err
			}
		}
		return nil
	})
	return err
}

func MGet(ctx context.Context, key []string) ([]interface{}, error) {
	if rdb == nil {
		return nil, rdbNilErr
	}
	return rdb.MGet(ctx, key...).Result()
}

func MDel(ctx context.Context, key []string) (int64, error) {
	if rdb == nil {
		return 0, rdbNilErr
	}
	return rdb.Del(ctx, key...).Result()
}

func HashSet(ctx context.Context, data map[string]map[string]interface{}, second uint) error {
	if rdb == nil {
		return rdbNilErr
	}
	_, err := rdb.Pipelined(ctx, func(pipe redis.Pipeliner) error {
		for k, vm := range data {
			err := rdb.HSet(ctx, k, vm).Err()
			if err != nil {
				return err
			}
			err = rdb.PExpire(ctx, k, time.Duration(second)*time.Second).Err()
			if err != nil {
				return err
			}
		}
		return nil
	})
	return err
}

func HashGet(ctx context.Context, key string) (time.Duration, map[string]string, error) {
	if rdb == nil {
		return 0, nil, rdbNilErr
	}
	t, err := rdb.TTL(ctx, key).Result()
	if err != nil {
		return 0, nil, err
	}
	data, err := rdb.HGetAll(ctx, key).Result()
	return t, data, err
}

func Keys(ctx context.Context, key string) ([]string, error) {
	if rdb == nil {
		return nil, rdbNilErr
	}
	return rdb.Keys(ctx, key).Result()
}
