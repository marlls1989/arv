package processor

type branchInput struct {
	a, b, c, pc uint32
	op          xuOperation
}

type branchOutput struct {
	link, taken      bool
	linkAddr, target uint32
}

func (s *Processor) branchUnit(
	input <-chan branchInput,
	output chan<- branchOutput) {

	buffer := make(chan branchOutput)
	s.pipeElement(buffer, output)

	go func() {
		defer close(buffer)

		for in := range input {
			out := branchOutput{
				taken:    false,
				target:   in.pc + in.c,
				linkAddr: in.pc,
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
