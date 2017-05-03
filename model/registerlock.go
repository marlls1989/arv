package model

import (
	"log"
)

func (s *Model) reglockStage(
	fifoIn <-chan uint32,
	fifoOut, lockOut chan<- uint32) {

	go func() {
		defer close(lockOut)
		defer close(fifoOut)

		<-s.start
		lockOut <- 0
		fifoOut <- 0
		for in := range fifoIn {
			lockOut <- in
			fifoOut <- in
		}
	}()
}

func (s *Model) registerLock(
	fifoIn <-chan uint32,
	fifoOut, lockedRegisters chan<- uint32,
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

	s.reglockStage(fifoIn, fifo[0], lock[0])

	for i := 1; i < stages-1; i++ {
		s.reglockStage(fifo[i-1], fifo[i], lock[i])
	}

	s.reglockStage(fifo[stages-2], fifoOut, lock[stages-1])

	go func() {
		var h uint32
		defer close(lockedRegisters)

		for {
			h = 0
			for _, l := range lock {
				a, va := <-l
				if !va {
					return
				}
				h |= a
			}
			lockedRegisters <- h & 0xFFFFFFFE
		}
	}()
}
