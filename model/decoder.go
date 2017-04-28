package model

import (
	"encoding/binary"
)

type xuSelector uint8
type opFormat uint8

const (
	xuBypassSel xuSelector = iota
	xuAdderSel
	xuLogicSel
	xuShiftSel
	xuMemorySel
)

const (
	opFormatR opFormat = iota
	opFormatI
	opFormatS
	opFormatSB
	opFormatU
	opFormatUJ
)

func (s *Model) decoderUnit(
	validIn <-chan bool,
	pcAddrIn <-chan uint32,
	instructionIn <-chan []byte,
	validOut chan<- bool,
	instructionOut, pcAddrOut chan<- uint32,
	xuSel chan<- xuSelector,
	opFmt chan<- opFormat) {

	s.pipeElement(validIn, validOut)
	s.pipeElement(pcAddrIn, pcAddrOut)

	go func() {
		defer close(instructionOut)
		defer close(xuSel)
		defer close(opFmt)

		for i := range instructionIn {
			ins := binary.LittleEndian.Uint32(i)
			instructionOut <- ins

			switch ins & 0x7F {
			case 0x1b:
				opFmt <- opFormatI
			case 0x3b:

			}
		}
	}()
}
