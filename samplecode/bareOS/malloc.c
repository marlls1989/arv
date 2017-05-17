#include "include/libos.h"

typedef uint32_t Align;

union header {
  struct {
    union header *ptr;
    unsigned int size;
  } s;
  Align x;
};
typedef union header Header;

static Header base = {.s = {.ptr = &base, .size = 0 }};
static Header *freep = &base;

void *malloc(size_t nbytes) {
  Header *p, *prevp;
  Header *moreroce(unsigned);
  unsigned nunits;
  int critical;
  nunits = (nbytes+sizeof(Header)-1)/sizeof(Header) + 1;

#ifdef DEBUG
  printf("malloc: Allocating %d bytes\n", nunits*sizeof(Header));
#endif

  critical = enter_critical();

  prevp = freep;
  for (p = prevp->s.ptr; ; prevp = p, p = p->s.ptr) {
    if (p->s.size >= nunits) {
      if (p->s.size == nunits)
	prevp->s.ptr = p->s.ptr;
      else {
	p->s.size -= nunits;
	p += p->s.size;
	p->s.size = nunits;
      }
      freep = prevp;
      leave_critical(critical);
      return (void *)(p+1);
    }
    if (p == freep) {
      leave_critical(critical);
      return NULL;
    }
  }
}

void free(void *ap) {
  Header *bp, *p;
  int critical;
  bp = (Header *)ap - 1;

#ifdef DEBUG
  printf("free: Freeing %d bytes\n", bp->s.size*sizeof(Header));
#endif

  critical = enter_critical();

  for (p = freep; !(bp > p && bp < p->s.ptr); p = p->s.ptr)
    if (p >= p->s.ptr && (bp > p || bp < p->s.ptr))
      break; /* freed block at start or end of arena */

  if (bp + bp->s.size == p->s.ptr) {
    /* join to upper nbr */
    bp->s.size += p->s.ptr->s.size;
    bp->s.ptr = p->s.ptr->s.ptr;
  } else
    bp->s.ptr = p->s.ptr;
  if (p + p->s.size == bp) {
    /* join to lower nbr */
    p->s.size += bp->s.size;
    p->s.ptr = bp->s.ptr;
  } else
    p->s.ptr = bp;
  freep = p;

  leave_critical(critical);
}

size_t bfree(void* ptr, size_t size) {
  Header *up = ptr;

  if(size < 2*sizeof(Header))
    return 0;
  
  up->s.size = size/sizeof(Header);
  free((void*)(up+1));

  return up->s.size * sizeof(Header);
}

void *calloc(size_t qty, size_t type_size){
  unsigned int *buf, *end, *a;
  unsigned int size;

  size = qty == 1 ? (type_size+3) & -4 : ((qty*type_size)+3) & -4;
  
  if((buf = malloc(size)))
    for(end = (unsigned int*)((unsigned int)buf + size), a = buf ; a < end ; a++)
      *a = 0;
  
  return (void *)buf;
}
