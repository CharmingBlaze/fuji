#ifndef FUJI_H
#define FUJI_H

#include <stdint.h>
#include <stdbool.h>
#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <time.h>

/**
 * KUJI RUNTIME HEADER
 * -------------------
 * Optimized for performance using NaN-boxing and a generational GC.
 */

typedef uint64_t FujiValue;

#define QNAN     ((uint64_t)0x7ffc000000000000)
#define SIGN_BIT ((uint64_t)0x8000000000000000)

#define TAG_NULL   1
#define TAG_FALSE  2
#define TAG_TRUE   3

// --- Value Encoding/Decoding ---

static inline double value_to_num(FujiValue v) {
    union { uint64_t u; double d; } cast;
    cast.u = v;
    return cast.d;
}

static inline FujiValue num_to_value(double n) {
    union { uint64_t u; double d; } cast;
    cast.d = n;
    return cast.u;
}

#define IS_NUMBER(v) (((v) & QNAN) != QNAN)
#define IS_NULL(v)   ((v) == (QNAN | TAG_NULL))
#define IS_BOOL(v)   (((v) | 1) == (QNAN | TAG_TRUE))
#define IS_OBJ(v)    (((v) & (QNAN | SIGN_BIT)) == (QNAN | SIGN_BIT))

#define AS_NUMBER(v) value_to_num(v)
#define AS_BOOL(v)   ((v) == (QNAN | TAG_TRUE))
#define AS_OBJ(v)    ((struct FujiObj*)(uintptr_t)((v) & ~(SIGN_BIT | QNAN)))
#define AS_CLOSURE(v) ((FujiClosure*)AS_OBJ(v))

#define NUMBER_VAL(n) num_to_value(n)
#define BOOL_VAL(b)   ((b) ? (QNAN | TAG_TRUE) : (QNAN | TAG_FALSE))
#define NULL_VAL      (QNAN | TAG_NULL)
#define OBJ_VAL(o)    (FujiValue)(SIGN_BIT | QNAN | (uintptr_t)(o))

// --- Object System ---

typedef enum {
    OBJ_STRING,
    OBJ_ARRAY,
    OBJ_OBJECT,
    OBJ_MAP,
    OBJ_SET,
    OBJ_TUPLE,
    OBJ_CLOSURE,
    OBJ_UPVALUE,
    OBJ_NATIVE,
    OBJ_BOUND_METHOD,
    OBJ_FFI,
} FujiObjType;

typedef struct FujiObj {
    FujiObjType type;
    uint8_t generation; // 0=young, 1=old
    uint8_t mark;       // 0=white, 1=gray/black
    uint8_t age;        // for promotion
    struct FujiObj* next;
} FujiObj;

typedef struct {
    FujiObj header;
    int length;
    uint32_t hash;
    char chars[];
} FujiString;

typedef struct {
    FujiObj header;
    int count;
    int capacity;
    FujiValue* elements;
} FujiArray;

typedef struct FujiUpvalue {
    FujiObj header;
    FujiValue* location;
    FujiValue closed;
    struct FujiUpvalue* next;
} FujiUpvalue;

typedef FujiValue (*FujiFn)(FujiValue thisVal, int argCount, FujiValue* args, FujiUpvalue** upvalues);

typedef struct {
    FujiObj header;
    FujiFn fn;
    FujiUpvalue** upvalues;
    int upvalueCount;
} FujiClosure;

// FFI wrapper
typedef void (*Finalizer)(void*);
typedef struct {
    FujiObj header;
    void* ptr;
    Finalizer finalizer;
    int ref_count;
} FujiFFI;
typedef FujiValue (*FujiNativeFn)(int argCount, FujiValue* args);

typedef struct {
    FujiObj header;
    FujiNativeFn fn;
} FujiNative;

typedef struct {
    FujiValue key;
    FujiValue value;
} FujiEntry;

typedef struct {
    FujiObj header;
    int count;
    int capacity;
    FujiEntry* entries;
} FujiObject;

typedef struct {
    FujiObj header;
    FujiValue receiver;
    FujiValue method; // Can be a Closure or NativeFn
} FujiBoundMethod;

