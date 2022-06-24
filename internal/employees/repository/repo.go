package repository

import (
	"github.com/BON4/employees/internal/models"
	kvStore "github.com/BON4/employees/internal/store"
)

type TreeRepo struct {
	store kvStore.Store
	repo  *models.EmpMapTree
}

func NewTreeRepo(s kvStore.Store) *TreeRepo {
	return &TreeRepo{store: s, repo: models.NewEmpMapTree()}
}
