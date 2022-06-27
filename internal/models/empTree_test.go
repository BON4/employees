package models

import (
	"errors"
	"fmt"
	"math/rand"
	"testing"

	kvStore "github.com/BON4/employees/internal/store"
)

var mapStore kvStore.Store

func TestMain(m *testing.M) {
	mapStore = kvStore.NewStore()
	m.Run()
}

func createDumbEmplMap() *EmpMapTree {
	return &EmpMapTree{(NewEmployeeMap(NewEmployee("admin", "", Admin),
		NewEmployeeMap(NewEmployee("1", "", Boss),
			NewEmployeeMap(NewEmployee("4", "", Boss),
				NewEmployeeMap(NewEmployee("9", "", Regular))),
			NewEmployeeMap(NewEmployee("5", "", Regular))),
		NewEmployeeMap(NewEmployee("2", "", Boss),
			NewEmployeeMap(NewEmployee("6", "", Regular)),
			NewEmployeeMap(NewEmployee("7", "", Regular))),
		NewEmployeeMap(NewEmployee("3", "", Boss),
			NewEmployeeMap(NewEmployee("8", "", Regular)))))}
}

func fuzzDumbEmplTree(n int, e *EmpMapTree) error {
	counter := 9
	for i := 0; i < n; i++ {
		v := fmt.Sprintf("%d", rand.Intn(counter)+1)
		emp, err := e.FindByUsername(v)
		if err == nil {
			err = e.Insert(emp.Payload.UUID, NewEmployee(fmt.Sprintf("%d", counter), "", Regular))
			if err != nil {
				return errors.New(fmt.Sprintf("%s:%s", err.Error(), emp.Payload.UUID))
			}
			counter++
		} else {
			i--
		}
	}
	return nil
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
		t.Logf("Saved Tree:\n%s", emps)
		t.Logf("Loaded Tree:\n%s", loadedTree)
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
	e := NewEmployee("10", "", Regular)

	if err := empTree.Insert(fEmp.Payload.UUID, e); err != nil {
		t.Error(err)
	}

	e = NewEmployee("8", "", Regular)

	if err := empTree.Insert(fEmp.Payload.UUID, e); err == nil {
		t.Error("Shuld be error")
	}

	t.Logf("\n%s", empTree)
}

func TestEmpMapFind(t *testing.T) {
	empTree := createDumbEmplMap()
	err := empTree.Save(mapStore)
	if err != nil {
		t.Error(err)
	}

	fEmp, err := empTree.FindByUsername("4")

	if err != nil {
		t.Error(err)
	}

	fEmp, err = empTree.FindById(fEmp.Payload.UUID)

	if err != nil {
		t.Error(err)
	}

	t.Logf("\n%s", fEmp)
}

func TestEmpMapTraverse(t *testing.T) {
	empTree := createDumbEmplMap()
	err := empTree.Save(mapStore)
	if err != nil {
		t.Error(err)
	}

	if err := empTree.root.Traverse(func(emp Employee) error {
		t.Logf("%+v\n", emp)
		return nil
	}); err != nil {
		t.Error(err)
	}

	if err := empTree.root.Traverse(func(emp Employee) error {
		if emp.Username == "4" {
			return errors.New("employee with this username already exists")
		}
		return nil
	}); err == nil {
		t.Error("Shoud be error")
	}
}

func TestJsonifyEmployeeTree(t *testing.T) {
	tree := createDumbEmplMap()
	t.Log(tree.root.Json())
}

func BenchmarkJsonifyEmployeeTree(b *testing.B) {
	b.Log(b.N)
	b.StopTimer()
	tree := createDumbEmplMap()
	err := fuzzDumbEmplTree(10, tree)
	if err != nil {
		b.Error(err)
	}
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		tree.root.Json()
	}
}
