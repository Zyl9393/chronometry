package chronometry_test

import (
	"testing"
	"time"

	"github.com/MMulthaupt/chronometry"
)

func TestTakeLapTime(t *testing.T) {
	sw := chronometry.NewStartedStopwatch()
	for i := 0; i < 3; i++ {
		time.Sleep(time.Second)
		duration := sw.TakeLapTime()
		maxError := time.Millisecond * 100
		durationsMatch(t, duration, time.Second, maxError, "")
	}
}

func TestTotalTime(t *testing.T) {
	sw := chronometry.NewStartedStopwatch()
	for i := 0; i < 3; i++ {
		time.Sleep(time.Second)
		duration := sw.TotalTime()
		maxError := time.Duration(i+1) * time.Millisecond * 100
		durationsMatch(t, duration, time.Second*time.Duration(i+1), maxError, "")
	}
}

func TestStoppedTotalTime(t *testing.T) {
	sw := chronometry.NewStoppedStopwatch()
	time.Sleep(time.Second)
	if sw.TotalTime() != 0 {
		t.Fatalf("sw.SplitTime() was not 0.")
	}
}

func TestStoppedTakeLapTime(t *testing.T) {
	sw := chronometry.NewStoppedStopwatch()
	time.Sleep(time.Second)
	if sw.TakeLapTime() != 0 {
		t.Fatalf("sw.TakeLapTime() was not 0.")
	}
}

func TestRestart(t *testing.T) {
	sw := chronometry.NewStartedStopwatch()
	time.Sleep(time.Second)
	sw.Restart()
	time.Sleep(time.Second)
	durationsMatch(t, sw.TotalTime(), time.Second, time.Millisecond*100, "")
}

func TestStopResumeTotalTime(t *testing.T) {
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
	durationsMatch(t, sw.TotalTime(), time.Second*3, time.Millisecond*100, "")
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
	durationsMatch(t, sw.TakeLapTime(), time.Second*3, time.Millisecond*100, "first try")
	time.Sleep(time.Second)
	sw.Stop()
	time.Sleep(time.Second)
	sw.Resume()
	time.Sleep(time.Second)
	sw.Stop()
	time.Sleep(time.Second)
	sw.Resume()
	time.Sleep(time.Second)
	durationsMatch(t, sw.TakeLapTime(), time.Second*3, time.Millisecond*100, "second try")
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
	durationsMatch(t, sw.TakeLapTime(), time.Second*3, time.Millisecond*100, "after stop")
}

func TestDocumentationClaims(t *testing.T) {
	sw := chronometry.NewStartedStopwatch()
	time.Sleep(time.Second)
	durationsMatch(t, sw.TotalTime(), time.Second, time.Millisecond*100, "test 1")
	durationsMatch(t, sw.TakeLapTime(), time.Second, time.Millisecond*100, "test 2")
	time.Sleep(time.Second)
	sw.Stop()
	durationsMatch(t, sw.TotalTime(), time.Second*2, time.Millisecond*100, "test 3")
	durationsMatch(t, sw.TakeLapTime(), time.Second, time.Millisecond*100, "test 4")
	time.Sleep(time.Second)
	durationsMatch(t, sw.TotalTime(), time.Second*2, time.Millisecond*100, "test 5")
	durationsMatch(t, sw.TakeLapTime(), 0, time.Millisecond*100, "test 6")
	sw.Resume()
	time.Sleep(time.Second)
	sw.Stop()
	time.Sleep(time.Second)
	sw.Resume()
	time.Sleep(time.Second)
	durationsMatch(t, sw.TotalTime(), time.Second*4, time.Millisecond*100, "test 7")
	durationsMatch(t, sw.TakeLapTime(), time.Second*2, time.Millisecond*100, "test 8")
	sw.Restart()
	durationsMatch(t, sw.TotalTime(), 0, time.Millisecond*100, "test 9")
}

func durationsMatch(t *testing.T, a, b, maxError time.Duration, context string) {
	if a > b+maxError {
		t.Fatalf("duration 'a' too long. (%s) a: %v; b: %v; maxError: %v", context, a, b, maxError)
	}
	if a < b-maxError {
		t.Fatalf("duration 'a' too short. (%s) a: %v; b: %v; maxError: %v", context, a, b, maxError)
	}
}
