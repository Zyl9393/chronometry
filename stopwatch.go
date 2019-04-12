package chronometry

import (
	"time"
)

// Stopwatch implements a stopwatch which can be stopped, resumed and reset and can report its current reading
// of time with SplitTime() as well as allow for convenient measuring of consecutive intervals using LapTime().
type Stopwatch struct {
	startTime, stopTime, readingTime              time.Time
	isStopped                                     bool
	passedDuration, accumulatedDifferenceDuration time.Duration
}

// NewStartedStopwatch returns a started stopwatch, ready to report SplitTime() and LapTime().
func NewStartedStopwatch() *Stopwatch {
	now := time.Now()
	return &Stopwatch{startTime: now, stopTime: now, readingTime: now}
}

// NewStoppedStopwatch returns a stopwatch which needs to have Restart() or Resume() called on it in order to
// begin counting time.
func NewStoppedStopwatch() *Stopwatch {
	now := time.Now()
	return &Stopwatch{startTime: now, stopTime: now, readingTime: now, isStopped: true}
}

// Restart resets and resumes the stopwatch.
func (sw *Stopwatch) Restart() {
	sw.Reset()
	sw.isStopped = false
}

// IsStopped reports whether the stopwatch is stopped, i.e. not running.
func (sw *Stopwatch) IsStopped() bool {
	return sw.isStopped
}

// IsRunning reports whether the stopwatch is running, i.e. not stopped.
func (sw *Stopwatch) IsRunning() bool {
	return !sw.IsStopped()
}

// Reset sets the stopwatch reading back to zero.
func (sw *Stopwatch) Reset() {
	now := time.Now()
	sw.startTime, sw.readingTime, sw.stopTime = now, now, now
	sw.passedDuration, sw.accumulatedDifferenceDuration = 0, 0
}

// Resume resumes the stopwatch when it is stopped.
func (sw *Stopwatch) Resume() {
	if sw.isStopped {
		now := time.Now()
		sw.startTime, sw.readingTime, sw.isStopped = now, now, false
	}
}

// Stop stops the stopwatch, causing it to stop observing the passing of time. It can be resumed with Resume().
func (sw *Stopwatch) Stop() time.Duration {
	if !sw.isStopped {
		now := time.Now()
		sw.stopTime, sw.isStopped = now, true
		sw.passedDuration += now.Sub(sw.startTime)
		sw.accumulatedDifferenceDuration += sw.currentSegmentDuration(now)
	}
	return sw.passedDuration
}

// CurrentSegmentDuration returns the duration since the stopwatch was last started/restarted, reset, resumed
// or had LapTime() called.
func (sw *Stopwatch) CurrentSegmentDuration() time.Duration {
	return sw.currentSegmentDuration(time.Now())
}

func (sw *Stopwatch) currentSegmentDuration(now time.Time) time.Duration {
	if sw.readingTime.Before(sw.startTime) {
		return now.Sub(sw.startTime)
	}
	return now.Sub(sw.readingTime)
}

// StartTime returns the time at which the stopwatch was last started/restarted, reset or resumed.
func (sw *Stopwatch) StartTime() time.Time {
	return sw.startTime
}

// StopTime returns the time at which the stopwatch was last stopped, started/restarted or reset.
func (sw *Stopwatch) StopTime() time.Time {
	return sw.stopTime
}

// ReadingTime returns the time at which the stopwatch was last started/restarted, reset or had LapTime() called.
func (sw *Stopwatch) ReadingTime() time.Time {
	return sw.readingTime
}

// LapTime makes a split, returning the change in SplitTime() since last calling this function
// or starting/restarting or resetting the stopwatch.
func (sw *Stopwatch) LapTime() time.Duration {
	now := time.Now()
	difference := sw.accumulatedDifferenceDuration
	if !sw.isStopped {
		difference += sw.currentSegmentDuration(now)
	}
	sw.readingTime, sw.accumulatedDifferenceDuration = now, 0
	return difference
}

// SplitTime returns the current stopwatch reading, i.e. the total passed time observed by the stopwatch while not stopped.
func (sw *Stopwatch) SplitTime() time.Duration {
	if sw.isStopped {
		return sw.passedDuration
	}
	return sw.passedDuration + time.Now().Sub(sw.startTime)
}
