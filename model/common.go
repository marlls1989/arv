package model

const Version string = "0.1"

type modelState struct {
	start, quit chan struct{}
	memory      *memory
	startPC     uint32
}
