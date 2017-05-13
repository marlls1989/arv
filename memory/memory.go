package memory

import (
	"launchpad.net/gommap"
	"os"
	"sync"
)

type Memory interface {
	ReadPort(addr, len <-chan uint32, data chan<- []byte)
	WritePort(addr <-chan uint32, data <-chan []byte, we <-chan bool)
}

type MemoryArray struct {
	mem gommap.MMap
	mux sync.Mutex
}

func MemoryArrayFromFile(file *os.File) (*MemoryArray, error) {
	ret := new(MemoryArray)

	memory, err :=
		gommap.Map(file.Fd(), gommap.PROT_READ|gommap.PROT_WRITE, gommap.MAP_SHARED)

	ret.mem = memory

	return ret, err
}

func (m *MemoryArray) ReadPort(addr, len <-chan uint32, data chan<- []byte) {
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

func (m *MemoryArray) WritePort(
	addr <-chan uint32,
	data <-chan []byte,
	we <-chan bool) {
	go func() {
		defer m.mem.Sync(gommap.MS_SYNC)

		for e := range we {
			a, da := <-addr
			d, dv := <-data
			if dv && da {
				if e {
					m.mux.Lock()
					for i, b := range d {
						m.mem[a+uint32(i)] = b
					}
					m.mux.Unlock()
				}
			} else {
				return
			}
		}
	}()
}
