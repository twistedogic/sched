package trigger

import (
	"context"
	"sync"
	"time"

	"github.com/pkg/errors"

	"github.com/twistedogic/sched/clock"
)

var (
	TaskTimeoutErr = errors.New("task timeout")
)

type triggerCtx struct {
	sync.Mutex
	context.Context
	err      error
	deadline time.Time
	cancel   context.CancelFunc
}

func withDeadline(ctx context.Context, deadline time.Time) context.Context {
	clk := clock.FromContext(ctx)
	now := clk.Now()
	cancelCtx, cancel := context.WithCancel(ctx)
	tCtx := &triggerCtx{
		Mutex:    sync.Mutex{},
		Context:  cancelCtx,
		deadline: deadline,
		cancel:   cancel,
	}
	dur := deadline.Sub(now)
	go tCtx.wait(clk.After(dur))
	return tCtx
}

func WithTrigger(ctx context.Context, t Trigger) context.Context {
	clk := clock.FromContext(ctx)
	deadline := t.Timeout(clk.Now())
	return withDeadline(ctx, deadline)
}

func (t *triggerCtx) wait(afterCh <-chan time.Time) {
	select {
	case <-t.Done():
	case <-afterCh:
		t.cancelAndSetErr()
	}
}

func (t *triggerCtx) setErr() {
	t.Lock()
	defer t.Unlock()
	t.err = TaskTimeoutErr
}

func (t *triggerCtx) cancelAndSetErr() {
	t.setErr()
	t.cancel()
}

func (t *triggerCtx) Deadline() (time.Time, bool) {
	return t.deadline, true
}

func (t *triggerCtx) Err() error {
	if t.err != nil {
		return t.err
	}
	return t.Context.Err()
}
