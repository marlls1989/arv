#include <libos.h>

#define NUM_RINGS 4
#define legal(f, t) ((table[t][top[t]] > table[f][top[f]]) && ((table[t][top[t]] & 1) ^ (table[f][top[f]] & 1)))

static unsigned int table[3][NUM_RINGS+1], top[3];

void move(unsigned int f, unsigned int t) {
	static unsigned int count = 0;

	printf("%d: Move from %d to %d\n", count++, f, t);
	table[t][++top[t]] = table[f][top[f]--];
}

void init_table() {
  unsigned int i;
  
  for (i = 0 ; i < 3 ; i++) {
    table[i][0] = 0xFFFFFFFF;
    top[i] = 0;
  }
  
  for(i = NUM_RINGS ; i ; i--)
    table[0][++top[0]] = i;
}

void hanoi_resolv() {
  unsigned int last, next;

  next = 0;
  last = NUM_RINGS & 1 ? 2 : 1;

  move(next, last);

  while((top[0]+top[1]) && (top[0]+top[2])) {
    switch(last) {
    default:
    case 0:
      next = (table[1][top[1]] < table[2][top[2]]) ? 1 : 2;
      break;
    case 1:
      next = table[0][top[0]] < table[2][top[2]] ? 0 : 2;
      break;
    case 2:
      next = table[0][top[0]] < table[1][top[1]] ? 0 : 1;
      break;
    }

    last = legal(next, last) ? last : 3-next-last;

    move(next, last);
  }
}

void kmain() {
    puts("init Hanoi\n");
    init_table();
    puts("Resolving Hanoi\n");
    hanoi_resolv();
    puts("Solved Hanoi :)\n");
		halt();
}
