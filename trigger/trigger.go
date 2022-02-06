package trigger

import (
	"time"
)

type Trigger interface {
	Timeout(time.Time) time.Time
	Next(time.Time) time.Time
}

type BaseTrigger struct {
	start             time.Time
	interval, timeout time.Duration
}

func NewIntervalTrigger(start time.Time, interval, timeout time.Duration) Trigger {
	return BaseTrigger{
		start:    start,
		interval: interval,
		timeout:  timeout,
	}
}

func (b BaseTrigger) lastRun(now time.Time) time.Time {
	num := now.Sub(b.start) / b.interval
	return b.start.Add(num * b.interval)
}

func (b BaseTrigger) Next(now time.Time) time.Time {
	return b.lastRun(now).Add(b.interval)
}

func (b BaseTrigger) Timeout(now time.Time) time.Time {
	return b.lastRun(now).Add(b.timeout)
}
