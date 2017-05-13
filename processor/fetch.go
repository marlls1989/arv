package processor

type branchCmd struct {
	taken  bool
	target uint32
}

func (s *Processor) nextPcUnit(
	currPC <-chan uint32,
	currValid <-chan bool,
	branch <-chan branchCmd,
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
			if br.taken {
				target = br.target
				valid = !valid
			} else {
				target = pc + 4
			}
			fetchAddr <- target
			nextPc <- target
			nextValid <- valid
		}
	}()
}

func (s *Processor) fetchUnit(
	branch <-chan branchCmd,

	pcAddr chan<- uint32,
	instruction chan<- []byte,
	valid chan<- bool) {

	nextPC := make(chan uint32)
	currPC := make(chan uint32)
	fetchAddr := make(chan uint32)
	nextValid := make(chan bool)
	currValid := make(chan bool)
	len := make(chan uint32)

	s.nextPcUnit(currPC, currValid, branch, fetchAddr, nextPC, nextValid)

	go func() {
		defer close(currPC)
		defer close(pcAddr)

		<-s.start
		currPC <- 0
		pcAddr <- 0
		for i := range nextPC {
			currPC <- i
			pcAddr <- i
		}
	}()

	go func() {
		defer close(currValid)
		defer close(valid)

		<-s.start
		currValid <- false
		valid <- false
		for i := range nextValid {
			currValid <- i
			valid <- i
		}
	}()

	go func() {
		defer close(len)
		for {
			select {
			case len <- 4:
			case <-s.quit:
				return
			}
		}
	}()

	s.Memory.ReadPort(fetchAddr, len, instruction)
}
