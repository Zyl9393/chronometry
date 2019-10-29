package chronometry_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/Zyl9393/chronometry"
)

func TestHPETResolution(t *testing.T) {
	zeroCount := 0
	const loopCount = 100
	for i := 0; i < loopCount; i++ {
		before := chronometry.Now()
		doSmallWorkLoad()
		after := chronometry.Now()
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
	before := chronometry.Now()
	time.Sleep(time.Second)
	after := chronometry.Now()
	diff := after.Sub(before)
	if diff < 900000000 {
		t.Fatalf("error too large. diff was %v.", diff)
	} else if diff > 10000000000 {
		t.Fatalf("error too large. diff was %v.", diff)
	} else if diff > 1100000000 {
		t.Logf("error looks large. diff was %v. Ignore for slow test machines.", diff)
	}
}

func TestBenchHPETSpeed(t *testing.T) {
	const loopCount = 1000
	startTime := chronometry.Now()
	for i := 0; i < loopCount; i++ {
		chronometry.Now()
	}
	hpNowTime := chronometry.Now()
	for i := 0; i < loopCount; i++ {
		chronometry.HPET()
	}
	hpetTime := chronometry.Now()
	for i := 0; i < loopCount; i++ {
		time.Now()
	}
	goNowTime := chronometry.Now()
	fmt.Printf("chronometry.Now() takes %v per call. chronometry.HPET() takes %v per call. time.Now() takes %v per call.\n",
		hpNowTime.Sub(startTime)/loopCount, hpetTime.Sub(hpNowTime)/loopCount, goNowTime.Sub(hpetTime)/loopCount)
}

func doSmallWorkLoad() {
	const dataSize = 1000
	data := make([]int, dataSize)
	for i := 0; i < dataSize; i++ {
		data[i] = i*7 + 42
	}
}
