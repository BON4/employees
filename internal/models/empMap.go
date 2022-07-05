package models

import (
	"bytes"
	"context"
	"encoding/gob"
	"encoding/json"
	"errors"
	"strings"
	"sync"

	kvStore "github.com/BON4/employees/pkg/store"
)

//TODO Change to SyncMap or RWMutex
type EmployeeMap struct {
	//Addition fields
	// ...
	rwMu    *sync.RWMutex
	Payload Employee                `json:"Employee"`
	Ords    map[string]*EmployeeMap `json:"Oredenates"`
}

func NewEmployeeMap(payload Employee, ords ...(*EmployeeMap)) *EmployeeMap {
	emp := &EmployeeMap{Payload: payload, Ords: make(map[string]*EmployeeMap, len(ords)), rwMu: &sync.RWMutex{}}
	for _, v := range ords {
		emp.Ords[v.Payload.UUID] = v
	}
	return emp
}

func (e *EmployeeMap) insert(newTree *EmployeeMap) {
	e.rwMu.Lock()
	e.Ords[newTree.Payload.UUID] = newTree
	e.rwMu.Unlock()
}

func (e *EmployeeMap) delete(UUID string) {
	e.rwMu.Lock()
	delete(e.Ords, UUID)
	e.rwMu.Unlock()
}

func (e *EmployeeMap) len() int {
	e.rwMu.Lock()
	defer e.rwMu.Unlock()
	return len(e.Ords)
}

func (e EmployeeMap) copy() EmployeeMap {
	e.rwMu.Lock()
	defer e.rwMu.Unlock()
	cEmp := *NewEmployeeMap(e.Payload)
	for k, v := range e.Ords {
		cEmp.Ords[k] = NewEmployeeMap(v.Payload)
	}
	return cEmp
}

//IsExists - checks if provided emp with this id is exists within of scoupe of current emp
func (e *EmployeeMap) IsExists(empUUID string) (Employee, bool) {
	var helper func(empUUID string, fe *EmployeeMap) (Employee, bool)
	helper = func(empUUID string, fe *EmployeeMap) (Employee, bool) {
		if fe.Payload.UUID == empUUID {
			return fe.Payload, true
		}

		for _, v := range fe.Ords {
			if v, ok := helper(empUUID, v); ok {
				return v, true
			}
		}

		return Employee{}, false
	}

	e.rwMu.Lock()
	defer e.rwMu.Unlock()

	return helper(empUUID, e)
}

func (e EmployeeMap) String() string {
	var helper func(e *EmployeeMap, offset int) string
	helper = func(e *EmployeeMap, offset int) string {
		accum := strings.Repeat("-", offset) + e.Payload.UUID + " : " + e.Payload.Username
		for _, c := range e.Ords {
			accum += "\n" + helper(c, offset+1)
		}
		return accum
	}

	e.rwMu.Lock()
	defer e.rwMu.Unlock()

	return helper(&e, 0)
}

func (e EmployeeMap) Traverse(f func(emp Employee) error) error {
	var helper func(e *EmployeeMap, f func(emp Employee) error) error
	helper = func(e *EmployeeMap, f func(emp Employee) error) error {
		if err := f(e.Payload); err != nil {
			return err
		}

		for _, v := range e.Ords {
			if err := helper(v, f); err != nil {
				return err
			}
		}

		return nil
	}

	e.rwMu.Lock()
	defer e.rwMu.Unlock()

	return helper(&e, f)
}

func (e EmployeeMap) marshal() ([]byte, error) {
	b := bytes.NewBuffer([]byte{})
	enc := gob.NewEncoder(b)

	newE := make(map[string]*EmployeeMap, len(e.Ords))
	for k, v := range e.Ords {
		newE[k] = &EmployeeMap{Payload: v.Payload}
	}

	e.Ords = newE
	err := enc.Encode(e)
	return b.Bytes(), err
}

func (e *EmployeeMap) unmarshal(bData []byte) error {
	b := bytes.NewBuffer(bData)
	dec := gob.NewDecoder(b)
	err := dec.Decode(e)
	return err
}

func (e EmployeeMap) writeToStore(ctx context.Context, s kvStore.Store) error {
	eGob, err := e.marshal()
	if err == nil {
		return s.Set(ctx, e.Payload.UUID, eGob)
	}
	return err
}

func (e *EmployeeMap) readFromStore(ctx context.Context, s kvStore.Store) error {
	bData := make([]byte, 0, 512)
	ok, err := s.Get(ctx, e.Payload.UUID, &bData)

	if ok {
		return e.unmarshal(bData)
	}

	return err
}

func (e *EmployeeMap) Json() (string, error) {
	e.rwMu.Lock()
	defer e.rwMu.Unlock()

	b, err := json.Marshal(e)
	if err != nil {
		return "", err
	}

	return string(b), nil
}

func buildMapFromStore(ctx context.Context, e *EmployeeMap, s kvStore.Store) error {
	if len(e.Payload.UUID) == 0 {
		return errors.New("employeeMap does not exists")
	}

	if err := e.readFromStore(ctx, s); err != nil {
		return err
	}

	for k, v := range e.Ords {
		newE := NewEmployeeMap(v.Payload)
		if err := buildMapFromStore(ctx, newE, s); err != nil {
			return err
		}
		e.Ords[k] = newE
	}

	return nil
}

func dumpMapToStore(ctx context.Context, e *EmployeeMap, s kvStore.Store) error {
	if err := e.writeToStore(ctx, s); err != nil {
		return err
	}

	for _, v := range e.Ords {
		if err := dumpMapToStore(ctx, v, s); err != nil {
			return err
		}
	}

	return nil
}
