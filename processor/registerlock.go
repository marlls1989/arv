package processor

import (
	"log"
)

// This function constructs a register lock queue logical stage element.
//
// The logical stage constructed reads a one-hot encoded register address,
// holds and outputs it to the two output channels.
//
// The first output channel is chained to the next stage or connected to the register write address output.
// The second output connects to a stage responsible for ORing the address and producing the lock mask.
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

// Constructs a variable length register lock queue and the mask generator logic element.
//
// The one-hot encoded destination register address is feed in the queue when the instruction is dispatched,
// the queue is the same depth as the instruction execution path.
// The register file write port or controller synchronises the destination address with the data.
// This imply that the register lock queue holds the destination register of all instructions dispatched
// that have not committed results to the register file.
//
// Each element of the queue have two output channels, one is chained to the next stage up to the output,
// the other channel is connected to the masking logic stage.
// The one-hot encoded address sent by each stage secondary channel are ORed by the masking logic stage,
// the output is the lockedRegs bitmask used by the operandFetch unit to identify locked register.
//
// It is important to notice that the lockedRegs channel must be read prior to any writes to the fifoIn channel.
// Otherwise a deadlock would occur as the queue has no bubbles to accommodate new tokens.
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

	// Mask generating logical stage code
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
