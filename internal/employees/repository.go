package employees

import (
	"context"

	"github.com/BON4/employees/internal/models"
)

type EmpRepository interface {
	Insert(ctx context.Context, emp *models.Employee) (*models.Employee, error)
	GetByID(ctx context.Context, empID uint) (*models.Employee, error)
	Delete(ctx context.Context, empID uint) error
}
