package repository

import (
	"context"
	"testing"

	"github.com/BON4/employees/internal/store"
)

func TestJson(t *testing.T) {
	store := store.NewStore()
	repo, err := NewTreeRepoDEBUG(store)
	if err != nil {
		t.Error(err)
	}

	empJson, err := repo.Json(context.Background(), "admin")

	if err != nil {
		t.Error(err)
	}

	t.Log(empJson)
}
