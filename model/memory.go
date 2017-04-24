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

func (s *modelState) memoryReadPort(addr, len <-chan uint32, data chan<- []byte) {
	go func() {
		defer close(data)

		for a := range addr {
			s.memory.mux.Lock()
			data <- s.memory.mem[a : a+<-len]
			s.memory.mux.Unlock()
		}
	}()

}

func (s *modelState) memoryWritePort(addr <-chan uint32, data <-chan []byte) {
	go func() {
		defer s.memory.mem.Sync(gommap.MS_ASYNC)

		for a := range addr {
			s.memory.mux.Lock()
			for i, b := range <-data {
				s.memory.mem[a+uint32(i)] = b
			}
			s.memory.mux.Unlock()
		}
	}()
}
