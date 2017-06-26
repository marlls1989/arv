package processor

import (
	"log"
	"sync/atomic"
)

func (s *processor) operandFetchUnit(
	decodedIn <-chan decoderOut,
	regLock <-chan uint32,

	dispatcherOut chan<- dispatcherInput,
	regDaddrOut chan<- regAddr) {

	chAaddr := make(chan regAddr)
	chAdata := make(chan uint32)

	s.regFile.ReadPort(chAaddr, chAdata)

	chBaddr := make(chan regAddr)
	chBdata := make(chan uint32)

	s.regFile.ReadPort(chBaddr, chBdata)

	go func() {

		var a, b, c uint32

		defer close(regDaddrOut)
		defer close(dispatcherOut)
		defer close(chAaddr)
		defer close(chBaddr)

		for decoded := range decodedIn {
			valid := decoded.valid
			pc := decoded.pc

			select {
			case _ = <-s.quit:
				return
			default:
			}

			if s.Debug {
				log.Printf("Decoded %+v pc:%X valid: %v", decoded, pc, valid)
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
			case opFormatI:
				for lock := range regLock {
					if uint32(regAaddr)&^lock == 0 {
						if s.Debug {
							log.Printf("Issuing bubble due to register lock %v", regAddr(lock))
						}
						atomic.AddUint64((&s.Bubbles), 1)
						dispatcherOut <- dispatcherInput{
							valid:  valid,
							pcAddr: pc,
							xuOper: bypassB,
							a:      0,
							b:      0,
							c:      0}
						regDaddrOut <- 1
					} else {
						break
					}
				}

				chAaddr <- regAaddr
				a = <-chAdata
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
					if (uint32(regAaddr)&^lock == 0) || (uint32(regBaddr)&^lock == 0) {
						if s.Debug {
							log.Printf("Issuing bubble due to register lock %v", regAddr(lock))
						}
						atomic.AddUint64((&s.Bubbles), 1)
						dispatcherOut <- dispatcherInput{
							valid:  valid,
							pcAddr: pc,
							xuOper: bypassB,
							a:      0,
							b:      0,
							c:      0}
						regDaddrOut <- 1
					} else {
						break
					}
				}

				chAaddr <- regAaddr
				chBaddr <- regBaddr
				a = <-chAdata
				b = (((uint32)((int32)(ins)>>20) & 0xFFFFFFE0) | ((ins >> 7) & 0x1F))
				c = <-chBdata

			case opFormatB, opFormatR:
				for lock := range regLock {
					if (uint32(regAaddr)&^lock == 0) || (uint32(regBaddr)&^lock == 0) {
						atomic.AddUint64((&s.Bubbles), 1)
						if s.Debug {
							log.Printf("Issuing bubble due to register lock %v", regAddr(lock))
						}
						dispatcherOut <- dispatcherInput{
							valid:  valid,
							pcAddr: pc,
							xuOper: bypassB,
							a:      0,
							b:      0,
							c:      0}
						regDaddrOut <- 1
					} else {
						break
					}
				}
				chAaddr <- regAaddr
				chBaddr <- regBaddr
				a = <-chAdata
				b = <-chBdata
				c = (((uint32)((int32)(ins)>>20) & 0xFFFFF000) |
					((ins << 4) & 0x800) |
					((ins >> 20) & 0x7E0) |
					((ins >> 7) & 0x1E))

			case opFormatNop:
				<-regLock
			}

			dispatcherOut <- dispatcherInput{
				valid:  valid,
				pcAddr: pc,
				xuOper: xuOp,
				a:      a,
				b:      b,
				c:      c}

			if fmt == opFormatB || fmt == opFormatS || fmt == opFormatNop {
				regDaddrOut <- 1
			} else {
				regDaddrOut <- regDaddr
			}
		}
	}()
}
