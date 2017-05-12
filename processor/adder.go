package processor

type adderInput struct {
	a, b uint32
	op   xuOperation
}

func (s *Processor) adderUnit(
	input <-chan adderInput,
	output chan<- uint32) {

	go func() {
		defer close(output)
		for i := range input {
			switch i.op {
			case adderSub:
				output <- i.a - i.b
			case adderSlt:
				if int32(i.a) < int32(i.b) {
					output <- 1
				} else {
					output <- 0
				}
			case adderSltu:
				if uint32(i.a) < uint32(i.b) {
					output <- 1
				} else {
					output <- 0
				}
			case adderSum:
				output <- i.a + i.b
			}
		}
	}()
}
