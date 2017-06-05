package processor

import (
	"log"
)

type regReadCmd struct {
	aaddr, baddr regAddr
}

type regDataRet struct {
	adata, bdata uint32
}

func (s *processor) registerBypass(
	regWcmd <-chan retireRegwCmd,
	regWaddrIn <-chan regAddr,
	regRcmd <-chan regReadCmd,

	regDataOut chan<- regDataRet) {

	regWaddr := make(chan regAddr)
	regWdata := make(chan uint32)

	s.regFile.WritePort(regWaddr, regWdata)

	regAaddr := make(chan regAddr)
	regAdata := make(chan uint32)

	s.regFile.ReadPort(regAaddr, regAdata)

	regBaddr := make(chan regAddr)
	regBdata := make(chan uint32)

	s.regFile.ReadPort(regBaddr, regBdata)

	go func() {
		defer close(regAaddr)
		defer close(regBaddr)
		defer close(regWaddr)
		defer close(regWdata)
		defer close(regDataOut)

		for wcmd := range regWcmd {
			wdata := wcmd.data

			waddr, vwaddr := <-regWaddrIn
			if !vwaddr {
				return
			}

			we := wcmd.we && (waddr != 0)

			if we {
				if s.Debug {
					log.Printf("Writing %X to register %v", wdata, waddr)
				}
				regWdata <- wdata
				regWaddr <- waddr
			}

			rcmd, vrcmd := <-regRcmd
			if !vrcmd {
				return
			}

			aaddr := rcmd.aaddr
			baddr := rcmd.baddr

			var adata, bdata uint32

			if aaddr != 0 {
				if we && (waddr&aaddr&0xFFFFFFFE != 0) {
					adata = wdata
				} else {
					regAaddr <- aaddr
					adata = <-regAdata
				}
			}

			if baddr != 0 {
				if we && (waddr&baddr&0xFFFFFFFE != 0) {
					bdata = wdata
				} else {
					regBaddr <- baddr
					bdata = <-regBdata
				}
			}

			if aaddr != 0 || baddr != 0 {
				regDataOut <- regDataRet{
					adata: adata,
					bdata: bdata}
			}
		}
	}()
}
