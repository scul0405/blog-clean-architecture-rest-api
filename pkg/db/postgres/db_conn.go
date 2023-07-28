package postgres

import (
	"fmt"
	_ "github.com/jackc/pgx/stdlib"
	"github.com/jmoiron/sqlx"
	"github.com/scul0405/blog-clean-architecture-rest-api/config"
	"time"
)

const (
	maxOpenConns    = 60
	maxIdleConns    = 30
	connMaxIdleTime = 20
	connMaxLifetime = 120
)

func NewPsqlDB(c *config.Config) (*sqlx.DB, error) {
	dataSourceName := fmt.Sprintf("host=%s port=%s user=%s dbname=%s sslmode=disable password=%s",
		c.Postgres.PostgresqlHost,
		c.Postgres.PostgresqlPort,
		c.Postgres.PostgresqlUser,
		c.Postgres.PostgresqlDbname,
		c.Postgres.PostgresqlPassword,
	)

	db, err := sqlx.Connect(c.Postgres.PgDriver, dataSourceName)
	if err != nil {
		return nil, err
	}

	db.SetMaxOpenConns(maxOpenConns)
	db.SetConnMaxLifetime(connMaxLifetime * time.Second)
	db.SetMaxIdleConns(maxIdleConns)
	db.SetConnMaxIdleTime(connMaxIdleTime * time.Second)
	if err = db.Ping(); err != nil {
		return nil, err
	}

	return db, nil
}
