// This package contains high-level functional memory models.
package memory

// Contains the logic stage constructor interfaces the processor model expects
// memory models implements.
//
//
// (i) ReadWritePort should implement a read-then-write memory port logic stage constructor.
// The logic stage constructed should exchange data through channels and respect the following interface:
// Input channel addr should provide the 32-bit memory base address;
// lng the lenght in bytes of the data to be read;
// dataIn the data to be optionally written;
// a true value read from boolen we input channel authorizes a memory write,
// the we channell is only read if the lenght of dataIn is greater than 0,
// otherwise the we channel should not be read.
// 
// The input channels are expected to be read in the following order:
// addr, dataIn, lng and optionally we.
// Not respecting this order may break compatibility with the processor model causing a deadlock.
// Real circuitry would avoid this deadlock condition by reading the three first channels in parallel.
//
// The data at the received base address should returned through the dataOut channel.
// A slice of bytes the length received should be written to the dataOut channel
// after successfully reading all required information from the addr and lng channels.
//
//
// (ii) ReadPort should implement a read-only memory port logic stage constructor
// The logic stage constructed should read the address and the length from the input channels addr and lng;
// writing a byte slice of the requested length containing the memory contents at the requested address to the output channel dataOut.
//
// The constructed stage should read first the addr channel and then the lng channel.
// Not respecting this order may break compatibility and cause deadlocks.
type Memory interface {
	ReadWritePort(
		addr, lng <-chan uint32,
		dataIn <-chan []byte,
		we <-chan bool,
		dataOut chan<- []byte)
	ReadPort(addr, lng <-chan uint32, data chan<- []byte)
}
