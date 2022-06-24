package models

import (
	"bytes"
	"encoding/gob"
	"errors"
	"strings"

	kvStore "github.com/BON4/employees/internal/store"
)

type EmployeeMap struct {
	//Addition fields
	// ...
	Payload Employee
	Ords    map[string]*EmployeeMap
}

func NewEmployeeMap(payload Employee, ords ...(*EmployeeMap)) *EmployeeMap {
	emp := &EmployeeMap{Payload: payload, Ords: make(map[string]*EmployeeMap, len(ords))}
	for _, v := range ords {
		emp.Ords[v.Payload.UUID] = v
	}
	return emp
}

func (e *EmployeeMap) insert(newTree *EmployeeMap) {
	e.Ords[newTree.Payload.UUID] = newTree
}

func (e *EmployeeMap) delete(UUID string) {
	delete(e.Ords, UUID)
}

func (e EmployeeMap) copy() EmployeeMap {
	cEmp := *NewEmployeeMap(e.Payload)
	for k, v := range e.Ords {
		cEmp.Ords[k] = NewEmployeeMap(v.Payload)
	}
	return cEmp
}

func (e *EmployeeMap) IsExists(empUUID string) bool {
	var helper func(empUUID string, fe *EmployeeMap) bool
	helper = func(empUUID string, fe *EmployeeMap) bool {
		if fe.Payload.UUID == empUUID {
			return true
		}

		for _, v := range fe.Ords {
			if helper(empUUID, v) {
				return true
			}
		}

		return false
	}
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
	return helper(&e, 0)
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

func (e EmployeeMap) writeToStore(s kvStore.Store) error {
	eGob, err := e.marshal()
	if err == nil {
		return s.Set(e.Payload.UUID, eGob)
	}
	return err
}

func (e *EmployeeMap) readFromStore(s kvStore.Store) error {
	bData := make([]byte, 0, 512)
	ok, err := s.Get(e.Payload.UUID, &bData)

	if ok {
		return e.unmarshal(bData)
	}

	return err
}

func buildMapFromStore(e *EmployeeMap, s kvStore.Store) error {
	if len(e.Payload.UUID) == 0 {
		return errors.New("employeeMap does not exists")
	}

	if err := e.readFromStore(s); err != nil {
		return err
	}

	for k, v := range e.Ords {
		newE := NewEmployeeMap(v.Payload)
		if err := buildMapFromStore(newE, s); err != nil {
			return err
		}
		e.Ords[k] = newE
	}

	return nil
}

func dumpMapToStore(e *EmployeeMap, s kvStore.Store) error {
	if err := e.writeToStore(s); err != nil {
		return err
	}

	for _, v := range e.Ords {
		if err := dumpMapToStore(v, s); err != nil {
			return err
		}
	}

	return nil
}
