package util

import (
	"context"
)

type Io[T any] func(context.Context)(T, error)

type Void struct{}
var Empty Void = struct{}{}

func Lift[T,U any](
	pure func(T) U,
) func(T) Io[U] {
	return func(t T) Io[U]{
		return func(_ context.Context)(U, error){
			var u U = pure(t)
			return u, nil
		}
	}
}

func Bind[T,U any](
	i Io[T],
	f func(T) Io[U],
) Io[U] {
	return func(ctx context.Context)(u U, e error){
		t, e := i(ctx)
		if nil != e {
			return u, e
		}
		return f(t)(ctx)
	}
}
