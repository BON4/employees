package jwtService

import (
	"errors"
	"time"

	jwt "github.com/golang-jwt/jwt/v4"
)

type JWTService struct {
	cfg *jwtConfig
}

func (j *JWTService) CreateAccess(val map[string]interface{}) (string, error) {
	claims := jwt.MapClaims{}
	for k, v := range val {
		claims[k] = v
	}

	claims["ExpiresAt"] = time.Now().Add(j.cfg.AccessExpireTime).Unix()
	claims["IssuedAt"] = time.Now().Unix()

	token := jwt.NewWithClaims(jwt.SigningMethodHS512, claims)

	out, err := token.SignedString(j.cfg.AccessKey)
	if err != nil {
		return "", err
	}

	return out, nil
}

func (j *JWTService) CreateRefresh(val map[string]interface{}) (string, error) {
	claims := jwt.MapClaims{}
	for k, v := range val {
		claims[k] = v
	}

	claims["ExpiresAt"] = time.Now().Add(j.cfg.RefreshExpireTime).Unix()
	claims["IssuedAt"] = time.Now().Unix()

	token := jwt.NewWithClaims(jwt.SigningMethodHS512, claims)

	out, err := token.SignedString(j.cfg.RefreshKey)
	if err != nil {
		return "", err
	}

	return out, nil
}

func (j *JWTService) VerifyAcess(token string) (map[string]interface{}, error) {
	claims := jwt.MapClaims{}
	parsedToken, err := jwt.ParseWithClaims(token, claims, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing algorithm: " + token.Method.Alg())
		}

		return j.cfg.AccessKey, nil
	})

	if err != nil {
		return claims, err
	}

	if parsedToken.Valid {
		return claims, nil
	}

	return claims, errors.New("invalid token")
}

func (j *JWTService) VerifyRefresh(token string) (map[string]interface{}, error) {
	claims := jwt.MapClaims{}
	parsedToken, err := jwt.ParseWithClaims(token, claims, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing algorithm: " + token.Method.Alg())
		}

		return j.cfg.RefreshKey, nil
	})

	if err != nil {
		return claims, err
	}

	if parsedToken.Valid {
		return claims, nil
	}

	return claims, errors.New("invalid token")
}

func NewJWTService(cfg *jwtConfig) *JWTService {
	return &JWTService{
		cfg: cfg,
	}
}
