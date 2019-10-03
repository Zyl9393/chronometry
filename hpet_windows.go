package chronometry

import (
	"fmt"
	"log"
	"time"

	"golang.org/x/sys/windows"

	"unsafe"
)

var tickFrequency int64
var queryPerformanceCounterProc *windows.Proc
var lpMonotonicReference time.Duration
var hpMonotonicReference time.Duration

func init() {
	kernel32 := windows.MustLoadDLL("kernel32.dll")
	queryPerformanceFrequencyProc := kernel32.MustFindProc("QueryPerformanceFrequency")
	queryPerformanceCounterProc = kernel32.MustFindProc("QueryPerformanceCounter")

	_, _, err := queryPerformanceFrequencyProc.Call(uintptr(unsafe.Pointer(&tickFrequency)))
	if err != nil {
		winErr, ok := err.(windows.Errno)
		if !ok || winErr != 0 {
			panic(fmt.Sprintf("QueryPerformanceFrequency() failed: %v", err))
		}
	}

	determineReferenceTimes()
}

func determineReferenceTimes() {
	var lpNow time.Time
	var before, after, diff, bestDiff time.Duration
	bestDiff = time.Duration(1<<63 - 1)
	const maxTries = 100
	for i := 0; i < maxTries; i++ {
		before = hpet()
		lpNow = time.Now()
		after = hpet()
		diff = after - before
		if diff <= bestDiff {
			lpMonotonicReference, hpMonotonicReference = getMonotonic(&lpNow), after
			bestDiff = diff
		}
	}
	if bestDiff >= time.Microsecond {
		log.Printf("Warning: large bestDiff: %v\n", bestDiff)
	}
}

func hpNow() time.Time {
	var now time.Time
	var hpNow time.Duration
	now = time.Now()
	hpNow = hpet()
	setMonotonic(&now, lpMonotonicReference+(hpNow-hpMonotonicReference))
	return now
}

func getMonotonic(t *time.Time) time.Duration {
	tint := (*[2]uint64)(unsafe.Pointer(t))
	if (*tint)[0]&(1<<63) != 0 {
		return time.Duration((*tint)[1])
	}
	panic("no monotonic time information in t")
}

func setMonotonic(t *time.Time, m time.Duration) {
	tint := (*[2]uint64)(unsafe.Pointer(t))
	if (*tint)[0]&(1<<63) != 0 {
		(*tint)[1] = uint64(m)
		return
	}
	panic("no monotonic time information in t")
}

func hpet() time.Duration {
	var tickCount int64
	queryPerformanceCounterProc.Call(uintptr(unsafe.Pointer(&tickCount)))
	return time.Duration((tickCount/tickFrequency*1e9)+(tickCount%tickFrequency*1e9/tickFrequency)) * time.Nanosecond
}
