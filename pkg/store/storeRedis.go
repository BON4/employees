package store

import (
	"context"

	redis "github.com/go-redis/redis/v8"
)

type RedisStore struct {
	cl *redis.Client
}

func NewRedisStore(cl *redis.Client) Store {
	return &RedisStore{cl: cl}
}

func (rS RedisStore) Get(ctx context.Context, k string, v *[]byte) (found bool, err error) {
	val, err := rS.cl.Get(ctx, k).Result()
	if err != nil {
		return false, err
	}

	*v = append(*v, []byte(val)...)

	return true, nil
}

func (rS RedisStore) Set(ctx context.Context, k string, v []byte) error {
	return rS.cl.Set(ctx, k, v, 0).Err()
}

func (rS RedisStore) Delete(ctx context.Context, k string) error {
	return rS.cl.Del(ctx, k).Err()
}

func (rS RedisStore) Close() error {
	return rS.cl.Close()
}
