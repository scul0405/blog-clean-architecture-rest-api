package repository

import (
	"context"
	"encoding/json"
	"github.com/opentracing/opentracing-go"
	"github.com/pkg/errors"
	"github.com/redis/go-redis/v9"
	"github.com/scul0405/blog-clean-architecture-rest-api/internal/auth"
	"github.com/scul0405/blog-clean-architecture-rest-api/internal/models"
	"time"
)

type authRedisRepo struct {
	rdb *redis.Client
}

func NewAuthRedisRepository(rdb *redis.Client) auth.RedisRepository {
	return &authRedisRepo{rdb: rdb}
}

func (r *authRedisRepo) GetByIDCtx(ctx context.Context, key string) (*models.User, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "authRedisRepo.GetByIDCtx")
	defer span.Finish()

	userBytes, err := r.rdb.Get(ctx, key).Bytes()
	if err != nil {
		return nil, errors.Wrap(err, "authRedisRepo.GetByIDCtx.redisClient.Get")
	}
	user := &models.User{}
	if err = json.Unmarshal(userBytes, user); err != nil {
		return nil, errors.Wrap(err, "authRedisRepo.GetByIDCtx.json.Unmarshal")
	}
	return user, nil
}

func (r *authRedisRepo) SetUserCtx(ctx context.Context, key string, seconds int, user *models.User) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "authRedisRepo.SetUserCtx")
	defer span.Finish()

	userBytes, err := json.Marshal(user)
	if err != nil {
		return errors.Wrap(err, "authRedisRepo.SetUserCtx.json.Unmarshal")
	}
	if err = r.rdb.Set(ctx, key, userBytes, time.Second*time.Duration(seconds)).Err(); err != nil {
		return errors.Wrap(err, "authRedisRepo.SetUserCtx.redisClient.Set")
	}
	return nil
}

func (r *authRedisRepo) DeleteUserCtx(ctx context.Context, key string) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "authRedisRepo.DeleteUserCtx")
	defer span.Finish()

	if err := r.rdb.Del(ctx, key).Err(); err != nil {
		return errors.Wrap(err, "authRedisRepo.DeleteUserCtx.redisClient.Del")
	}
	return nil
}