// --- GC & Memory ---

void fuji_init(void);
void fuji_shutdown(void);
void fuji_runtime_init(void);
void fuji_runtime_shutdown(void);
FujiObj* fuji_alloc(FujiObjType type, size_t size);
FujiString* fuji_copy_string(const char* chars, int length);
/** Boxed string for compiler string literals (length, UTF-8 bytes). */
FujiValue fuji_allocate_string(int length, const char* chars);
FujiArray* fuji_new_array();
FujiObject* fuji_alloc_object();
FujiUpvalue* fuji_new_upvalue(FujiValue* slot);
FujiClosure* fuji_new_closure(FujiFn fn, int upvalueCount);

void fuji_upvalue_set(FujiUpvalue* up, FujiValue val);
FujiValue fuji_upvalue_get(FujiUpvalue* up);
void fuji_closure_set_upvalue(FujiValue closure, int index, FujiUpvalue* u);
FujiValue fuji_obj_as_value(void* obj);

/** Fatal arity mismatch (stderr + exit). kind: 0=min, 1=max, 2=exact, 3=missing default index (b unused). */
void fuji_assert(FujiValue cond, FujiValue msg);
void fuji_abort_arg_error(int kind, int a, int b);

/** Build a Fuji array from argv[start : start+count). count may be 0 (argv may be NULL). */
FujiValue fuji_argv_slice_to_array(FujiValue* argv, int start, int count);

void fuji_array_push(FujiValue array, FujiValue value);
void fuji_object_set(FujiValue obj, FujiValue key, FujiValue value);
/** Remove entry whose key compares equal; returns true if a key was removed. */
FujiValue fuji_object_delete(FujiValue obj, FujiValue key);
FujiValue fuji_map_new(int argCount, FujiValue* args);
FujiValue fuji_set_new(int argCount, FujiValue* args);
FujiValue fuji_tuple_new(int argCount, FujiValue* args);
FujiValue fuji_input(int argCount, FujiValue* args);
FujiValue fuji_get_index(FujiValue obj, FujiValue index);
FujiValue fuji_call(FujiValue callee, int argCount, FujiValue* args);

FujiValue fuji_len(FujiValue val);
FujiValue fuji_type(FujiValue val);
FujiValue fuji_clock();
/** Wall-clock seconds since Unix epoch (maps to `time()` in the language reference). */
FujiValue fuji_wall_time(void);
/** Sleep for a wall duration in milliseconds (NaN-boxed number). */
void fuji_sleep_ms(FujiValue ms);
FujiValue fuji_abs_num(FujiValue v);
FujiValue fuji_sqrt_num(FujiValue v);
FujiValue fuji_cbrt_num(FujiValue v);
FujiValue fuji_sin_num(FujiValue v);
FujiValue fuji_cos_num(FujiValue v);
FujiValue fuji_tan_num(FujiValue v);
FujiValue fuji_asin_num(FujiValue v);
FujiValue fuji_acos_num(FujiValue v);
FujiValue fuji_atan_num(FujiValue v);
FujiValue fuji_atan2_num(FujiValue y, FujiValue x);
FujiValue fuji_log_num(FujiValue v);
FujiValue fuji_log2_num(FujiValue v);
FujiValue fuji_log10_num(FujiValue v);
FujiValue fuji_exp_num(FujiValue v);
FujiValue fuji_floor_num(FujiValue v);
FujiValue fuji_ceil_num(FujiValue v);
FujiValue fuji_round_num(FujiValue v);
FujiValue fuji_trunc_num(FujiValue v);
FujiValue fuji_min_num(int argCount, FujiValue* args);
FujiValue fuji_max_num(int argCount, FujiValue* args);
FujiValue fuji_clamp_num(FujiValue v, FujiValue min, FujiValue max);
FujiValue fuji_lerp_num(FujiValue a, FujiValue b, FujiValue t);
FujiValue fuji_smoothstep_num(FujiValue a, FujiValue b, FujiValue t);
FujiValue fuji_map_num(FujiValue v, FujiValue inMin, FujiValue inMax, FujiValue outMin, FujiValue outMax);
FujiValue fuji_sign_num(FujiValue v);
FujiValue fuji_hypot_num(FujiValue a, FujiValue b);
FujiValue fuji_distance_num(FujiValue x1, FujiValue y1, FujiValue x2, FujiValue y2);
FujiValue fuji_distance_sq_num(FujiValue x1, FujiValue y1, FujiValue x2, FujiValue y2);
FujiValue fuji_angle_between_num(FujiValue x1, FujiValue y1, FujiValue x2, FujiValue y2);
FujiValue fuji_normalize_num(FujiValue x, FujiValue y);

