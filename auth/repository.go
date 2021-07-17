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
	Exists(tokenId, userId string) (bool, error)
	Delete(tokenId, userId string) error
}

type repository struct {
	rdb *redis.Client
}

func NewRepository(rdb *redis.Client) *repository {
	return &repository{rdb}
}

func buildString(tId, uId string) string {
	return fmt.Sprintf("%s:%s", uId, tId)
}

func (r *repository) Create(t *entity.RedisToken) error {
	utc := time.Unix(t.Expiry, 0)
	now := time.Now()

	str := buildString(t.Id, t.UserId)

	err := r.rdb.Set(context.TODO(), str, 1, utc.Sub(now)).Err()
	if err != nil {
		return err
	}

	return nil
}

func (r *repository) Exists(tokenId, userId string) (bool, error) {
	str := buildString(tokenId, userId)

	err := r.rdb.Get(context.TODO(), str).Err()
	if err != nil {
		return false, err
	}

	return true, nil
}

func (r *repository) Delete(tokenId, userId string) error {
	str := buildString(tokenId, userId)

	if err := r.rdb.Del(context.TODO(), str).Err(); err != nil {
		return err
	}

	return nil
}
