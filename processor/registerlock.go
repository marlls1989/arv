package processor

import (
	"log"
)

func (s *processor) reglockEl(
	fifoIn <-chan regAddr,
	fifoOut chan<- regAddr,
	lockOut chan<- uint32) {

	go func() {
		defer close(lockOut)
		defer close(fifoOut)

		<-s.start
		fifoOut <- 0
		lockOut <- 0
		for in := range fifoIn {
			fifoOut <- in
			lockOut <- uint32(in)
		}
	}()
}

func (s *processor) registerLock(
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

	s.reglockEl(fifoIn, fifo[0], lock[0])

	for i := 1; i < stages-1; i++ {
		s.reglockEl(fifo[i-1], fifo[i], lock[i])
	}

	s.reglockEl(fifo[stages-2], fifoOut, lock[stages-1])

	go func() {
		defer close(lockedRegs)

		for {
			var out uint32 = 0
			for i := len(lock); i > 0; i-- {
				in, vin := <-lock[i-1]
				if !vin {
					return
				}
				out |= in & 0xFFFFFFFE
			}
			lockedRegs <- out
		}
	}()
}
