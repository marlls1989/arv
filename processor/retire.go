package processor

import (
	"log"
	"sync/atomic"
)

type regWCmd struct {
	we   bool
	data uint32
}

func (s *processor) retireUnit(
	qIn <-chan programElement,
	bypassIn <-chan uint32,
	adderIn <-chan uint32,
	logicIn <-chan uint32,
	shifterIn <-chan uint32,
	memoryIn <-chan memoryUnitOutput,
	branchIn <-chan branchOutput,

	regWcmd chan<- regWCmd,
	memoryWe chan<- bool,
	branchOut chan<- uint32) {

	currValid := make(chan uint8)
	nextValid := make(chan uint8)

	nextBr := make(chan uint32)

	go func() {
		defer close(branchOut)

		for i := range nextBr {
			branchOut <- i
		}
	}()

	// utilizes a loop to keep track of the current flow validity flag
	go func() {
		defer close(currValid)

		for i := range nextValid {
			currValid <- i
		}
	}()

	go func() {
		defer close(regWcmd)
		defer close(memoryWe)
		defer close(nextBr)
		defer close(nextValid)

		<-s.start
		nextValid <- 0
		regWcmd <- regWCmd{we: false}

		for q := range qIn {
			select {
			case _ = <-s.quit:
				return
			default:
			}

			var data, brTarget uint32
			valid := <-currValid
			rwe := true
			memWe := false
			brTaken := false

			/* Selectively handshake the instruction unit
			 * acording to program order as defined by the queue */
			switch q.unit {
			case xuBypassSel:
				data = <-bypassIn
			case xuAdderSel:
				data = <-adderIn
			case xuLogicSel:
				data = <-logicIn
			case xuShiftSel:
				data = <-shifterIn
			case xuMemorySel:
				meminfo := <-memoryIn
				memWe = meminfo.writeRequest
				data = meminfo.value
			case xuBranchSel:
				br := <-branchIn
				brTaken = br.taken
				brTarget = br.target
				rwe = br.link
				data = br.linkAddr
			}

			/* Case the validty flag of the current instruction
			 * and the validity flag of the current flow mismatch
			 * invalidate the instruction */
			if valid != q.valid {
				atomic.AddUint64(&(s.Cancelled), 1)
				if s.Debug {
					log.Printf("Canceling Instruction [q: %+v br: %v brTarget:%x data: %x rwe: %v mwe: %v]", q, brTaken, brTarget, data, rwe, memWe)
				}
				brTaken = false
				rwe = false

				/* memoryWe is peculiar, it should not be handshake unless the memory unit
				 * has signed it is expecting a handshake to complete a memory write */
				if memWe {
					memoryWe <- false
				}
			} else {
				atomic.AddUint64(&(s.Retired), 1)
				if s.Debug {
					log.Printf("Retiring Instruction [q: %+v br: %v brTarget:%x data: %x rwe: %v mwe: %v]", q, brTaken, brTarget, data, rwe, memWe)
				}
				if memWe {
					memoryWe <- true
				}
			}

			/* if a branch is taken, increment the validity flag
			 * to execute only the valid flow  */
			if brTaken {
				nextValid <- valid + 1
			} else {
				nextValid <- valid
			}

			regWcmd <- regWCmd{
				we:   rwe,
				data: data}
			if brTaken {
				nextBr <- brTarget
			}
		}
	}()

}
