package generator

import (
	"errors"
	"time"
)

var Generator *generator

var EPOCH_TIME time.Time = time.Unix(0, 1609459200000*int64(time.Millisecond))
var TIME_MASK uint64 = 0x1FFFFFFFFFF

var tooManyRequests = errors.New("Too many requests in the current ms")
var incorrectSystemTime = errors.New("The current time is less than the last generated time! Check the system time.")

type generator struct {
	LastGeneratedTime   time.Duration
	Counter             uint64
	LastCounterRollover time.Duration
	WorkerID            uint64
	RequestChan         chan chan uint64
}

func (g *generator) GenerateSnowflake() (uint64, error) {
	time := time.Since(EPOCH_TIME)
	if time < g.LastGeneratedTime {
		return 0, incorrectSystemTime
	}
	if g.Counter == 4095 {
		if g.LastCounterRollover >= time {
			return 0, tooManyRequests
		}
		g.LastCounterRollover = time
		g.Counter = 0
	} else {
		g.Counter += 1
	}
	g.LastGeneratedTime = time
	return (((uint64(time.Milliseconds()) & TIME_MASK) << 22) + g.WorkerID + g.Counter), nil
}

func CreateGenerator(workerID uint64) {
	Generator = &generator{
		LastGeneratedTime:   time.Since(EPOCH_TIME),
		Counter:             0,
		LastCounterRollover: time.Since(EPOCH_TIME),
		WorkerID:            workerID,
		RequestChan:         make(chan chan uint64),
	}
}

func (g *generator) Run() {
	for output := range g.RequestChan {
		id, err := g.GenerateSnowflake()
		if err != nil {
			if err == tooManyRequests {
				time.Sleep(time.Millisecond)
				continue
			} else if err == incorrectSystemTime {
				panic("System time is incorrect")
			}
		}
		output <- id
	}
}
