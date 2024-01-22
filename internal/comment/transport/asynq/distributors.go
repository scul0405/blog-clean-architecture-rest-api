package asynq

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/hibiken/asynq"
	"github.com/scul0405/blog-clean-architecture-rest-api/pkg/logger"
)

type CommentTaskDistributor interface {
	DistributeTaskLikeComment(ctx context.Context, payload *LikeCommentPayload, opts ...asynq.Option) error
	DistributeTaskDislikeComment(ctx context.Context, payload *DislikeCommentPayload, opts ...asynq.Option) error
}

type commentTaskDistributor struct {
	client *asynq.Client
	logger logger.Logger
}

func NewCommentTaskDistributor(client *asynq.Client, logger logger.Logger) CommentTaskDistributor {
	return &commentTaskDistributor{
		client: client,
		logger: logger,
	}
}

func (distributor *commentTaskDistributor) DistributeTaskLikeComment(ctx context.Context, payload *LikeCommentPayload, opts ...asynq.Option) error {
	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("failed to marshal task payload: %w", err)
	}

	task := asynq.NewTask(TypeLikeCommentTask, jsonPayload, opts...)
	info, err := distributor.client.EnqueueContext(ctx, task)
	if err != nil {
		return fmt.Errorf("failed to enqueue task: %w", err)
	}

	distributor.logger.Infof("type=%s, payload=%s, queue=%s, maxRetry=%d enqueued task", info.Type, info.Payload, info.Queue, info.MaxRetry)

	return nil
}

func (distributor *commentTaskDistributor) DistributeTaskDislikeComment(ctx context.Context, payload *DislikeCommentPayload, opts ...asynq.Option) error {
	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("failed to marshal task payload: %w", err)
	}

	task := asynq.NewTask(TypeDislikeCommentTask, jsonPayload, opts...)
	info, err := distributor.client.EnqueueContext(ctx, task)
	if err != nil {
		return fmt.Errorf("failed to enqueue task: %w", err)
	}

	distributor.logger.Infof("type=%s, payload=%s, queue=%s, maxRetry=%d enqueued task", info.Type, info.Payload, info.Queue, info.MaxRetry)

	return nil
}
