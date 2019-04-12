package chronometry_test

import (
	"testing"
	"time"

	"github.com/MMulthaupt/chronometry"
)

func TestLapTime(t *testing.T) {
	sw := chronometry.NewStartedStopwatch()
	for i := 0; i < 3; i++ {
		time.Sleep(time.Second)
		duration := sw.LapTime()
		maxError := time.Millisecond * 100
		durationsMatch(t, duration, time.Second, maxError, "")
	}
}

func TestElapsed(t *testing.T) {
	sw := chronometry.NewStartedStopwatch()
	for i := 0; i < 3; i++ {
		time.Sleep(time.Second)
		duration := sw.SplitTime()
		maxError := time.Duration(i+1) * time.Millisecond * 100
		durationsMatch(t, duration, time.Second*time.Duration(i+1), maxError, "")
	}
}

func TestStoppedElapsed(t *testing.T) {
	sw := chronometry.NewStoppedStopwatch()
	time.Sleep(time.Second)
	if sw.SplitTime() != 0 {
		t.Fatalf("sw.SplitTime() was not 0.")
	}
}

func TestStoppedLapTime(t *testing.T) {
	sw := chronometry.NewStoppedStopwatch()
	time.Sleep(time.Second)
	if sw.LapTime() != 0 {
		t.Fatalf("sw.LapTime() was not 0.")
	}
}

func TestRestart(t *testing.T) {
	sw := chronometry.NewStartedStopwatch()
	time.Sleep(time.Second)
	sw.Restart()
	time.Sleep(time.Second)
	durationsMatch(t, sw.SplitTime(), time.Second, time.Millisecond*100, "")
}

func TestStopResumeSplitTime(t *testing.T) {
	sw := chronometry.NewStartedStopwatch()
	time.Sleep(time.Second)
	sw.Stop()
	time.Sleep(time.Second)
	sw.Resume()
	time.Sleep(time.Second)
	sw.Stop()
	time.Sleep(time.Second)
	sw.Resume()
	time.Sleep(time.Second)
	durationsMatch(t, sw.SplitTime(), time.Second*3, time.Millisecond*100, "")
}

func TestStopResumeLapTime(t *testing.T) {
	sw := chronometry.NewStartedStopwatch()
	time.Sleep(time.Second)
	sw.Stop()
	time.Sleep(time.Second)
	sw.Resume()
	time.Sleep(time.Second)
	sw.Stop()
	time.Sleep(time.Second)
	sw.Resume()
	time.Sleep(time.Second)
	durationsMatch(t, sw.LapTime(), time.Second*3, time.Millisecond*100, "first try")
	time.Sleep(time.Second)
	sw.Stop()
	time.Sleep(time.Second)
	sw.Resume()
	time.Sleep(time.Second)
	sw.Stop()
	time.Sleep(time.Second)
	sw.Resume()
	time.Sleep(time.Second)
	durationsMatch(t, sw.LapTime(), time.Second*3, time.Millisecond*100, "second try")
	time.Sleep(time.Second)
	sw.Stop()
	time.Sleep(time.Second)
	sw.Resume()
	time.Sleep(time.Second)
	sw.Stop()
	time.Sleep(time.Second)
	sw.Resume()
	time.Sleep(time.Second)
	sw.Stop()
	durationsMatch(t, sw.LapTime(), time.Second*3, time.Millisecond*100, "after stop")
}

func TestDocumentationClaims(t *testing.T) {
	sw := chronometry.NewStartedStopwatch()
	time.Sleep(time.Second)
	durationsMatch(t, sw.SplitTime(), time.Second, time.Millisecond*100, "test 1")
	durationsMatch(t, sw.LapTime(), time.Second, time.Millisecond*100, "test 2")
	time.Sleep(time.Second)
	sw.Stop()
	durationsMatch(t, sw.SplitTime(), time.Second*2, time.Millisecond*100, "test 3")
	durationsMatch(t, sw.LapTime(), time.Second, time.Millisecond*100, "test 4")
	time.Sleep(time.Second)
	durationsMatch(t, sw.SplitTime(), time.Second*2, time.Millisecond*100, "test 5")
	durationsMatch(t, sw.LapTime(), 0, time.Millisecond*100, "test 6")
	sw.Resume()
	time.Sleep(time.Second)
	sw.Stop()
	time.Sleep(time.Second)
	sw.Resume()
	time.Sleep(time.Second)
	durationsMatch(t, sw.SplitTime(), time.Second*4, time.Millisecond*100, "test 7")
	durationsMatch(t, sw.LapTime(), time.Second*2, time.Millisecond*100, "test 8")
	sw.Restart()
	durationsMatch(t, sw.SplitTime(), 0, time.Millisecond*100, "test 9")
}

func durationsMatch(t *testing.T, a, b, maxError time.Duration, context string) {
	if a > b+maxError {
		t.Fatalf("duration 'a' too long. (%s) a: %v; b: %v; maxError: %v", context, a, b, maxError)
	}
	if a < b-maxError {
		t.Fatalf("duration 'a' too short. (%s) a: %v; b: %v; maxError: %v", context, a, b, maxError)
	}
}
