package repository

import (
	"context"
	"testing"

	"github.com/BON4/employees/internal/models"
	"github.com/BON4/employees/pkg/store"
)

func TestJson(t *testing.T) {
	store := store.NewMapStore()
	eT := models.NewEmpMapTreeDEBUG()
	repo := NewTreeRepo(eT, store)

	emp, err := repo.GetByUsername(context.Background(), "1")
	if err != nil {
		t.Error(err)
	}

	empJson, err := repo.Json(context.Background(), emp.UUID)

	if err != nil {
		t.Error(err)
	}

	t.Log(empJson)
}
