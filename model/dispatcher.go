package model

func (s *Model) dispatcherUnit(
	validIn <-chan bool,
	instructionIn, pcAddrIn <-chan uint32,
	xuOper <-chan xuOperation,
	opFmt <-chan opFormat,
	regLock <-chan uint32,
	regA, regB, regD <-chan uint32) {

}
