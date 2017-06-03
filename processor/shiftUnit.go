package processor

type shifterInput struct {
	op   xuOperation
	a, b uint32
}

func (s *processor) shifterUnit(
	input <-chan shifterInput,
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

		for in := range input {
			switch in.op {
			case shifterLl:
				buffer <- in.a << (in.b & 0x1f)
			case shifterRl:
				buffer <- uint32(in.a) >> (in.b & 0x1F)
			case shifterRa:
				buffer <- uint32(int32(in.a) >> in.b & 0x1F)
			}
		}
	}()
}
