package chronometry

import "time"

// BenchExecutionTime returns the minimum expectable execution time of the given function f.
// f is called multiple times to determine this.
func BenchExecutionTime(f func()) time.Duration {
	viableSampleCount := 10
	const maximumSampleCount = 100000
	const runCount = 4
	var averages [runCount]time.Duration
	var bests [runCount]time.Duration
	reachedViableSampleCount := false
	for r := 0; r < runCount; r++ {
		averages[r] = benchAverage(f, viableSampleCount)
		if !reachedViableSampleCount && viableSampleCount < maximumSampleCount {
			viableSampleCount *= 10
			reachedViableSampleCount = averages[r] > 0
			if viableSampleCount > maximumSampleCount {
				viableSampleCount = maximumSampleCount
			}
			r = -1
			continue
		}
		bests[r] = benchBest(f, viableSampleCount)
	}
	bestAverage := smallestNonZero(averages[:])
	best := smallestNonZero(bests[:])
	if best < bestAverage && best != 0 {
		return best
	}
	return bestAverage
}

func benchAverage(f func(), sampleCount int) time.Duration {
	before := HPET()
	for s := 0; s < sampleCount; s++ {
		f()
	}
	return (HPET() - before) / time.Duration(sampleCount)
}

func benchBest(f func(), sampleCount int) (best time.Duration) {
	var delta time.Duration
	var before, after time.Duration
	for s := 0; s < sampleCount; s++ {
		before = HPET()
		f()
		after = HPET()
		delta = after - before
		if delta < best {
			best = delta
		}
	}
	return best
}

func smallestNonZero(durations []time.Duration) (smallest time.Duration) {
	for _, d := range durations {
		if d > 0 && (d < smallest || smallest == 0) {
			smallest = d
		}
	}
	return smallest
}

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
