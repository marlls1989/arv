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

func (s *Model) memoryReadPort(addr, len <-chan uint32, data chan<- []byte) {
	go func() {
		defer close(data)
		for a := range addr {
			l, lv := <-len
			if lv {
				s.memory.mux.Lock()
				d := s.memory.mem[a : a+l]
				s.memory.mux.Unlock()
				data <- d
			} else {
				return
			}
		}
	}()

}

func (s *Model) memoryWritePort(addr <-chan uint32, data <-chan []byte) {
	go func() {
		defer s.memory.mem.Sync(gommap.MS_ASYNC)

		for a := range addr {
			d, dv := <-data
			if dv {
				s.memory.mux.Lock()
				for i, b := range d {
					s.memory.mem[a+uint32(i)] = b
				}
				s.memory.mux.Unlock()
			} else {
				return
			}
		}
	}()
}
