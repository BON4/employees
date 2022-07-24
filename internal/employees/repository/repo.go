package repository

import (
	"context"

	"github.com/BON4/employees/internal/employees"
	"github.com/BON4/employees/internal/models"
	kvStore "github.com/BON4/employees/pkg/store"
)

type TreeRepo struct {
	store kvStore.Store
	repo  *models.EmpMapTree
}

// Delete implements employees.EmpRepository
func (t *TreeRepo) Delete(ctx context.Context, empID string) error {
	if err := t.repo.Delete(ctx, empID); err != nil {
		return err
	}

	if err := t.store.Delete(ctx, empID); err != nil {
		return err
	}

	return t.repo.Save(ctx, t.store)
}

func (t *TreeRepo) Move(ctx context.Context, bossID, empID, toID string) error {
	emp, err := t.repo.FindById(ctx, bossID)
	if err != nil {
		return err
	}

	if _, ok := emp.IsExists(ctx, empID); ok {
		if err := t.repo.Move(ctx, empID, toID); err != nil {
			return err
		}

		return t.repo.Save(ctx, t.store)
	}

	return nil
}

// GetByID implements employees.EmpRepository
func (t *TreeRepo) GetByID(ctx context.Context, empID string) (models.Employee, error) {
	emp, err := t.repo.FindById(ctx, empID)
	if err != nil {
		return models.Employee{}, err
	}

	return emp.Payload, nil
}

func (t *TreeRepo) GetByUsername(ctx context.Context, empUsername string) (models.Employee, error) {
	emp, err := t.repo.FindByUsername(ctx, empUsername)
	if err != nil {
		return models.Employee{}, err
	}

	return emp.Payload, nil
}

// Insert implements employees.EmpRepository
func (t *TreeRepo) Insert(ctx context.Context, empID string, emp models.Employee) error {
	err := t.repo.Insert(ctx, empID, emp)
	if err != nil {
		return err
	}

	return t.repo.Save(ctx, t.store)
}

// Traverse implements employees.EmpRepository
func (t *TreeRepo) Traverse(ctx context.Context, emp *models.Employee, f func(emp models.Employee) error) error {
	empMap, err := t.repo.FindById(ctx, emp.UUID)
	if err != nil {
		return err
	}

	//TODO somehow implement ctx
	return empMap.Traverse(ctx, f)
}

func (t *TreeRepo) Json(ctx context.Context, empID string) (string, error) {
	emp, err := t.repo.FindById(ctx, empID)
	if err != nil {
		return "", err
	}

	return emp.Json()
}

func NewTreeRepo(mapTree *models.EmpMapTree, s kvStore.Store) employees.EmpRepository {
	return &TreeRepo{store: s, repo: mapTree}
}
