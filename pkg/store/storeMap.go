package store

import "sync"

type MapStore struct {
	db *sync.Map
}

func NewStore() Store {
	store := &MapStore{db: &sync.Map{}}
	return store
}

func (mS MapStore) Get(k string, v *[]byte) (found bool, err error) {

	bData, ok := mS.db.Load(k)
	if ok {
		*v = append(*v, bData.([]byte)...)
	}

	return ok, nil
}

func (mS MapStore) Set(k string, v []byte) error {
	mS.db.Store(k, v)
	return nil
}

func (mS MapStore) Delete(k string) error {
	mS.db.Delete(k)
	return nil
}

func (mS MapStore) Close() error {
	return nil
}
