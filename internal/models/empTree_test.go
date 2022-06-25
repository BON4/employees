package models

import (
	"testing"

	kvStore "github.com/BON4/employees/internal/store"
)

var mapStore kvStore.Store

func TestMain(m *testing.M) {
	mapStore = kvStore.NewMapStore()
	m.Run()
}

func createDumbEmplMap() *EmpMapTree {
	return &EmpMapTree{(NewEmployeeMap(NewEmployee("admin", "", Admin),
		[]*EmployeeMap{
			NewEmployeeMap(NewEmployee("1", "", Boss),
				NewEmployeeMap(NewEmployee("4", "", Boss),
					NewEmployeeMap(NewEmployee("9", "", Regular))),
				NewEmployeeMap(NewEmployee("5", "", Regular))),
			NewEmployeeMap(NewEmployee("2", "", Boss),
				NewEmployeeMap(NewEmployee("6", "", Regular)),
				NewEmployeeMap(NewEmployee("7", "", Regular))),
			NewEmployeeMap(NewEmployee("3", "", Boss),
				NewEmployeeMap(NewEmployee("8", "", Regular)))}...))}
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

	loadedTree := NewEmpMapTree()
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
	e := NewEmployee("8", "", Regular)

	if err := empTree.Insert(fEmp.Payload.UUID, e); err != nil {
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
