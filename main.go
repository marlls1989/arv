package main

import (
	"bitbucket.org/marcos_sartori/qdi-riscv/model"
	"fmt"
	"os"
	"runtime"
)

func main() {
	// Initialize the runtime for best using the available cores
	runtime.GOMAXPROCS(runtime.NumCPU())

	fmt.Printf("RISC-V QDI Model version %s\n", model.Version)

	file, err := os.OpenFile("memdump", os.O_RDWR, 0)
	if err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}

	mem, err := model.InitializeMemoryFromFile(file)

	if err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}

	fmt.Println(mem)
}
