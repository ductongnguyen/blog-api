package utils

import (
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

// // Validate is user from owner of content
// func ValidateIsOwner(ctx context.Context, creatorID string, logger logger.Logger) error {
// 	user, err := GetUserFromCtx(ctx)
// 	if err != nil {
// 		return err
// 	}

// 	if user.Id.(String) != creatorID {
// 		logger.Errorf(
// 			ctx,
// 			"ValidateIsOwner, userID: %v, creatorID: %v",
// 			user.UserID.String(),
// 			creatorID,
// 		)
// 		return errors.Forbidden
// 	}

// 	return nil
// }

func HashPasswordBcrypt(password string) (string, error) {
	hashedPasswordBytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", fmt.Errorf("failed to hash password: %w", err)
	}

	return string(hashedPasswordBytes), nil
}
