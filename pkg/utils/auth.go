package utils

import (
	"context"

	"github.com/ductong169z/blog-api/pkg/errors"
	"github.com/ductong169z/blog-api/pkg/logger"
)

// Validate is user from owner of content
func ValidateIsOwner(ctx context.Context, creatorID string, logger logger.Logger) error {
	user, err := GetUserFromCtx(ctx)
	if err != nil {
		return err
	}

	if user.UserID.String() != creatorID {
		logger.Errorf(
			ctx,
			"ValidateIsOwner, userID: %v, creatorID: %v",
			user.UserID.String(),
			creatorID,
		)
		return errors.Forbidden
	}

	return nil
}
