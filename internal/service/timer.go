package service

import (
	"time"
)

type Timer struct {
	duration   time.Duration
	ticker     *time.Ticker
	paused     bool
	done       chan bool
	pause      chan bool
	timePassed int
	startTime  time.Time
}

func newTimer() *Timer {
	return &Timer{
		paused: true,
		done:   make(chan bool),
		pause:  make(chan bool),
	}
}

func (t *Timer) Start() {
	t.ticker = time.NewTicker(t.duration)
	t.startTime = time.Now()
	t.timePassed = 0
	t.paused = false

	go func() {
		defer t.ticker.Stop()

		for range t.ticker.C {
			t.done <- true

			return
		}
	}()
}

func (t *Timer) Paused() <-chan bool {
	return t.pause
}

func (t *Timer) Pause() {
	if !t.IsPaused() {
		t.paused = true
		t.timePassed = int(time.Now().Unix() - t.startTime.Unix())
		t.pause <- true
		t.ticker.Stop()
	}
}

func (t *Timer) Stop() {
	t.ticker.Stop()
}

func (t *Timer) IsPaused() bool {
	return t.paused
}

func (t *Timer) Done() <-chan bool {
	return t.done
}
