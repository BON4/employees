package employeesTree

import (
	"testing"

	"github.com/BON4/employees/models"
	kvStore "github.com/BON4/employees/store"
)

var mapStore kvStore.Store

func TestMain(m *testing.M) {
	mapStore = kvStore.NewMapStore()
	m.Run()
}

func createDumbEmplMap() *EmpMapTree {
	return NewEmpMapTree(NewEmployeeMap(models.NewEmployee("admin", "", models.Admin),
		[]*EmployeeMap{
			NewEmployeeMap(models.NewEmployee("1", "", models.Boss),
				NewEmployeeMap(models.NewEmployee("4", "", models.Boss),
					NewEmployeeMap(models.NewEmployee("9", "", models.Regular))),
				NewEmployeeMap(models.NewEmployee("5", "", models.Regular))),
			NewEmployeeMap(models.NewEmployee("2", "", models.Boss),
				NewEmployeeMap(models.NewEmployee("6", "", models.Regular)),
				NewEmployeeMap(models.NewEmployee("7", "", models.Regular))),
			NewEmployeeMap(models.NewEmployee("3", "", models.Boss),
				NewEmployeeMap(models.NewEmployee("8", "", models.Regular)))}...))
}

func TestEmpMapNew(t *testing.T) {
	emps := createDumbEmplMap()
	t.Logf("\n%s", emps)
}

func TestMapTreeLoadSave(t *testing.T) {
	emps := createDumbEmplMap()

	err := emps.Save(mapStore)
	if err != nil {
		t.Error(err)
	}

	loadedTree := NewEmpMapTree(NewEmployeeMap(models.NewEmployee("admin", "", models.Admin)))
	if err := loadedTree.Load(mapStore); err != nil {
		t.Error(err)
	}

	if len(emps.String()) != len(loadedTree.String()) {
		t.Errorf("Trees dont match:")
		t.Logf("Trees len: %d - %d", len(emps.String()), len(loadedTree.String()))
	}
}

func TestEmpMapDelete(t *testing.T) {
	empTree := createDumbEmplMap()
	err := empTree.Save(mapStore)
	if err != nil {
		t.Error(err)
	}
	t.Logf("\n%s", empTree)

	fEmp, err := empTree.FindByUsername("8")

	if err != nil {
		t.Error(err)
	}

	if err := empTree.Delete(fEmp.Payload.UUID); err != nil {
		t.Error(err)
	}

	fAdmin, err := empTree.FindByUsername("admin")

	if err != nil {
		t.Error(err)
	}

	if err := empTree.Delete(fAdmin.Payload.UUID); err == nil {
		t.Error("Cant delete admin")
	}

	t.Logf("\n%s", empTree)
}

func TestEmpMapInsert(t *testing.T) {
	empTree := createDumbEmplMap()
	err := empTree.Save(mapStore)
	if err != nil {
		t.Error(err)
	}

	fEmp, err := empTree.FindByUsername("admin")

	if err != nil {
		t.Error(err)
	}

	if err := empTree.Insert(fEmp.Payload.UUID, NewEmployeeMap(models.NewEmployee("8", "", models.Regular))); err != nil {
		t.Error(err)
	}

	t.Logf("\n%s", empTree)
}

func TestEmpMapFind(t *testing.T) {
	empTree := createDumbEmplMap()
	err := empTree.Save(mapStore)
	if err != nil {
		t.Error(err)
	}

	fEmp, err := empTree.FindByUsername("1")

	if err != nil {
		t.Error(err)
	}

	t.Logf("\n%s", fEmp)
}
