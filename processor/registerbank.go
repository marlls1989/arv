package processor

import (
	"sync"
)

type regFile struct {
	reg [32]uint32
	mux sync.Mutex
}

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

func (r *regFile) WritePort(addr <-chan regAddr, regWCmd <-chan regWCmd) {

	go func() {
		for wcmd := range regWCmd {
			a, av := <-addr
			if !av {
				return
			}

			d := wcmd.data
			we := wcmd.we

			ad := decodeOneHot32(uint32(a))
			if we && len(ad) != 0 && ad[0] != 0 {
				r.mux.Lock()
				r.reg[ad[0]-1] = d
				r.mux.Unlock()
			}
		}
	}()
}
