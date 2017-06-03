package processor

import (
	"log"
)

func (s *processor) bypassEl(input <-chan uint32, output chan<- uint32) {
	go func() {
		defer close(output)

		<-s.start
		output <- 0
		for i := range input {
			output <- i
		}
	}()
}

func (s *processor) bypassUnit(input <-chan uint32, output chan<- uint32, depth int) {
	if depth < 2 {
		log.Panic("bypassunit queue depth must be at least 2")
	}

	internal := make([]chan uint32, depth-1)

	for i := range internal {
		internal[i] = make(chan uint32)
	}

	s.bypassEl(input, internal[0])
	for i := 0; i < depth-2; i++ {
		s.bypassEl(internal[i], internal[i+1])
	}
	s.bypassEl(internal[depth-2], output)
}
