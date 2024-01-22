package asynq

import asynqPkg "github.com/scul0405/blog-clean-architecture-rest-api/pkg/asynq"

func MapHandlers(tp *asynqPkg.RedisTaskProcessor, cp CommentProcessor) {
	tp.RegisterHandler(TypeLikeCommentTask, cp.ProcessTaskLikeComment)
	tp.RegisterHandler(TypeDislikeCommentTask, cp.ProcessTaskDislikeComment)
}
