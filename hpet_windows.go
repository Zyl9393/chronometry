package chronometry

import (
	"fmt"
	"log"
	"time"

	"golang.org/x/sys/windows"

	"unsafe"
)

/*
#include <windows.h>

LARGE_INTEGER hpetTicks() {
	LARGE_INTEGER tickCount;
	QueryPerformanceCounter(&tickCount);
	return tickCount;
}
*/
import "C"

var tickFrequency int64
var lpMonotonicReference time.Duration
var hpMonotonicReference time.Duration

func init() {
	kernel32 := windows.NewLazySystemDLL("kernel32.dll")
	queryPerformanceFrequencyProc := kernel32.NewProc("QueryPerformanceFrequency")

	tickFrequencyC := C.calloc(1, 8)
	r1, _, err := queryPerformanceFrequencyProc.Call(uintptr(unsafe.Pointer(tickFrequencyC)))
	tickFrequency = *((*int64)(unsafe.Pointer(tickFrequencyC)))
	C.free(tickFrequencyC)
	if r1 == 0 || err != nil {
		winErr, ok := err.(windows.Errno)
		if !ok || winErr != 0 {
			panic(fmt.Sprintf("QueryPerformanceFrequency() failed: return code %d: %v", r1, err))
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
	now, hpNow = time.Now(), hpet()
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
	tickCountC := C.hpetTicks()
	tickCount := *((*int64)(unsafe.Pointer(&tickCountC)))
	return time.Duration((tickCount/tickFrequency*1e9)+(tickCount%tickFrequency*1e9/tickFrequency)) * time.Nanosecond
}
