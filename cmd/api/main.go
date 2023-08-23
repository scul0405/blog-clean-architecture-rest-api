package main

import (
	"github.com/hibiken/asynq"
	"github.com/opentracing/opentracing-go"
	"github.com/scul0405/blog-clean-architecture-rest-api/config"
	_ "github.com/scul0405/blog-clean-architecture-rest-api/docs"
	"github.com/scul0405/blog-clean-architecture-rest-api/internal/server"
	asynqPkg "github.com/scul0405/blog-clean-architecture-rest-api/pkg/asynq"
	"github.com/scul0405/blog-clean-architecture-rest-api/pkg/db/minio"
	"github.com/scul0405/blog-clean-architecture-rest-api/pkg/db/postgres"
	"github.com/scul0405/blog-clean-architecture-rest-api/pkg/db/redis"
	"github.com/scul0405/blog-clean-architecture-rest-api/pkg/jaeger"
	"github.com/scul0405/blog-clean-architecture-rest-api/pkg/logger"
	"log"
)

// @version 1.0
// @title Blog Clean Architecture Rest API Server
// @description Simple server written by Golang
// @contact.name Duy Truong
// @contact.url https://github.com/scul0405
// @contact.email vldtruong1221@gmail.com
// @BasePath /api/v1
// @securityDefinitions.apikey Access Token
// @in header
// @name Authorization
func main() {
	log.Println("Starting api server")

	cfgFile, err := config.LoadConfig("./config/config")
	if err != nil {
		log.Fatalf("LoadConfig: %v", err)
	}

	cfg, err := config.ParseConfig(cfgFile)
	if err != nil {
		log.Fatalf("ParseConfig: %v", err)
	}

	// Logger
	appLogger := logger.NewApiLogger(cfg)
	appLogger.InitLogger()
	appLogger.Infof("AppVersion: %s, LogLevel: %s, Mode: %s", cfg.Server.AppVersion, cfg.Logger.Level, cfg.Server.Mode)

	// Database
	psqlDB, err := postgres.NewPsqlDB(cfg)
	if err != nil {
		appLogger.Fatalf("Postgresql init: %s", err)
	} else {
		appLogger.Infof("Postgres connected, Status: %#v", psqlDB.Stats())
	}
	defer psqlDB.Close()

	redisClient := redis.NewRedisClient(cfg)
	defer redisClient.Close()
	appLogger.Info("Redis connected")

	minioClient, err := minio.NewMinioClient(cfg)
	if err != nil {
		appLogger.Infof("Minio client init: %v", err)
	}

	// Jaeger
	tracer, closer, err := jaeger.InitJaeger(cfg)
	if err != nil {
		appLogger.Fatal("cannot create tracer", err)
	}
	appLogger.Info("Jaeger connected")

	opentracing.SetGlobalTracer(tracer)
	defer closer.Close()
	appLogger.Info("Opentracing connected")

	// Asynq
	asynqClient := asynqPkg.NewAsynqClient(asynq.RedisClientOpt{
		Addr: cfg.Asynq.AsynqEndpoint,
	})

	taskProcessor := asynqPkg.NewRedisTaskProcessor(asynq.RedisClientOpt{
		Addr: cfg.Asynq.AsynqEndpoint,
	}, appLogger)

	s := server.NewServer(cfg, psqlDB, redisClient, minioClient, asynqClient, taskProcessor, appLogger)
	if err = s.Run(); err != nil {
		log.Fatal(err)
	}
}
