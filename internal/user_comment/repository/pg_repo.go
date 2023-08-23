package repository

import (
	"context"
	"database/sql"
	"github.com/jmoiron/sqlx"
	"github.com/opentracing/opentracing-go"
	"github.com/pkg/errors"
	"github.com/scul0405/blog-clean-architecture-rest-api/internal/models"
	"github.com/scul0405/blog-clean-architecture-rest-api/internal/user_comment"
)

type userCommentRepo struct {
	db *sqlx.DB
}

func NewUserCommentRepository(db *sqlx.DB) user_comment.Repository {
	return &userCommentRepo{db: db}
}

func (r *userCommentRepo) Create(ctx context.Context, userComment *models.UserComments) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "userCommentRepo.Create")
	defer span.Finish()

	result, err := r.db.ExecContext(ctx, createUserCommentQuery, userComment.UserID, userComment.CommentID)
	if err != nil {
		return errors.Wrap(err, "userCommentRepo.Create.StructScan")
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return errors.Wrap(err, "userCommentRepo.Create.RowsAffected")
	}
	if rowsAffected == 0 {
		return errors.Wrap(sql.ErrNoRows, "userCommentRepo.Create.rowsAffected")
	}

	return nil
}

func (r *userCommentRepo) GetByID(ctx context.Context, userComment *models.UserComments) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "userCommentRepo.GetByID")
	defer span.Finish()

	var c models.UserComments
	if err := r.db.QueryRowxContext(ctx, getUserCommentQuery, userComment.UserID, userComment.CommentID).StructScan(&c); err != nil {
		return errors.Wrap(err, "userCommentRepo.GetByID.StructScan")
	}

	return nil
}

func (r *userCommentRepo) Delete(ctx context.Context, userComment *models.UserComments) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "userCommentRepo.Delete")
	defer span.Finish()

	result, err := r.db.ExecContext(ctx, deleteUserCommentQuery, userComment.UserID, userComment.CommentID)
	if err != nil {
		return errors.Wrap(err, "userCommentRepo.Delete.StructScan")
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return errors.Wrap(err, "userCommentRepo.Delete.RowsAffected")
	}
	if rowsAffected == 0 {
		return errors.Wrap(sql.ErrNoRows, "userCommentRepo.Delete.rowsAffected")
	}

	return nil
}
