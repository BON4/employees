package store

type MapStore struct {
	db map[string][]byte
}

func NewStore() Store {
	store := &MapStore{db: map[string][]byte{}}
	return store
}

func (mS MapStore) Get(k string, v *[]byte) (found bool, err error) {
	bData, ok := mS.db[k]
	if ok {
		*v = append(*v, bData...)
	}

	return ok, nil
}

func (mS MapStore) Set(k string, v []byte) error {
	mS.db[k] = v
	return nil
}

func (mS MapStore) Delete(k string) error {
	delete(mS.db, k)
	return nil
}

func (mS MapStore) Close() error {
	return nil
}

func NewMapStore() Store {
	return MapStore{db: map[string][]byte{}}
}
