#ifndef _HAL_H
#define _HAL_H

#include "prototypes.h"

// Stubs, need real implementation
static inline int enter_critical() { return 0; }
static inline void leave_critical(int i) {}
void halt() __attribute__ ((noreturn));

#endif
