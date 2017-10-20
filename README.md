# ARV Go High-level Functional Model #

The Asynchronous RISC-V (ARV) is a 7-stage superscalar asynchronous processor model.
It employs principles based on Communicating Sequential Processes (CSP), a paradigm suitable for modelling asynchronous circuits, while abstracting handshake protocols and data encoding details. In this way, the model is suitable as a specification for either quasi-delay insensitive (QDI) or bundled-data (BD) implementations.

This repository contains the ARV high-level executable model written in [The Go Programming Language](https://golang.org/), which features channels and CSP-based goroutine constructs.
The model describes the overall processor organisation, abstracting implementation details.

## Installation Instructions

Intallation requires the following packages:
+ Go (minimal version 1.7)
+ Git (to fetch dependencies and the latest version)
+ GNU Toolchain for RISC-V (for building example code)

Instructions assume the model runs under some current Unix-like environment, e.g. Linux or MacOS.

Retrieve the latest version from the repository:

	go get bitbucket.org/marcos_sartori/qdi-riscv

If installing from a tarball release, untar the contents in `src/bitbucket.org/marcos_sartori/qdi-riscv` in your `$GOPATH` directory. An example set of commands to perform this is:

	mkdir -p ${GOPATH}/src/bitbucket.org/marcos_sartori
	tar -xvf qdi-riscv-1.0.tar.bz2 -C ${GOPATH}/src/bitbucket.org/marcos_sartori
	go get bitbucket.org/marcos_sartori/qdi-riscv
	
In order to compile the code in `samplecode` and `riscv-tests` directories it is first necessary to get and install [The GNU Toolchain for RISC-V](https://github.com/riscv/riscv-gnu-toolchain), as described in:

	git clone --recursive https://github.com/riscv/riscv-gnu-toolchain
	cd riscv-gnu-toolchain
	mkdir build
	cd build
	../configure --with-arch=rv32ima --with-abi=ilp32
	sudo make
	
## Running Code in the Model

After completing the above installation, it is possible to compile and run the included sample code:

	cd samplecode
	make
	qdi-riscv -memfile hanoi.bin

The code provided in the `samplecode` directory is useful to write programs targeting the simulation platform.

To run the RISC-V Unit Test Suite, use the command set below:

	cd riscv-tests
	make
	qdi-riscv -memfile test.bin
	
For further options, including debug flags, use:

	qdi-riscv -v

## Version History

* v1.0 - First release described in Marcos Sartori End of Term Work (In Brazilian Portuguese, this is called a "Trabalho de Conclusao de Curso" or TCC). The text of the TCC is in English.

## Authors and License

The project started as an End of Term Work by Marcos Luiggi Lemos Sartori <marcos.sartori@acad.pucrs.br> at the Pontifical Catholic University of Rio Grande do Sul (PUCRS), RS, Brazil. It was developed within the Hardware Design Support Research Group (GAPH), advised by Prof. Ney Laert Vilar Calazans <ney.calazans@pucrs.br>.

This work is licensed under GPLv2.
