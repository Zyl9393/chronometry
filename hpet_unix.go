// +build !windows

package chronometry

import "time"

var referenceTime = time.Now()

func hpNow() time.Time {
	return time.Now()
}

func hpet() time.Duration {
	return time.Now().Sub(referenceTime)
}
