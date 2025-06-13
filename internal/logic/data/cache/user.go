package cache

import (
	"context"
	"github.com/zhixunjie/im-fun/internal/logic/data"
	k "github.com/zhixunjie/im-fun/pkg/goredis/key"
	"time"
)

type UserCache struct {
	userRepo *data.UserRepo
}

func NewUserCache(userRepo *data.UserRepo) *UserCache {
	return &UserCache{userRepo: userRepo}
}

func (b *UserCache) GetToken(ctx context.Context, uid uint64) string {
	key := data.UserToken.Format(k.M{"uid": uid})
	mem := b.userRepo.RedisClient

	token, _ := mem.Get(ctx, key).Result()

	return token
}

func (b *UserCache) SetToken(ctx context.Context, uid uint64, token string) error {
	key := data.UserToken.Format(k.M{"uid": uid})
	mem := b.userRepo.RedisClient

	return mem.Set(ctx, key, token, time.Second*86400*10).Err()
}

func (b *UserCache) DelToken(ctx context.Context, uid uint64) error {
	key := data.UserToken.Format(k.M{"uid": uid})
	mem := b.userRepo.RedisClient

	return mem.Del(ctx, key).Err()
}
