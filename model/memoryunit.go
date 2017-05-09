package model

type memoryUnitInput struct {
	op      xuOperation
	a, b, c uint32
}

type memoryUnitOutput struct {
	writeRequest bool
	readData     uint32
}

func (s *Model) memoryUnit(
	input <-chan memoryUnitInput,
	we <-chan bool,
	output chan<- memoryUnitOutput) {

}
