package repository

import (
	"context"
	"database/sql"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/opentracing/opentracing-go"
	"github.com/pkg/errors"
	"github.com/scul0405/blog-clean-architecture-rest-api/internal/blog"
	"github.com/scul0405/blog-clean-architecture-rest-api/internal/models"
	"github.com/scul0405/blog-clean-architecture-rest-api/pkg/utils"
)

type blogRepo struct {
	db *sqlx.DB
}

func NewBlogRepository(db *sqlx.DB) blog.Repository {
	return &blogRepo{db: db}
}

func (r *blogRepo) Create(ctx context.Context, blog *models.Blog) (*models.BlogBase, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "blogRepo.Create")
	defer span.Finish()

	var b models.BlogBase
	if err := r.db.QueryRowxContext(ctx, createBlogQuery, &blog.AuthorID, &blog.Title, &blog.Content, &blog.ImageURL, &blog.Category).StructScan(&b); err != nil {
		return nil, errors.Wrap(err, "blogRepo.Create.StructScan")
	}

	return &b, nil
}

func (r *blogRepo) GetByID(ctx context.Context, id uuid.UUID) (*models.BlogBase, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "blogRepo.GetByID")
	defer span.Finish()

	var b models.BlogBase
	if err := r.db.QueryRowxContext(ctx, getBlogByIDQuery, id).StructScan(&b); err != nil {
		return nil, errors.Wrap(err, "blogRepo.GetByID.StructScan")
	}

	return &b, nil
}

func (r *blogRepo) Update(ctx context.Context, blog *models.BlogBase) (*models.BlogBase, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "blogRepo.Update")
	defer span.Finish()

	var b models.BlogBase
	if err := r.db.QueryRowxContext(ctx, updateBlogQuery, &blog.Title, &blog.Content, &blog.ImageURL, &blog.Category, &blog.BlogID).StructScan(&b); err != nil {
		return nil, errors.Wrap(err, "blogRepo.Update.StructScan")
	}

	return &b, nil
}

func (r *blogRepo) Delete(ctx context.Context, id uuid.UUID) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "blogRepo.Delete")
	defer span.Finish()

	result, err := r.db.ExecContext(ctx, deleteBlogQuery, id)
	if err != nil {
		return errors.Wrap(err, "blogRepo.Delete.StructScan")
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return errors.Wrap(err, "blogRepo.Delete.RowsAffected")
	}
	if rowsAffected == 0 {
		return errors.Wrap(sql.ErrNoRows, "blogRepo.Delete.rowsAffected")
	}

	return nil
}

func (r *blogRepo) List(ctx context.Context, pq *utils.PaginationQuery) (*models.BlogsList, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "blogRepo.List")
	defer span.Finish()

	var totalCount int
	if err := r.db.GetContext(ctx, &totalCount, getTotalCountQuery); err != nil {
		return nil, errors.Wrap(err, "blogRepo.List.GetContext.totalCount")
	}

	if totalCount == 0 {
		return &models.BlogsList{
			TotalCount: totalCount,
			TotalPages: utils.GetTotalPages(totalCount, pq.GetSize()),
			Page:       pq.GetPage(),
			Size:       pq.GetSize(),
			HasMore:    utils.GetHasMore(pq.GetPage(), totalCount, pq.GetSize()),
			Blogs:      make([]*models.BlogBase, 0),
		}, nil
	}

	// TODO: update order by
	var blogsList = make([]*models.BlogBase, 0, pq.GetSize())
	rows, err := r.db.QueryxContext(ctx, listBlogsQuery, pq.GetOffset(), pq.GetLimit())
	if err != nil {
		return nil, errors.Wrap(err, "blogRepo.List.QueryxContext")
	}
	defer rows.Close()

	for rows.Next() {
		n := &models.BlogBase{}
		if err = rows.StructScan(n); err != nil {
			return nil, errors.Wrap(err, "blogRepo.List.StructScan")
		}
		blogsList = append(blogsList, n)
	}

	if err = rows.Err(); err != nil {
		return nil, errors.Wrap(err, "blogRepo.List.rows.Err")
	}

	return &models.BlogsList{
		TotalCount: totalCount,
		TotalPages: utils.GetTotalPages(totalCount, pq.GetSize()),
		Page:       pq.GetPage(),
		Size:       pq.GetSize(),
		HasMore:    utils.GetHasMore(pq.GetPage(), totalCount, pq.GetSize()),
		Blogs:      blogsList,
	}, nil
}
