package testutil

import (
	"context"
	"sync"
	"testing"
	"time"
)

type clock interface {
	Now() time.Time
	After(time.Duration) <-chan time.Time
}

type MockTask struct {
	sync.Mutex
	clock
	duration time.Duration
	value    int
}

func SetupMockTask(clk clock, dur time.Duration) *MockTask {
	return &MockTask{
		Mutex:    sync.Mutex{},
		clock:    clk,
		duration: dur,
		value:    0,
	}
}

func (m *MockTask) inc() {
	m.Lock()
	defer m.Unlock()
	m.value += 1
}

func (m *MockTask) Run(ctx context.Context) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-m.After(m.duration):
		m.inc()
		return nil
	}
}

func (m *MockTask) Check(t *testing.T, want int) {
	if m.value != want {
		t.Fatalf("at %v, want: %d, got: %d", m.Now(), want, m.value)
	}
}
