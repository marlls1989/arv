package processor

type logicInput struct {
	a, b uint32
	op   xuOperation
}

func (s *Processor) logicUnit(
	input <-chan logicInput,
	output chan<- uint32) {

	go func() {
		defer close(output)
		for i := range input {
			switch i.op {
			case logicXor:
				output <- i.a ^ i.b
			case logicOr:
				output <- i.a | i.b
			case logicAnd:
				output <- i.a & i.b
			}
		}
	}()
}
