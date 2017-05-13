package processor

import (
	"bitbucket.org/marcos_sartori/qdi-riscv/memory"
	"log"
)

type Processor struct {
	start, quit chan struct{}
	Memory      memory.Memory
	regFile     regFile
	startPC     uint32
}

func encodeOneHot32(val ...uint) (ret uint32) {
	ret = 0
	for _, v := range val {
		ret |= 1 << (v & 0x1F)
	}

	return
}

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

func (s *Processor) Start() {
	if s.Memory != nil {
		close(s.start)
	} else {
		log.Panic("Processor has no memory attached")
	}
}

func (s *Processor) Stop() {
	close(s.quit)
}
