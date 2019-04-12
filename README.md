# Chronometry
Package for time-keeping purposes.

## Features

### Stopwatch
A monotonically stable stopwatch implementation which resembles a real-life stopwatch, best explained with some sample code:

    sw := chronometry.NewStartedStopwatch()
    time.Sleep(time.Second)
    sw.SplitTime() // 1 second
    sw.LapTime() // 1 second
    time.Sleep(time.Second)
    sw.Stop()
    sw.SplitTime() // 2 seconds
    sw.LapTime() // 1 second
    time.Sleep(time.Second)
    sw.SplitTime() // 2 seconds
    sw.LapTime() // 0 seconds
    sw.Resume()
    time.Sleep(time.Second)
    sw.Stop()
    time.Sleep(time.Second)
    sw.Resume()
    time.Sleep(time.Second)
    sw.SplitTime() // 4 seconds
    sw.LapTime() // 2 seconds
    sw.Restart()
    sw.SplitTime() // 0 seconds



## TODO
Mock time in tests without making the module more difficult to use.
