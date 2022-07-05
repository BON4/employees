package repository

import (
	"context"
	"testing"
	"time"

	"github.com/BON4/employees/internal/authJWT"
	"github.com/BON4/employees/internal/employees"
	"github.com/BON4/employees/internal/employees/repository"
	"github.com/BON4/employees/internal/models"
	"github.com/BON4/employees/pkg/jwtService"
	"github.com/BON4/employees/pkg/store"
)

var accessKey = []byte("123456789")
var refreshKey = []byte("987654321")

var repo authJWT.AuthRepo
var treeRepo employees.EmpRepository
var jwtM *jwtService.JWTService

func TestMain(m *testing.M) {
	s := store.NewMapStore()
	eT := models.NewEmpMapTreeDEBUG()
	treeRepo = repository.NewTreeRepo(eT, s)
	jwtConf := jwtService.NewJWTConfig(accessKey, refreshKey, time.Minute, time.Hour)
	jwtM = jwtService.NewJWTService(jwtConf)
	repo = NewJWTRepo(treeRepo, s, jwtM)
	m.Run()
}

func TestLogin(t *testing.T) {
	emp, err := treeRepo.GetByUsername(context.Background(), "3")
	if err != nil {
		t.Error(err)
	}

	acess, refresh, err := repo.Login(context.Background(), emp.Username, emp.Password)
	if err != nil {
		t.Error(err)
	}

	_, err = jwtM.VerifyAcess(acess)
	if err != nil {
		t.Error(err)
	}

	_, err = jwtM.VerifyRefresh(refresh)
	if err != nil {
		t.Error(err)
	}

	t.Logf("acess:%s\nrefresh:%s\n", acess, refresh)
}

func TestRefresh(t *testing.T) {
	emp, err := treeRepo.GetByUsername(context.Background(), "3")
	if err != nil {
		t.Error(err)
	}

	acess, refresh, err := repo.Login(context.Background(), emp.Username, emp.Password)
	if err != nil {
		t.Error(err)
	}

	_, err = jwtM.VerifyAcess(acess)
	if err != nil {
		t.Error(err)
	}

	_, err = jwtM.VerifyRefresh(refresh)
	if err != nil {
		t.Error(err)
	}

	t.Logf("OLDacess:%s\nOLDrefresh:%s\n", acess, refresh)

	time.Sleep(time.Second * 14)

	acess, refresh, err = repo.Refresh(context.Background(), refresh)
	if err != nil {
		t.Error(err)
	}

	_, err = jwtM.VerifyAcess(acess)
	if err != nil {
		t.Error(err)
	}

	_, err = jwtM.VerifyRefresh(refresh)
	if err != nil {
		t.Error(err)
	}

	t.Logf("NEWacess:%s\nNEWrefresh:%s\n", acess, refresh)
}

func TestRegister(t *testing.T) {
	acess, refresh, err := repo.Register(context.Background(), "test", "testPassword")
	if err != nil {
		t.Error(err)
	}

	_, err = jwtM.VerifyAcess(acess)
	if err != nil {
		t.Error(err)
	}

	_, err = jwtM.VerifyRefresh(refresh)
	if err != nil {
		t.Error(err)
	}

	t.Logf("acess:%s\nrefresh:%s\n", acess, refresh)
}
