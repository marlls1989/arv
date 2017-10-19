# ARV Go High-level Functional Model #

The Asynchronous RISC-V (ARV) is a 7-stage superscalar asynchronous processor design.
It is designed using Communicating Sequential Processes (CSP), a paradigm suitable for modelling asynchronous circuits abstracting handshake protocol and encoding details.

This repository contains the ARV high-level model written in [The Go Programming Language](https://golang.org/), which features channels and goroutine CSP-based constructs.
This model is used to develop the processor organisation, abstracting asynchronous specific design complexities.

## Installation Instructions

You will need the following packages installed before you start:
+ Go (minimal version 1.7)
+ Git (to fetch dependencies and the latest version)
+ GNU Toolchain for RISC-V (for building example code)

The instructions assume you are running a modern Unix-like environment, e.g. Linux or MacOS.

Retrieve the latest version from the repository:

	go get bitbucket.org/marcos_sartori/qdi-riscv

If you are installing from a tarball release, untar the contents in `src/bitbucket.org/marcos_sartori/qdi-riscv` in your `$GOPATH` directory:

	mkdir -p ${GOPATH}/src/bitbucket.org/marcos_sartori
	tar -xvf qdi-riscv-1.0.tar.bz2 -C ${GOPATH}/src/bitbucket.org/marcos_sartori
	go get bitbucket.org/marcos_sartori/qdi-riscv
	
In order to compile the code in `samplecode` and `riscv-tests` directories you first need to get [The GNU Toolchain for RISC-V](https://github.com/riscv/riscv-gnu-toolchain) dependencies before installing it:

	git clone --recursive https://github.com/riscv/riscv-gnu-toolchain
	cd riscv-gnu-toolchain
	mkdir build
	cd build
	../configure --with-arch=rv32ima --with-abi=ilp32
	sudo make
	
## Running Code in the Model

After installing, you are ready to compile and run the included sample code:

	cd samplecode
	make
	qdi-riscv -memfile hanoi.bin

You can use the code provided in the `samplecode` directory to write your own programs targeting the simulation platform.
	
To run the RISC-V Unit Test Suite use:

	cd riscv-tests
	make
	qdi-riscv -memfile test.bin
	
For further options, including debug flags use:

	qdi-riscv -v

## Version History

* v1.0 - First release described in Marcos Sartori End of Term Project (In Brazilian Portuguese this is called a "Trabalho de Conclusão de Curso" or TCC). The text of the TCC is in English.

## Authors and License

This project started as an End of Term Project by Marcos Luiggi Lemos Sartori <marcos.sartori@acad.pucrs.br> at the Pontifical Catholic University of Rio Grande do Sul (PUCRS). It was developed withing the Hardware Design Support Research Group (GAPH), and was advised by Prof. Dr. Ney Laert Vilar Calazans <ney.calazans@pucrs.br>.

This work is licensed under GPLv2.
