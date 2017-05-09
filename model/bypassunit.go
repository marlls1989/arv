package model

import (
	"log"
)

func (s *Model) bypassUnit(input <-chan uint32, output chan<- uint32, depth int) {
	if depth < 2 {
		log.Panic("bypassunit queue depth must be at least 2")
	}

	internal := make([]chan uint32, depth-1)

	for i := range internal {
		internal[i] = make(chan uint32)
	}

	s.pipeElementWithInitization(input, uint32(0), internal[0])
	for i := 0; i < depth-2; i++ {
		s.pipeElementWithInitization(internal[i], uint32(0), internal[i+1])
	}
	s.pipeElementWithInitization(internal[depth-2], uint32(0), output)
}
