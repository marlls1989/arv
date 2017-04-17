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

func InitializeMemoryFromFile(file *os.File) (*memory, error) {
	ret := new(memory)

	memory, err :=
		gommap.Map(file.Fd(), gommap.PROT_READ|gommap.PROT_WRITE, gommap.MAP_SHARED)

	ret.mem = memory

	return ret, err
}

func (s *modelState) memoryReadPort(addr chan uint32, data chan []byte, len uint32) {
	go func() {
		for {
			select {
			case a := <-addr:
				s.memory.mux.Lock()
				data <- s.memory.mem[a : a+len]
				s.memory.mux.Unlock()
			case <-s.quit:
				return
			}
		}
	}()
}

func (s *modelState) memoryWritePort(addr chan uint32, data chan []byte) {
	go func() {
		defer s.memory.mem.Sync(gommap.MS_ASYNC)
		for {
			select {
			case a := <-addr:
				s.memory.mux.Lock()
				d := <-data
				for i, b := range d {
					s.memory.mem[a+uint32(i)] = b
				}
				s.memory.mux.Unlock()
			case <-s.quit:
				return
			}
		}
	}()
}
