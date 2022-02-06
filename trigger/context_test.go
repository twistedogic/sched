package trigger

import (
	"context"
	"testing"
	"time"

	"github.com/twistedogic/sched/clock"
	"github.com/twistedogic/sched/testutil"
)

func Test_triggerCtx(t *testing.T) {
	start := time.Date(2022, 1, 1, 0, 0, 0, 0, time.UTC)
	cases := map[string]struct {
		now, deadline time.Time
		done          bool
		wantErr       error
	}{
		"base": {
			now:      time.Date(2022, 1, 1, 0, 0, 0, 0, time.UTC),
			deadline: time.Date(2022, 1, 1, 0, 0, 1, 0, time.UTC),
		},
		"done": {
			now:      time.Date(2022, 1, 1, 0, 0, 3, 0, time.UTC),
			deadline: time.Date(2022, 1, 1, 0, 0, 2, 0, time.UTC),
			done:     true,
			wantErr:  TaskTimeoutErr,
		},
	}
	for name := range cases {
		tc := cases[name]
		t.Run(name, func(t *testing.T) {
			clk := testutil.NewClockAt(start)
			ctx, cancel := context.WithCancel(context.TODO())
			defer cancel()
			cctx := clock.WithClock(ctx, clk)
			tCtx := withDeadline(cctx, tc.deadline)
			clk.AdvanceTo(tc.now)
			time.Sleep(10 * time.Millisecond) // wait for go routine
			select {
			case <-tCtx.Done():
				if !tc.done {
					t.Fatalf("ctx shouldn't be done")
				}
				if got := tCtx.Err(); got != tc.wantErr {
					t.Fatalf("err, want: %v, got: %v", tc.wantErr, got)
				}
			default:
				if tc.done {
					t.Fatalf("ctx should be done")
				}
			}
		})
	}
}
