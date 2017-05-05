package model

type adderInput struct {
	a, b      uint32
	operation xuOperation
}

func (s *Model) adderUnit(
	input <-chan adderInput) {
}
