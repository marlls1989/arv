package model

import (
	"launchpad.net/gommap"
	"os"
	"sync"
)

type memory struct {
	mem gommap.MMap
	mux sync.Mutex
}

func initializeMemoryFromFile(file *os.File) (*memory, error) {
	ret := new(memory)

	memory, err :=
		gommap.Map(file.Fd(), gommap.PROT_READ|gommap.PROT_WRITE, gommap.MAP_SHARED)

	ret.mem = memory

	return ret, err
}

func (m *memory) ReadPort(addr, len <-chan uint32, data chan<- []byte) {
	go func() {
		defer close(data)
		for a := range addr {
			l, lv := <-len
			if lv {
				m.mux.Lock()
				d := m.mem[a : a+l]
				m.mux.Unlock()
				data <- d
			} else {
				return
			}
		}
	}()

}

func (m *memory) WritePort(addr <-chan uint32, data <-chan []byte) {
	go func() {
		defer m.mem.Sync(gommap.MS_ASYNC)

		for a := range addr {
			d, dv := <-data
			if dv {
				m.mux.Lock()
				for i, b := range d {
					m.mem[a+uint32(i)] = b
				}
				m.mux.Unlock()
			} else {
				return
			}
		}
	}()
}
