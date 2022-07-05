package store

import (
	"context"
	"sync"
)

type MapStore struct {
	db *sync.Map
}

func NewMapStore() Store {
	store := &MapStore{db: &sync.Map{}}
	return store
}

func (mS MapStore) Get(_ context.Context, k string, v *[]byte) (found bool, err error) {

	bData, ok := mS.db.Load(k)
	if ok {
		*v = append(*v, bData.([]byte)...)
	}

	return ok, nil
}

func (mS MapStore) Set(_ context.Context, k string, v []byte) error {
	mS.db.Store(k, v)
	return nil
}

func (mS MapStore) Delete(_ context.Context, k string) error {
	mS.db.Delete(k)
	return nil
}

func (mS MapStore) Close() error {
	return nil
}
