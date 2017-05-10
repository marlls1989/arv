package model

func (s *Model) operandFetchUnit(
	validIn <-chan bool,
	instructionIn, pcAddrIn <-chan uint32,
	xuOperIn <-chan xuOperation,
	opFmt <-chan opFormat,
	regLock <-chan uint32,
	regAaddrIn, regBaddrIn, regDaddrIn <-chan uint32,
	regAdata, regBdata <-chan uint32,

	dispatcherOut chan<- dispatcherInput,
	regAaddrOut, regBaddrOut, regDaddrOut chan<- uint32) {

	go func() {

		var a, b, c uint32

		defer close(regAaddrOut)
		defer close(regBaddrOut)
		defer close(regDaddrOut)
		defer close(dispatcherOut)

		<-s.start
		dispatcherOut <- dispatcherNOP

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
						dispatcherOut <- dispatcherNOP
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
						dispatcherOut <- dispatcherNOP
						regDaddrOut <- 0
					} else {
						break
					}
				}
				regAaddrOut <- regAaddr
				regBaddrOut <- 0
				a = <-regAdata
				b = (uint32)((int32)(ins) >> 20)
			case opFormatU:
				<-regLock
				a = pc
				b = ins & 0xFFFFF000
				regAaddrOut <- 0
				regBaddrOut <- 0
			case opFormatJ:
				<-regLock
				a = pc
				b = (((uint32)((int32)(ins)>>12) & 0xFFF00000) |
					(ins & 0x000FF000) |
					((ins >> 9) & 0x800) |
					((ins >> 20) & 0x7FE))
				regAaddrOut <- 0
				regBaddrOut <- 0
			case opFormatS:
				for lock := range regLock {
					if (regAaddr&^lock == 0) || (regBaddr&^lock == 0) {
						dispatcherOut <- dispatcherNOP
						regDaddrOut <- 0
					} else {
						break
					}
				}
				regAaddrOut <- regAaddr
				regBaddrOut <- regBaddr
				a = <-regAdata
				b = (((uint32)((int32)(ins)>>20) & 0xFFFFFFE0) | ((ins >> 7) & 0x1F))
				c = <-regBdata
			case opFormatB:
				for lock := range regLock {
					if (regAaddr&^lock == 0) || (regBaddr&^lock == 0) {
						dispatcherOut <- dispatcherNOP
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
			dispatcherOut <- dispatcherInput{
				valid:  valid,
				pcAddr: pc,
				xuOper: xuOp,
				a:      a,
				b:      b,
				c:      c}

			if fmt == opFormatB || fmt == opFormatS {
				regDaddrOut <- 0
			} else {
				regDaddrOut <- regDaddr
			}
		}
	}()
}
