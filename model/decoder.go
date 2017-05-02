package model

import (
	"encoding/binary"
)

type xuSelector uint8
type regAddr uint32
type opFormat uint8
type xuOperation uint8

const (
	bypassA   xuOperation = 0x00
	adderSum              = 0x10
	adderSub              = 0x11
	adderSlt              = 0x12
	adderSltu             = 0x13
	shifterLl             = 0x20
	shifterRl             = 0x21
	shifterRa             = 0x22
	logicXor              = 0x30
	logicOr               = 0x31
	logicAnd              = 0x32
	memoryLB              = 0x40
	memoryLH              = 0x41
	memoryLW              = 0x42
	memoryLBU             = 0x44
	memoryLHU             = 0x45
	memorySB              = 0x48
	memorySH              = 0x49
	memorySW              = 0x4A
	branchEQ              = 0x50
	branchNE              = 0x51
	branchLT              = 0x54
	branchGE              = 0x55
	branchLTU             = 0x56
	branchGEU             = 0x57
	branchJL              = 0x58
)

const (
	xuBypassSel xuSelector = 0x00
	xuAdderSel             = 0x01
	xuLogicSel             = 0x02
	xuShiftSel             = 0x03
	xuMemorySel            = 0x04
	xuBranchSel            = 0x05
)

const (
	opFormatR opFormat = iota
	opFormatI
	opFormatS
	opFormatB
	opFormatU
	opFormatJ
)

func (s *Model) decoderUnit(
	validIn <-chan bool,
	pcAddrIn <-chan uint32,
	instructionIn <-chan []byte,
	validOut chan<- bool,
	instructionOut, pcAddrOut chan<- uint32,
	xuOper chan<- xuOperation,
	opFmt chan<- opFormat,
	regAaddr chan<- uint32,
	rebBaddr chan<- uint32) {

	go func() {
		defer close(validOut)

		<-s.start
		validOut <- false
		for v := range validIn {
			validOut <- v
		}
	}()

	go func() {
		defer close(pcAddrOut)

		<-s.start
		pcAddrOut <- 0
		for p := range pcAddrIn {
			pcAddrOut <- p
		}
	}()

	go func() {
		defer close(instructionOut)
		defer close(opFmt)
		defer close(xuOper)

		<-s.start
		instructionOut <- 0x13
		opFmt <- opFormatI
		for i := range instructionIn {
			ins := binary.LittleEndian.Uint32(i)
			fmt := opFormatI
			op := bypassA

			switch ins & 0x7F {
			case 0x13, 0x67, 0x03:
				fmt = opFormatI
			case 0x6f:
				fmt = opFormatJ
			case 0x33:
				fmt = opFormatR
			case 0x63:
				fmt = opFormatB
			case 0x37, 0x17:
				fmt = opFormatU
			case 0x23:
				fmt = opFormatS
			}

			switch ins & 0x7F {
			case 0x37: //LUI
				op = bypassA
			case 0x1F: //AUIPC
				op = adderSum
			case 0x6F: //JAL
				op = branchJL
			default:
				switch ins & 0x707F {
				case 0x67: //JALR
					op = branchJL
				case 0x63: //BEQ
					op = branchEQ
				case 0x1063: //BNE
					op = branchNE
				case 0x4063: //BLT
					op = branchLT
				case 0x5063: //BGE
					op = branchGE
				case 0x6063: //BLTU
					op = branchLTU
				case 0x7063: //BGEU
					op = branchGEU
				case 0x03: //LB
					op = memoryLB
				case 0x1003: //LH
					op = memoryLH
				case 0x2003: //LW
					op = memoryLW
				case 0x4003: //LBU
					op = memoryLBU
				case 0x5003: //LHU
					op = memoryLHU
				case 0x43: //SB
					op = memorySB
				case 0x1043: //SH
					op = memorySH
				case 0x2043: //SW
					op = memorySW
				case 0x13: //ADDI
					op = adderSum
				case 0x2013: //SLTI
					op = adderSlt
				case 0x3013: //SLTIU
					op = adderSltu
				case 0x4013: //XORI
					op = logicXor
				case 0x6013: //ORI
					op = logicOr
				case 0x7013: //ANDI
					op = logicAnd
				default:
					switch ins & 0xFE00707F {
					case 0x1033, 0x1013: //SLL, SLLI
						op = shifterLl
					case 0x5033, 0x5013: //SRL, SRLI
						op = shifterRl
					case 0x40005033, 0x40005013: //SRA, SRAI
						op = shifterRa
					case 0x33: //ADD
						op = adderSum
					case 0x40000033: //SUB
						op = adderSub
					case 0x2033: //SLT
						op = adderSlt
					case 0x3033: //SLTU
						op = adderSltu
					case 0x4033: //XOR
						op = logicXor
					case 0x6033: //OR
						op = logicOr
					case 0x7033: //AND
						op = logicAnd
					}
				}
			}

			opFmt <- fmt
			xuOper <- op
			instructionOut <- ins
		}
	}()
}
