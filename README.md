# Chronometry
Package for time-keeping purposes.

## Features

### Stopwatch
A monotonically stable stopwatch implementation which resembles a real-life stopwatch, best explained with some sample code:

    sw := chronometry.NewStartedStopwatch()
    time.Sleep(time.Second)
    sw.Elapsed() // 1 second
    sw.ReadDifference() // 1 second
    time.Sleep(time.Second)
    sw.Stop()
    sw.Elapsed() // 2 seconds
    sw.ReadDifference() // 1 second
    time.Sleep(time.Second)
    sw.Elapsed() // 2 seconds
    sw.ReadDifference() // 0 seconds
    sw.Resume()
    time.Sleep(time.Second)
    sw.Stop()
    time.Sleep(time.Second)
    sw.Resume()
    time.Sleep(time.Second)
    sw.Elapsed() // 4 seconds
    sw.ReadDifference() // 2 seconds
    sw.Restart()
    sw.Elapsed() // 0 seconds



## TODO
Mock time in tests without making the module more difficult to use.
