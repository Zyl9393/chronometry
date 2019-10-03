package chronometry

import "time"

// HPNow returns the current time with high precision monotonic reading.
func HPNow() time.Time {
	return hpNow()
}

// HPET returns the reading of the system's monotonic High Precision Event Timer as a time.Duration.
func HPET() time.Duration {
	return hpet()
}
