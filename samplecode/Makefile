TRIPLET = riscv32-unknown-elf
CC = $(TRIPLET)-gcc
CFLAGS = -IbareOS/include -Os -fno-builtin -march=rv32i -nostdinc
AS = $(TRIPLET)-as
ASFLAGS = -march=rv32i
CPP = $(TRIPLET)-g++
CCFLAGS = $(CFLAGS)
LD = $(TRIPLET)-gcc
LDFLAGS = -nostdlib -Os -march=rv32i -Tqdi-riscv.ld
OBJCOPY = $(TRIPLET)-objcopy
OBJDUMP = $(TRIPLET)-objdump

PROGNAME = helloworld dummy test_hanoi

BIN_FILES = $(addsuffix .bin, $(PROGNAME))
LST_FILES = $(addsuffix .lst, $(PROGNAME))

BARE_OS = bareOS/crt0.o \
	bareOS/arch.o \
	bareOS/libos.o \
	bareOS/malloc.o \
	bareOS/softmul.o

all : $(BIN_FILES) $(LST_FILES)

%.bin : %.elf
	$(OBJCOPY) -O binary $< $@

%.lst : %.elf
	$(OBJDUMP) -d $< > $@

%.o: %.c
	$(CC) $(CFLAGS) -o $@ -c $<

%.o: %.cpp
	$(CPP) $(CCFLAGS) -o $@ -c $<

%.o: %.s
	$(AS) $(ASFLAGS) -o $@ $<

%.elf: $(BARE_OS) %.o
	$(LD) -o $@ $+ $(LDFLAGS)

clean :
	rm -f $(BIN_FILES) $(LST_FILES) $(BARE_OS) $(wildcard *~) $(wildcard *.o) 
