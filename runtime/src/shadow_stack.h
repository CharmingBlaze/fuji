#ifndef FUJI_SHADOW_STACK_H
#define FUJI_SHADOW_STACK_H

#include "value.h"
#include <stdbool.h>

#define FUJI_SHADOW_STACK_INITIAL_CAPACITY 4096
#define FUJI_SHADOW_STACK_MAX_CAPACITY 65536

typedef struct {
    Value** slot_ptrs;
    int count;
} FujiShadowFrame;

extern FujiShadowFrame* fuji_shadow_stack;
extern int fuji_shadow_depth;
extern int fuji_shadow_capacity;

void fuji_push_frame(Value** slot_ptrs, int count);
void fuji_pop_frame(void);
int fuji_get_shadow_depth(void);

#endif
