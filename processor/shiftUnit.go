package processor

type shifterInput struct {
	op   xuOperation
	a, b uint32
}

// Constructs the shift execution unit logic stages
// This unit is used by SLL(I), SRL(I), SRA(I) instructions.
//
// The unit is two-stage long to match the the expected depth of execution units.
// Due to the simple logic, use of buffer may be carried to the silicon design.
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
				buffer <- uint32(int32(in.a) >> (in.b & 0x1F))
			}
		}
	}()
}
