package trigger

import (
	"time"
)

type Trigger interface {
	Timeout(time.Time) time.Time
	Next(time.Time) time.Time
}

type baseTrigger struct {
	start             time.Time
	interval, timeout time.Duration
}

func newIntervalTrigger(start time.Time, interval, timeout time.Duration) Trigger {
	return baseTrigger{
		start:    start,
		interval: interval,
		timeout:  timeout,
	}
}

func (b baseTrigger) lastRun(now time.Time) time.Time {
	num := now.Sub(b.start) / b.interval
	return b.start.Add(num * b.interval)
}

func (b baseTrigger) Next(now time.Time) time.Time {
	return b.lastRun(now).Add(b.interval)
}

func (b baseTrigger) Timeout(now time.Time) time.Time {
	return b.lastRun(now).Add(b.timeout)
}

func EveryAt(at time.Time, interval, timeout time.Duration) Trigger {
	return newIntervalTrigger(at, interval, timeout)
}

func Every(interval, timeout time.Duration) Trigger {
	return EveryAt(time.Now(), interval, timeout)
}
