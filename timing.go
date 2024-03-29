package client

import (
	"math"
	"sync"
	"time"

	"github.com/blushft/jitsuclient/event/contexts"
)

type timer struct {
	lock  sync.Mutex
	start time.Time

	timing *contexts.Timing
}

func (t *Client) TimingStart(cat, label, variable string) *timer {
	return &timer{
		start: time.Now(),
		timing: &contexts.Timing{
			Category: cat,
			Label:    label,
			Variable: variable,
			Unit:     "seconds",
			Value:    -1,
		},
	}
}

func (timer *timer) End() *contexts.Timing {
	timer.lock.Lock()
	defer timer.lock.Unlock()
	dur := time.Since(timer.start).Seconds()

	timer.timing.Value = math.Floor(dur*100) / 100
	return timer.timing
}
