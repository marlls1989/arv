package model

func (s *Model) nextPcUnit(
	currPC, branchAddr <-chan uint32,
	currValid, branch <-chan bool,
	fetchAddr, nextPc chan<- uint32,
	nextValid chan<- bool) {

	go func() {

		defer close(fetchAddr)
		defer close(nextPc)
		defer close(nextValid)

		<-s.start
		nextPc <- s.startPC
		fetchAddr <- s.startPC
		for br := range branch {
			var target uint32
			valid := <-currValid
			pc := <-currPC
			if br {
				brt, vbrt := <-branchAddr
				if vbrt {
					target = brt
					valid = !valid
				} else {
					return
				}
			} else {
				target = pc + 4
			}
			fetchAddr <- target
			nextPc <- target
			nextValid <- valid
		}
	}()
}

func (s *Model) fetchUnit(
	branch <-chan bool,
	branchAddr <-chan uint32,
	pcAddr chan<- uint32,
	instruction chan<- []byte,
	valid chan<- bool) {

	nextPC := make(chan uint32)
	currPC := make(chan uint32)
	fetchAddr := make(chan uint32)
	nextValid := make(chan bool)
	currValid := make(chan bool)
	len := make(chan uint32)

	s.nextPcUnit(currPC, branchAddr, currValid, branch, fetchAddr, nextPC, nextValid)

	s.pipeElement(nextPC, currPC, pcAddr)
	s.pipeElement(nextValid, currValid, valid)
	s.pipeElement(uint32(4), len)

	s.memory.ReadPort(fetchAddr, len, instruction)
}
