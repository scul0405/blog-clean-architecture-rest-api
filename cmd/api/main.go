package main

import (
	"github.com/opentracing/opentracing-go"
	"github.com/scul0405/blog-clean-architecture-rest-api/config"
	"github.com/scul0405/blog-clean-architecture-rest-api/internal/server"
	"github.com/scul0405/blog-clean-architecture-rest-api/pkg/db/minio"
	"github.com/scul0405/blog-clean-architecture-rest-api/pkg/db/postgres"
	"github.com/scul0405/blog-clean-architecture-rest-api/pkg/db/redis"
	"github.com/scul0405/blog-clean-architecture-rest-api/pkg/jaeger"
	"github.com/scul0405/blog-clean-architecture-rest-api/pkg/logger"
	"log"
)

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

	s := server.NewServer(cfg, psqlDB, redisClient, minioClient, appLogger)
	if err = s.Run(); err != nil {
		log.Fatal(err)
	}
}
