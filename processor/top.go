package processor

import (
	"bitbucket.org/marcos_sartori/qdi-riscv/memory"
)

/* This function instanciate a new processor unit
 * calling each module constructor and constructing the interconect */

func ConstructProcessor(mem memory.Memory) *Processor {
	proc := new(Processor)

	proc.start = make(chan struct{})
	proc.quit = make(chan struct{})
	proc.Memory = mem

	/* Creating an unitialised queue for branch commands
	 * to acommodate NOPs inserted by the execution loop
	 * into the control loop */
	brCmd := make(chan uint32, 1)
	instruction := make(chan []byte)
	fetchValid := make(chan uint8)
	fetchPcAddr := make(chan uint32)

	proc.fetchUnit(brCmd, fetchPcAddr, instruction, fetchValid)

	decoderValid := make(chan uint8)
	decoderPcAddr := make(chan uint32)
	decoderOutput := make(chan decoderOut)

	proc.decoderUnit(fetchValid, fetchPcAddr, instruction,
		decoderValid, decoderPcAddr, decoderOutput)

	regLock := make(chan uint32)
	regRData := make(chan regDataRet)
	dispatcherCmd := make(chan dispatcherInput)
	regRcmd := make(chan regReadCmd)
	regDaddrIn := make(chan regAddr)

	proc.operandFetchUnit(decoderValid, decoderPcAddr, decoderOutput,
		regRData, regLock, dispatcherCmd, regRcmd, regDaddrIn)

	regDaddrOut := make(chan regAddr)
	proc.registerLock(regDaddrIn, regDaddrOut, regLock, 4)

	regWcmd := make(chan retireRegwCmd)

	proc.registerBypass(regWcmd, regDaddrOut, regRcmd, regRData)

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
