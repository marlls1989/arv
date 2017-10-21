package processor

import (
	"github.com/marlls1989/arv/memory"
)

// This function instanciate a new processor unit
// calling each logical stage constructor and creating the interconect channels.
// The function receives the memory model as argument and returns a processor struct.
func ConstructProcessor(mem memory.Memory) *processor {
	proc := new(processor)

	proc.start = make(chan struct{})
	proc.quit = make(chan struct{})
	proc.Memory = mem

	brCmd := make(chan uint32)
	instruction := make(chan []byte)
	fetchValid := make(chan uint8)
	fetchPcAddr := make(chan uint32)

	proc.fetchUnit(brCmd, fetchPcAddr, instruction, fetchValid)

	decoderOutput := make(chan decoderOut)

	proc.decoderUnit(fetchValid, fetchPcAddr, instruction, decoderOutput)

	regLock := make(chan uint32)
	dispatcherCmd := make(chan dispatcherInput)
	regDaddrIn := make(chan regAddr)

	proc.operandFetchUnit(decoderOutput, regLock, dispatcherCmd, regDaddrIn)

	regDaddrOut := make(chan regAddr)
	proc.registerLock(regDaddrIn, regDaddrOut, regLock, 4)

	regWcmd := make(chan regWCmd)

	proc.regFile.WritePort(regDaddrOut, regWcmd)

	ctrlQIn := make(chan programElement)
	bypassIn := make(chan uint32)
	adderIn := make(chan adderInput)
	logicIn := make(chan logicInput)
	shifterIn := make(chan shifterInput)
	memoryIn := make(chan memoryUnitInput)
	branchIn := make(chan branchInput)

	proc.dispatcherUnit(dispatcherCmd,
		ctrlQIn,
		bypassIn,
		adderIn,
		logicIn,
		shifterIn,
		memoryIn,
		branchIn)

	bypassOut := make(chan uint32)
	proc.bypassUnit(bypassIn, bypassOut, 2)

	ctrlQOut := make(chan programElement)
	proc.programQueue(ctrlQIn, ctrlQOut, 2)

	adderOut := make(chan uint32)
	proc.adderUnit(adderIn, adderOut)

	logicOut := make(chan uint32)
	proc.logicUnit(logicIn, logicOut)

	shifterOut := make(chan uint32)
	proc.shifterUnit(shifterIn, shifterOut)

	memoryOut := make(chan memoryUnitOutput)
	memoryWe := make(chan bool)
	proc.memoryUnit(memoryIn, memoryWe, memoryOut)

	branchOut := make(chan branchOutput)
	proc.branchUnit(branchIn, branchOut)

	proc.retireUnit(ctrlQOut,
		bypassOut,
		adderOut,
		logicOut,
		shifterOut,
		memoryOut,
		branchOut,

		regWcmd,
		memoryWe,
		brCmd)

	return proc
}
