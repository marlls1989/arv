package processor

import (
	"encoding/binary"
	"fmt"
	"log"
	"sync/atomic"
)

type opFormat uint8
type xuOperation uint8

//go:generate stringer -type=xuOperation
const (
	bypassB xuOperation = xuOperation(xuBypassSel<<4) + iota
)

const (
	adderSum xuOperation = xuOperation(xuAdderSel<<4) + iota
	adderSub
	adderSlt
	adderSltu
)

const (
	logicXor xuOperation = xuOperation(xuLogicSel<<4) + iota
	logicOr
	logicAnd
)

const (
	shifterLl xuOperation = xuOperation(xuShiftSel<<4) + iota
	shifterRl
	shifterRa
)

const (
	memoryLB xuOperation = xuOperation(xuMemorySel<<4) + iota
	memoryLH
	memoryLW
	memoryLBU
	memoryLHU
	memorySB
	memorySH
	memorySW
)

const (
	branchEQ xuOperation = xuOperation(xuBranchSel<<4) + iota
	branchNE
	branchLT
	branchGE
	branchLTU
	branchGEU
	branchJL
)

//go:generate stringer -type=opFormat
const (
	opFormatR opFormat = iota
	opFormatI
	opFormatS
	opFormatB
	opFormatU
	opFormatJ
	opFormatNop
)

type regAddr uint32

//go:generate stringer -type=regAddr
const (
	zero regAddr = 1 << iota
	ra
	sp
	gp
	tp
	t0
	t1
	t2
	s0
	s1
	a0
	a1
	a2
	a3
	a4
	a5
	a6
	a7
	s2
	s3
	s4
	s5
	s6
	s7
	s8
	s9
	s10
	s11
	t3
	t4
	t5
	t6
)

type decoderOut struct {
	pc               uint32
	valid            uint8
	ins              uint32
	regA, regB, regD regAddr
	op               xuOperation
	fmt              opFormat
}

func (d decoderOut) String() string {
	return fmt.Sprintf("{ins:%08X regA:%v regB:%v regD:%v op:%v fmt:%v}", d.ins, d.regA, d.regB, d.regD, d.op, d.fmt)
}

func (s *processor) decoderUnit(
	validIn <-chan uint8,
	pcAddrIn <-chan uint32,
	instructionIn <-chan []byte,

	output chan<- decoderOut) {

	go func() {
		defer close(output)

		for i := range instructionIn {

			select {
			case _ = <-s.quit:
				return
			default:
			}

			atomic.AddUint64((&s.Decoded), 1)
			pc, pe := <-pcAddrIn
			valid, ve := <-validIn

			if !(pe || ve) {
				return
			}

			var fmt opFormat
			var op xuOperation

			ins := binary.LittleEndian.Uint32(i)

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
			case 0xFF:
				fmt = opFormatNop
			default:
				if s.Debug {
					log.Printf("Decoding unknown instruction at %X as NOP", pc)
				}
				fmt = opFormatNop
			}

			switch ins & 0x7F {
			case 0x0F: //FENCE (decoded as NOP)
				op = bypassB
			case 0x37: //LUI
				op = bypassB
			case 0x17: //AUIPC
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
				case 0x23: //SB
					op = memorySB
				case 0x1023: //SH
					op = memorySH
				case 0x2023: //SW
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
					default:
						op = bypassB
					}
				}
			}

			ra := regAddr(encodeOneHot32((uint)((ins >> 15) & 0x1F)))
			rb := regAddr(encodeOneHot32((uint)((ins >> 20) & 0x1F)))
			rd := regAddr(encodeOneHot32((uint)((ins >> 07) & 0x1F)))

			output <- decoderOut{
				pc:    pc,
				valid: valid,
				ins:   ins,
				fmt:   fmt,
				op:    op,
				regA:  ra,
				regB:  rb,
				regD:  rd}
		}
	}()
}
