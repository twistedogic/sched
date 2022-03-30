// Package sched implements context based task scheduling
package sched

import (
	"context"

	"github.com/twistedogic/sched/clock"
	"github.com/twistedogic/sched/task"
	"github.com/twistedogic/sched/trigger"
)

func Run(ctx context.Context, t trigger.Trigger, f task.Func) error {
	errCh := make(chan error)
	clk := clock.FromContext(ctx)
	for {
		tCtx := trigger.WithTrigger(ctx, t)
		next := t.Next(clk.Now())
		go func() { errCh <- f(tCtx) }()
		select {
		case <-ctx.Done():
			return ctx.Err()
		case err := <-errCh:
			if err != nil {
				return err
			}
			waitDur := next.Sub(clk.Now())
			<-clk.After(waitDur)
		}
	}
	return nil
}
