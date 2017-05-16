package processor

import (
	"fmt"
	"log"
)

type dispatcherInput struct {
	valid   uint8
	pcAddr  uint32
	xuOper  xuOperation
	a, b, c uint32
}

func (d dispatcherInput) String() string {
	return fmt.Sprintf("{valid:%v pcAddr:%X xuOper:%v a:%x b:%x c:%x}", d.valid, d.pcAddr, d.xuOper, d.a, d.b, d.c)
}

func (s *Processor) dispatcherUnit(
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

		for in := range dispatcherIn {
			xuSel := xuSelector(in.xuOper >> 4)

			log.Printf("Dispatching instruction %v", in)

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
