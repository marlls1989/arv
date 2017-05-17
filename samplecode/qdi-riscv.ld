ENTRY(_entry)

/* specify HF-RISC memory areas */


/* define a global symbol _stack  */

/* now define the output sections  */

SECTIONS {
	. = 0;

	.text : {
			_text = .;
			*(.text.init)	
			*(.text)
			*(.text.*)
			_etext = .;
	}

	.data : {
			_rodata = .;
			*(.rdata)
			*(.rodata)
			*(.rodata.*)
			_erodata = .;
			_data = .;
			*(.data)
			*(.data.*)
			_gp = . + 0x800;
			*(.sdata)
			*(.sdata.*)
			. = ALIGN(16); 
			_edata = .;
			_bss_start = .;
			*(.sbss) 
			*(.sbss.*) 
			*(.bss) 
			*(.bss.*)
			. = ALIGN(4);
			_bss_end = . ;
			. += 0x2000;
			_stack = .;
	}

	_end = .;
}

