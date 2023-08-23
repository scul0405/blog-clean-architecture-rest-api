package asynq

import (
	"github.com/google/uuid"
)

const (
	TypeLikeCommentTask    = "comment:like"
	TypeDislikeCommentTask = "comment:dislike"
)

type LikeCommentPayload struct {
	UserUID   uuid.UUID
	CommentID uuid.UUID
}

type DislikeCommentPayload struct {
	UserUID   uuid.UUID
	CommentID uuid.UUID
}
