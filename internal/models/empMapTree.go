package models

import (
	"errors"

	kvStore "github.com/BON4/employees/internal/store"
)

//TODO UPTODATE SAVING
type EmpMapTree struct {
	root *EmployeeMap
}

func NewEmpMapTree() *EmpMapTree {
	return &EmpMapTree{(NewEmployeeMap(NewEmployee("admin", "", Admin)))}
}

func NewEmpMapTreeDEBUG() *EmpMapTree {
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

func (e EmpMapTree) String() string {
	return e.root.String()
}

func (e *EmpMapTree) findParent(uuid string) (*EmployeeMap, error) {
	var helper func(em *EmployeeMap, uuid string) *EmployeeMap
	helper = func(em *EmployeeMap, uuid string) *EmployeeMap {
		if _, ok := em.Ords[uuid]; ok {
			return em
		}

		for _, v := range em.Ords {
			if found := helper(v, uuid); found != nil {
				return found
			}
		}

		return nil
	}

	fEmp := helper(e.root, uuid)
	if fEmp == nil {
		return nil, errors.New("employee does not exists")
	}

	return fEmp, nil
}

func (e *EmpMapTree) FindById(UUID string) (*EmployeeMap, error) {
	var helper func(e *EmployeeMap, UUID string) *EmployeeMap
	helper = func(e *EmployeeMap, UUID string) *EmployeeMap {
		if e.Payload.UUID == UUID {
			return e
		}

		for _, c := range e.Ords {
			if found := helper(c, UUID); found != nil {
				return found
			}
		}

		return nil
	}

	if e := helper(e.root, UUID); e != nil {
		return e, nil
	}
	return nil, errors.New("employee does not exists")
}

func (e *EmpMapTree) FindByUsername(Username string) (*EmployeeMap, error) {
	var helper func(e *EmployeeMap, Username string) *EmployeeMap
	helper = func(e *EmployeeMap, Username string) *EmployeeMap {
		if e.Payload.Username == Username {
			return e
		}

		for _, c := range e.Ords {
			if found := helper(c, Username); found != nil {
				return found
			}
		}

		return nil
	}

	if e := helper(e.root, Username); e != nil {
		return e, nil
	}
	return nil, errors.New("employee does not exists")
}

//TODO (Maby) if parent of inserted was at Role Regular change it to Boss
func (e *EmpMapTree) Insert(uuid string, newEmp Employee) error {
	if e.root.IsExists(newEmp.UUID) {
		return errors.New("employee with this UUID already exists")
	}

	p, err := e.FindById(uuid)
	if err != nil {
		return err
	}

	p.insert(NewEmployeeMap(newEmp))

	if p.Payload.Role == Regular {
		p.Payload.Role = Boss
	}

	return nil
}

//TODO (Maby) if parent of deleted was at Role Boss, and after deletion it have no childs, change it to Regular
func (e *EmpMapTree) Delete(childUUID string) error {
	p, err := e.findParent(childUUID)
	if err != nil {
		return err
	}
	p.delete(childUUID)
	return nil
}

func (e *EmpMapTree) Load(s kvStore.Store) error {
	return buildMapFromStore(e.root, s)
}

func (e *EmpMapTree) Save(s kvStore.Store) error {
	return dumpMapToStore(e.root, s)
}
