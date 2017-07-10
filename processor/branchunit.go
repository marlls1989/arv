package processor

type branchInput struct {
	a, b, c, pc uint32
	op          xuOperation
}

// Used by the branch unit to signal the retire unit intention to branch
//
// Link signals that the contents of linkAddr should be written to the destination register;
// taken signals if the branch instruction overwrites the program counter with the content of target.
type branchOutput struct {
	link, taken      bool
	linkAddr, target uint32
}

// This function constructs the branch unit two logical stages.
//
// In the real implementation adder logic may be distributed across the two stages,
// in this model the second stage is a just a dummy delay stage.
//
//
func (s *processor) branchUnit(
	input <-chan branchInput,
	output chan<- branchOutput) {

	buffer := make(chan branchOutput)

	go func() {
		defer close(output)
		for i := range buffer {
			output <- i
		}
	}()

	go func() {
		defer close(buffer)

		for in := range input {
			out := branchOutput{
				taken:    false,
				target:   in.pc + in.c,
				linkAddr: in.pc + 4,
				link:     false}

			switch in.op {
			case branchEQ:
				out.taken = (in.a == in.b)
			case branchNE:
				out.taken = (in.a != in.b)
			case branchLT:
				out.taken = (int32(in.a) < int32(in.b))
			case branchGE:
				out.taken = (int32(in.a) >= int32(in.b))
			case branchLTU:
				out.taken = (in.a < in.b)
			case branchGEU:
				out.taken = (in.a >= in.b)
			case branchJL:
				out.taken = true
				out.target = in.a + in.b
				out.link = true
			}

			buffer <- out
		}
	}()
}
