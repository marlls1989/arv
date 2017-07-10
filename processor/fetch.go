package processor

import (
	"log"
)

// nextPc and streamTag logic unit constructor
//
// receives the currPC and currValid values from the delay logic stage in the fetchUnit
// and computes the nextPc, fetchaddr and nextValid (streamtag) values based on the presence of a branch address
func (s *processor) nextPcUnit(
	currPC <-chan uint32,
	currValid <-chan uint8,
	branch <-chan uint32,
	fetchAddr, nextPc chan<- uint32,
	nextValid chan<- uint8) {

	go func() {

		defer close(fetchAddr)
		defer close(nextPc)
		defer close(nextValid)

		<-s.start
		nextPc <- s.StartPC
		fetchAddr <- s.StartPC
		nextValid <- 0
		for pc := range currPC {
			select {
			case _ = <-s.quit:
				return
			default:
			}
			var target uint32
			valid := <-currValid
			select {
			case br := <-branch:
				target = br
				valid = valid + 1
				if s.Debug {
					log.Printf("Branching to 0x%X", target)
				}
			default: //Uncouple the fetch loop by taking branch completeness as a cue
				target = pc + 4
			}
			fetchAddr <- target
			nextPc <- target
			nextValid <- valid
		}
	}()
}

// the fetchUnit constructor instantiates the nextPcUnit,
// a delay logic element and a memory port.
//
// The delay logic element delays in one handshake cycle the value of nextPc and nextValid
// sending them back to the nextPcUnit to form the pc/streamtag loop
// and forwarding them to the decoder unit.
//
// In parallel with the delay element, the memory port fetches the instruction from the fetchAddr
// outputting the fetched instruction directly to the decoderUnit.
func (s *processor) fetchUnit(
	branch <-chan uint32,

	pcAddr chan<- uint32,
	instruction chan<- []byte,
	valid chan<- uint8) {

	nextPC := make(chan uint32)
	currPC := make(chan uint32)
	fetchAddr := make(chan uint32)
	nextValid := make(chan uint8)
	currValid := make(chan uint8)
	len := make(chan uint32)

	s.nextPcUnit(currPC, currValid, branch, fetchAddr, nextPC, nextValid)

	go func() {
		defer close(currPC)
		defer close(pcAddr)

		for i := range nextPC {
			currPC <- i
			pcAddr <- i
		}
	}()

	go func() {
		defer close(currValid)
		defer close(valid)

		for i := range nextValid {
			currValid <- i
			valid <- i
		}
	}()

	go func() {
		defer close(len)
		for {
			select {
			case len <- 4:
			case <-s.quit:
				return
			}
		}
	}()

	s.Memory.ReadPort(fetchAddr, len, instruction)
}
