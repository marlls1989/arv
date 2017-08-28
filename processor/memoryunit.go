package processor

import (
	"encoding/binary"
)

// operands and control information received from the dispatcher.
//
// a and b contains the base + offset pair used to calculate the memory address;
// c contains the data to be optionally written to memory write operations.
type memoryUnitInput struct {
	op      xuOperation
	a, b, c uint32
}

// Used by the memoryUnit to send results and communicates with the retire unit.
//
// writeRequest signals a pending memory write on the memory port waiting on the we channel
// and memoryRead signals if the contents of value should be written to the register bank.
type memoryUnitOutput struct {
	writeRequest bool
	memoryRead   bool
	value        uint32
}

// Constructs the memory access controller as a execution unit stage.
// It instantiate the first stage, a memory read-then-write port and the second stage.
//
// This unit is used for LW, LH(U), LB(U), SW, SH, SB.
//
// Data in memory is encoded as little-endian.
// The first stage computes the base address, encodes the data to be written
// and generate control signals to access the memory;
// the memory port then retrieves the data, if a write is issued it waits on the we channel
// for confirmation from the retire unit;
// the second stage decodes and sign extends the data read and sends to the retire unit,
// if the operation is a write the second stage signals the retire unit to control the we channel.
//
// The we channel is only used if a write is issued to the memory port,
// otherwise it is left idle and is ignored by the memory port.
func (s *processor) memoryUnit(
	input <-chan memoryUnitInput,
	we <-chan bool,
	output chan<- memoryUnitOutput) {

	addr := make(chan uint32)
	rdata := make(chan []byte)
	wdata := make(chan []byte)
	rlen := make(chan uint32)

	s.Memory.ReadWritePort(addr, rlen, wdata, we, rdata)

	operation := make(chan xuOperation)

	go func() {
		defer close(addr)
		defer close(rlen)
		defer close(wdata)
		defer close(operation)

		emptyData := make([]byte, 0)

		for in := range input {
			addr <- in.a + in.b
			operation <- in.op

			switch in.op {
			case memoryLB, memoryLBU:
				wdata <- emptyData
				rlen <- 1
			case memoryLH, memoryLHU:
				wdata <- emptyData
				rlen <- 2
			case memoryLW:
				wdata <- emptyData
				rlen <- 4
			case memorySB:
				data := make([]byte, 1)
				data[0] = byte(in.c & 0xFF)
				wdata <- data
				rlen <- 1
			case memorySH:
				data := make([]byte, 2)
				binary.LittleEndian.PutUint16(data, uint16(in.c&0xFFFF))
				wdata <- data
				rlen <- 2
			case memorySW:
				data := make([]byte, 4)
				binary.LittleEndian.PutUint32(data, in.c)
				wdata <- data
				rlen <- 4
			}
		}
	}()

	go func() {
		defer close(output)

		var out memoryUnitOutput

		for op := range operation {
			out.writeRequest = false
			out.memoryRead = true
			val, vval := <-rdata

			if !vval {
				return
			}

			switch op {
			case memorySB, memorySH, memorySW:
				out.writeRequest = true
				out.memoryRead = false
			case memoryLB:
				out.value = uint32(int8(val[0]))
			case memoryLBU:
				out.value = uint32(val[0])
			case memoryLH:
				out.value = uint32(int16(binary.LittleEndian.Uint16(val)))
			case memoryLHU:
				out.value = uint32(binary.LittleEndian.Uint16(val))
			case memoryLW:
				out.value = binary.LittleEndian.Uint32(val)
			}

			output <- out
		}
	}()
}
