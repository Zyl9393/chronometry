# Chronometry
Package for time-keeping purposes.

## Features

### Stopwatch
A monotonically stable stopwatch implementation in Golang which resembles a real-life stopwatch, best explained with some sample code:

    sw := chronometry.NewStartedStopwatch()
    time.Sleep(time.Second)
    sw.TotalTime() // 1 second
    sw.TakeLapTime() // 1 second
    time.Sleep(time.Second)
    sw.Stop()
    sw.TotalTime() // 2 seconds
    sw.TakeLapTime() // 1 second
    time.Sleep(time.Second)
    sw.TotalTime() // 2 seconds
    sw.TakeLapTime() // 0 seconds
    sw.Resume()
    time.Sleep(time.Second)
    sw.Stop()
    time.Sleep(time.Second)
    sw.Resume()
    time.Sleep(time.Second)
    sw.TotalTime() // 4 seconds
    sw.TakeLapTime() // 2 seconds
    sw.Restart()
    sw.TotalTime() // 0 seconds
