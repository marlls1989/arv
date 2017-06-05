#include "include/libos.h"

#define UART_REG (*(volatile char*)0x80004000)

void putchar(char c) {
  UART_REG = c;
}

void puts(const char *s) {
  while (*s) putchar(*(s++));
}

char *utoa(unsigned int i, char *s, int base) {
  char c;
  char *p = s;
  char *q = s;
  int b, shift;

  for(b = base, shift = 0 ; !(b & 1) ; shift++, b >>= 1);
  
  if(b == 1) {
    b = (1 << shift)-1;
    do {
      *q++ = '0' + (i & b);
    } while(i >>= shift);
  } else
    do {
      *q++ = '0' + (i % base);
    } while (i /= base);
  
  for (*q = 0; p <= --q; p++){
    (*p > '9')?(c = *p + 39):(c = *p);
    (*q > '9')?(*p = *q + 39):(*p = *q);
    *q = c;
  }
  
  return s;
}

char *itoa(int i, char *s, int base){
  char c;
  char *p = s;
  char *q = s;

  if (i >= 0)
    return utoa(i, s, base);

  *q++ = '-';
  p++;
  do{
    *q++ = '0' - (i % base);
  } while (i /= base);
  
  for (*q = 0; p <= --q; p++){
    (*p > '9')?(c = *p + 39):(c = *p);
    (*q > '9')?(*p = *q + 39):(*p = *q);
    *q = c;
  }
  
  return s;
}

void *memset(void *dst, int ic, size_t bytes) {
  uint8_t *Dst = dst;
  uint32_t f, *Dst32; 
  size_t b;

	uint8_t c = ic;

  // Fills the lower bytes till aligned
  while(((size_t)Dst & 0x3) && (bytes--)) *(Dst++) = c;
  
  // Fills the aligned region using faster memory access instruction
	f = (c << 24) | (c << 16) | (c << 8) | c;
	Dst32 = (uint32_t*)Dst;
	
  for(b = bytes & -4; b ; b -=4)
		*(Dst32++) = f;

  // Fills the upper unligned bytes
  bytes &= 0x3;
  Dst = (uint8_t*)Dst32;
  while(bytes--) *(Dst++) = c;

  return dst;
}

void *memcpy(void *dst, const void *src, size_t bytes) {
  const uint8_t *Src;
  const uint32_t *Src32;
  uint8_t *Dst;
  uint32_t *Dst32;
  size_t b;
  
  Dst = dst;
  Src = src;

	// Fills the aligned region using faster memory access instruction
	for(Dst32 = (uint32_t*)Dst, b = (bytes & -4), Src32 = (uint32_t*)Src ; b ; b -=4)
		*(Dst32++) = *(Src32++);
	
	bytes &= 0x3;
	Src = (uint8_t*)Src32;
	Dst = (uint8_t*)Dst32;

  // Fills the upper unligned bytes
  while(bytes--) *(Dst++) = *(Src++);
  
  return dst;
}

int printf(const char *fmt, ...) {
  va_list ap;
  int ret;

  va_start(ap, fmt);
  ret = vprintf(fmt, ap);
  va_end(ap);

  return ret;
}

int vprintf(const char *fmt, va_list ap){
  char *s;
  int i,j;
  char buf[30];

  while (*fmt){
    if (*fmt != '%')
      putchar(*fmt++);
    else{
      j = 0;
      switch (*++fmt){
      case 'i':
      case 'd':
	i = va_arg(ap, int);
	itoa(i,buf,10);
	j=0;
	while (buf[j]) putchar(buf[j++]);
	break;
      case 'u':
	i = va_arg(ap, int);
	utoa(i, buf, 10);
	j=0;
	while (buf[j]) putchar(buf[j++]);
	break;
      case 'o':
	i = va_arg(ap, int);
	utoa(i,buf,8);
	j=0;
	while (buf[j]) putchar(buf[j++]);
	break;
      case 'p':
	puts("0x");
      case 'X':
	i = va_arg(ap, int);
	utoa(i,buf,16);
	j=0;
	while (buf[j])
	  islower(buf[j])?putchar(buf[j++] - 0x20):putchar(buf[j++]);
	break;
      case 'x':
	i = va_arg(ap, int);
	itoa(i,buf,16);
	j=0;
	while (buf[j]) putchar(buf[j++]);
	break;
      case 'c':
	putchar(va_arg(ap, int));
	break;
      case 's':
	s = va_arg(ap, char*);
	if (!s) s = "(null)";
	while (*s) putchar(*s++);
	break;
      case '%' :
	putchar('%');
	break;
      }
      fmt++;
    }
  }

  return 0;
}

void panic(const char *fmt, ...) {
  va_list ap;

  va_start(ap, fmt);
  printf("KERNEL PANIC\n");
  vprintf(fmt, ap);
  va_end(ap);
  puts("\nCPU Halted!");
  
  for(;;) halt();
}
