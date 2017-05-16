package processor

import (
	"fmt"
	"log"
)

type xuSelector uint8

//go:generate stringer -type=xuSelector
const (
	xuBypassSel xuSelector = 0x00
	xuAdderSel             = 0x01
	xuLogicSel             = 0x02
	xuShiftSel             = 0x03
	xuMemorySel            = 0x04
	xuBranchSel            = 0x05
)

type programElement struct {
	valid uint8
	unit  xuSelector
}

func (p programElement) String() string {
	return fmt.Sprintf("{valid:%v unit:%v}", p.valid, p.unit)
}

func (s *Processor) prgQElement(
	fifoIn <-chan programElement,
	fifoOut chan<- programElement) {

	go func() {
		defer close(fifoOut)

		<-s.start
		fifoOut <- programElement{
			valid: 255,
			unit:  xuBypassSel}
		for in := range fifoIn {
			fifoOut <- in
		}
	}()
}

func (s *Processor) programQueue(
	fifoIn <-chan programElement,
	fifoOut chan<- programElement,
	depth int) {

	if depth < 2 {
		log.Panic("Program ordering queue depth must be at least 2")
	}

	fifo := make([]chan programElement, depth-1)

	for i := range fifo {
		fifo[i] = make(chan programElement)
	}

	s.prgQElement(fifoIn, fifo[0])

	for i := 1; i < depth-1; i++ {
		s.prgQElement(fifo[i-1], fifo[i])
	}

	s.prgQElement(fifo[depth-2], fifoOut)
}
