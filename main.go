package main

import (
	"bitbucket.org/marcos_sartori/qdi-riscv/memory"
	"bitbucket.org/marcos_sartori/qdi-riscv/processor"
	"fmt"
	"os"
	"runtime"
)

func main() {
	// Initialize the runtime for best using the available cores
	runtime.GOMAXPROCS(runtime.NumCPU())

	file, err := os.OpenFile("memdump", os.O_RDWR|os.O_CREATE, 0666)
	if err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}

	proc := processor.ConstructProcessor()

	file.Truncate(4096)

	mem, err := memory.MemoryArrayFromFile(file)

	if err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}

	proc.Memory = mem

}
