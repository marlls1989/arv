package model

func (s *Model) retireUnit(
	qIn <-chan programElement,
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

	go func() {

	}()

}
