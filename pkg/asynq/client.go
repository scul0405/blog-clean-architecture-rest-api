package asynq

import (
	"github.com/hibiken/asynq"
)

func NewAsynqClient(redisOpt asynq.RedisClientOpt) *asynq.Client {
	client := asynq.NewClient(redisOpt)

	return client
}
