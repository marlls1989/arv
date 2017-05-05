package model

import (
	"log"
)

func (s *Model) reglockStage(
	fifoIn, lockIn <-chan uint32,
	fifoOut chan<- uint32,
	lockOut ...chan<- uint32) {

	go func() {
		for _, c := range lockOut {
			defer close(c)
		}
		defer close(fifoOut)

		<-s.start
		for _, c := range lockOut {
			c <- 0
		}
		fifoOut <- 0
		for in := range fifoIn {
			l, lv := <-lockIn
			if !lv {
				return
			}
			fifoOut <- in
			for _, c := range lockOut {
				c <- (in | l) & 0xFFFFFFFE
			}
		}
	}()
}

func (s *Model) registerLock(
	fifoIn <-chan uint32,
	fifoOut, lockedRegs chan<- uint32,
	stages int) {

	if stages < 2 {
		log.Panic("registerLock should have at least two stages")
	}

	fifo := make([]chan uint32, stages-1)
	lock := make([]chan uint32, stages)

	for i := range fifo {
		fifo[i] = make(chan uint32)
	}

	for i := range lock {
		lock[i] = make(chan uint32)
	}

	s.pipeElement(uint32(0), lock[0])

	s.reglockStage(fifoIn, lock[0], fifo[0], lock[1])

	for i := 1; i < stages-1; i++ {
		s.reglockStage(fifo[i-1], lock[i], fifo[i], lock[i+1])
	}

	s.reglockStage(fifo[stages-2], lock[stages-1], fifoOut, lockedRegs)
}
