package chronometry_test

import (
	"testing"
	"time"

	"github.com/MMulthaupt/chronometry"
)

type clock struct {
	clockTime time.Time
}

var testClock = &clock{clockTime: time.Unix(60*60*24*365*70, 345678)}

func (c *clock) now() time.Time {
	return c.clockTime
}

func (c *clock) advance(duration time.Duration) {
	c.clockTime = c.clockTime.Add(duration)
}

func init() {
	chronometry.NowFunc = testClock.now
}

func TestTakeLapTime(t *testing.T) {
	sw := chronometry.NewStartedStopwatch()
	for i := 0; i < 3; i++ {
		testClock.advance(time.Second)
		duration := sw.TakeLapTime()
		durationsMatch(t, duration, time.Second, "")
	}
}

func TestTotalTime(t *testing.T) {
	sw := chronometry.NewStartedStopwatch()
	for i := 0; i < 3; i++ {
		testClock.advance(time.Second)
		duration := sw.TotalTime()
		durationsMatch(t, duration, time.Second*time.Duration(i+1), "")
	}
}

func TestStoppedTotalTime(t *testing.T) {
	sw := chronometry.NewStoppedStopwatch()
	testClock.advance(time.Second)
	if sw.TotalTime() != 0 {
		t.Fatalf("sw.SplitTime() was not 0.")
	}
}

func TestStoppedTakeLapTime(t *testing.T) {
	sw := chronometry.NewStoppedStopwatch()
	testClock.advance(time.Second)
	if sw.TakeLapTime() != 0 {
		t.Fatalf("sw.TakeLapTime() was not 0.")
	}
}

func TestRestart(t *testing.T) {
	sw := chronometry.NewStartedStopwatch()
	testClock.advance(time.Second)
	sw.Restart()
	testClock.advance(time.Second)
	durationsMatch(t, sw.TotalTime(), time.Second, "")
}

func TestStopResumeTotalTime(t *testing.T) {
	sw := chronometry.NewStartedStopwatch()
	testClock.advance(time.Second)
	sw.Stop()
	testClock.advance(time.Second)
	sw.Resume()
	testClock.advance(time.Second)
	sw.Stop()
	testClock.advance(time.Second)
	sw.Resume()
	testClock.advance(time.Second)
	durationsMatch(t, sw.TotalTime(), time.Second*3, "")
}

func TestStopResumeLapTime(t *testing.T) {
	sw := chronometry.NewStartedStopwatch()
	testClock.advance(time.Second)
	sw.Stop()
	testClock.advance(time.Second)
	sw.Resume()
	testClock.advance(time.Second)
	sw.Stop()
	testClock.advance(time.Second)
	sw.Resume()
	testClock.advance(time.Second)
	durationsMatch(t, sw.TakeLapTime(), time.Second*3, "first try")
	testClock.advance(time.Second)
	sw.Stop()
	testClock.advance(time.Second)
	sw.Resume()
	testClock.advance(time.Second)
	sw.Stop()
	testClock.advance(time.Second)
	sw.Resume()
	testClock.advance(time.Second)
	durationsMatch(t, sw.TakeLapTime(), time.Second*3, "second try")
	testClock.advance(time.Second)
	sw.Stop()
	testClock.advance(time.Second)
	sw.Resume()
	testClock.advance(time.Second)
	sw.Stop()
	testClock.advance(time.Second)
	sw.Resume()
	testClock.advance(time.Second)
	sw.Stop()
	durationsMatch(t, sw.TakeLapTime(), time.Second*3, "after stop")
}

func TestDocumentationClaims(t *testing.T) {
	sw := chronometry.NewStartedStopwatch()
	testClock.advance(time.Second)
	durationsMatch(t, sw.TotalTime(), time.Second, "test 1")
	durationsMatch(t, sw.TakeLapTime(), time.Second, "test 2")
	testClock.advance(time.Second)
	sw.Stop()
	durationsMatch(t, sw.TotalTime(), time.Second*2, "test 3")
	durationsMatch(t, sw.TakeLapTime(), time.Second, "test 4")
	testClock.advance(time.Second)
	durationsMatch(t, sw.TotalTime(), time.Second*2, "test 5")
	durationsMatch(t, sw.TakeLapTime(), 0, "test 6")
	sw.Resume()
	testClock.advance(time.Second)
	sw.Stop()
	testClock.advance(time.Second)
	sw.Resume()
	testClock.advance(time.Second)
	durationsMatch(t, sw.TotalTime(), time.Second*4, "test 7")
	durationsMatch(t, sw.TakeLapTime(), time.Second*2, "test 8")
	sw.Restart()
	durationsMatch(t, sw.TotalTime(), 0, "test 9")
}

func durationsMatch(t *testing.T, a, b time.Duration, context string) {
	if a > b {
		t.Fatalf("duration 'a' too long. (%s) a: %v; b: %v", context, a, b)
	}
	if a < b {
		t.Fatalf("duration 'a' too short. (%s) a: %v; b: %v", context, a, b)
	}
}
