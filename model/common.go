package model

const Version string = "0.1"

type modelState struct {
	memory *memory
	quit   chan bool
}
