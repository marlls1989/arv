package model

type dispatcherInput struct {
	valid   bool
	pcAddr  uint32
	xuOper  xuOperation
	a, b, c uint32
}

var dispatcherNOP = dispatcherInput{
	valid:  false,
	pcAddr: 0,
	xuOper: bypassB,
	a:      0,
	b:      0,
	c:      0}

func (s *Model) dispatcherUnit(
	dispatcherIn <-chan dispatcherInput,

	programQOut chan<- programElement,
	bypassOut chan<- uint32,
	adderOut chan<- adderInput,
	logicOut chan<- logicInput,
	shifterOut chan<- shifterInput,
	memoryOut chan<- memoryUnitInput,
	branchOut chan<- branchInput) {

	go func() {
		defer close(programQOut)
		defer close(bypassOut)
		defer close(adderOut)
		defer close(logicOut)
		defer close(memoryOut)
		defer close(branchOut)

		<-s.start
		bypassOut <- 0
		programQOut <- programQNOP
		for in := range dispatcherIn {

			xuSel := xuSelector(in.xuOper >> 4)

			switch xuSel {
			case xuBypassSel:
				bypassOut <- in.b
			case xuAdderSel:
				adderOut <- adderInput{
					op: in.xuOper,
					a:  in.a,
					b:  in.b}
			case xuLogicSel:
				logicOut <- logicInput{
					op: in.xuOper,
					a:  in.a,
					b:  in.b}
			case xuShiftSel:
				shifterOut <- shifterInput{
					op: in.xuOper,
					a:  in.a,
					b:  in.b}
			case xuMemorySel:
				memoryOut <- memoryUnitInput{
					op: in.xuOper,
					a:  in.a,
					b:  in.b,
					c:  in.c}
			case xuBranchSel:
				branchOut <- branchInput{
					op: in.xuOper,
					a:  in.a,
					b:  in.b,
					c:  in.c,
					pc: in.pcAddr}
			}

			programQOut <- programElement{
				valid: in.valid,
				unit:  xuSel}
		}

	}()
}
