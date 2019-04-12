package chronometry

import (
	"time"
)

// Stopwatch implements a stopwatch which can be stopped and resumed and can report its total observed
// passed time as well as time intervals between "readings".
type Stopwatch struct {
	startTime, stopTime, readingTime              time.Time
	isStopped                                     bool
	passedDuration, accumulatedDifferenceDuration time.Duration
}

// NewStartedStopwatch returns a started stopwatch which will report an ever increasing duration with Elapsed() as well
// as time intervals between calls to ReadDifference().
func NewStartedStopwatch() *Stopwatch {
	now := time.Now()
	return &Stopwatch{startTime: now, stopTime: now, readingTime: now}
}

// NewStoppedStopwatch returns a stopwatch which needs to have Restart() or Resume() called on it in order to count time.
func NewStoppedStopwatch() *Stopwatch {
	now := time.Now()
	return &Stopwatch{startTime: now, stopTime: now, readingTime: now, isStopped: true}
}

// Restart resets and resumes the stopwatch.
func (sw *Stopwatch) Restart() {
	now := time.Now()
	sw.startTime, sw.stopTime, sw.readingTime, sw.isStopped, sw.passedDuration, sw.accumulatedDifferenceDuration = now, now, now, false, 0, 0
}

// IsStopped reports whether the stopwatch is stopped, i.e. not running.
func (sw *Stopwatch) IsStopped() bool {
	return sw.isStopped
}

// IsRunning reports whether the stopwatch is running, i.e. not stopped.
func (sw *Stopwatch) IsRunning() bool {
	return !sw.IsStopped()
}

// Reset sets alls durations tracked by the stopwatch back to zero.
func (sw *Stopwatch) Reset() {
	now := time.Now()
	sw.startTime, sw.readingTime, sw.stopTime, sw.passedDuration, sw.accumulatedDifferenceDuration = now, now, now, 0, 0
}

// Resume resumes the stopwatch when it was stopped.
func (sw *Stopwatch) Resume() {
	if sw.isStopped {
		now := time.Now()
		sw.startTime, sw.readingTime, sw.isStopped = now, now, false
	}
}

// Stop causes the stopwatch to stop observing the passing of time. It can be resumed with Resume().
func (sw *Stopwatch) Stop() time.Duration {
	if !sw.isStopped {
		now := time.Now()
		sw.stopTime, sw.passedDuration, sw.accumulatedDifferenceDuration, sw.isStopped =
			now, sw.passedDuration+now.Sub(sw.startTime), sw.accumulatedDifferenceDuration+sw.GetSegmentDifference(now), true
	}
	return sw.passedDuration
}

// GetSegmentDifference returns the shortest duration since the stopwatch was last started, restarted, resumed or read with ReadDifference().
func (sw *Stopwatch) GetSegmentDifference(now time.Time) time.Duration {
	if sw.readingTime.Before(sw.startTime) {
		return now.Sub(sw.startTime)
	}
	return now.Sub(sw.readingTime)
}

// StartTime returns the time at which the stopwatch was last started, restarted or resumed.
func (sw *Stopwatch) StartTime() time.Time {
	return sw.startTime
}

// StopTime returns the time of the most recent call to Stop().
func (sw *Stopwatch) StopTime() time.Time {
	return sw.stopTime
}

// ReadingTime returns the time of the most recent call to ReadDifference().
func (sw *Stopwatch) ReadingTime() time.Time {
	return sw.readingTime
}

// ReadDifference returns the difference between the current stopwatch reading and the one from the previous call to ReadDifference().
func (sw *Stopwatch) ReadDifference() time.Duration {
	now := time.Now()
	difference := sw.accumulatedDifferenceDuration
	if !sw.isStopped {
		difference += sw.GetSegmentDifference(now)
	}
	sw.readingTime, sw.accumulatedDifferenceDuration = now, 0
	return difference
}

// Elapsed returns the total observed passed time of the stopwatch.
func (sw *Stopwatch) Elapsed() time.Duration {
	if sw.isStopped {
		return sw.passedDuration
	}
	return sw.passedDuration + time.Now().Sub(sw.startTime)
}
