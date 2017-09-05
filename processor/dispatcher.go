package processor

import (
	"fmt"
	"log"
	"sync/atomic"
)

// This structure is used to transfer instruction data from the operand fetch unit to the dispatcher
//
// valid is the stream tag
// pcAddr requires no introductions
// xuOper is the control signal selecting which operation to be performed
// a, b and c are operands retrived by the operand fetch unit acording to the instruction format
type dispatcherInput struct {
	valid   uint8
	pcAddr  uint32
	xuOper  xuOperation
	a, b, c uint32
}

func (d dispatcherInput) String() string {
	return fmt.Sprintf("{valid:%v pcAddr:%X xuOper:%v a:%x b:%x c:%x}", d.valid, d.pcAddr, d.xuOper, d.a, d.b, d.c)
}

// This function constructs the dispatcher logic stage
//
// The instruction control signals is received through the the dispatcherIn,
// the dispatcher then selects the correct execution unit using the upper bits of the xuOper,
// sending the instruction to the correct execution unit.
//
// The selected execution unit is sent to the program ordering queue along with the stream tag.
//
// The stage is initiated to a bubble
func (s *processor) dispatcherUnit(
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
		defer close(shifterOut)
		defer close(memoryOut)
		defer close(branchOut)

		<-s.start
		programQOut <- programElement{
			unit:  xuBypassSel,
			valid: 255}
		bypassOut <- 0

		for in := range dispatcherIn {

			select {
			case <-s.quit:
				return
			default:
			}

			xuSel := xuSelector(in.xuOper >> 4)
			if s.Debug {
				log.Printf("Dispatching instruction %v", in)
			}

			switch xuSel {
			case xuBypassSel:
				atomic.AddUint64((&s.Stats.Unit.Bypass), 1)
				bypassOut <- in.b
			case xuAdderSel:
				atomic.AddUint64((&s.Stats.Unit.Adder), 1)
				adderOut <- adderInput{
					op: in.xuOper,
					a:  in.a,
					b:  in.b}
			case xuLogicSel:
				atomic.AddUint64((&s.Stats.Unit.Logic), 1)
				logicOut <- logicInput{
					op: in.xuOper,
					a:  in.a,
					b:  in.b}
			case xuShiftSel:
				atomic.AddUint64((&s.Stats.Unit.Shifter), 1)
				shifterOut <- shifterInput{
					op: in.xuOper,
					a:  in.a,
					b:  in.b}
			case xuMemorySel:
				atomic.AddUint64((&s.Stats.Unit.Memory), 1)
				memoryOut <- memoryUnitInput{
					op: in.xuOper,
					a:  in.a,
					b:  in.b,
					c:  in.c}
			case xuBranchSel:
				atomic.AddUint64((&s.Stats.Unit.Branch), 1)
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
