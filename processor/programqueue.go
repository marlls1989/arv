package processor

import (
	"fmt"
	"log"
)

type xuSelector uint8

//go:generate stringer -type=xuSelector
const (
	xuBypassSel xuSelector = iota
	xuAdderSel
	xuLogicSel
	xuShiftSel
	xuMemorySel
	xuBranchSel
)

// Holds the stream tag and execution unit selector
// by the program queue to carry control information between
// the dispatcher and retire unit
type programElement struct {
	valid uint8
	unit  xuSelector
}

func (p programElement) String() string {
	return fmt.Sprintf("{valid:%v unit:%v}", p.valid, p.unit)
}

// construct a program queue stage
//
// the stage is initiated to a bubble
func (s *processor) prgQElement(
	fifoIn <-chan programElement,
	fifoOut chan<- programElement) {

	go func() {
		defer close(fifoOut)

		<-s.start
		fifoOut <- programElement{
			valid: 255,
			unit:  xuBypassSel}

		for in := range fifoIn {

			select {
			case <-s.quit:
				return
			default:
			}

			fifoOut <- in
		}
	}()
}

// Construct parameterisabled length program queue.
func (s *processor) programQueue(
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
