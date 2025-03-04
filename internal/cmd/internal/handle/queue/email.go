package queue

import (
	"context"

	"github.com/redis/go-redis/v9"
)

type EmailHandle struct {
	rds *redis.Client
}

func NewEmailHandle(rds *redis.Client) *EmailHandle {
	return &EmailHandle{rds: rds}
}

func (e *EmailHandle) Handle(ctx context.Context) error {
	return nil
}
