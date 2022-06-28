package authJWT

import (
	"context"
)

type AuthRepo interface {
	// Auth user with credentials and send access and refresh token
	Login(ctx context.Context, username, password string) (string, string, error)

	// Only with using session manager
	// Logout()

	// Refresh takes refresh token and send new refresh and access tokens
	Refresh(ctx context.Context, refreshToken string) (string, string, error)

	// Create new object with provided username and password and send access and refresh token
	Register(ctx context.Context, username, password string) (string, string, error)
}
