package middleware

import (
	"github.com/scul0405/blog-clean-architecture-rest-api/config"
	"github.com/scul0405/blog-clean-architecture-rest-api/internal/auth"
	"github.com/scul0405/blog-clean-architecture-rest-api/pkg/logger"
)

type MiddlewareManager struct {
	authUC auth.UseCase
	cfg    *config.Config
	logger logger.Logger
}

func NewMiddlewareManager(authUC auth.UseCase, cfg *config.Config, logger logger.Logger) *MiddlewareManager {
	return &MiddlewareManager{
		authUC: authUC,
		cfg:    cfg,
		logger: logger,
	}
}
