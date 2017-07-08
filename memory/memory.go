// This package contains high-level functional memory models.
package memory

// Contains the logic stage constructor interfaces the processor model expects
// memory models implements.
//
// ReadWritePort should implement a read-then-write memory port logic stage constructor.
// The logic stage constructed should exchange data through channels and respect the following interface:
// Input channel addr should provide the 32-bit memory base address;
// lng the lenght in bytes of the data to be read;
// dataIn the data to be optionally written;
// and the we channel is read if the lenght of dataIn is greater than 0, a true value 
// otherwise the we channel should not be read.

type Memory interface {
	ReadWritePort(
		addr, lng <-chan uint32,
		dataIn <-chan []byte,
		we <-chan bool,
		dataOut chan<- []byte)
	ReadPort(addr, len <-chan uint32, data chan<- []byte)
}
