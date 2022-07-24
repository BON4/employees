package repository

import (
	"context"

	"github.com/BON4/employees/internal/authJWT"
	"github.com/BON4/employees/internal/employees"
	uerrors "github.com/BON4/employees/internal/errors"
	"github.com/BON4/employees/internal/models"
	"github.com/BON4/employees/pkg/jwtService"
	kvStore "github.com/BON4/employees/pkg/store"
)

type JWTRepo struct {
	store      kvStore.Store
	repo       employees.EmpRepository
	jwtManager *jwtService.JWTService
}

// Login implements authJWT.AuthRepo
func (j *JWTRepo) Login(ctx context.Context, username, password string) (string, string, error) {
	emp, err := j.repo.GetByUsername(ctx, username)
	if err != nil {
		return "", "", err
	}

	if emp.Password != password {
		return "", "", uerrors.NewError[uerrors.UserError]("wrong password")
	}

	access, err := j.jwtManager.CreateAccess(map[string]interface{}{
		"Emp": emp,
	})
	if err != nil {
		return "", "", err
	}

	refresh, err := j.jwtManager.CreateRefresh(map[string]interface{}{
		"Emp": emp,
	})
	if err != nil {
		return "", "", err
	}

	return access, refresh, nil
}

// Refresh implements authJWT.AuthRepo
func (j *JWTRepo) Refresh(ctx context.Context, refreshToken string) (string, string, error) {
	tokData, err := j.jwtManager.VerifyRefresh(refreshToken)
	if err != nil {
		return "", "", err
	}

	if emp, ok := tokData["Emp"]; ok {
		access, err := j.jwtManager.CreateAccess(map[string]interface{}{
			"Emp": emp,
		})
		if err != nil {
			return "", "", err
		}

		refresh, err := j.jwtManager.CreateRefresh(map[string]interface{}{
			"Emp": emp,
		})
		if err != nil {
			return "", "", err
		}

		return access, refresh, nil
	}
	return "", "", uerrors.NewError[uerrors.UserError]("invalid refresh token")
}

// Register implements authJWT.AuthRepo
func (j *JWTRepo) Register(ctx context.Context, username, password string) (string, string, error) {
	emp := models.NewEmployee(username, password, models.Regular)
	if err := j.repo.Insert(ctx, "admin", emp); err != nil {
		return "", "", err
	}

	access, err := j.jwtManager.CreateAccess(map[string]interface{}{
		"Emp": emp,
	})
	if err != nil {
		return "", "", err
	}

	refresh, err := j.jwtManager.CreateRefresh(map[string]interface{}{
		"Emp": emp,
	})
	if err != nil {
		return "", "", err
	}

	return access, refresh, nil
}

func NewJWTRepo(repo employees.EmpRepository, store kvStore.Store, jwtManager *jwtService.JWTService) authJWT.AuthRepo {
	return &JWTRepo{
		store:      store,
		repo:       repo,
		jwtManager: jwtManager,
	}
}
