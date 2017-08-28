package processor

type logicInput struct {
	a, b uint32
	op   xuOperation
}

// Constructs the logic execution unit logic stages
// This unit is used by XOR(I), OR(I), AND(I) instructions.
//
// The unit is two-stage long to match the the expected depth of execution units.
// Due to the simple logic, use of buffer may be carried to the silicon design.
func (s *processor) logicUnit(
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
