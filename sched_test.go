package sched

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/twistedogic/sched/clock"
	"github.com/twistedogic/sched/testutil"
	"github.com/twistedogic/sched/trigger"
)

type checkpoint struct {
	step time.Duration
	want int
}

func Test_Run(t *testing.T) {
	cases := map[string]struct {
		interval, timeout, duration time.Duration
		steps                       []checkpoint
		err                         error
	}{
		"base": {
			interval: 3 * time.Second,
			timeout:  2 * time.Second,
			duration: time.Second,
			steps: []checkpoint{
				{step: 3 * time.Second, want: 1},
				{step: 3 * time.Second, want: 2},
				{step: 3 * time.Second, want: 3},
			},
		},
		"timeout": {
			interval: time.Second,
			timeout:  time.Second,
			duration: 3 * time.Second,
			steps: []checkpoint{
				{step: time.Second, want: 0},
				{step: time.Second, want: 0},
				{step: time.Second, want: 0},
			},
			err: trigger.TaskTimeoutErr,
		},
	}
	for name := range cases {
		tc := cases[name]
		t.Run(name, func(t *testing.T) {
			var err error
			errCh := make(chan error)
			ctx, cancel := context.WithCancel(context.TODO())
			defer cancel()
			tr := trigger.NewIntervalTrigger(time.Now(), tc.interval, tc.timeout)
			task := testutil.SetupMockTask(clock.New(), tc.duration)
			go func() {
				if err := Run(ctx, tr, task.Run); err != nil {
					errCh <- err
				}
			}()
			for _, c := range tc.steps {
				time.Sleep(c.step)
				select {
				case e := <-errCh:
					err = e
				default:
				}
				task.Check(t, c.want)
			}
			if !errors.Is(tc.err, err) {
				t.Fatalf("err, want: %v, got: %v", tc.err, err)
			}
		})
	}
}
