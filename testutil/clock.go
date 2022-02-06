package testutil

import (
	"sync"
	"time"
)

type block struct {
	end time.Time
	ch  chan time.Time
}

func (b block) notify(t time.Time) bool {
	if t.Before(b.end) {
		return false
	}
	b.ch <- b.end
	return true
}

type Clock struct {
	sync.RWMutex
	current time.Time
	blocks  []block
}

func NewClockAt(t time.Time) *Clock {
	return &Clock{RWMutex: sync.RWMutex{}, current: t, blocks: make([]block, 0)}
}

func (c *Clock) Now() time.Time {
	c.RLock()
	defer c.RUnlock()
	return c.current
}

func (c *Clock) addBlock(b block) {
	c.Lock()
	defer c.Unlock()
	c.blocks = append(c.blocks, b)
}

func (c *Clock) After(d time.Duration) <-chan time.Time {
	b := block{
		end: c.Now().Add(d),
		ch:  make(chan time.Time, 1),
	}
	c.addBlock(b)
	return b.ch
}

func (c *Clock) AdvanceTo(t time.Time) {
	c.Lock()
	defer c.Unlock()
	c.current = t
	newBlocks := make([]block, 0, len(c.blocks))
	for _, b := range c.blocks {
		if !b.notify(t) {
			newBlocks = append(newBlocks, b)
		}
	}
	c.blocks = newBlocks
}

func (c *Clock) Advance(d time.Duration) { c.AdvanceTo(c.Now().Add(d)) }
