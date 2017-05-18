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
	sizePtr := flag.Int64("memsize", -1, "Truncate memory file to `size` in kbytes")
	nProcs := flag.Int("jobs", runtime.NumCPU(), "Number of concurent execution `threads`")
	memfile := flag.String("memfile", "", "Memory dump `file name` containing the binary file name")

	flag.Parse()

	// Initialize the runtime for best using the available cores
	runtime.GOMAXPROCS(*nProcs)

	if *memfile == "" {
		fmt.Fprintln(os.Stderr, "Required parameter -memfile not defined")
		os.Exit(-1)
	}

	file, err := os.OpenFile(*memfile, os.O_RDWR|os.O_CREATE, 0666)
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
