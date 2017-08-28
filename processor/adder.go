package processor

// used by the dispatcher to send instructions to the adderUnit.
//
// holds two operands and the operation selector signal.
type adderInput struct {
	a, b uint32
	op   xuOperation
}

// Constructs a two stages adder execution unit.
// This unit is used by ADD(I), SUB(I), SLT(I), SLT(I)U instructions.
//
// In the real implementation adder logic may be distributed across the two stages,
// in this model the second stage is a just a dummy delay stage.
func (s *processor) adderUnit(
	input <-chan adderInput,
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
			case adderSub:
				buffer <- i.a - i.b
			case adderSlt:
				if int32(i.a) < int32(i.b) {
					buffer <- 1
				} else {
					buffer <- 0
				}
			case adderSltu:
				if uint32(i.a) < uint32(i.b) {
					buffer <- 1
				} else {
					buffer <- 0
				}
			case adderSum:
				buffer <- i.a + i.b
			}
		}
	}()
}
