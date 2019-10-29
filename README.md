# Chronometry

<p align="left">
  <a href="https://github.com/Zyl9393/chronometry/tree/master"><img alt="Status of tests" src="https://github.com/Zyl9393/chronometry/workflows/Tests/badge.svg?branch=master"></a>
</p>

Package for time-keeping purposes.

## Features

### Stopwatch
A monotonically stable stopwatch with high precision, best explained with some sample code:

```go
sw := chronometry.NewStartedStopwatch()
time.Sleep(time.Second)
sw.Total() // 1 second
sw.Lap() // 1 second
time.Sleep(time.Second)
sw.Stop()
sw.Total() // 2 seconds
sw.Lap() // 1 second
time.Sleep(time.Second)
sw.Total() // 2 seconds
sw.Lap() // 0 seconds
sw.Resume()
time.Sleep(time.Second)
sw.Stop()
time.Sleep(time.Second)
sw.Resume()
time.Sleep(time.Second)
sw.Total() // 4 seconds
sw.Lap() // 2 seconds
sw.Restart()
sw.Total() // 0 seconds
```

### High Precision Event Timer and Clock
Can call `chronometry.Now()` for an equivalent of `time.Now()` with a high-resolution monotonic timestamp. Can call `chronometry.HPET()` for just the high-resolution monotonic timestamp.

### Known issues
Because the library needs to call C functions for the increased precision on Windows, there, each invocation of `chronometry.Now()` takes about 20 times as long as `time.Now()` (on my test machine, 80 nanoseconds instead of 4). This is not fixable.
