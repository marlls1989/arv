package model

import (
	"log"
)

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

func (s *Model) prgQElement(
	fifoIn <-chan programAtom,
	fifoOut chan<- programAtom) {

	go func() {
		defer close(fifoOut)

		<-s.start
		fifoOut <- programAtom{valid: false, unit: xuBypassSel}
		for in := range fifoIn {
			fifoOut <- in
		}
	}()
}

func (s *Model) programQueue(
	fifoIn <-chan programAtom,
	fifoOut chan<- programAtom,
	depth int) {

	if depth < 2 {
		log.Panic("Program ordering queue depth must be at least 2")
	}

	fifo := make([]chan programAtom, depth-1)

	for i := range fifo {
		fifo[i] = make(chan programAtom)
	}

	s.prgQElement(fifoIn, fifo[0])

	for i := 1; i < depth-1; i++ {
		s.prgQElement(fifo[i-1], fifo[i])
	}

	s.prgQElement(fifo[depth-2], fifoOut)
}
