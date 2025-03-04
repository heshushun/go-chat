package cron

import (
	"context"
	"fmt"

	"go-chat/internal/repository/cache"
)

type ClearWsCache struct {
	storage *cache.ServerStorage
}

func NewClearWsCache(storage *cache.ServerStorage) *ClearWsCache {
	return &ClearWsCache{storage: storage}
}

// Spec 配置定时任务规则
// 每隔30分钟处理 websocket 缓存
func (c *ClearWsCache) Spec() string {
	return "* * * * *"
}

func (c *ClearWsCache) Enable() bool {
	return true
}

func (c *ClearWsCache) Handle(ctx context.Context) error {

	for _, sid := range c.storage.GetExpireServerAll(ctx) {
		c.clear(ctx, sid)
	}

	return nil
}

func (c *ClearWsCache) clear(ctx context.Context, sid string) {
	var cursor uint64
	for {
		var keys []string
		var err error
		keys, cursor, err = c.storage.Redis().Scan(ctx, cursor, fmt.Sprintf("ws:%s:*", sid), 200).Result()
		if err != nil {
			return
		}

		c.storage.Redis().Del(ctx, keys...)

		if cursor == 0 {
			_ = c.storage.DelExpireServer(ctx, sid)
			break
		}
	}
}
