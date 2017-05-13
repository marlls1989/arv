package main

import (
	"bitbucket.org/marcos_sartori/qdi-riscv/memory"
	"bitbucket.org/marcos_sartori/qdi-riscv/processor"
	"log"
	"os"
	"runtime"
)

func main() {
	// Initialize the runtime for best using the available cores
	runtime.GOMAXPROCS(runtime.NumCPU())

	file, err := os.OpenFile("memdump", os.O_RDWR|os.O_CREATE, 0666)
	if err != nil {
		log.Fatal(err)
		os.Exit(-1)
	}
	log.Print("Memdump file opened")

	proc := processor.ConstructProcessor()
	log.Print("Processor model instantiated")

	file.Truncate(1024 * 1024)
	mem, err := memory.MemoryArrayFromFile(file)

	if err != nil {
		log.Fatal(err)
		os.Exit(-1)
	}
	proc.Memory = mem
	log.Print("Memory model created from file")

	proc.Start()
	log.Print("Simulation started")

	for {
	}

}
