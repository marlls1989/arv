package model

func (s *Model) registerBypass(
	writeEnable <-chan bool,
	regWaddrIn, regWdataIn <-chan uint32,
	regAaddrIn, regBaddrIn <-chan uint32,

	regAdataOut, regBdataOut chan<- uint32) {

	regWaddr := make(chan uint32)
	regWdata := make(chan uint32)

	s.regFile.WritePort(regWaddr, regWdata)

	regAaddr := make(chan uint32)
	regAdata := make(chan uint32)

	s.regFile.ReadPort(regAaddr, regAdata)

	regBaddr := make(chan uint32)
	regBdata := make(chan uint32)

	s.regFile.ReadPort(regBaddr, regBdata)

	go func() {
		defer close(regAaddr)
		defer close(regBaddr)
		defer close(regWaddr)
		defer close(regAdataOut)
		defer close(regBdataOut)
		defer close(regWdata)

		for we := range writeEnable {

			wdata, vwdata := <-regWdataIn
			if !vwdata {
				return
			}

			waddr, vwaddr := <-regWaddrIn
			if !vwaddr {
				return
			}

			aaddr, vaaddr := <-regAaddrIn
			if !vaaddr {
				return
			}

			baddr, vbaddr := <-regBaddrIn
			if !vbaddr {
				return
			}

			we = we && (waddr != 0)

			if we {
				regWdata <- wdata
				regWaddr <- waddr
			}

			if aaddr != 0 {
				if we && (waddr&aaddr&0xFFFFFFFE != 0) {
					regAdataOut <- wdata
				} else {
					regAaddr <- aaddr
					adata := <-regAdata
					regAdataOut <- adata
				}
			}

			if baddr != 0 {
				if we && (waddr&baddr&0xFFFFFFFE != 0) {
					regBdataOut <- wdata
				} else {
					regBaddr <- baddr
					bdata := <-regBdata
					regBdataOut <- bdata
				}
			}
		}
	}()
}
