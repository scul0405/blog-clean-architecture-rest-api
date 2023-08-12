package repository

import (
	"context"
	"database/sql"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/opentracing/opentracing-go"
	"github.com/pkg/errors"
	"github.com/scul0405/blog-clean-architecture-rest-api/internal/comment"
	"github.com/scul0405/blog-clean-architecture-rest-api/internal/models"
	"github.com/scul0405/blog-clean-architecture-rest-api/pkg/utils"
)

type commentRepo struct {
	db *sqlx.DB
}

func NewCommentRepository(db *sqlx.DB) comment.Repository {
	return &commentRepo{db: db}
}

func (r *commentRepo) Create(ctx context.Context, comment *models.Comment) (*models.Comment, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "commentRepo.Create")
	defer span.Finish()

	var c models.Comment
	if err := r.db.QueryRowxContext(ctx, createCommentQuery, &comment.AuthorID, &comment.BlogID, &comment.Message).StructScan(&c); err != nil {
		return nil, errors.Wrap(err, "commentRepo.Create.StructScan")
	}

	return &c, nil
}

func (r *commentRepo) GetByID(ctx context.Context, id uuid.UUID) (*models.CommentBase, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "commentRepo.GetByID")
	defer span.Finish()

	var c models.CommentBase
	if err := r.db.QueryRowxContext(ctx, getCommentByIDQuery, id).StructScan(&c); err != nil {
		return nil, errors.Wrap(err, "commentRepo.GetByID.StructScan")
	}

	return &c, nil
}

func (r *commentRepo) Update(ctx context.Context, comment *models.CommentBase) (*models.CommentBase, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "commentRepo.Update")
	defer span.Finish()

	var c models.CommentBase
	if err := r.db.QueryRowxContext(ctx, updateCommentQuery, &comment.Message, &comment.CommentID).StructScan(&c); err != nil {
		return nil, errors.Wrap(err, "commentRepo.Update.StructScan")
	}

	return &c, nil
}

func (r *commentRepo) Delete(ctx context.Context, id uuid.UUID) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "commentRepo.Delete")
	defer span.Finish()

	result, err := r.db.ExecContext(ctx, deleteCommentQuery, id)
	if err != nil {
		return errors.Wrap(err, "commentRepo.Delete.StructScan")
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return errors.Wrap(err, "commentRepo.Delete.RowsAffected")
	}
	if rowsAffected == 0 {
		return errors.Wrap(sql.ErrNoRows, "commentRepo.Delete.rowsAffected")
	}

	return nil
}

func (r *commentRepo) List(ctx context.Context, blogID uuid.UUID, pq *utils.PaginationQuery) (*models.CommentsList, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "commentRepo.List")
	defer span.Finish()

	var totalCount int
	if err := r.db.GetContext(ctx, &totalCount, getTotalCountByBlogIDQuery, blogID); err != nil {
		return nil, errors.Wrap(err, "commentRepo.List.GetContext.totalCount")
	}

	if totalCount == 0 {
		return &models.CommentsList{
			TotalCount: totalCount,
			TotalPages: utils.GetTotalPages(totalCount, pq.GetSize()),
			Page:       pq.GetPage(),
			Size:       pq.GetSize(),
			HasMore:    utils.GetHasMore(pq.GetPage(), totalCount, pq.GetSize()),
			Comments:   make([]*models.CommentBase, 0),
		}, nil
	}

	var commentsList = make([]*models.CommentBase, 0, pq.GetSize())
	rows, err := r.db.QueryxContext(ctx, listCommentsByBlogIDQuery, blogID, pq.GetOffset(), pq.GetLimit())
	if err != nil {
		return nil, errors.Wrap(err, "commentRepo.List.QueryxContext")
	}
	defer rows.Close()

	for rows.Next() {
		n := &models.CommentBase{}
		if err = rows.StructScan(n); err != nil {
			return nil, errors.Wrap(err, "commentRepo.List.StructScan")
		}
		commentsList = append(commentsList, n)
	}

	if err = rows.Err(); err != nil {
		return nil, errors.Wrap(err, "commentRepo.List.rows.Err")
	}

	return &models.CommentsList{
		TotalCount: totalCount,
		TotalPages: utils.GetTotalPages(totalCount, pq.GetSize()),
		Page:       pq.GetPage(),
		Size:       pq.GetSize(),
		HasMore:    utils.GetHasMore(pq.GetPage(), totalCount, pq.GetSize()),
		Comments:   commentsList,
	}, nil
}
