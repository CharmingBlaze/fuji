#include "value.h"
#include "object.h"
#include <stdio.h>
#include <string.h>

bool values_equal(Value a, Value b) {
    if (IS_NUMBER(a) && IS_NUMBER(b)) {
        return AS_NUMBER(a) == AS_NUMBER(b);
    }
    return a == b;
}

void print_value(Value v) {
    if (IS_NIL(v)) {
        printf("nil");
    } else if (IS_FALSE(v)) {
        printf("false");
    } else if (IS_TRUE(v)) {
        printf("true");
    } else if (IS_NUMBER(v)) {
        printf("%g", AS_NUMBER(v));
    } else if (IS_OBJ(v)) {
        Obj* o = AS_OBJ(v);
        if (o->type == OBJ_STRING) {
            ObjString* s = (ObjString*)o;
            printf("%.*s", s->length, s->chars);
        } else {
            printf("[object]");
        }
    } else {
        printf("?");
    }
}
