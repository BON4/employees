package models

import (
	"context"

	uerrors "github.com/BON4/employees/internal/errors"
	kvStore "github.com/BON4/employees/pkg/store"
)

//TODO UPTODATE SAVING
type EmpMapTree struct {
	root *EmployeeMap
}

func (e EmpMapTree) String() string {
	return e.root.String()
}

func (e *EmpMapTree) findParent(ctx context.Context, uuid string) (*EmployeeMap, error) {
	var helper func(em *EmployeeMap, uuid string) *EmployeeMap
	helper = func(em *EmployeeMap, uuid string) *EmployeeMap {
		if err := ctx.Err(); err != nil {
			return nil
		}

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
		return nil, uerrors.NewError[uerrors.UserError]("employee does not exists")
	}

	return fEmp, nil
}

func (e *EmpMapTree) FindById(ctx context.Context, UUID string) (*EmployeeMap, error) {
	var helper func(e *EmployeeMap, UUID string) *EmployeeMap
	helper = func(e *EmployeeMap, UUID string) *EmployeeMap {
		if err := ctx.Err(); err != nil {
			return nil
		}

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
	return nil, uerrors.NewError[uerrors.UserError]("employee does not exists")
}

func (e *EmpMapTree) FindByUsername(ctx context.Context, Username string) (*EmployeeMap, error) {
	var helper func(e *EmployeeMap, Username string) *EmployeeMap
	helper = func(e *EmployeeMap, Username string) *EmployeeMap {
		if err := ctx.Err(); err != nil {
			return nil
		}

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
	return nil, uerrors.NewError[uerrors.UserError]("employee does not exists")
}

func (e *EmpMapTree) Move(ctx context.Context, empUUID, toUUID string) error {
	empMapTo, err := e.FindById(ctx, toUUID)
	if err != nil {
		return err
	}

	empMap, err := e.FindById(ctx, empUUID)
	if err != nil {
		return err
	}

	empMapParrent, err := e.findParent(ctx, empUUID)
	if err != nil {
		return err
	}

	empMapParrent.delete(empUUID)

	if empMapParrent.len() == 0 {
		if empMapParrent.Payload.Role == Boss {
			empMapParrent.Payload.Role = Regular
		}
	}

	empMapTo.insert(empMap)

	if empMapTo.Payload.Role == Regular {
		empMapTo.Payload.Role = Boss
	}

	return nil
}

//TODO (Maby) if parent of inserted was at Role Regular change it to Boss
func (e *EmpMapTree) Insert(ctx context.Context, uuid string, newEmp Employee) error {
	if _, ok := e.root.IsExists(ctx, newEmp.UUID); ok {
		return uerrors.NewError[uerrors.UserError]("employee with this UUID already exists")
	}

	if err := e.root.Traverse(ctx, func(emp Employee) error {
		if emp.Username == newEmp.Username {
			return uerrors.NewError[uerrors.UserError]("employee with this username already exists")
		}
		return nil
	}); err != nil {
		return err
	}

	p, err := e.FindById(ctx, uuid)
	if err != nil {
		return err
	}

	p.insert(newEmployeeMap(newEmp))

	if p.Payload.Role == Regular {
		p.Payload.Role = Boss
	}

	return nil
}

//TODO (Maby) if parent of deleted was at Role Boss, and after deletion it have no childs, change it to Regular
func (e *EmpMapTree) Delete(ctx context.Context, childUUID string) error {
	p, err := e.findParent(ctx, childUUID)
	if err != nil {
		return err
	}
	p.delete(childUUID)
	return nil
}

func (e *EmpMapTree) Load(ctx context.Context, s kvStore.Store) error {
	return buildMapFromStore(ctx, e.root, s)
}

func (e *EmpMapTree) Save(ctx context.Context, s kvStore.Store) error {
	return dumpMapToStore(ctx, e.root, s)
}

func NewEmpMapTree() *EmpMapTree {
	return &EmpMapTree{(newEmployeeMap(NewEmployee("admin", "adminadmin", Admin)))}
}

func NewEmpMapTreeDEBUG() *EmpMapTree {
	return &EmpMapTree{(newEmployeeMap(NewEmployee("admin", "adminadmin", Admin),
		newEmployeeMap(NewEmployee("1", "1boss", Boss),
			newEmployeeMap(NewEmployee("4", "4boss", Boss),
				newEmployeeMap(NewEmployee("9", "9regular", Regular))),
			newEmployeeMap(NewEmployee("5", "5regular", Regular))),
		newEmployeeMap(NewEmployee("2", "2boss", Boss),
			newEmployeeMap(NewEmployee("6", "6reglar", Regular)),
			newEmployeeMap(NewEmployee("7", "7regular", Regular))),
		newEmployeeMap(NewEmployee("3", "3boss", Boss),
			newEmployeeMap(NewEmployee("8", "8regular", Regular)))))}
}
