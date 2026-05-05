#ifndef FUJI_SHADOW_STACK_H
#define FUJI_SHADOW_STACK_H

#include "value.h"
#include <stdbool.h>

#define FUJI_SHADOW_STACK_MAX 4096

typedef struct {
    Value** slot_ptrs;
    int count;
} FujiShadowFrame;

extern FujiShadowFrame fuji_shadow_stack[FUJI_SHADOW_STACK_MAX];
extern int fuji_shadow_depth;

void fuji_push_frame(Value** slot_ptrs, int count);
void fuji_pop_frame(void);

#endif
