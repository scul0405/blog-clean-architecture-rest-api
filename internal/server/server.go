package server

import (
	"context"
	"github.com/hibiken/asynq"
	"github.com/jmoiron/sqlx"
	"github.com/labstack/echo/v4"
	"github.com/minio/minio-go/v7"
	"github.com/redis/go-redis/v9"
	"github.com/scul0405/blog-clean-architecture-rest-api/config"
	asynqPkg "github.com/scul0405/blog-clean-architecture-rest-api/pkg/asynq"
	"github.com/scul0405/blog-clean-architecture-rest-api/pkg/logger"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

const (
	ctxTimeout     = 5
	maxHeaderBytes = 1 << 20
)

type Server struct {
	echo          *echo.Echo
	cfg           *config.Config
	db            *sqlx.DB
	rdb           *redis.Client
	minioClient   *minio.Client
	asynqClient   *asynq.Client
	taskProcessor *asynqPkg.RedisTaskProcessor
	logger        logger.Logger
}

func NewServer(
	cfg *config.Config,
	db *sqlx.DB,
	rdb *redis.Client,
	minioClient *minio.Client,
	asynqClient *asynq.Client,
	taskProcessor *asynqPkg.RedisTaskProcessor,
	logger logger.Logger) *Server {
	return &Server{echo: echo.New(), cfg: cfg,
		db:            db,
		rdb:           rdb,
		minioClient:   minioClient,
		asynqClient:   asynqClient,
		taskProcessor: taskProcessor,
		logger:        logger}
}

func (s *Server) Run() error {

	server := &http.Server{
		Addr:           s.cfg.Server.Port,
		ReadTimeout:    time.Second * s.cfg.Server.ReadTimeout,
		WriteTimeout:   time.Second * s.cfg.Server.WriteTimeout,
		MaxHeaderBytes: maxHeaderBytes,
	}

	go func() {
		s.logger.Infof("Server is listening on PORT: %s", s.cfg.Server.Port)
		if err := s.echo.StartServer(server); err != nil {
			s.logger.Fatalf("Error starting Server: ", err)
		}
	}()

	if err := s.MapHandlers(s.echo); err != nil {
		return err
	}

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)

	<-quit

	ctx, shutdown := context.WithTimeout(context.Background(), ctxTimeout*time.Second)
	defer shutdown()

	s.logger.Info("Server Exited Properly")
	return s.echo.Server.Shutdown(ctx)
}
