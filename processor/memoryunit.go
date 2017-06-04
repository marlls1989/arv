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
			case memorySH:
				data := make([]byte, 2)
				binary.LittleEndian.PutUint16(data, uint16(in.c&0xFFFF))
				wdata <- data
			case memorySW:
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
				out.value = uint32(int8(val[0]))
			case memoryLBU:
				val := <-rdata
				out.value = uint32(val[0])
			case memoryLH:
				val := <-rdata
				out.value = uint32(int16(binary.LittleEndian.Uint16(val)))
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
