package model

type xuSelector uint8

const (
	xuBypassSel xuSelector = 0x00
	xuAdderSel             = 0x01
	xuLogicSel             = 0x02
	xuShiftSel             = 0x03
	xuMemorySel            = 0x04
	xuBranchSel            = 0x05
)

type programAtom struct {
	valid bool
	unit  xuSelector
}

func (s *Model) programQueue(
	fifoIn <-chan programAtom,
	fifoOut chan<- programAtom,
	depth int) {

	fifo := make(chan programAtom, depth)

	go func() {
		defer close(fifo)

		for in := range fifoIn {
			fifo <- in
		}
	}()

	go func() {
		defer close(fifoOut)

		for out := range fifo {
			fifoOut <- out
		}
	}()
}
