package repository

import (
	"context"
	"encoding/json"
	"github.com/opentracing/opentracing-go"
	"github.com/pkg/errors"
	"github.com/redis/go-redis/v9"
	"github.com/scul0405/blog-clean-architecture-rest-api/internal/blog"
	"github.com/scul0405/blog-clean-architecture-rest-api/internal/models"
	"time"
)

type blogRedisRepo struct {
	rdb *redis.Client
}

func NewBlogRedisRepository(rdb *redis.Client) blog.RedisRepository {
	return &blogRedisRepo{rdb: rdb}
}

func (r *blogRedisRepo) GetBlogByIDCtx(ctx context.Context, key string) (*models.BlogBase, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "blogRedisRepo.GetBlogByIDCtx")
	defer span.Finish()

	blogBytes, err := r.rdb.Get(ctx, key).Bytes()
	if err != nil {
		return nil, errors.Wrap(err, "blogRedisRepo.GetBlogByIDCtx.redisClient.Get")
	}

	blog := &models.BlogBase{}
	if err = json.Unmarshal(blogBytes, blog); err != nil {
		return nil, errors.Wrap(err, "blogRedisRepo.GetBlogByIDCtx.json.Unmarshal")
	}

	return blog, nil
}

func (r *blogRedisRepo) SetBlogCtx(ctx context.Context, key string, seconds int, blog *models.BlogBase) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "blogRedisRepo.SetBlogCtx")
	defer span.Finish()

	blogBytes, err := json.Marshal(blog)
	if err != nil {
		return errors.Wrap(err, "blogRedisRepo.SetBlogCtx.json.Unmarshal")
	}

	if err = r.rdb.Set(ctx, key, blogBytes, time.Second*time.Duration(seconds)).Err(); err != nil {
		return errors.Wrap(err, "blogRedisRepo.SetBlogCtx.redisClient.Set")
	}

	return nil
}

func (r *blogRedisRepo) DeleteBlogCtx(ctx context.Context, key string) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "blogRedisRepo.DeleteBlogCtx")
	defer span.Finish()

	if err := r.rdb.Del(ctx, key).Err(); err != nil {
		return errors.Wrap(err, "blogRedisRepo.DeleteBlogCtx.redisClient.Del")
	}
	return nil
}
