package chronometry_test

import (
	"testing"
	"time"

	"github.com/Zyl9393/chronometry"
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

func TestLap(t *testing.T) {
	sw := chronometry.NewStartedStopwatch()
	for i := 0; i < 3; i++ {
		testClock.advance(time.Second)
		duration := sw.Lap()
		durationsMatch(t, duration, time.Second, "")
	}
}

func TestTotal(t *testing.T) {
	sw := chronometry.NewStartedStopwatch()
	for i := 0; i < 3; i++ {
		testClock.advance(time.Second)
		duration := sw.Total()
		durationsMatch(t, duration, time.Second*time.Duration(i+1), "")
	}
}

func TestStoppedTotal(t *testing.T) {
	sw := chronometry.NewStoppedStopwatch()
	testClock.advance(time.Second)
	if sw.Total() != 0 {
		t.Fatalf("sw.SplitTime() was not 0.")
	}
}

func TestStoppedLap(t *testing.T) {
	sw := chronometry.NewStoppedStopwatch()
	testClock.advance(time.Second)
	if sw.Lap() != 0 {
		t.Fatalf("sw.Lap() was not 0.")
	}
}

func TestRestart(t *testing.T) {
	sw := chronometry.NewStartedStopwatch()
	testClock.advance(time.Second)
	sw.Restart()
	testClock.advance(time.Second)
	durationsMatch(t, sw.Total(), time.Second, "")
}

func TestStopResumeTotal(t *testing.T) {
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
	durationsMatch(t, sw.Total(), time.Second*3, "")
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
	durationsMatch(t, sw.Lap(), time.Second*3, "first try")
	testClock.advance(time.Second)
	sw.Stop()
	testClock.advance(time.Second)
	sw.Resume()
	testClock.advance(time.Second)
	sw.Stop()
	testClock.advance(time.Second)
	sw.Resume()
	testClock.advance(time.Second)
	durationsMatch(t, sw.Lap(), time.Second*3, "second try")
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
	durationsMatch(t, sw.Lap(), time.Second*3, "after stop")
}

func TestDocumentationClaims(t *testing.T) {
	sw := chronometry.NewStartedStopwatch()
	testClock.advance(time.Second)
	durationsMatch(t, sw.Total(), time.Second, "test 1")
	durationsMatch(t, sw.Lap(), time.Second, "test 2")
	testClock.advance(time.Second)
	sw.Stop()
	durationsMatch(t, sw.Total(), time.Second*2, "test 3")
	durationsMatch(t, sw.Lap(), time.Second, "test 4")
	testClock.advance(time.Second)
	durationsMatch(t, sw.Total(), time.Second*2, "test 5")
	durationsMatch(t, sw.Lap(), 0, "test 6")
	sw.Resume()
	testClock.advance(time.Second)
	sw.Stop()
	testClock.advance(time.Second)
	sw.Resume()
	testClock.advance(time.Second)
	durationsMatch(t, sw.Total(), time.Second*4, "test 7")
	durationsMatch(t, sw.Lap(), time.Second*2, "test 8")
	sw.Restart()
	durationsMatch(t, sw.Total(), 0, "test 9")
}

func TestPeekLap(t *testing.T) {
	sw := chronometry.NewStartedStopwatch()
	testClock.advance(time.Hour)
	durationsMatch(t, sw.Lap(), time.Hour, "test 1")
	testClock.advance(time.Second)
	durationsMatch(t, sw.PeekLap(), time.Second, "test 2")
	testClock.advance(time.Second)
	durationsMatch(t, sw.PeekLap(), time.Second*2, "test 3")
}

func durationsMatch(t *testing.T, actual, expected time.Duration, context string) {
	if actual > expected {
		t.Fatalf("duration too long; %s: expected: %v; actual: %v", context, expected, actual)
	}
	if actual < expected {
		t.Fatalf("duration too short; %s: expected: %v; actual: %v", context, expected, actual)
	}
}
