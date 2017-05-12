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

	file, err := os.OpenFile("memdump", os.O_RDWR|os.O_CREATE, 0666)
	if err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}

	file.Truncate(4096)
	if err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}

}
