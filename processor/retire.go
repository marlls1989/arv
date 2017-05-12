package processor

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

	regWData chan<- uint32,
	regWe chan<- bool,
	memoryWe chan<- bool,
	branchOut chan<- branchCmd) {

	currValid := make(chan bool)
	nextValid := make(chan bool)

	// utilizes a loop to keep track of the current flow validity flag
	s.pipeElementWithInitization(nextValid, true, currValid)

	go func() {
		defer close(regWData)
		defer close(regWe)
		defer close(memoryWe)
		defer close(branchOut)
		defer close(nextValid)

		for q := range qIn {

			var data uint32
			valid := <-currValid
			nValid := valid
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
				/* if a branch is taken, flip the validity flag
				 * of the instruction flow for the next instructions */
				nValid = brCmd.taken != valid // XORing to negate

			}

			/* Case the validty flag of the current instruction
			 * and the validity flag of the current flow mismatch
			 * invalidate the instruction */
			if valid != q.valid {
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

			branchOut <- brCmd
			regWe <- rwe
			regWData <- data
			nextValid <- nValid
		}

	}()

}
