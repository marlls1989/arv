package main

import (
	"bitbucket.org/marcos_sartori/qdi-riscv/memory"
	"bitbucket.org/marcos_sartori/qdi-riscv/processor"
	"flag"
	"fmt"
	"log"
	"os"
	"runtime"
)

func main() {
	// Initialize the runtime for best using the available cores
	runtime.GOMAXPROCS(runtime.NumCPU())

	sizePtr := flag.Int64("memsize", -1, "Truncate memory file to size in kbytes")

	flag.Parse()

	args := flag.Args()

	if len(args) < 1 {
		fmt.Fprintf(os.Stderr, "USAGE: %s [options] <memorydump file>", os.Args[0])
	}

	file, err := os.OpenFile(args[0], os.O_RDWR|os.O_CREATE, 0666)
	if err != nil {
		log.Fatal(err)
		os.Exit(-1)
	}
	log.Print("Memdump file opened")

	if *sizePtr >= 0 {
		file.Truncate((*sizePtr) * 1024)
		log.Printf("Memdump file truncated to %dkb", *sizePtr)
	}

	mem, err := memory.MemoryArrayFromFile(file)

	if err != nil {
		log.Fatal(err)
	}

	log.Print("Memory model created from file")

	proc := processor.ConstructProcessor(mem)
	log.Print("Processor model instantiated")

	proc.Start()
	log.Print("Simulation started")

	<-mem.EndSimulation
	proc.Stop()
	log.Print("Finishing Simulation")
	file.Close()
}
