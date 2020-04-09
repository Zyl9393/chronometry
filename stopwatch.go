package chronometry

import (
	"time"
)

// NowFunc specifies the function used to retrieve the current time. It is changed only during tests.
var NowFunc = hpNow

// Stopwatch implements a stopwatch which can be stopped, resumed and reset and can report its current reading
// of time with Total() as well as allow for convenient measuring of consecutive intervals using Lap().
type Stopwatch struct {
	startTime, stopTime, splitTime time.Time
	isStopped                      bool
	totalDuration, splitDuration   time.Duration
}

// NewStartedStopwatch returns a started stopwatch, ready to report Total() and Lap().
func NewStartedStopwatch() *Stopwatch {
	now := NowFunc()
	return &Stopwatch{startTime: now, stopTime: now, splitTime: now}
}

// NewStoppedStopwatch returns a stopwatch which needs to have Restart() or Resume() called on it in order to
// begin counting time.
func NewStoppedStopwatch() *Stopwatch {
	now := NowFunc()
	return &Stopwatch{startTime: now, stopTime: now, splitTime: now, isStopped: true}
}

// Restart resets and resumes the stopwatch.
func (sw *Stopwatch) Restart() {
	sw.Reset()
	sw.isStopped = false
}

// IsStopped reports whether the stopwatch is stopped, i.e. not started.
func (sw *Stopwatch) IsStopped() bool {
	return sw.isStopped
}

// IsStarted reports whether the stopwatch is started, i.e. not stopped.
func (sw *Stopwatch) IsStarted() bool {
	return !sw.IsStopped()
}

// Reset sets all stopwatch readings back to zero without affecting whether it is currently started or stopped.
func (sw *Stopwatch) Reset() {
	now := NowFunc()
	sw.startTime, sw.splitTime, sw.stopTime = now, now, now
	sw.totalDuration, sw.splitDuration = 0, 0
}

// Resume resumes the stopwatch if it is stopped. The returned boolean is false if the stopwatch was already started; true otherwise.
func (sw *Stopwatch) Resume() bool {
	if sw.isStopped {
		now := NowFunc()
		sw.startTime, sw.splitTime, sw.isStopped = now, now, false
		return true
	}
	return false
}

// Stop stops the stopwatch such that it pauses observing the passing of time. It can be resumed with Resume().
func (sw *Stopwatch) Stop() time.Duration {
	if !sw.isStopped {
		now := NowFunc()
		sw.stopTime, sw.isStopped = now, true
		sw.totalDuration += now.Sub(sw.startTime)
		sw.splitDuration += sw.segment(now)
	}
	return sw.totalDuration
}

// Toggle stops the stopwatch if it is started and stops it otherwise, returning the result of IsStarted() at the end.
func (sw *Stopwatch) Toggle() bool {
	if sw.isStopped {
		sw.Resume()
	} else {
		sw.Stop()
	}
	return !sw.isStopped
}

// Segment returns the duration into the latest contiguous measuring interval of the current lap.
func (sw *Stopwatch) Segment() time.Duration {
	return sw.segment(NowFunc())
}

func (sw *Stopwatch) segment(now time.Time) time.Duration {
	if sw.isStopped {
		if sw.splitTime.After(sw.startTime) {
			return sw.stopTime.Sub(sw.splitTime)
		}
		return sw.stopTime.Sub(sw.startTime)
	}
	if sw.splitTime.Before(sw.startTime) {
		return now.Sub(sw.startTime)
	}
	return now.Sub(sw.splitTime)
}

// StartTime returns the time at which the stopwatch was last started/restarted, reset or resumed.
func (sw *Stopwatch) StartTime() time.Time {
	return sw.startTime
}

// StopTime returns the time at which the stopwatch was last stopped or reset. If neither was ever done, it returns the same as StartTime().
func (sw *Stopwatch) StopTime() time.Time {
	return sw.stopTime
}

// SplitTime returns the time at which the current split has begun, i.e. the start time of the stopwatch's latest contiguous measuring interval.
func (sw *Stopwatch) SplitTime() time.Time {
	return sw.splitTime
}

// Lap makes a split, returning the change in Total() since last calling this function or resetting or first starting the stopwatch.
func (sw *Stopwatch) Lap() time.Duration {
	return sw.lap(NowFunc(), true)
}

// PeekLap returns the change in Total() since last calling Lap() or resetting or first starting the stopwatch.
func (sw *Stopwatch) PeekLap() time.Duration {
	return sw.lap(NowFunc(), false)
}

func (sw *Stopwatch) lap(now time.Time, split bool) time.Duration {
	lapDuration := sw.splitDuration
	if !sw.isStopped {
		lapDuration += sw.segment(now)
	}
	if split {
		sw.splitTime, sw.splitDuration = now, 0
	}
	return lapDuration
}

// Total returns the current stopwatch reading, i.e. the total passed time observed by the stopwatch while not stopped.
func (sw *Stopwatch) Total() time.Duration {
	return sw.totalTime(NowFunc())
}

func (sw *Stopwatch) totalTime(now time.Time) time.Duration {
	if sw.isStopped {
		return sw.totalDuration
	}
	return sw.totalDuration + now.Sub(sw.startTime)
}

// Split does the same as Lap(), but also returns the exact according total time.
func (sw *Stopwatch) Split() (lapTime, totalTime time.Duration) {
	now := NowFunc()
	return sw.lap(now, true), sw.totalTime(now)
}

// PeekSplit does the same as PeekLap(), but also returns the exact according total time.
func (sw *Stopwatch) PeekSplit() (lapTime, totalTime time.Duration) {
	now := NowFunc()
	return sw.lap(now, false), sw.totalTime(now)
}
