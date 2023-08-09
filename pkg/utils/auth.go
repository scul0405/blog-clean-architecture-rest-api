package utils

import (
	"context"
	httpErrors "github.com/scul0405/blog-clean-architecture-rest-api/pkg/http_errors"
	"github.com/scul0405/blog-clean-architecture-rest-api/pkg/logger"
)

// ValidateIsOwner validate is user from owner of content
func ValidateIsOwner(ctx context.Context, creatorID string, logger logger.Logger) error {
	userUID, err := GetUserUIDFromCtx(ctx)
	if err != nil {
		return err
	}

	if userUID.String() != creatorID {
		logger.Errorf(
			"ValidateIsOwner, userID: %v, creatorID: %v",
			userUID.String(),
			creatorID,
		)
		return httpErrors.Forbidden
	}

	return nil
}
