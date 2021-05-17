package memory

import (
  "fmt"
  "github.com/edsrzf/mmap-go"
  "log"
  "os"
  "sync"
  "time"
  "unsafe"
)

// Memory model instance created using MemoryArrayFromFile function.
//
// The file is mapped to the lower portion of the memory model starting at 0.
// EndSimulation is a channel closed when data is written to address range 0x80000000-0x80001000,
// it can be used to sign end of simulation.
// Any byte written to memory location above 0x80001000 are printed to stdout
// Reads and writes to memory locations not backed by file bellow the special upper regions are ignored.
//
// Invalid memory access warnings are printed to stderr if the Debug flag is enabled.
type memoryArray struct {
  mem           mmap.MMap
  mux           sync.Mutex
  EndSimulation chan struct{}
  Debug         bool
}

// Construct a memoryArray memory model by mmapping the contents of file received as argument.
//
// The memoryArray instance returned implements the Memory interface presenting the following memory map:
// the file is mapped to the lower portion of the memory model starting at 0;
// memoryArray.EndSimulation is a struct {} channel closed when data is written to address range 0x80000000-0x80001000,
// it can be used to sign end of simulation;
// any byte written to memory location above 0x80001000 are printed to stdout;
// reads and writes to memory locations not backed by file bellow the special upper regions are ignored.

func MemoryArrayFromFile(file *os.File) (*memoryArray, error) {
  ret := new(memoryArray)

  memory, err :=
    mmap.Map(file, mmap.RDWR, 0)

  ret.mem = memory
  ret.EndSimulation = make(chan struct{})

  return ret, err
}

// Constructs a read-only memory port logic stage as described in the Memory interface.
//
// Memory locations corresponding to the valid range in the file are read normally,
// non mapped memory and special memory regions returns all zero data.
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

// Constructs a read-then-write memory port logic stage as described in the Memory interface.
//
// Memory locations backed by file are read normally,
// non mapped memory and special memory regions returns all zero data.
//
// Memory regions backed by the file are written normally,
// writes to special memory regions trigger actions as documented in memoryArray struct,
// writes to memory locations bellow 0x80000000 not backed by file are discarded.
func (m *memoryArray) ReadWritePort(
  addr, lng <-chan uint32,
  dataIn <-chan []byte,
  we <-chan bool,

  dataOut chan<- []byte) {
  startTime := time.Now()
  go func() {
    defer m.mem.Flush()
    defer close(dataOut)
    for a := range addr {
      di, dv := <-dataIn
      l, lv := <-lng

      if !(dv || lv) {
        return
      }

      // Memory Read portion
      var do []byte
      if a >= 0x80006000 && a < 0x80006008 {
        x := a % 8
        elapsed := time.Now().Sub(startTime).Milliseconds()
        do = ((*[8]byte)(unsafe.Pointer(&elapsed)))[x:x+l]
      } else if a+l < uint32(len(m.mem)) {
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
