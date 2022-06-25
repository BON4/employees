package employees

import (
	"context"

	"github.com/BON4/employees/internal/models"
)

type EmpRepository interface {
	Insert(ctx context.Context, empID string, emp models.Employee) error
	GetByID(ctx context.Context, empID string) (models.Employee, error)
	Delete(ctx context.Context, empID string) error
	//TODO Cretare proper imutable algorithm for traverse
	//Or method To HTML/String tree
	Traverse(ctx context.Context, emp *models.Employee, f func(emp models.Employee)) error
}
