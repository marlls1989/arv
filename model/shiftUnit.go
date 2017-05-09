package model

type shifterInput struct {
	op   xuOperation
	a, b uint32
}

func (s *Model) shifterUnit(
	input <-chan shifterInput,
	output chan<- uint32) {

	go func() {
		defer close(output)

		for in := range input {
			switch in.op {
			case shifterLl:
				output <- in.a << (in.b & 0x1f)
			case shifterRl:
				output <- uint32(in.a) >> (in.b & 0x1F)
			case shifterRa:
				output <- uint32(int32(in.a) >> in.b & 0x1F)
			}
		}
	}()
}
