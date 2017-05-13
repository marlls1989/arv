package processor

type logicInput struct {
	a, b uint32
	op   xuOperation
}

func (s *Processor) logicUnit(
	input <-chan logicInput,
	output chan<- uint32) {

	buffer := make(chan uint32)

	go func() {
		defer close(output)
		for i := range buffer {
			output <- i
		}
	}()

	go func() {
		defer close(buffer)
		for i := range input {
			switch i.op {
			case logicXor:
				buffer <- i.a ^ i.b
			case logicOr:
				buffer <- i.a | i.b
			case logicAnd:
				buffer <- i.a & i.b
			}
		}
	}()
}
