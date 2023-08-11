package repository

import (
	"context"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/scul0405/blog-clean-architecture-rest-api/internal/models"
	"github.com/scul0405/blog-clean-architecture-rest-api/pkg/utils"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestBlogRepo_Create(t *testing.T) {
	t.Parallel()

	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	require.NoError(t, err)
	defer db.Close()

	sqlxDB := sqlx.NewDb(db, "sqlmock")
	defer sqlxDB.Close()

	blogRepo := NewBlogRepository(sqlxDB)

	t.Run("Create", func(t *testing.T) {
		authorUID := uuid.New()
		title := "title"
		content := "content"

		rows := sqlmock.NewRows([]string{"author_id", "title", "content"}).
			AddRow(authorUID, title, content)

		blog := &models.Blog{
			AuthorID: authorUID,
			Title:    title,
			Content:  content,
		}

		mock.ExpectQuery(createBlogQuery).
			WithArgs(blog.AuthorID,
				blog.Title,
				blog.Content,
				blog.ImageURL,
				blog.Category).
			WillReturnRows(rows)

		createdBlog, err := blogRepo.Create(context.Background(), blog)

		require.NoError(t, err)
		require.NotNil(t, createdBlog)
		require.Equal(t, createdBlog.AuthorID, blog.AuthorID)
	})
}

func TestBlogRepo_GetByID(t *testing.T) {
	t.Parallel()

	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	require.NoError(t, err)
	defer db.Close()

	sqlxDB := sqlx.NewDb(db, "sqlmock")
	defer sqlxDB.Close()

	blogRepo := NewBlogRepository(sqlxDB)

	t.Run("GetByID", func(t *testing.T) {
		blogUID := uuid.New()
		authorUID := uuid.New()
		title := "title"
		content := "content"

		rows := sqlmock.NewRows([]string{"blog_id", "author_id", "title", "content"}).
			AddRow(blogUID, authorUID, title, content)

		blogTest := &models.Blog{
			BlogID:   blogUID,
			AuthorID: authorUID,
			Title:    title,
			Content:  content,
		}

		mock.ExpectQuery(getBlogByIDQuery).WithArgs(blogUID).WillReturnRows(rows)

		blog, err := blogRepo.GetByID(context.Background(), blogUID)

		require.NoError(t, err)
		require.NotNil(t, blog)
		require.Equal(t, blogTest.BlogID, blog.BlogID)
		require.Equal(t, blogTest.AuthorID, blog.AuthorID)
	})
}

func TestBlogRepo_Update(t *testing.T) {
	t.Parallel()

	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	require.NoError(t, err)
	defer db.Close()

	sqlxDB := sqlx.NewDb(db, "sqlmock")
	defer sqlxDB.Close()

	blogRepo := NewBlogRepository(sqlxDB)

	t.Run("Update", func(t *testing.T) {
		blogUID := uuid.New()
		authorUID := uuid.New()
		title := "update title"
		content := "content"

		rows := sqlmock.NewRows([]string{"blog_id", "author_id", "title", "content"}).
			AddRow(blogUID, authorUID, title, content)

		blog := &models.BlogBase{
			BlogID:   blogUID,
			AuthorID: authorUID,
			Title:    title,
			Content:  content,
		}

		mock.ExpectQuery(updateBlogQuery).WithArgs(
			blog.Title,
			blog.Content,
			blog.ImageURL,
			blog.Category,
			blog.BlogID).
			WillReturnRows(rows)

		updatedBlog, err := blogRepo.Update(context.Background(), blog)

		require.NoError(t, err)
		require.NotNil(t, updatedBlog)
		require.Equal(t, updatedBlog.BlogID, blog.BlogID)
		require.Equal(t, updatedBlog.Title, blog.Title)
	})
}

func TestBlogRepo_Delete(t *testing.T) {
	t.Parallel()

	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	require.NoError(t, err)
	defer db.Close()

	sqlxDB := sqlx.NewDb(db, "sqlmock")
	defer sqlxDB.Close()

	blogRepo := NewBlogRepository(sqlxDB)

	t.Run("Delete", func(t *testing.T) {
		blogUID := uuid.New()
		mock.ExpectExec(deleteBlogQuery).WithArgs(blogUID).WillReturnResult(sqlmock.NewResult(1, 1))

		err := blogRepo.Delete(context.Background(), blogUID)

		require.NoError(t, err)
	})
}

func TestBlogRepo_List(t *testing.T) {
	t.Parallel()

	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	require.NoError(t, err)
	defer db.Close()

	sqlxDB := sqlx.NewDb(db, "sqlmock")
	defer sqlxDB.Close()

	blogRepo := NewBlogRepository(sqlxDB)

	t.Run("List", func(t *testing.T) {
		expectedCount := 5
		expectedPage := 1
		expectedSize := 2
		expectedTotalPage := 3
		mockBlogs := make([]models.BlogBase, expectedCount, expectedCount)

		for i := 0; i < expectedCount; i++ {
			mockBlogs[i] = models.BlogBase{
				BlogID:   uuid.New(),
				AuthorID: uuid.New(),
				Title:    "title",
				Content:  "content",
			}
		}

		rows := sqlmock.NewRows([]string{"blog_id", "author_id", "title", "content"})
		for _, blog := range mockBlogs {
			rows.AddRow(blog.BlogID, blog.AuthorID, blog.Title, blog.Content)
		}

		countRows := sqlmock.NewRows([]string{"count"}).AddRow(expectedCount)
		mock.ExpectQuery(getTotalCountQuery).WillReturnRows(countRows)

		pq := utils.PaginationQuery{
			Size: expectedSize,
			Page: expectedPage,
		}
		mock.ExpectQuery(listBlogsQuery).WithArgs(pq.GetOffset(), pq.GetLimit()).WillReturnRows(rows)

		listBlogs, err := blogRepo.List(context.Background(), &pq)

		require.NoError(t, err)
		require.NotNil(t, listBlogs)
		require.Equal(t, expectedCount, listBlogs.TotalCount)
		require.Equal(t, expectedTotalPage, listBlogs.TotalPages)
		require.Equal(t, expectedPage, listBlogs.Page)
		require.Equal(t, expectedSize, listBlogs.Size)
		require.True(t, listBlogs.HasMore)
		require.NotNil(t, listBlogs.Blogs)
	})
}
