// This package contains the high-level functional model of the ARV core.
package processor

import (
	"github.com/marlls1989/arv/memory"
	"log"
)

// Processor model instace created using ConstructProcessor function.
//
// It is used by logical stage constructors to refer to the appropriated register bank and memory model.
// It is also used to keep execution performance statistics.
//
// Debug flag enabled pipeline debugging verbosity;
// and StartPC points to the first instruction to be executed,
// should be set before calling the Start function.
//
// Cancelled, Bubbles, Retired and Decoded are respectively counters of instructions and bubbles
// cancelled due to branching; bubbles inserted due to pipeline hazards; instructions and bubbles
// successfully retired and instructions fetched and decoded.

type Stats struct {
	Cancelled uint64
	Bubbles   uint64
	Retired   uint64
	Decoded   uint64
	Unit      struct {
		Bypass  uint64
		Adder   uint64
		Logic   uint64
		Shifter uint64
		Memory  uint64
		Branch  uint64
	}
}

type processor struct {
	start, quit chan struct{}
	Memory      memory.Memory
	regFile     regFile
	StartPC     uint32
	Debug       bool
	Stats       Stats
}

// Receives one or more values between 0 and 31 and encodes into a 32-bit
// bitfield by setting the respective numbered bits to 1
func encodeOneHot32(val ...uint) (ret uint32) {
	ret = 0
	for _, v := range val {
		ret |= 1 << (v & 0x1F)
	}

	return
}

// Receives a 32-bit bitfield decomposing into a ordered vector of uint
// containing the number of bits set as 1.
func decodeOneHot32(val uint32) (ret []uint) {
	var i uint
	for i = 0; i < 32; i++ {
		if (val & 1) != 0 {
			ret = append(ret, i)
		}

		val >>= 1
	}

	return
}

// Starts the processor model execution.
func (s *processor) Start() {
	if s.Memory != nil {
		close(s.start)
	} else {
		log.Panic("Processor has no memory attached")
	}
}

// Terminates the processor model execution.
// BUG(msartori): Untested
func (s *processor) Stop() {
	close(s.quit)
}
