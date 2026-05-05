#ifndef FUJI_WRAPGEN_ABI_H
#define FUJI_WRAPGEN_ABI_H

/*
 * ABI shim for cmd/wrapgen-generated wrapper.c when linked with libfuji_runtime.a.
 * Wrapgen historically emitted Fuji* names from the embedded alternate header; the
 * shipped runtime uses Value / ObjString / NIL_VAL instead.
 */
#include "value.h"
#include "object.h"
#include "fuji_runtime.h"

typedef Value FujiValue;
typedef ObjString FujiString;

#ifndef NULL_VAL
#define NULL_VAL NIL_VAL
#endif

#endif /* FUJI_WRAPGEN_ABI_H */
