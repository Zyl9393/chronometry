package chronometry

import "time"

// Now returns the current time with high precision monotonic reading.
func Now() time.Time {
	return hpNow()
}

// Since returns the precise time elapsed since t. It is shorthand for chronometry.Now().Sub(t).
func Since(t time.Time) time.Duration {
	return Now().Sub(t)
}

// Until returns the precise duration until t. It is shorthand for t.Sub(chronometry.Now()).
func Until(t time.Time) time.Duration {
	return t.Sub(Now())
}

// HPET returns the reading of the system's monotonic High Precision Event Timer as a time.Duration.
func HPET() time.Duration {
	return hpet()
}
