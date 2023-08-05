package repository

import (
	"context"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/scul0405/blog-clean-architecture-rest-api/internal/models"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestAuthRepo_Register(t *testing.T) {
	t.Parallel()

	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	require.NoError(t, err)
	defer db.Close()

	sqlxDB := sqlx.NewDb(db, "sqlmock")
	defer sqlxDB.Close()

	authRepo := NewAuthRepository(sqlxDB)

	t.Run("Register", func(t *testing.T) {
		gender := "male"
		role := "user"

		rows := sqlmock.NewRows([]string{"first_name", "last_name", "password", "email", "role", "gender"}).AddRow(
			"Liem", "Le", "123123", "liemledeptrai@gmail.com", "user", &gender)

		user := &models.User{
			FirstName: "Liem",
			LastName:  "Le",
			Email:     "liemledeptrai@gmail.com",
			Password:  "123123",
			Role:      &role,
			Gender:    &gender,
		}

		mock.ExpectQuery(createUserQuery).WithArgs(&user.FirstName, &user.LastName, &user.Email,
			&user.Password, &user.Role, &user.About, &user.Avatar, &user.PhoneNumber, &user.Address, &user.City,
			&user.Gender, &user.Postcode, &user.Birthday).WillReturnRows(rows)

		createdUser, err := authRepo.Register(context.Background(), user)

		require.NoError(t, err)
		require.NotNil(t, createdUser)
		require.Equal(t, createdUser, user)
	})
}

func TestAuthRepo_GetByID(t *testing.T) {
	t.Parallel()

	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	require.NoError(t, err)
	defer db.Close()

	sqlxDB := sqlx.NewDb(db, "sqlmock")
	defer sqlxDB.Close()

	authRepo := NewAuthRepository(sqlxDB)

	t.Run("GetByID", func(t *testing.T) {
		uID := uuid.New()

		rows := sqlmock.NewRows([]string{"user_id", "first_name", "last_name", "email"}).AddRow(
			uID, "Liem", "Le", "liemledeptrai@gmail.com")

		testUser := &models.User{
			UserID:    uID,
			FirstName: "Liem",
			LastName:  "Le",
			Email:     "liemledeptrai@gmail.com",
		}

		mock.ExpectQuery(getUserQuery).WithArgs(uID).WillReturnRows(rows)

		user, err := authRepo.GetByID(context.Background(), uID)

		require.NoError(t, err)
		require.Equal(t, testUser, user)
	})
}

func TestAuthRepo_FindByEmail(t *testing.T) {
	t.Parallel()

	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	require.NoError(t, err)
	defer db.Close()

	sqlxDB := sqlx.NewDb(db, "sqlmock")
	defer sqlxDB.Close()

	authRepo := NewAuthRepository(sqlxDB)

	t.Run("FindByEmail", func(t *testing.T) {
		uID := uuid.New()

		rows := sqlmock.NewRows([]string{"user_id", "first_name", "last_name", "email"}).AddRow(
			uID, "Liem", "Le", "liemledeptrai@gmail.com")

		testUser := &models.User{
			UserID:    uID,
			FirstName: "Liem",
			LastName:  "Le",
			Email:     "liemledeptrai@gmail.com",
		}

		mock.ExpectQuery(getUserByEmailQuery).WithArgs(testUser.Email).WillReturnRows(rows)

		user, err := authRepo.FindByEmail(context.Background(), testUser.Email)

		require.NoError(t, err)
		require.NotNil(t, user)
		require.Equal(t, user, testUser)
	})
}
