package model

type xuSelector uint8
type opSelector uint8

const (
	xuBypassSel xuSelector = iota
	xuAdderSel
	xuLogicSel
	xuShiftSel
	xuMemorySel
)

const (
	opRSelector opSelector = iota
	opISelector
	opSSelector
	opSBSelector
	opUSelector
	opUJSelector
)

func (s *Model) decoderUnit(
	validIn <-chan bool,
	pcAddrIn <-chan uint32,
	instructionIn <-chan []byte,
	validOut chan<- bool,
	instructionOut, pcAddrOut chan<- uint32,
	xuSel chan<- xuSelector) {

	go func() {
	}()
}
