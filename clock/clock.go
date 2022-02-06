package clock

import (
	"context"
	"time"
)

type contextKey uint

const ctxKey contextKey = iota

type Clock interface {
	Now() time.Time
	After(time.Duration) <-chan time.Time
}

type RealClock struct{}

func New() Clock                                           { return RealClock{} }
func (r RealClock) Now() time.Time                         { return time.Now() }
func (r RealClock) After(d time.Duration) <-chan time.Time { return time.After(d) }

func FromContext(ctx context.Context) Clock {
	if v, ok := ctx.Value(ctxKey).(Clock); ok {
		return v
	}
	return New()
}

func WithClock(ctx context.Context, c Clock) context.Context {
	return context.WithValue(ctx, ctxKey, c)
}
