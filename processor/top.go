package processor

/* This function instanciate a new processor unit
 * calling each module constructor and constructing the interconect */

func ConstructProcessor() *Processor {
	proc := new(Processor)

	proc.start = make(chan struct{})
	proc.quit = make(chan struct{})

	brCmd := make(chan branchCmd)
	instruction := make(chan []byte)
	fetchValid := make(chan bool)
	fetchPcAddr := make(chan uint32)

	proc.fetchUnit(brCmd, pcAddr, instruction, fetchValid)

	decoderValid := make(chan bool)
	decoderPcAddr := make(chan uint32)
	decoderOutput := make(chan decoderOut)

	proc.decoderUnit(fetchValid, fetchPcAddr, instruction,
		decoderValid, decoderPcAddr, decoderOutput)

	regLock := make(chan uint32)
	regRData := make(chan regDataRet)

}
