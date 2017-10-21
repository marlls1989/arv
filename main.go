package main

import (
	"flag"
	"fmt"
	"github.com/marlls1989/arv/memory"
	"github.com/marlls1989/arv/processor"
	"gopkg.in/yaml.v2"
	"log"
	"os"
	"runtime"
	"sync/atomic"
)

func main() {
	sizePtr := flag.Int64("memsize", -1, "Truncate memory file to `size` in kbytes")
	nProcs := flag.Int("jobs", runtime.NumCPU(), "Number of concurent execution `threads`")
	memfile := flag.String("memfile", "", "Memory dump `file name` containing the binary file name")
	memdebug := flag.Bool("memdebug", false, "Logs memory writes to stderr")
	coreedebug := flag.Bool("coredebug", false, "Logs instructions flow to stderr")
	statsfile := flag.String("statsfile", "stats.yaml", "Record execution statistics to file")

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
	}
	log.Print("Memdump file opened")
	defer file.Close()

	if *sizePtr >= 0 {
		file.Truncate((*sizePtr) * 1024)
		log.Printf("Memdump file truncated to %dkb", *sizePtr)
	}

	mem, err := memory.MemoryArrayFromFile(file)

	mem.Debug = *memdebug

	if err != nil {
		log.Fatal(err)
	}

	log.Print("Memory model created from file")

	proc := processor.ConstructProcessor(mem)
	log.Print("Processor model instantiated")

	proc.Debug = *coreedebug

	proc.Start()
	log.Print("Simulation started")

	<-mem.EndSimulation
	proc.Stop()
	log.Print("Finishing Simulation")

	file, err = os.OpenFile(*statsfile, os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		log.Printf("Could not open statsfile %s, %v", *statsfile, err)
	} else {
		defer file.Close()
		stats, err := yaml.Marshal(&proc.Stats)
		if err == nil {
			file.Write(stats)
			log.Printf("Execution statistics written to file %s", *statsfile)
		}
	}

	if err != nil {
		log.Printf("Decoded: %v instructions", atomic.LoadUint64(&proc.Stats.Decoded))
		log.Printf("Inserted: %v bubbles", atomic.LoadUint64(&proc.Stats.Bubbles))
		log.Printf("Retired: %v instructions", atomic.LoadUint64(&proc.Stats.Retired))
		log.Printf("Cancelled: %v instructions", atomic.LoadUint64(&proc.Stats.Cancelled))
	}
}
