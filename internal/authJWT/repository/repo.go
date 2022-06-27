package repository

import (
	"github.com/BON4/employees/internal/authJWT"
	"github.com/BON4/employees/internal/models"
	kvStore "github.com/BON4/employees/internal/store"
	"github.com/BON4/employees/pkg/jwtService"
)

type JWTRepo struct {
	store      kvStore.Store
	repo       *models.EmpMapTree
	jwtService *jwtService.JWTService
}

// Login implements authJWT.AuthRepo
func (*JWTRepo) Login(username, password string) (string, string, error) {
	panic("unimplemented")
}

// Refresh implements authJWT.AuthRepo
func (*JWTRepo) Refresh(refreshToken string) (string, string, error) {
	panic("unimplemented")
}

// Register implements authJWT.AuthRepo
func (*JWTRepo) Register(username, password string) (string, string, error) {
	panic("unimplemented")
}

func NewJWTRepo(repo *models.EmpMapTree, store kvStore.Store) authJWT.AuthRepo {
	return &JWTRepo{
		store: store,
		repo:  repo,
	}
}
