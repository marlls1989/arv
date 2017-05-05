package model

func (s *Model) dispatcherUnit(
	validIn <-chan bool,
	pcAddrIn <-chan uint32,
	xuOperIn <-chan xuOperation,
	operandA, operandB, operandC <-chan uint32,

	programQOut chan<- programElement,
	bypassOut chan<- uint32,
	adderOut chan<- adderInput,
	logicOut chan<- logicInput,
	branchOut chan<- branchInput,
) {

}
