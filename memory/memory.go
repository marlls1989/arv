package memory

import (
	"fmt"
	"launchpad.net/gommap"
	"log"
	"os"
	"sync"
)

type Memory interface {
	ReadPort(addr, len <-chan uint32, data chan<- []byte)
	WritePort(addr <-chan uint32, data <-chan []byte, we <-chan bool)
}

type MemoryArray struct {
	mem           gommap.MMap
	mux           sync.Mutex
	EndSimulation chan struct{}
}

func MemoryArrayFromFile(file *os.File) (*MemoryArray, error) {
	ret := new(MemoryArray)

	memory, err :=
		gommap.Map(file.Fd(), gommap.PROT_READ|gommap.PROT_WRITE, gommap.MAP_SHARED)

	ret.mem = memory
	ret.EndSimulation = make(chan struct{})

	return ret, err
}

func (m *MemoryArray) ReadPort(addr, lng <-chan uint32, data chan<- []byte) {
	go func() {
		defer close(data)
		for a := range addr {
			var d []byte
			l, lv := <-lng
			if lv {
				if a+l < uint32(len(m.mem)) {
					m.mux.Lock()
					d = m.mem[a : a+l]
					m.mux.Unlock()
				} else {
					d = make([]byte, l)
				}
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
					log.Printf("Writing %v to addr 0x%X", d, a)
					if a < 0x80000000 {
						m.mux.Lock()
						for i, b := range d {
							m.mem[a+uint32(i)] = b
						}
						m.mux.Unlock()
					} else if a < 0x80001000 {
						close(m.EndSimulation)
					} else {
						for _, c := range d {
							fmt.Printf("%c", c)
						}
					}
				}
			} else {
				return
			}
		}
	}()
}
