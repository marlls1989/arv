#ifndef _LIBOS_H
#define _LIBOS_H

#include "hal.h"

#define isprint(c)		(' '<=(c)&&(c)<='~')
#define isspace(c)		((c)==' '||(c)=='\t'||(c)=='\n'||(c)=='\r')
#define isdigit(c)		('0'<=(c)&&(c)<='9')
#define islower(c)		('a'<=(c)&&(c)<='z')
#define isupper(c)		('A'<=(c)&&(c)<='Z')
#define isalpha(c)		(islower(c)||isupper(c))
#define isalnum(c)		(isalpha(c)||isdigit(c))
#define min(a,b)		((a)<(b)?(a):(b))
#define ntohs(A)		(((A)>>8) | (((A)&0xff)<<8))
#define ntohl(A)		(((A)>>24) | (((A)&0xff0000)>>8) | (((A)&0xff00)<<8) | ((A)<<24))

#define va_start(ap, argN) __builtin_va_start(ap, argN)
#define va_arg(ap, type) __builtin_va_arg(ap, type)
#define va_end(ap) __builtin_va_end(ap)
#define va_copy(dst, src) __builtin_va_copy(dst, src)

#define new(type) calloc(1, sizeof(type))

#define ASSERT_CONCAT_(a, b) a##b
#define ASSERT_CONCAT(a, b) ASSERT_CONCAT_(a, b)
#define ct_assert(e) enum { ASSERT_CONCAT(assert_line_, __LINE__) = 1/(!!(e)) }

typedef __builtin_va_list __gnuc_va_list;
typedef __gnuc_va_list va_list;

int vprintf(const char *fmt, va_list ap);
int printf(const char *fmt, ...)
  __attribute__((format(printf, 1, 2)));

void putchar(char);
void puts(const char*);
char *utoa(unsigned int i, char *s, int base);
char *itoa(int i, char *s, int base);
void *memset(void *dst, int c, size_t bytes);
void *memcpy(void *dst, const void *src, size_t bytes);
void *malloc(size_t size);
void *calloc(size_t qty, size_t type_size);
void free(void *ap);
size_t bfree(void *ptr, size_t size);

void panic(const char *fmt, ...)
  __attribute__ ((noreturn))
  __attribute__((format(printf, 1, 2)));

#endif /* !_LIBOS_H */
