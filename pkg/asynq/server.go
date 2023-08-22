package asynq

import (
	"context"
	"github.com/hibiken/asynq"
	"github.com/scul0405/blog-clean-architecture-rest-api/pkg/logger"
)

const (
	QueueCritical = "critical"
	QueueDefault  = "default"
)

type RedisTaskProcessor struct {
	server   *asynq.Server
	handlers map[string]asynq.HandlerFunc
}

func NewRedisTaskProcessor(redisOpt asynq.RedisClientOpt, logger logger.Logger) *RedisTaskProcessor {
	server := asynq.NewServer(
		redisOpt,
		asynq.Config{
			Queues: map[string]int{
				QueueCritical: 10,
				QueueDefault:  5,
			},
			ErrorHandler: asynq.ErrorHandlerFunc(func(ctx context.Context, task *asynq.Task, err error) {
				logger.Errorf("process task failed: type=%s, payload=%s, err=%v", task.Type(), task.Payload(), err)
			}),
			Logger: logger,
		},
	)

	return &RedisTaskProcessor{
		server:   server,
		handlers: make(map[string]asynq.HandlerFunc),
	}
}

func (p *RedisTaskProcessor) Start() error {
	mux := asynq.NewServeMux()

	p.mapHandlersToMux(mux)
	return p.server.Start(mux)
}

func (p *RedisTaskProcessor) RegisterHandler(taskName string, handler asynq.HandlerFunc) {
	p.handlers[taskName] = handler
}

func (p *RedisTaskProcessor) mapHandlersToMux(mux *asynq.ServeMux) {
	for taskName, handler := range p.handlers {
		mux.HandleFunc(taskName, handler)
	}
}
