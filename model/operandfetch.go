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
		xuOperOut <- bypassB
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
						xuOperOut <- bypassB
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
						xuOperOut <- bypassB
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
			case opFormatU:
				<-regLock
				a = pc
				b = ins & 0xFFFFF000
			case opFormatJ:
				<-regLock
				a = pc
				b = (((uint32)((int32)(ins)>>12) & 0xFFF00000) |
					(ins & 0x000FF000) |
					((ins >> 9) & 0x800) |
					((ins >> 20) & 0x7FE))
			case opFormatS:
				for lock := range regLock {
					if (regAaddr&^lock == 0) || (regBaddr&^lock == 0) {
						validOut <- valid
						pcAddrOut <- 0
						xuOperOut <- bypassB
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
				c = (((uint32)((int32)(ins)>>20) & 0xFFFFFFE0) | ((ins >> 7) & 0x1F))
			case opFormatB:
				for lock := range regLock {
					if (regAaddr&^lock == 0) || (regBaddr&^lock == 0) {
						validOut <- valid
						pcAddrOut <- 0
						xuOperOut <- bypassB
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
				c = (((uint32)((int32)(ins)>>20) & 0xFFFFF000) |
					((ins << 4) & 0x800) |
					((ins >> 20) & 0x7E0) |
					((ins >> 7) & 0x1E))
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
