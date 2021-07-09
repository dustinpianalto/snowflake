package generator

import (
	"errors"
	"time"
)

var EPOCH_TIME time.Time = time.Unix(0, 1609459200000*int64(time.Millisecond))
var TIME_MASK uint64 = 0x1FFFFFFFFFF

var tooManyRequests = errors.New("Too many requests in the current ms")
var incorrectSystemTime = errors.New("The current time is less than the last generated time! Check the system time.")

type Generator struct {
	LastGeneratedTime   time.Duration
	Counter             uint64
	LastCounterRollover time.Duration
	WorkerID            uint64
}

func (g *Generator) GenerateSnowflake() (uint64, error) {
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

func CreateGenerator(workerID uint64) (*Generator, error) {
	if workerID <= 0 {
		return nil, errors.New("Worker ID must be greater than 0")
	}
	return &Generator{
		LastGeneratedTime:   time.Since(EPOCH_TIME),
		Counter:             0,
		LastCounterRollover: time.Since(EPOCH_TIME),
		WorkerID:            workerID,
	}, nil
}

func (g *Generator) Run(request chan chan uint64) {
	for output := range request {
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
