package model

type branchCommandT struct {
	address uint32
	branch  bool
}

func (s *modelState) nextPcUnit(
	cmd <-chan branchCommandT,
	currPC <-chan uint32,
	currValid <-chan bool,
	fetchAddr, nextPc, fetchLen chan<- uint32,
	nextValid chan<- bool) {

	go func() {

		defer close(fetchAddr)
		defer close(nextPc)
		defer close(nextValid)
		defer close(fetchLen)

		<-s.start
		nextPc <- s.startPC
		fetchAddr <- s.startPC
		fetchLen <- 4
		for pc := range currPC {
			var target uint32
			c := <-cmd
			valid := <-currValid
			if c.branch {
				target = c.address
				valid = !valid
			} else {
				target = pc + 4
			}
			fetchAddr <- target
			nextPc <- target
			nextValid <- valid
			fetchLen <- 4
		}
	}()
}

func (s *modelState) pcUnit(
	nextPC <-chan uint32,
	nextValid <-chan bool,
	currPc, pcOut chan<- uint32,
	currValid, validOut chan<- bool) {

}

func (s *modelState) fetchUnit(cmd <-chan branchCommandT, fetchAddr chan<- uint32) {
}
