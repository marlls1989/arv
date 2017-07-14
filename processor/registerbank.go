package processor

import (
	"sync"
)

// Represents a instance of the simulated register file model.
type regFile struct {
	reg [31]uint32
	mux sync.Mutex
}

// This function constructs a register read port logical stage associated to a register file model.
//
// The constructed logic stage reads the address channel;
// if the address points register zero, the port writes the value zero on the data channel;
// otherwise it reads the content of register pointed by address and write to the data channel.
func (r *regFile) ReadPort(
	addr <-chan regAddr,
	data chan<- uint32) {

	go func() {
		defer close(data)

		for a := range addr {
			ad := decodeOneHot32(uint32(a))
			if len(ad) != 0 {
				if ad[0] == 0 {
					data <- 0
				} else {
					r.mux.Lock()
					d := r.reg[ad[0]-1]
					r.mux.Unlock()
					data <- d
				}
			}
		}
	}()
}

// This function constructs a register write port logical stage associated to a register file model.
//
// The Constructed logic stage first reads the data channel, then the address channel;
// if the address points to register zero, the data is discarded;
// otherwise it is written to the register pointed by address before reading the next data-address pair.
//
// It is important to notice that the if the data and address channels are written by the same logic stage
// they must be written in the order they are read, i.e. data first, address last, otherwise a deadlock occurs.
// This limitation is due to the way go operates on channels, real hardware or a more precise model would read both channels in parallel.
func (r *regFile) WritePort(addr <-chan regAddr, data <-chan uint32) {

	go func() {
		for d := range data {
			a, av := <-addr
			if !av {
				return
			}
			ad := decodeOneHot32(uint32(a))
			if len(ad) != 0 && ad[0] != 0 {
				r.mux.Lock()
				r.reg[ad[0]-1] = d
				r.mux.Unlock()
			}
		}
	}()
}
