# ARV Go High-level Functional Model #

The Asynchronous RISC-V (ARV) is a 7-stage superscalar asynchronous processor design.
It is designed using Communicating Sequential Processes (CSP), a paradigm suitable for modelling asynchronous circuits abstracting handshake protocol and encoding details.

This repository contains ARV high-level model written in [The Go Programming Language](https://golang.org/), which features channels and goroutine CSP-based constructs.
This model is used to develop the processor organisation, abstracting asynchronous specific design complexities.

## Installation Instructions

The memory model requires [launchpad.net/gommap](https://godoc.org/launchpad.net/gommap) as a consequence it only runs on Linux machines.

## Version History

* v1.0 - First release described in my [End of Term Project](https://link.pro/tcc.pdf)

## Related Publications
