package model

func (s *Model) operandFetchUnit(
	validIn <-chan bool,
	instructionIn, pcAddrIn <-chan uint32,
	xuOperIn <-chan xuOperation,
	opFmt <-chan opFormat,
	regLock <-chan uint32,
	regAaddrIn, regBaddrIn, regDaddrIn <-chan uint32,
	regAdata, regBdata <-chan uint32,

	validOut chan<- bool,
	pcAddrOut chan<- uint32,
	xuOperOut chan<- xuOperation,
	regAaddrOut, regBaddrOut, regDaddrOut chan<- uint32,
	operandA, operandB, operandC chan<- uint32) {

	go func() {

		var a, b, c uint32

		defer close(validOut)
		defer close(xuOperOut)
		defer close(regAaddrOut)
		defer close(regBaddrOut)
		defer close(regDaddrOut)
		defer close(operandA)
		defer close(operandB)
		defer close(operandC)

		<-s.start
		validOut <- false
		pcAddrOut <- 0
		xuOperOut <- bypassA
		operandA <- 0
		operandB <- 0
		operandC <- 0
		regDaddrOut <- 0
		for ins := range instructionIn {
			valid, vvalid := <-validIn
			pc, vpc := <-pcAddrIn

			if !vvalid || !vpc {
				return
			}

			fmt := <-opFmt
			xuOp := <-xuOperIn
			regAaddr := <-regAaddrIn
			regBaddr := <-regBaddrIn
			regDaddr := <-regDaddrIn

			a = 0
			b = 0
			c = 0

			switch fmt {
			case opFormatR:
				for lock := range regLock {
					if (regAaddr&^lock == 0) || (regBaddr&^lock == 0) {
						validOut <- valid
						pcAddrOut <- 0
						xuOperOut <- bypassA
						operandA <- 0
						operandB <- 0
						operandC <- 0
						regDaddrOut <- 0
					} else {
						break
					}
				}
				regAaddrOut <- regAaddr
				regBaddrOut <- regBaddr
				a = <-regAdata
				b = <-regBdata
			case opFormatI:
				for lock := range regLock {
					if regAaddr&^lock == 0 {
						validOut <- valid
						pcAddrOut <- 0
						xuOperOut <- bypassA
						operandA <- 0
						operandB <- 0
						operandC <- 0
						regDaddrOut <- 0
					} else {
						break
					}
				}
				regAaddrOut <- regAaddr
				a = <-regAdata
				b = (uint32)((int32)(ins) >> 20)
			case opFormatB:
			case opFormatJ:
			case opFormatS:
			case opFormatU:
			}
			validOut <- valid
			pcAddrOut <- pc
			xuOperOut <- xuOp
			operandA <- a
			operandB <- b
			operandC <- c
			if fmt == opFormatB || fmt == opFormatS {
				regDaddrOut <- 0
			} else {
				regDaddrOut <- regDaddr
			}
		}
	}()
}
