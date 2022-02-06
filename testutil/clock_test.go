package testutil

import (
	"testing"
	"time"
)

func Test_Clock(t *testing.T) {
	c := NewClockAt(time.Date(2022, 1, 1, 0, 0, 0, 0, time.UTC))
	ch1 := c.After(time.Hour)
	ch2 := c.After(2 * time.Hour)
	ch3 := c.After(3 * time.Hour)
	c.Advance(10 * time.Minute)
	select {
	case <-ch1:
		t.Fatalf("ch1 fired at %v", c.Now())
	case <-ch2:
		t.Fatalf("ch2 fired at %v", c.Now())
	case <-ch3:
		t.Fatalf("ch3 fired at %v", c.Now())
	default:
	}
	c.Advance(50 * time.Minute)
	select {
	case <-ch1:
	case <-ch2:
		t.Fatalf("ch2 fired at %v", c.Now())
	case <-ch3:
		t.Fatalf("ch3 fired at %v", c.Now())
	default:
		t.Fatalf("default fired at %v", c.Now())
	}
	c.Advance(5 * time.Minute)
	select {
	case <-ch1:
		t.Fatalf("ch1 fired at %v", c.Now())
	case <-ch2:
		t.Fatalf("ch2 fired at %v", c.Now())
	case <-ch3:
		t.Fatalf("ch3 fired at %v", c.Now())
	default:
	}
	c.Advance(time.Hour)
	select {
	case <-ch1:
		t.Fatalf("ch1 fired at %v", c.Now())
	case <-ch2:
	case <-ch3:
		t.Fatalf("ch3 fired at %v", c.Now())
	default:
		t.Fatalf("default fired at %v", c.Now())
	}
	c.Advance(time.Hour)
	select {
	case <-ch1:
		t.Fatalf("ch1 fired at %v", c.Now())
	case <-ch2:
		t.Fatalf("ch2 fired at %v", c.Now())
	case <-ch3:
	default:
		t.Fatalf("default fired at %v", c.Now())
	}
	c.Advance(time.Hour)
	select {
	case <-ch1:
		t.Fatalf("ch1 fired at %v", c.Now())
	case <-ch2:
		t.Fatalf("ch2 fired at %v", c.Now())
	case <-ch3:
		t.Fatalf("ch3 fired at %v", c.Now())
	default:
	}
}
