package trigger

import (
	"testing"
	"time"

	"github.com/twistedogic/sched/testutil"
)

func Test_BaseTrigger(t *testing.T) {
	cases := map[string]struct {
		interval, timeout          time.Duration
		start, now, deadline, next time.Time
	}{
		"base": {
			interval: 10 * time.Second,
			timeout:  5 * time.Second,
			start:    time.Date(2022, 1, 1, 0, 0, 0, 0, time.UTC),
			now:      time.Date(2022, 1, 1, 0, 0, 1, 0, time.UTC),
			deadline: time.Date(2022, 1, 1, 0, 0, 5, 0, time.UTC),
			next:     time.Date(2022, 1, 1, 0, 0, 10, 0, time.UTC),
		},
		"2nd run": {
			interval: 10 * time.Second,
			timeout:  5 * time.Second,
			start:    time.Date(2022, 1, 1, 0, 0, 0, 0, time.UTC),
			now:      time.Date(2022, 1, 1, 0, 0, 11, 0, time.UTC),
			deadline: time.Date(2022, 1, 1, 0, 0, 15, 0, time.UTC),
			next:     time.Date(2022, 1, 1, 0, 0, 20, 0, time.UTC),
		},
	}
	for name := range cases {
		tc := cases[name]
		t.Run(name, func(t *testing.T) {
			clk := testutil.NewClockAt(tc.start)
			trigger := newIntervalTrigger(tc.start, tc.interval, tc.timeout)
			clk.AdvanceTo(tc.now)
			if got := trigger.Next(tc.now); !got.Equal(tc.next) {
				t.Fatalf("next, want: %v, got: %v", tc.next, got)
			}
			if got := trigger.Timeout(tc.now); !got.Equal(tc.deadline) {
				t.Fatalf("timeout, want: %v, got: %v", tc.deadline, got)
			}
		})
	}
}
