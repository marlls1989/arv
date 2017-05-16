package processor

import (
	"encoding/binary"
)

type memoryUnitInput struct {
	op      xuOperation
	a, b, c uint32
}

type memoryUnitOutput struct {
	writeRequest bool
	value        uint32
}

func (s *Processor) memoryUnit(
	input <-chan memoryUnitInput,
	we <-chan bool,
	output chan<- memoryUnitOutput) {

	raddr := make(chan uint32)
	rdata := make(chan []byte)
	rlen := make(chan uint32)

	s.Memory.ReadPort(raddr, rlen, rdata)

	waddr := make(chan uint32)
	wdata := make(chan []byte)

	s.Memory.WritePort(waddr, wdata, we)

	operation := make(chan xuOperation)

	go func() {
		defer close(raddr)
		defer close(rlen)
		defer close(waddr)
		defer close(wdata)
		defer close(operation)

		var addr uint32

		for in := range input {
			addr = in.a + in.b

			operation <- in.op

			switch in.op {
			case memoryLB, memoryLBU:
				raddr <- addr
				rlen <- 1
			case memoryLH, memoryLHU:
				raddr <- addr
				rlen <- 2
			case memoryLW:
				raddr <- addr
				rlen <- 4
			case memorySB:
				waddr <- addr
				data := make([]byte, 1)
				data[0] = byte(in.c & 0xFF)
				wdata <- data
			case memorySH:
				waddr <- addr
				data := make([]byte, 2)
				binary.LittleEndian.PutUint16(data, uint16(in.c&0xFFFF))
				wdata <- data
			case memorySW:
				waddr <- addr
				data := make([]byte, 4)
				binary.LittleEndian.PutUint32(data, in.c)
				wdata <- data
			}
		}
	}()

	go func() {
		defer close(output)

		var out memoryUnitOutput

		for op := range operation {
			out.writeRequest = false
			out.value = 0
			switch op {
			case memorySB, memorySH, memorySW:
				out.writeRequest = true
			case memoryLB:
				val := <-rdata
				out.value = uint32(int32(val[0]))
			case memoryLBU:
				val := <-rdata
				out.value = uint32(val[0])
			case memoryLH:
				val := <-rdata
				out.value = uint32(int32(binary.LittleEndian.Uint16(val)))
			case memoryLHU:
				val := <-rdata
				out.value = uint32(binary.LittleEndian.Uint16(val))
			case memoryLW:
				val := <-rdata
				out.value = binary.LittleEndian.Uint32(val)
			}

			output <- out
		}
	}()
}
