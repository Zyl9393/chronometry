package chronometry_test

import (
	"testing"
	"time"

	"github.com/MMulthaupt/chronometry"
)

func TestHPETResolution(t *testing.T) {
	zeroCount := 0
	const loopCount = 100
	for i := 0; i < loopCount; i++ {
		before := chronometry.HPNow()
		doSmallWorkLoad()
		after := chronometry.HPNow()
		diff := after.Sub(before)
		if diff == 0 {
			zeroCount++
		}
	}
	if zeroCount > 0 {
		t.Fatalf("diff was zero in %d out of %d instances", zeroCount, loopCount)
	}
}

func TestHPETAccuracy(t *testing.T) {
	before := chronometry.HPNow()
	time.Sleep(time.Second)
	after := chronometry.HPNow()
	diff := after.Sub(before)
	if diff > 1100000000 || diff < 900000000 {
		t.Fatalf("error too large. diff was %v", diff)
	}
}

func TestBenchHPETSpeed(t *testing.T) {
	const loopCount = 1000
	first := chronometry.HPNow()
	for i := 0; i < loopCount; i++ {
		chronometry.HPNow()
	}
	last := chronometry.HPNow()
	t.Logf("chronometry.HPNow() takes %v per call.", last.Sub(first)/loopCount)
}

func doSmallWorkLoad() {
	const dataSize = 10000
	data := make([]int, dataSize)
	for i := 0; i < dataSize; i++ {
		data[i] = i*7 + 42
	}
}
