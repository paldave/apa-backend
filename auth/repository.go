package auth

import (
	"apa-backend/entity"
	"context"
	"fmt"
	"time"

	"github.com/go-redis/redis/v8"
)

type Repository interface {
	Create(*entity.RedisToken) error
}

type repository struct {
	rdb *redis.Client
}

func NewRepository(rdb *redis.Client) *repository {
	return &repository{rdb}
}

func (r *repository) Create(t *entity.RedisToken) error {
	utc := time.Unix(t.Expiry, 0)
	now := time.Now()

	str := fmt.Sprintf("%s:%s", t.UserId, t.Id)

	err := r.rdb.Set(context.TODO(), str, 1, utc.Sub(now)).Err()
	if err != nil {
		return err
	}

	return nil
}
