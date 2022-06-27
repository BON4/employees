package authJWT

type AuthRepo interface {
	// Auth user with credentials and send access and refresh token
	Login(username, password string) (string, string, error)

	// Only with using session manager
	// Logout()

	// Refresh takes refresh token and send new refresh and access tokens
	Refresh(refreshToken string) (string, string, error)

	// Create new object with provided username and password and send access and refresh token
	Register(username, password string) (string, string, error)
}