/** Time functions */
FujiValue fuji_delta_time(void);
FujiValue fuji_time(void);
FujiValue fuji_timestamp(void);
void fuji_sleep_ms(FujiValue ms);

/** Random functions */
FujiValue fuji_random(int argCount, FujiValue* args);
FujiValue fuji_random_int(int argCount, FujiValue* args);
FujiValue fuji_random_choice(FujiValue array);
FujiValue fuji_set_timeout(FujiValue closure, FujiValue ms);
FujiValue fuji_set_interval(FujiValue closure, FujiValue ms);
void fuji_poll_tasks(void);

FujiValue fuji_io_read_file(int argCount, FujiValue* args);
FujiValue fuji_io_write_file(int argCount, FujiValue* args);
FujiValue fuji_module_init_io();
FujiValue fuji_json_parse(int argCount, FujiValue* args);
FujiValue fuji_json_stringify(int argCount, FujiValue* args);
FujiValue fuji_module_init_json();

void fuji_push(FujiValue value);
FujiValue fuji_pop();
void fuji_collect();

FujiValue fuji_add(FujiValue a, FujiValue b);
FujiValue fuji_to_string_val(FujiValue v);
/** Parse a number from a string, or pass through numeric values; else null. */
FujiValue fuji_parse_number(FujiValue v);
FujiValue fuji_range(FujiValue from, FujiValue to);
FujiValue fuji_sub(FujiValue a, FujiValue b);
FujiValue fuji_mul(FujiValue a, FujiValue b);
FujiValue fuji_div(FujiValue a, FujiValue b);
FujiValue fuji_print(FujiValue val);
void fuji_print_no_newline(FujiValue val);
void fuji_print_space(void);
void fuji_print_newline(void);

FujiValue fuji_is_truthy(FujiValue v);
FujiValue fuji_gt(FujiValue a, FujiValue b);
FujiValue fuji_lt(FujiValue a, FujiValue b);
FujiValue fuji_eq(FujiValue a, FujiValue b);
FujiValue fuji_neq(FujiValue a, FujiValue b);
FujiValue fuji_le(FujiValue a, FujiValue b);
FujiValue fuji_ge(FujiValue a, FujiValue b);
FujiValue fuji_mod(FujiValue a, FujiValue b);
FujiValue fuji_pow(FujiValue a, FujiValue b);
FujiValue fuji_bit_and(FujiValue a, FujiValue b);
FujiValue fuji_bit_or(FujiValue a, FujiValue b);
FujiValue fuji_bit_xor(FujiValue a, FujiValue b);
FujiValue fuji_bit_not(FujiValue a);
FujiValue fuji_shl(FujiValue a, FujiValue b);
FujiValue fuji_shr(FujiValue a, FujiValue b);
FujiValue fuji_ushr(FujiValue a, FujiValue b);
void fuji_set_index(FujiValue obj, FujiValue index, FujiValue value);
int fuji_for_in_len(FujiValue iterable);
FujiValue fuji_for_in_get(FujiValue iterable, int index);
FujiValue fuji_slice(FujiValue obj, FujiValue start, FujiValue end);

FujiValue fuji_negate(FujiValue a);

// Write Barrier
void fuji_write_barrier(FujiObj* obj, FujiValue value);

#define FUJI_WRITE_BARRIER(obj, val) \
    do { if (IS_OBJ(val)) fuji_write_barrier((FujiObj*)(obj), (val)); } while(0)

#endif // FUJI_H
