package employeesTree

import (
	"testing"
)

var mapStore Store

func TestMain(m *testing.M) {
	mapStore = NewMapStore()
	m.Run()
}

func createDumbEmplMap() *EmpMapTree {
	return NewEmpMapTree(NewEmployeeMap("admin", []*EmployeeMap{NewEmployeeMap("1",
		NewEmployeeMap("4",
			NewEmployeeMap("9")),
		NewEmployeeMap("5")),
		NewEmployeeMap("2",
			NewEmployeeMap("6"),
			NewEmployeeMap("7")),
		NewEmployeeMap("3",
			NewEmployeeMap("8"))}...))
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

	loadedTree := NewEmpMapTree(NewEmployeeMap("admin"))
	if err := loadedTree.Load(mapStore); err != nil {
		t.Error(err)
	}

	if len(emps.String()) != len(loadedTree.String()) {
		t.Errorf("Trees dont match:")
	}
}

func TestEmpMapDelete(t *testing.T) {
	empTree := createDumbEmplMap()
	err := empTree.Save(mapStore)
	if err != nil {
		t.Error(err)
	}
	t.Logf("\n%s", empTree)

	if err := empTree.Delete("8"); err != nil {
		t.Error(err)
	}

	if err := empTree.Delete("admin"); err == nil {
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

	if err := empTree.Insert("admin", NewEmployeeMap("10")); err != nil {
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

	fEmp, err := empTree.Find("1")

	if err != nil {
		t.Error(err)
	}

	t.Logf("\n%s", fEmp)
}
