package processor

import (
	"log"
)

type retireRegwCmd struct {
	we   bool
	data uint32
}

func (s *Processor) retireUnit(
	qIn <-chan programElement,
	bypassIn <-chan uint32,
	adderIn <-chan uint32,
	logicIn <-chan uint32,
	shifterIn <-chan uint32,
	memoryIn <-chan memoryUnitOutput,
	branchIn <-chan branchOutput,

	regWcmd chan<- retireRegwCmd,
	memoryWe chan<- bool,
	branchOut chan<- branchCmd) {

	currValid := make(chan uint8)
	nextValid := make(chan uint8)

	// utilizes a loop to keep track of the current flow validity flag
	go func() {
		defer close(currValid)

		<-s.start
		currValid <- 0
		for i := range nextValid {
			currValid <- i
		}
	}()

	go func() {
		defer close(regWcmd)
		defer close(memoryWe)
		defer close(branchOut)
		defer close(nextValid)

		for q := range qIn {
			var data uint32
			valid := <-currValid
			rwe := true
			memWe := false
			brCmd := branchCmd{
				taken: false}

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
				brCmd.taken = br.taken
				brCmd.target = br.target
				rwe = br.link
				data = br.linkAddr
				/* if a branch is taken, increment the validity flag
				 * to execute only the valid flow  */
			}

			/* Case the validty flag of the current instruction
			 * and the validity flag of the current flow mismatch
			 * invalidate the instruction */
			if valid != q.valid {
				log.Printf("Canceling Instruction [q: %+v br: %+v data: %d rwe: %v mwe: %v]", q, brCmd, data, rwe, memWe)
				brCmd.taken = false
				rwe = false

				/* memoryWe is peculiar, it should not be handshake unless the memory unit
				 * has signed it is expecting a handshake to complete a memory write */
				if memWe {
					memoryWe <- false
				}
			} else if memWe {
				memoryWe <- true
			}

			if brCmd.taken {
				nextValid <- valid + 1
			} else {
				nextValid <- valid
			}

			regWcmd <- retireRegwCmd{
				we:   rwe,
				data: data}
			if brCmd.taken {
				branchOut <- brCmd
			}
		}

	}()

}
