package models

import (
	"context"
	"errors"
	"fmt"
	"math/rand"
	"testing"

	kvStore "github.com/BON4/employees/pkg/store"
)

var mapStore kvStore.Store

func TestMain(m *testing.M) {
	mapStore = kvStore.NewMapStore()
	m.Run()
}

func createDumbEmplMap() *EmpMapTree {
	return NewEmpMapTreeDEBUG()
}

func fuzzDumbEmplTree(n int, e *EmpMapTree) error {
	counter := 9
	ctx := context.Background()
	for i := 0; i < n; i++ {
		v := fmt.Sprintf("%d", rand.Intn(counter)+1)
		emp, err := e.FindByUsername(ctx, v)
		if err == nil {
			err = e.Insert(ctx, emp.Payload.UUID, NewEmployee(fmt.Sprintf("%d", counter), "", Regular))
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

	err := emps.Save(context.Background(), mapStore)
	if err != nil {
		t.Error(err)
	}

	loadedTree := NewEmpMapTree()
	if err := loadedTree.Load(context.Background(), mapStore); err != nil {
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
	ctx := context.Background()

	empTree := createDumbEmplMap()
	err := empTree.Save(context.Background(), mapStore)
	if err != nil {
		t.Error(err)
	}
	t.Logf("\n%s", empTree)

	fEmp, err := empTree.FindByUsername(ctx, "8")

	if err != nil {
		t.Error(err)
	}

	if err := empTree.Delete(ctx, fEmp.Payload.UUID); err != nil {
		t.Error(err)
	}

	fAdmin, err := empTree.FindByUsername(ctx, "admin")

	if err != nil {
		t.Error(err)
	}

	if err := empTree.Delete(ctx, fAdmin.Payload.UUID); err == nil {
		t.Error("Cant delete admin")
	}

	t.Logf("\n%s", empTree)
}

func TestEmpMapInsert(t *testing.T) {
	empTree := createDumbEmplMap()
	ctx := context.Background()
	err := empTree.Save(context.Background(), mapStore)
	if err != nil {
		t.Error(err)
	}

	fEmp, err := empTree.FindByUsername(ctx, "admin")

	if err != nil {
		t.Error(err)
	}
	e := NewEmployee("10", "", Regular)

	if err := empTree.Insert(ctx, fEmp.Payload.UUID, e); err != nil {
		t.Error(err)
	}

	e = NewEmployee("8", "", Regular)

	if err := empTree.Insert(ctx, fEmp.Payload.UUID, e); err == nil {
		t.Error("Shuld be error")
	}

	t.Logf("\n%s", empTree)
}

func TestEmpMapFind(t *testing.T) {
	ctx := context.Background()
	empTree := createDumbEmplMap()
	err := empTree.Save(context.Background(), mapStore)
	if err != nil {
		t.Error(err)
	}

	fEmp, err := empTree.FindByUsername(ctx, "4")

	if err != nil {
		t.Error(err)
	}

	fEmp, err = empTree.FindById(ctx, fEmp.Payload.UUID)

	if err != nil {
		t.Error(err)
	}

	t.Logf("\n%s", fEmp)
}

func TestEmpMapTraverse(t *testing.T) {
	empTree := createDumbEmplMap()
	ctx := context.Background()
	err := empTree.Save(context.Background(), mapStore)
	if err != nil {
		t.Error(err)
	}

	if err := empTree.root.Traverse(ctx, func(emp Employee) error {
		t.Logf("%+v\n", emp)
		return nil
	}); err != nil {
		t.Error(err)
	}

	if err := empTree.root.Traverse(ctx, func(emp Employee) error {
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
