package task

import (
	"context"
	"log"

	"github.com/pkg/errors"
)

type Func func(context.Context) error

type Wrapper func(f Func) Func

func ComposeWrapper(wraps ...Wrapper) Wrapper {
	return func(f Func) Func {
		wrapped := f
		for _, wrap := range wraps {
			wrapped = wrap(wrapped)
		}
		return wrapped
	}
}

func LogErrorWrapper(f Func) Func {
	return func(ctx context.Context) error {
		if err := f(ctx); err != nil {
			log.Println(errors.Wrap(err, "task failure"))
		}
		return nil
	}
}
