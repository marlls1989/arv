package model

func (s *Model) operandFetchUnit(
	validIn <-chan bool,
	pcAddrIn <-chan uint32,
	regLock <-chan uint32,
	decodedIn <-chan decoderOut,
	regDataIn <-chan regDataRet,

	dispatcherOut chan<- dispatcherInput,
	regRcmdOut chan<- regReadCmd,
	regDaddrOut chan<- uint32) {

	go func() {

		var a, b, c uint32

		defer close(regRcmdOut)
		defer close(regDaddrOut)
		defer close(dispatcherOut)

		<-s.start
		dispatcherOut <- dispatcherNOP

		regDaddrOut <- 0
		for decoded := range decodedIn {
			valid, vvalid := <-validIn
			pc, vpc := <-pcAddrIn

			if !vvalid || !vpc {
				return
			}

			ins := decoded.ins
			fmt := decoded.fmt
			xuOp := decoded.op
			regAaddr := decoded.regA
			regBaddr := decoded.regB
			regDaddr := decoded.regD

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
				regRcmdOut <- regReadCmd{
					aaddr: regAaddr,
					baddr: regBaddr}
				read := <-regDataIn
				a = read.adata
				b = read.bdata
			case opFormatI:
				for lock := range regLock {
					if regAaddr&^lock == 0 {
						dispatcherOut <- dispatcherNOP
						regDaddrOut <- 0
					} else {
						break
					}
				}
				regRcmdOut <- regReadCmd{
					aaddr: regAaddr}
				read := <-regDataIn
				a = read.adata
				b = (uint32)((int32)(ins) >> 20)
			case opFormatU:
				<-regLock
				a = pc
				b = ins & 0xFFFFF000
				regRcmdOut <- regReadCmd{
					aaddr: 0,
					baddr: 0}
			case opFormatJ:
				<-regLock
				a = pc
				b = (((uint32)((int32)(ins)>>12) & 0xFFF00000) |
					(ins & 0x000FF000) |
					((ins >> 9) & 0x800) |
					((ins >> 20) & 0x7FE))
				regRcmdOut <- regReadCmd{
					aaddr: 0,
					baddr: 0}
			case opFormatS:
				for lock := range regLock {
					if (regAaddr&^lock == 0) || (regBaddr&^lock == 0) {
						dispatcherOut <- dispatcherNOP
						regDaddrOut <- 0
					} else {
						break
					}
				}
				regRcmdOut <- regReadCmd{
					aaddr: regAaddr,
					baddr: regBaddr}
				read := <-regDataIn
				a = read.adata
				b = (((uint32)((int32)(ins)>>20) & 0xFFFFFFE0) | ((ins >> 7) & 0x1F))
				c = read.bdata
			case opFormatB:
				for lock := range regLock {
					if (regAaddr&^lock == 0) || (regBaddr&^lock == 0) {
						dispatcherOut <- dispatcherNOP
						regDaddrOut <- 0
					} else {
						break
					}
				}
				regRcmdOut <- regReadCmd{
					aaddr: regAaddr,
					baddr: regBaddr}
				read := <-regDataIn
				a = read.adata
				b = read.bdata
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
