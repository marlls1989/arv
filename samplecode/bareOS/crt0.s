	.section .text.init
	.global _entry
_entry:
	la a0, _bss_start
	la a1, _bss_end
	la sp, _stack
.clr_bss:
	bgeu a0, a1, .clr_bss_done
	sw zero, 0(a0)
	addi a0, a0, 4
	beq zero, zero, .clr_bss
.clr_bss_done:
	la  ra, halt
	jal zero, kmain

