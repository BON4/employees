package store

type Store interface {
	Set(k string, v []byte) error
	Get(k string, v *[]byte) (found bool, err error)
	Delete(k string) error
	Close() error
}
