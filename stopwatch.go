package chronometry

import (
	"time"
)

// NowFunc specifies the function used to retrieve the current time. It is changed from time.Now only during tests.
var NowFunc = time.Now

// Stopwatch implements a stopwatch which can be stopped, resumed and reset and can report its current reading
// of time with TotalTime() as well as allow for convenient measuring of consecutive intervals using TakeLapTime().
type Stopwatch struct {
	startTime, stopTime, readingTime      time.Time
	isStopped                             bool
	totalDuration, lapAccumulatorDuration time.Duration
}

// NewStartedStopwatch returns a started stopwatch, ready to report TotalTime() and TakeLapTime().
func NewStartedStopwatch() *Stopwatch {
	now := NowFunc()
	return &Stopwatch{startTime: now, stopTime: now, readingTime: now}
}

// NewStoppedStopwatch returns a stopwatch which needs to have Restart() or Resume() called on it in order to
// begin counting time.
func NewStoppedStopwatch() *Stopwatch {
	now := NowFunc()
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
	now := NowFunc()
	sw.startTime, sw.readingTime, sw.stopTime = now, now, now
	sw.totalDuration, sw.lapAccumulatorDuration = 0, 0
}

// Resume resumes the stopwatch when it is stopped.
func (sw *Stopwatch) Resume() {
	if sw.isStopped {
		now := NowFunc()
		sw.startTime, sw.readingTime, sw.isStopped = now, now, false
	}
}

// Stop stops the stopwatch, causing it to stop observing the passing of time. It can be resumed with Resume().
func (sw *Stopwatch) Stop() time.Duration {
	if !sw.isStopped {
		now := NowFunc()
		sw.stopTime, sw.isStopped = now, true
		sw.totalDuration += now.Sub(sw.startTime)
		sw.lapAccumulatorDuration += sw.currentSegmentDuration(now)
	}
	return sw.totalDuration
}

// CurrentSegmentDuration returns the duration since the stopwatch was last started/restarted, reset, resumed
// or had TakeLapTime() called.
func (sw *Stopwatch) CurrentSegmentDuration() time.Duration {
	return sw.currentSegmentDuration(NowFunc())
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

// ReadingTime returns the time at which the stopwatch was last started/restarted, reset, resumed or had TakeLapTime() called.
func (sw *Stopwatch) ReadingTime() time.Time {
	return sw.readingTime
}

// TakeLapTime makes a split, returning the change in TotalTime() since last calling this function
// or starting/restarting or resetting the stopwatch.
func (sw *Stopwatch) TakeLapTime() time.Duration {
	return sw.takeLapTime(NowFunc())
}

func (sw *Stopwatch) takeLapTime(now time.Time) time.Duration {
	difference := sw.lapAccumulatorDuration
	if !sw.isStopped {
		difference += sw.currentSegmentDuration(now)
	}
	sw.readingTime, sw.lapAccumulatorDuration = now, 0
	return difference
}

// TotalTime returns the current stopwatch reading, i.e. the total passed time observed by the stopwatch while not stopped.
func (sw *Stopwatch) TotalTime() time.Duration {
	return sw.totalTime(NowFunc())
}

func (sw *Stopwatch) totalTime(now time.Time) time.Duration {
	if sw.isStopped {
		return sw.totalDuration
	}
	return sw.totalDuration + now.Sub(sw.startTime)
}

// MakeSplit does the same as TakeLapTime(), but also returns the exact according total time.
func (sw *Stopwatch) MakeSplit() (lapTime, totalTime time.Duration) {
	now := NowFunc()
	return sw.takeLapTime(now), sw.totalTime(now)
}
