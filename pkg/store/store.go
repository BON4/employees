package store

import "context"

type Store interface {
	Set(ctx context.Context, k string, v []byte) error
	Get(ctx context.Context, k string, v *[]byte) (found bool, err error)
	Delete(ctx context.Context, k string) error
	Close() error
}
