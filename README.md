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

	go get -u github.com/marlls1989/arv
	
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
	arv -memfile hanoi.bin

The code provided in the `samplecode` directory is useful to write programs targeting the simulation platform.

To run the RISC-V Unit Test Suite, use the command set below:

	cd riscv-tests
	make
	arv -memfile test.bin
	
For further options, including debug flags, use:

	arv -h

## Version History

* v1.0 - First release described in Sartori's End of Term Work (In Brazilian Portuguese, this is called a "Trabalho de Conclusao de Curso" or TCC). The text of the TCC is in English.
* v1.1 - First public release on Github, the register write-back mechanism is simplified and new performance counters are introduced. A conference paper containing results extracted from this version was approved for publication in the [24th IEEE International Conference on Electronics, Circuits and Systems](http://icecs2017.org/).

## Publications
* SARTORI, M. L. L. ARV: Towards an Asynchronous Implementation of the RISC-V Architecture. End of Term Work. Computer Engineering - PUCRS, July 2017. 57p. (Presented and approved. Advisor: Ney Laert Vilar Calazans).
* SARTORI, M. L. L.; CALAZANS, N. L. V. Go Functional Model for a RISC-V Asynchronous Organization - ARV. In: IEEE International Conference on Electronics, Circuits and Systems (ICECS'17), Batumi, 2017. Accepted for publication, final version in preparation.

## Authors and License

The project started as an End of Term Work by Marcos Luiggi Lemos Sartori <marcos.sartori@acad.pucrs.br> at the Pontifical Catholic University of Rio Grande do Sul (PUCRS), RS, Brazil. It was developed within the [Hardware Design Support Research Group (GAPH)](http://www.inf.pucrs.br/gaph), advised by Prof. Ney Laert Vilar Calazans <ney.calazans@pucrs.br>.

This work is licensed under GPLv2.
