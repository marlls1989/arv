	.global halt
halt:
	li s10, 0x80000000
	sw zero, 0(s10)
	j halt
