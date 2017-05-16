package processor

import (
	"log"
)

func (s *Processor) reglockEl(
	fifoIn <-chan regAddr,
	lockIn <-chan uint32,
	fifoOut chan<- regAddr,
	lockOut chan<- uint32) {

	go func() {
		defer close(lockOut)
		defer close(fifoOut)

		<-s.start
		fifoOut <- 0
		lockOut <- 0
		for in := range fifoIn {
			l, lv := <-lockIn
			if !lv {
				return
			}
			lock := (uint32(in) | l) & 0xFFFFFFFE
			fifoOut <- in
			lockOut <- lock
		}
	}()
}

func (s *Processor) registerLock(
	fifoIn <-chan regAddr,
	fifoOut chan<- regAddr,
	lockedRegs chan<- uint32,
	stages int) {

	if stages < 2 {
		log.Panic("registerLock should have at least two stages")
	}

	fifo := make([]chan regAddr, stages-1)
	lock := make([]chan uint32, stages)

	for i := range fifo {
		fifo[i] = make(chan regAddr)
	}

	for i := range lock {
		lock[i] = make(chan uint32)
	}

	go func() {
		defer close(lock[0])

		for {
			select {
			case lock[0] <- 0:
			case <-s.quit:
				return
			}
		}
	}()

	s.reglockEl(fifoIn, lock[0], fifo[0], lock[1])

	for i := 1; i < stages-1; i++ {
		s.reglockEl(fifo[i-1], lock[i], fifo[i], lock[i+1])
	}

	s.reglockEl(fifo[stages-2], lock[stages-1], fifoOut, lockedRegs)
}
