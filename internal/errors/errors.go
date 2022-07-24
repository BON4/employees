package errors

type UserError interface{}

type CustomError[T UserError] struct {
	err string
}

func NewError[T UserError](err string) CustomError[T] {
	return CustomError[T]{
		err: err,
	}
}

func (c CustomError[T]) Error() string {
	return ""
}
