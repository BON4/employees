package employees

import (
	"context"

	"github.com/BON4/employees/internal/models"
)

type EmpRepository interface {
	Insert(ctx context.Context, empID string, emp models.Employee) error
	Move(ctx context.Context, bossID, empID, toID string) error
	GetByID(ctx context.Context, empID string) (models.Employee, error)
	GetByUsername(ctx context.Context, empID string) (models.Employee, error)
	Delete(ctx context.Context, empID string) error
	Json(ctx context.Context, empID string) (string, error)
	Traverse(ctx context.Context, emp *models.Employee, f func(emp models.Employee) error) error
}
