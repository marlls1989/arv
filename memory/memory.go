package memory

import (
	"fmt"
	"launchpad.net/gommap"
	"log"
	"os"
	"sync"
)

type Memory interface {
	ReadWritePort(
		addr, lng <-chan uint32,
		dataIn <-chan []byte,
		we <-chan bool,
		dataOut chan<- []byte)
	ReadPort(addr, len <-chan uint32, data chan<- []byte)
}

type memoryArray struct {
	mem           gommap.MMap
	mux           sync.Mutex
	EndSimulation chan struct{}
	Debug         bool
}

func MemoryArrayFromFile(file *os.File) (*memoryArray, error) {
	ret := new(memoryArray)

	memory, err :=
		gommap.Map(file.Fd(), gommap.PROT_READ|gommap.PROT_WRITE, gommap.MAP_SHARED)

	ret.mem = memory
	ret.EndSimulation = make(chan struct{})

	if err != nil {
		//Force initialization
		memory.Sync(gommap.MS_SYNC)
	}

	return ret, err
}

func (m *memoryArray) ReadPort(addr, lng <-chan uint32, data chan<- []byte) {
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
					if m.Debug {
						log.Printf("Reading %d bytes from out of bounds memory location %x", l, a)
					}
					d = make([]byte, l)
				}
				data <- d
			} else {
				return
			}
		}
	}()
}

// Constructs a read-then-write memory port
func (m *memoryArray) ReadWritePort(
	addr, lng <-chan uint32,
	dataIn <-chan []byte,
	we <-chan bool,

	dataOut chan<- []byte) {
	go func() {
		defer m.mem.Sync(gommap.MS_SYNC)
		defer close(dataOut)
		for a := range addr {
			di, dv := <-dataIn
			l, lv := <-lng

			if !(dv || lv) {
				return
			}

			// Memory Read portion
			var do []byte
			if a+l < uint32(len(m.mem)) {
				m.mux.Lock()
				do = m.mem[a : a+l]
				m.mux.Unlock()
				if m.Debug {
					log.Printf("Reading %d from memory address %X", do, a)
				}
			} else {
				if m.Debug {
					log.Printf("Reading %d bytes from out of bounds memory location %x", l, a)
				}
				do = make([]byte, l)
			}
			dataOut <- do

			// Memory Write portion
			if len(di) > 0 {
				e := <-we
				if e {
					if m.Debug {
						log.Printf("Writing %v to memory address %X", di, a)
					}
					if a < 0x80000000 {
						m.mux.Lock()
						for i, b := range di {
							m.mem[a+uint32(i)] = b
						}
						m.mux.Unlock()
					} else if a < 0x80001000 {
						if m.Debug {
							log.Print("Simulation End invoked")
						}
						close(m.EndSimulation)
					} else {
						for _, c := range di {
							fmt.Printf("%c", c)
						}
					}
				}
			}
		}
	}()
}
