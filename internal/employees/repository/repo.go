package repository

import (
	"context"

	"github.com/BON4/employees/internal/employees"
	"github.com/BON4/employees/internal/models"
	kvStore "github.com/BON4/employees/internal/store"
)

type TreeRepo struct {
	store kvStore.Store
	repo  *models.EmpMapTree
}

// Delete implements employees.EmpRepository
func (t *TreeRepo) Delete(ctx context.Context, empID string) error {
	err := t.repo.Delete(empID)
	if err != nil {
		return err
	}

	return t.repo.Save(t.store)
}

// GetByID implements employees.EmpRepository
func (t *TreeRepo) GetByID(ctx context.Context, empID string) (models.Employee, error) {
	emp, err := t.repo.FindById(empID)
	if err != nil {
		return models.Employee{}, err
	}

	return emp.Payload, nil
}

// Insert implements employees.EmpRepository
func (t *TreeRepo) Insert(ctx context.Context, empID string, emp models.Employee) error {
	err := t.repo.Insert(empID, emp)
	if err != nil {
		return err
	}

	return t.repo.Save(t.store)
}

// Traverse implements employees.EmpRepository
func (t *TreeRepo) Traverse(ctx context.Context, emp *models.Employee, f func(emp models.Employee)) error {
	empMap, err := t.repo.FindById(emp.UUID)
	if err != nil {
		return err
	}

	//TODO somehow implement ctx
	empMap.Traverse(f)
	return nil
}

func NewTreeRepo(s kvStore.Store) (employees.EmpRepository, error) {
	repo := models.NewEmpMapTree()
	err := repo.Load(s)
	return &TreeRepo{store: s, repo: repo}, err
}
