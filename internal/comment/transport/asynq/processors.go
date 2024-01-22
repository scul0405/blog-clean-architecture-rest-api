package asynq

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/hibiken/asynq"
	"github.com/scul0405/blog-clean-architecture-rest-api/internal/comment"
	"github.com/scul0405/blog-clean-architecture-rest-api/internal/models"
	"github.com/scul0405/blog-clean-architecture-rest-api/pkg/logger"
)

type CommentProcessor interface {
	ProcessTaskLikeComment(ctx context.Context, t *asynq.Task) error
	ProcessTaskDislikeComment(ctx context.Context, t *asynq.Task) error
}

type commentProcessor struct {
	commentUC comment.UseCase
	logger    logger.Logger
}

func NewCommentProcessor(commentUC comment.UseCase, logger logger.Logger) CommentProcessor {
	return &commentProcessor{
		commentUC: commentUC,
		logger:    logger,
	}
}

func (p *commentProcessor) ProcessTaskLikeComment(ctx context.Context, t *asynq.Task) error {
	var payload LikeCommentPayload
	if err := json.Unmarshal(t.Payload(), &payload); err != nil {
		return fmt.Errorf("json.Unmarshal failed: %v: %w", err, asynq.SkipRetry)
	}

	userComment := &models.UserComments{
		CommentID: payload.CommentID,
		UserID:    payload.UserUID,
	}

	err := p.commentUC.Like(ctx, userComment)
	if err != nil {
		return err
	}

	return nil
}

func (p *commentProcessor) ProcessTaskDislikeComment(ctx context.Context, t *asynq.Task) error {
	var payload DislikeCommentPayload
	if err := json.Unmarshal(t.Payload(), &payload); err != nil {
		return fmt.Errorf("json.Unmarshal failed: %v: %w", err, asynq.SkipRetry)
	}

	userComment := &models.UserComments{
		CommentID: payload.CommentID,
		UserID:    payload.UserUID,
	}

	err := p.commentUC.Dislike(ctx, userComment)
	if err != nil {
		return err
	}

	return nil
}
