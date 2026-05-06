#include "fuji_runtime.h"
#include "value.h"
#include "object.h"
#include "gc.h"
#include "shadow_stack.h"
#include <stdbool.h>
#include <stddef.h>
#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <time.h>
#include <math.h>
#include <stdint.h>
#include <errno.h>
#include <ctype.h>

// Platform-specific sleep and monotonic clock
#ifdef _WIN32
#include <synchapi.h>
#include <windows.h>
#else
#include <unistd.h>
#endif

#ifdef _WIN32
static double fuji_monotonic_seconds(void) {
    return (double)GetTickCount64() * 0.001;
}
#else
static double fuji_monotonic_seconds(void) {
    struct timespec ts;
    if (clock_gettime(CLOCK_MONOTONIC, &ts) != 0) {
        return (double)time(NULL);
    }
    return (double)ts.tv_sec + (double)ts.tv_nsec * 1e-9;
}
#endif

void* gc_stack_base = NULL;

Value fuji_bool(int argc, Value* args);

FujiShadowFrame* fuji_shadow_stack = NULL;
int fuji_shadow_depth = 0;
int fuji_shadow_capacity = 0;
static int fuji_shadow_depth_high_water = 0;

typedef struct {
    ObjString** entries;
    int count;
    int capacity;
} FujiStringInternTable;

static FujiStringInternTable fuji_string_intern = { NULL, 0, 0 };

static int fuji_gc_debug_enabled(void) {
    const char* env = getenv("FUJI_GC_DEBUG");
    return env != NULL && env[0] != '\0' && env[0] != '0';
}

static int fuji_stack_base_plausible(void* stack_base) {
    if (stack_base == NULL) {
        return 0;
    }
    uintptr_t a = (uintptr_t)stack_base;
    uintptr_t b = (uintptr_t)&a;
    uintptr_t dist = (a > b) ? (a - b) : (b - a);
    return dist <= ((uintptr_t)1u << 30);
}

static void fuji_shadow_stack_init(void) {
    if (fuji_shadow_stack != NULL) {
        return;
    }
    fuji_shadow_capacity = FUJI_SHADOW_STACK_INITIAL_CAPACITY;
    fuji_shadow_stack = (FujiShadowFrame*)malloc(sizeof(FujiShadowFrame) * (size_t)fuji_shadow_capacity);
    if (fuji_shadow_stack == NULL) {
        fuji_panic_str("out of memory allocating shadow stack");
    }
}

void fuji_push_frame(Value** slot_ptrs, int count) {
    fuji_shadow_stack_init();
    if (fuji_shadow_depth >= fuji_shadow_capacity) {
        if (fuji_shadow_capacity >= FUJI_SHADOW_STACK_MAX_CAPACITY) {
            fuji_panic_str("stack overflow — maximum recursion depth reached");
        }
        int next_capacity = fuji_shadow_capacity * 2;
        if (next_capacity > FUJI_SHADOW_STACK_MAX_CAPACITY) {
            next_capacity = FUJI_SHADOW_STACK_MAX_CAPACITY;
        }
        FujiShadowFrame* next = (FujiShadowFrame*)realloc(fuji_shadow_stack, sizeof(FujiShadowFrame) * (size_t)next_capacity);
        if (next == NULL) {
            fuji_panic_str("out of memory growing shadow stack");
        }
        fuji_shadow_stack = next;
        fuji_shadow_capacity = next_capacity;
    }
    fuji_shadow_stack[fuji_shadow_depth].slot_ptrs = slot_ptrs;
    fuji_shadow_stack[fuji_shadow_depth].count = count;
    fuji_shadow_depth++;
    if (fuji_shadow_depth > fuji_shadow_depth_high_water) {
        fuji_shadow_depth_high_water = fuji_shadow_depth;
    }
}

void fuji_pop_frame(void) {
    if (fuji_shadow_depth <= 0) {
        fprintf(stderr, "fuji: shadow stack underflow (push/pop mismatch)\n");
        abort();
    }
    fuji_shadow_depth--;
}

int fuji_get_shadow_depth(void) {
    return fuji_shadow_depth;
}

int fuji_shadow_stack_high_water(void) {
    return fuji_shadow_depth_high_water;
}

#define FUJI_CALL_STACK_MAX 256

typedef struct {
    const char* function_name;
    const char* file_name;
    int line;
} FujiCallFrame;

static FujiCallFrame fuji_call_stack[FUJI_CALL_STACK_MAX];
static int fuji_call_stack_depth = 0;

static ObjString* fuji_intern_find(const char* chars, int length) {
    if (chars == NULL || length < 0) {
        return NULL;
    }
    for (int i = 0; i < fuji_string_intern.count; i++) {
        ObjString* s = fuji_string_intern.entries[i];
        if (s == NULL) {
            continue;
        }
        if (s->length == length && memcmp(s->chars, chars, (size_t)length) == 0) {
            return s;
        }
    }
    return NULL;
}

static void fuji_intern_add(ObjString* str) {
    if (str == NULL) {
        return;
    }
    if (fuji_string_intern.count >= fuji_string_intern.capacity) {
        int next_capacity = fuji_string_intern.capacity == 0 ? 256 : fuji_string_intern.capacity * 2;
        ObjString** next = (ObjString**)realloc(fuji_string_intern.entries, sizeof(ObjString*) * (size_t)next_capacity);
        if (next == NULL) {
            fuji_panic_str("out of memory growing string intern table");
        }
        fuji_string_intern.entries = next;
        fuji_string_intern.capacity = next_capacity;
    }
    fuji_string_intern.entries[fuji_string_intern.count++] = str;
}

void fuji_sweep_intern_table(void) {
    int w = 0;
    for (int r = 0; r < fuji_string_intern.count; r++) {
        ObjString* s = fuji_string_intern.entries[r];
        if (s == NULL) {
            continue;
        }
        /* Keep only strings that the mark phase proved reachable from real roots. */
        if (s->obj.is_marked) {
            fuji_string_intern.entries[w++] = s;
        }
    }
    fuji_string_intern.count = w;
}

void fuji_push_call(const char* fn_name, const char* file_name, int line) {
    if (fuji_call_stack_depth < FUJI_CALL_STACK_MAX) {
        fuji_call_stack[fuji_call_stack_depth].function_name = fn_name ? fn_name : "?";
        fuji_call_stack[fuji_call_stack_depth].file_name = file_name ? file_name : "?";
        fuji_call_stack[fuji_call_stack_depth].line = line;
        fuji_call_stack_depth++;
    }
}

void fuji_pop_call(void) {
    if (fuji_call_stack_depth > 0) {
        fuji_call_stack_depth--;
    }
}

void fuji_print_stack_trace(void) {
    fprintf(stderr, "\nstack trace:\n");
    for (int i = fuji_call_stack_depth - 1; i >= 0; i--) {
        fprintf(stderr, "  at %s (%s:%d)\n",
            fuji_call_stack[i].function_name,
            fuji_call_stack[i].file_name,
            fuji_call_stack[i].line);
    }
    fflush(stderr);
}

static void fuji_print_value_stderr(Value v) {
    if (IS_NIL(v)) {
        fprintf(stderr, "nil");
    } else if (IS_FALSE(v)) {
        fprintf(stderr, "false");
    } else if (IS_TRUE(v)) {
        fprintf(stderr, "true");
    } else if (IS_NUMBER(v)) {
        fprintf(stderr, "%g", AS_NUMBER(v));
    } else if (IS_OBJ(v)) {
        Obj* o = AS_OBJ(v);
        if (o->type == OBJ_STRING) {
            ObjString* s = (ObjString*)o;
            fprintf(stderr, "%.*s", s->length, s->chars);
        } else {
            fprintf(stderr, "[object]");
        }
    } else {
        fprintf(stderr, "?");
    }
}

Value fuji_ok(int argc, Value* argv) {
    Value val = (argc > 0) ? argv[0] : NIL_VAL;
    Value objv = fuji_allocate_object(3);
    Value key_ok = fuji_copy_string("ok", 2);
    Value key_value = fuji_copy_string("value", 5);
    Value key_error = fuji_copy_string("error", 5);
    fuji_object_set(objv, key_ok, TRUE_VAL);
    fuji_object_set(objv, key_value, val);
    fuji_object_set(objv, key_error, NIL_VAL);
    return objv;
}

Value fuji_err(int argc, Value* argv) {
    Value msg = (argc > 0) ? argv[0] : fuji_copy_string("error", 5);
    Value objv = fuji_allocate_object(3);
    Value key_ok = fuji_copy_string("ok", 2);
    Value key_value = fuji_copy_string("value", 5);
    Value key_error = fuji_copy_string("error", 5);
    fuji_object_set(objv, key_ok, FALSE_VAL);
    fuji_object_set(objv, key_value, NIL_VAL);
    fuji_object_set(objv, key_error, msg);
    return objv;
}

Value fuji_err_str(const char* msg) {
    if (msg == NULL) {
        msg = "error";
    }
    Value v = fuji_copy_string(msg, (int)strlen(msg));
    Value args[1] = { v };
    return fuji_err(1, args);
}

void fuji_panic_str(const char* msg) {
    if (msg == NULL) {
        msg = "(null message)";
    }
    Value v = fuji_copy_string(msg, (int)strlen(msg));
    Value args[1] = { v };
    fuji_panic(1, args);
}

void fuji_panic(int argc, Value* argv) {
    fprintf(stderr, "\nfuji panic: ");
    if (argc > 0) {
        fuji_print_value_stderr(argv[0]);
    } else {
        fprintf(stderr, "(no message)");
    }
    fprintf(stderr, "\n");
    fuji_print_stack_trace();
    fflush(stderr);
    exit(1);
}

Value fuji_assert(int argc, Value* argv) {
    if (argc < 1) {
        fuji_panic_str("assert() requires at least one argument");
    }
    Value b = fuji_bool(1, argv);
    if (IS_TRUE(b)) {
        return NIL_VAL;
    }
    if (argc >= 2) {
        Value pargs[1] = { argv[1] };
        fuji_panic(1, pargs);
    }
    fuji_panic_str("assertion failed");
    return NIL_VAL;
}

void fuji_gc_use_shadow_stack(bool enable) {
    gc_set_use_shadow_stack(enable);
}

Value* fuji_globals = NULL;
int fuji_globals_count = 0;
int fuji_globals_capacity = 0;
Value** fuji_global_slots = NULL;
int fuji_global_slots_count = 0;
int fuji_global_slots_capacity = 0;

void fuji_mark_module_cache(void) {}

void fuji_mark_open_upvalues(void) {}

static int fuji_frame_clock_inited = 0;
static double fuji_frame_clock_start = 0;
static double fuji_frame_clock_last = 0;

static void fuji_frame_clock_ensure_init(void) {
    if (fuji_frame_clock_inited) {
        return;
    }
    double now = fuji_monotonic_seconds();
    fuji_frame_clock_start = now;
    fuji_frame_clock_last = now;
    fuji_frame_clock_inited = 1;
}

void fuji_runtime_set_stack_base(void* base) {
    gc_stack_base = base;
}

void fuji_globals_init(void) {
    if (fuji_globals != NULL && fuji_globals_capacity > 0) {
        return;
    }
    fuji_globals_capacity = 256;
    fuji_globals_count = 0;
    fuji_globals = (Value*)malloc(sizeof(Value) * (size_t)fuji_globals_capacity);
    if (fuji_globals == NULL) {
        fprintf(stderr, "fuji: out of memory allocating globals array\n");
        exit(1);
    }
    fuji_global_slots_capacity = 256;
    fuji_global_slots_count = 0;
    fuji_global_slots = (Value**)malloc(sizeof(Value*) * (size_t)fuji_global_slots_capacity);
    if (fuji_global_slots == NULL) {
        fprintf(stderr, "fuji: out of memory allocating global slots array\n");
        exit(1);
    }
}

void fuji_globals_free(void) {
    if (fuji_globals != NULL) {
        free(fuji_globals);
        fuji_globals = NULL;
    }
    if (fuji_global_slots != NULL) {
        free(fuji_global_slots);
        fuji_global_slots = NULL;
    }
    fuji_globals_count = 0;
    fuji_globals_capacity = 0;
    fuji_global_slots_count = 0;
    fuji_global_slots_capacity = 0;
}

void fuji_register_global(Value v) {
    if (fuji_globals == NULL || fuji_globals_capacity <= 0) {
        fuji_globals_init();
    }
    if (fuji_globals_count >= fuji_globals_capacity) {
        int new_cap = fuji_globals_capacity * 2;
        if (new_cap < 256) {
            new_cap = 256;
        }
        Value* next = (Value*)realloc(fuji_globals, sizeof(Value) * (size_t)new_cap);
        if (next == NULL) {
            fprintf(stderr, "fuji: out of memory growing globals array\n");
            exit(1);
        }
        fuji_globals = next;
        fuji_globals_capacity = new_cap;
    }
    fuji_globals[fuji_globals_count++] = v;
}

void fuji_register_global_slot(Value* slot) {
    if (slot == NULL) {
        return;
    }
    if (fuji_global_slots == NULL || fuji_global_slots_capacity <= 0) {
        fuji_globals_init();
    }
    if (fuji_global_slots_count >= fuji_global_slots_capacity) {
        int new_cap = fuji_global_slots_capacity * 2;
        if (new_cap < 256) {
            new_cap = 256;
        }
        Value** next = (Value**)realloc(fuji_global_slots, sizeof(Value*) * (size_t)new_cap);
        if (next == NULL) {
            fprintf(stderr, "fuji: out of memory growing global slots array\n");
            exit(1);
        }
        fuji_global_slots = next;
        fuji_global_slots_capacity = new_cap;
    }
    fuji_global_slots[fuji_global_slots_count++] = slot;
}

void fuji_runtime_init_ex(void* stack_base) {
    if (!fuji_stack_base_plausible(stack_base)) {
        fprintf(stderr, "fuji: invalid stack base passed to fuji_runtime_init_ex\n");
        abort();
    }
    gc_stack_base = stack_base;
    gc_init();
    fuji_globals_init();
    fuji_shadow_stack_init();
    fuji_gc_use_shadow_stack(true);
}

void fuji_runtime_init(void) {
#if defined(__clang__) || defined(__GNUC__)
    fuji_runtime_init_ex(__builtin_frame_address(0));
#else
    uintptr_t anchor = 0;
    fuji_runtime_init_ex(&anchor);
#endif
}

void fuji_gc_set_threshold(size_t bytes) {
    gc_set_next_threshold(bytes);
}

void fuji_gc_disable(void) {
    gc_set_disabled(true);
}

void fuji_gc_enable(void) {
    gc_set_disabled(false);
}

void fuji_gc_collect(void) {
    gc_collect();
}

void fuji_gc_frame_step(double budget_ms) {
    if (gc_incremental_is_idle()) {
        size_t used = gc_nursery_used_bytes();
        size_t half = gc_nursery_capacity_bytes() / 2u;
        if (used >= half) {
            gc_collect_minor();
        }
    }
    if (budget_ms > 0.0) {
        double ms = budget_ms;
        if (ms < 0.0) {
            ms = 0.0;
        }
        uint64_t budget_us = (uint64_t)(ms * 1000.0 + 0.999);
        if (budget_us < 64u) {
            budget_us = 64u;
        }
        gc_frame_step_incremental(budget_us);
    }
}

Value fuji_gc_stats(int argc, Value* argv) {
    (void)argc;
    (void)argv;
    GCStats st = gc_get_stats();
    Value objv = fuji_allocate_object(5);
    /* Property keys are ASCII lowercase: dot access (e.g. s.bytesAllocated) compiles to s["bytesallocated"]. */
    fuji_object_set(objv, fuji_copy_string("collections", 11), NUMBER_VAL((double)st.collections));
    fuji_object_set(objv, fuji_copy_string("bytesallocated", 14), NUMBER_VAL((double)st.bytes_allocated));
    fuji_object_set(objv, fuji_copy_string("bytesfreed", 10), NUMBER_VAL((double)st.bytes_freed));
    fuji_object_set(objv, fuji_copy_string("maxpausetimeus", 14), NUMBER_VAL((double)st.max_pause_time_us));
    fuji_object_set(objv, fuji_copy_string("totalpausetimeus", 16), NUMBER_VAL((double)st.total_pause_time_us));
    return objv;
}

Value* fuji_alloc_cell(void) {
    ObjCell* c = allocate_cell();
    return &c->value;
}

Value fuji_cell_read(Value* cell) {
    return *cell;
}

void fuji_cell_write(Value* cell, Value v) {
    ObjCell* cellObj = (ObjCell*)((uint8_t*)cell - offsetof(ObjCell, value));
    gc_write_barrier((Obj*)&cellObj->obj, v);
    *cell = v;
}

void fuji_runtime_shutdown(void) {
    gc_collect();
    if (fuji_gc_debug_enabled()) {
        fprintf(stderr,
            "fuji gc debug: remembered_overflow=%llu shadow_depth_hwm=%d globals=%d/%d global_slots=%d/%d\n",
            (unsigned long long)gc_debug_remembered_overflow_count(),
            fuji_shadow_depth_high_water,
            fuji_globals_count,
            fuji_globals_capacity,
            fuji_global_slots_count,
            fuji_global_slots_capacity);
    }
    fuji_globals_free();
    if (fuji_string_intern.entries != NULL) {
        free(fuji_string_intern.entries);
        fuji_string_intern.entries = NULL;
    }
    fuji_string_intern.count = 0;
    fuji_string_intern.capacity = 0;
    if (fuji_shadow_stack != NULL) {
        free(fuji_shadow_stack);
        fuji_shadow_stack = NULL;
    }
    fuji_shadow_capacity = 0;
    fuji_shadow_depth = 0;
}

Value fuji_print_val(Value v) {
    print_value(v);
    printf("\n");
    fflush(stdout);
    return NIL_VAL;
}

void fuji_print_newline(void) {
    printf("\n");
    fflush(stdout);
}

Value fuji_print_argv(int arg_count, Value* args) {
    for (int i = 0; i < arg_count; i++) {
        print_value(args[i]);
        if (i < arg_count - 1) printf(" ");
    }
    printf("\n");
    return NIL_VAL;
}

Value fuji_typeof(int arg_count, Value* args) {
    if (arg_count == 0) return NIL_VAL;
    
    Value v = args[0];
    if (IS_NIL(v)) {
        ObjString* str = allocate_string(3);
        memcpy(str->chars, "nil", 3);
        return OBJ_VAL((Obj*)str);
    } else if (IS_BOOL(v)) {
        ObjString* str = allocate_string(7);
        memcpy(str->chars, "boolean", 7);
        return OBJ_VAL((Obj*)str);
    } else if (IS_NUMBER(v)) {
        ObjString* str = allocate_string(6);
        memcpy(str->chars, "number", 6);
        return OBJ_VAL((Obj*)str);
    } else if (IS_OBJ(v)) {
        Obj* obj = AS_OBJ(v);
        const char* type_str = "object";
        switch (obj->type) {
            case OBJ_STRING: type_str = "string"; break;
            case OBJ_ARRAY: type_str = "array"; break;
            case OBJ_TABLE: type_str = "table"; break;
            case OBJ_FUNCTION: type_str = "function"; break;
            case OBJ_CLOSURE: type_str = "closure"; break;
            case OBJ_NATIVE: type_str = "native"; break;
            case OBJ_CELL: type_str = "cell"; break;
        }
        int len = strlen(type_str);
        ObjString* str = allocate_string(len);
        memcpy(str->chars, type_str, len);
        return OBJ_VAL((Obj*)str);
    }
    
    return NIL_VAL;
}

Value fuji_allocate_object(int property_count) {
    ObjTable* table = allocate_table(property_count);
    return OBJ_VAL((Obj*)table);
}

Value fuji_allocate_struct(int field_count) {
    ObjTable* table = allocate_struct_table(field_count);
    return OBJ_VAL((Obj*)table);
}

static uint32_t fnv1a_bytes(const uint8_t* data, int len) {
    uint32_t h = 2166136261u;
    for (int i = 0; i < len; i++) {
        h ^= (uint32_t)data[i];
        h *= 16777619u;
    }
    return h == 0u ? 1u : h;
}

static uint32_t fuji_string_hash(ObjString* s) {
    if (s->hash != 0u) {
        return s->hash;
    }
    s->hash = fnv1a_bytes((const uint8_t*)s->chars, s->length);
    if (s->hash == 0u) {
        s->hash = 1u;
    }
    return s->hash;
}

static uint32_t fuji_hash_value(Value key) {
    if (IS_OBJ(key)) {
        Obj* ob = AS_OBJ(key);
        if (ob != NULL && ob->type == OBJ_STRING) {
            return fuji_string_hash((ObjString*)ob);
        }
    }
    if (IS_NUMBER(key)) {
        union {
            Value v;
            uint64_t u;
        } u;
        u.v = key;
        uint32_t hi = (uint32_t)(u.u >> 32);
        uint32_t lo = (uint32_t)u.u;
        uint32_t h = hi ^ lo;
        return h == 0u ? 1u : h;
    }
    return 1u;
}

Value fuji_struct_get(Value obj, int64_t index) {
    if (!IS_OBJ(obj)) {
        return NIL_VAL;
    }
    Obj* o = AS_OBJ(obj);
    if (o->type != OBJ_TABLE) {
        return NIL_VAL;
    }
    ObjTable* t = (ObjTable*)o;
    if (!t->is_struct_layout) {
        return NIL_VAL;
    }
    if (index < 0 || index >= t->count) {
        return NIL_VAL;
    }
    return t->values[(int)index];
}

Value fuji_struct_set(Value obj, int64_t index, Value val) {
    if (!IS_OBJ(obj)) {
        return NIL_VAL;
    }
    Obj* o = AS_OBJ(obj);
    if (o->type != OBJ_TABLE) {
        return NIL_VAL;
    }
    ObjTable* t = (ObjTable*)o;
    if (!t->is_struct_layout) {
        return NIL_VAL;
    }
    if (index < 0 || index >= t->count) {
        return NIL_VAL;
    }
    gc_write_barrier(&t->obj, val);
    t->values[(int)index] = val;
    return val;
}

static Value fuji_object_get_linear(ObjTable* table, Value key) {
    for (int i = 0; i < table->count; i++) {
        if (values_equal(table->keys[i], key)) {
            return table->values[i];
        }
    }
    return NIL_VAL;
}

static bool table_key_is_tombstone(Value k) {
    return IS_BOOL(k) && AS_BOOL(k);
}

static Value fuji_object_get_hash(ObjTable* table, Value key) {
    uint32_t cap = (uint32_t)table->capacity;
    if (cap == 0u) {
        return NIL_VAL;
    }
    uint32_t start = fuji_hash_value(key) % cap;
    for (uint32_t p = 0u; p < cap; p++) {
        uint32_t i = (start + p) % cap;
        Value k = table->keys[i];
        if (IS_NIL(k)) {
            return NIL_VAL;
        }
        if (table_key_is_tombstone(k)) {
            continue;
        }
        if (values_equal(k, key)) {
            return table->values[i];
        }
    }
    return NIL_VAL;
}

Value fuji_object_get(Value obj, Value key) {
    if (IS_NIL(obj)) {
        const char* prop = "?";
        if (IS_OBJ(key) && AS_OBJ(key)->type == OBJ_STRING) {
            prop = ((ObjString*)AS_OBJ(key))->chars;
        }
        char msg[256];
        snprintf(msg, sizeof(msg), "cannot read property '%s' of null", prop);
        fuji_panic_str(msg);
    }
    if (!IS_OBJ(obj)) {
        return NIL_VAL;
    }
    Obj* o = AS_OBJ(obj);
    if (o->type != OBJ_TABLE) {
        return NIL_VAL;
    }

    ObjTable* table = (ObjTable*)o;
    if (table->is_struct_layout) {
        return fuji_object_get_linear(table, key);
    }
    if (table->hashes != NULL) {
        return fuji_object_get_hash(table, key);
    }
    return fuji_object_get_linear(table, key);
}

static bool table_remove_pair(ObjTable* table, Value key) {
    if (table->hashes != NULL && !table->is_struct_layout) {
        uint32_t cap = (uint32_t)table->capacity;
        if (cap == 0u) {
            return false;
        }
        uint32_t start = fuji_hash_value(key) % cap;
        for (uint32_t p = 0u; p < cap; p++) {
            uint32_t i = (start + p) % cap;
            Value k = table->keys[i];
            if (IS_NIL(k)) {
                return false;
            }
            if (table_key_is_tombstone(k)) {
                continue;
            }
            if (values_equal(k, key)) {
                table->keys[i] = TRUE_VAL;
                table->values[i] = NIL_VAL;
                table->hashes[i] = 0u;
                table->count--;
                return true;
            }
        }
        return false;
    }
    for (int i = 0; i < table->count; i++) {
        if (values_equal(table->keys[i], key)) {
            for (int j = i; j < table->count - 1; j++) {
                table->keys[j] = table->keys[j + 1];
                table->values[j] = table->values[j + 1];
            }
            table->count--;
            return true;
        }
    }
    return false;
}

/** Remove own property `key` from table object; returns boxed boolean. */
Value fuji_object_remove(Value obj, Value key) {
    if (!IS_OBJ(obj)) return NIL_VAL;
    Obj* o = AS_OBJ(obj);
    if (o->type != OBJ_TABLE) return NIL_VAL;
    ObjTable* table = (ObjTable*)o;
    bool ok = table_remove_pair(table, key);
    return BOOL_VAL(ok);
}

static void table_rehash_open(Obj* parent, ObjTable* table);

static bool fuji_object_set_hash(Obj* parent, ObjTable* table, Value key, Value value) {
    uint32_t cap = (uint32_t)table->capacity;
    if (cap == 0u) {
        return false;
    }
    if ((table->count + 1) * 10 > (int)cap * 7) {
        table_rehash_open(parent, table);
        cap = (uint32_t)table->capacity;
        if (cap == 0u) {
            return false;
        }
    }
    uint32_t start = fuji_hash_value(key) % cap;
    int first_tomb = -1;
    for (uint32_t p = 0u; p < cap; p++) {
        uint32_t i = (start + p) % cap;
        Value k = table->keys[i];
        if (IS_NIL(k)) {
            int ins = (first_tomb >= 0) ? first_tomb : (int)i;
            gc_write_barrier(parent, key);
            gc_write_barrier(parent, value);
            table->keys[ins] = key;
            table->values[ins] = value;
            table->hashes[ins] = fuji_hash_value(key);
            table->count++;
            return true;
        }
        if (table_key_is_tombstone(k)) {
            if (first_tomb < 0) {
                first_tomb = (int)i;
            }
            continue;
        }
        if (values_equal(k, key)) {
            gc_write_barrier(parent, key);
            gc_write_barrier(parent, value);
            table->keys[i] = key;
            table->values[i] = value;
            return true;
        }
    }
    table_rehash_open(parent, table);
    return fuji_object_set_hash(parent, table, key, value);
}

static void table_rehash_open(Obj* parent, ObjTable* table) {
    int old_cap = table->capacity;
    Value* old_keys = table->keys;
    Value* old_vals = table->values;
    uint32_t* old_hashes = table->hashes;
    int new_cap = old_cap < 8 ? 16 : old_cap * 2;
    if (new_cap < old_cap) {
        return;
    }
    table->capacity = new_cap;
    table->count = 0;
    table->keys = (Value*)gc_alloc(sizeof(Value) * (size_t)new_cap);
    table->values = (Value*)gc_alloc(sizeof(Value) * (size_t)new_cap);
    table->hashes = (uint32_t*)calloc((size_t)new_cap, sizeof(uint32_t));
    for (int i = 0; i < new_cap; i++) {
        table->keys[i] = NIL_VAL;
        table->values[i] = NIL_VAL;
    }
    for (int i = 0; i < old_cap; i++) {
        Value k = old_keys[i];
        if (IS_NIL(k) || table_key_is_tombstone(k)) {
            continue;
        }
        (void)fuji_object_set_hash(parent, table, k, old_vals[i]);
    }
    if (old_keys != NULL) {
        gc_free(old_keys, (size_t)old_cap * sizeof(Value));
    }
    if (old_vals != NULL) {
        gc_free(old_vals, (size_t)old_cap * sizeof(Value));
    }
    if (old_hashes != NULL) {
        gc_free(old_hashes, (size_t)old_cap * sizeof(uint32_t));
    }
}

Value fuji_object_set(Value obj, Value key, Value value) {
    if (!IS_OBJ(obj)) {
        return NIL_VAL;
    }
    Obj* o = AS_OBJ(obj);
    if (o->type != OBJ_TABLE) {
        return NIL_VAL;
    }

    ObjTable* table = (ObjTable*)o;
    if (table->is_struct_layout) {
        for (int i = 0; i < table->count; i++) {
            if (values_equal(table->keys[i], key)) {
                gc_write_barrier(o, key);
                gc_write_barrier(o, value);
                table->keys[i] = key;
                table->values[i] = value;
                return value;
            }
        }
        return NIL_VAL;
    }
    if (table->hashes != NULL) {
        if (fuji_object_set_hash(o, table, key, value)) {
            return value;
        }
        return NIL_VAL;
    }

    for (int i = 0; i < table->count; i++) {
        if (values_equal(table->keys[i], key)) {
            gc_write_barrier(o, key);
            gc_write_barrier(o, value);
            table->keys[i] = key;
            table->values[i] = value;
            return value;
        }
    }

    if (table->count < table->capacity) {
        gc_write_barrier(o, key);
        gc_write_barrier(o, value);
        table->keys[table->count] = key;
        table->values[table->count] = value;
        table->count++;
        return value;
    }

    return NIL_VAL;
}

Value fuji_allocate_string(int length, char* chars) {
    ObjString* interned = fuji_intern_find(chars, length);
    if (interned != NULL) {
        return OBJ_VAL((Obj*)interned);
    }
    ObjString* str = allocate_string(length);
    memcpy(str->chars, chars, length);
    str->chars[length] = '\0';
    fuji_intern_add(str);
    return OBJ_VAL((Obj*)str);
}

Value fuji_copy_string(const char* chars, int length) {
    if (chars == NULL || length < 0) {
        return NIL_VAL;
    }
    ObjString* interned = fuji_intern_find(chars, length);
    if (interned != NULL) {
        return OBJ_VAL((Obj*)interned);
    }
    ObjString* str = allocate_string(length);
    memcpy(str->chars, chars, (size_t)length);
    str->chars[length] = '\0';
    fuji_intern_add(str);
    return OBJ_VAL((Obj*)str);
}

Value fuji_allocate_array(int length) {
    ObjArray* arr = allocate_array(length);
    return OBJ_VAL((Obj*)arr);
}

Value fuji_array_get(Value arr, int index) {
    if (!IS_OBJ(arr)) return NIL_VAL;
    Obj* obj = AS_OBJ(arr);
    if (obj->type != OBJ_ARRAY) return NIL_VAL;
    
    ObjArray* array = (ObjArray*)obj;
    if (index < 0 || index >= array->count) {
        char msg[256];
        snprintf(msg, sizeof(msg), "index %d out of bounds (array length %d)", index, array->count);
        fuji_panic_str(msg);
    }
    
    return array->elements[index];
}

double fuji_unbox_number(Value v) {
    if (!IS_NUMBER(v)) {
        fprintf(stderr, "fuji_unbox_number: value is not a number\n");
        return 0.0;
    }
    return AS_NUMBER(v);
}

Value fuji_box_number(double d) {
    return NUMBER_VAL(d);
}

Value fuji_get(Value obj, Value key) {
    if (!IS_OBJ(obj)) {
        fprintf(stderr, "Cannot index a non-object value\n");
        return NIL_VAL;
    }
    Obj* o = AS_OBJ(obj);
    if (o->type == OBJ_ARRAY) {
        if (!IS_NUMBER(key)) {
            fprintf(stderr, "Array index must be a number\n");
            return NIL_VAL;
        }
        int idx = (int)AS_NUMBER(key);
        return fuji_array_get(obj, idx);
    }
    if (o->type == OBJ_STRING) {
        if (!IS_NUMBER(key)) {
            fprintf(stderr, "String index must be a number\n");
            return NIL_VAL;
        }
        int idx = (int)AS_NUMBER(key);
        ObjString* s = (ObjString*)o;
        if (idx < 0 || idx >= s->length) return NIL_VAL;
        return fuji_copy_string(&s->chars[idx], 1);
    }
    if (o->type == OBJ_TABLE) {
        return fuji_object_get(obj, key);
    }
    fprintf(stderr, "Cannot index this object type\n");
    return NIL_VAL;
}

Value fuji_get_index(Value obj, Value key) {
    return fuji_get(obj, key);
}

void fuji_array_set(Value arr, int64_t index, Value value) {
    if (!IS_OBJ(arr)) return;
    Obj* obj = AS_OBJ(arr);
    if (obj->type != OBJ_ARRAY) return;
    
    ObjArray* array = (ObjArray*)obj;
    if (index < 0 || index >= array->capacity) return;

    gc_write_barrier(obj, value);
    array->elements[(int)index] = value;
    if (index + 1 > array->count) {
        array->count = (int)index + 1;
    }
}

Value fuji_set(Value obj, Value key, Value val) {
    if (!IS_OBJ(obj)) {
        fprintf(stderr, "Cannot index-assign a non-object value\n");
        return NIL_VAL;
    }
    Obj* o = AS_OBJ(obj);
    if (o->type == OBJ_ARRAY) {
        if (!IS_NUMBER(key)) {
            fprintf(stderr, "Array index must be a number\n");
            return NIL_VAL;
        }
        int64_t i = (int64_t)AS_NUMBER(key);
        fuji_array_set(obj, i, val);
        return val;
    }
    if (o->type == OBJ_TABLE) {
        return fuji_object_set(obj, key, val);
    }
    fprintf(stderr, "Cannot assign into this object type\n");
    return NIL_VAL;
}

void fuji_array_push(Value arr, Value value) {
    if (!IS_OBJ(arr)) return;
    Obj* obj = AS_OBJ(arr);
    if (obj->type != OBJ_ARRAY) return;
    
    ObjArray* array = (ObjArray*)obj;
    if (array->count < array->capacity) {
        gc_write_barrier(obj, value);
        array->elements[array->count] = value;
        array->count++;
    }
}

Value fuji_array_pop(Value arr) {
    if (!IS_OBJ(arr)) return NIL_VAL;
    Obj* obj = AS_OBJ(arr);
    if (obj->type != OBJ_ARRAY) return NIL_VAL;
    
    ObjArray* array = (ObjArray*)obj;
    if (array->count > 0) {
        array->count--;
        return array->elements[array->count];
    }
    return NIL_VAL;
}

Value fuji_array_length(Value arr) {
    if (!IS_OBJ(arr)) return NIL_VAL;
    Obj* obj = AS_OBJ(arr);
    if (obj->type != OBJ_ARRAY) return NIL_VAL;
    
    ObjArray* array = (ObjArray*)obj;
    return NUMBER_VAL(array->count);
}

/** Slot iteration length for native `for-of` / `for-in` (arrays + tables). */
Value fuji_forof_length(Value v) {
    if (!IS_OBJ(v)) return NUMBER_VAL(0);
    Obj* o = AS_OBJ(v);
    if (o->type == OBJ_ARRAY) {
        return NUMBER_VAL((double)((ObjArray*)o)->count);
    }
    if (o->type == OBJ_TABLE) {
        return NUMBER_VAL((double)((ObjTable*)o)->count);
    }
    return NUMBER_VAL(0);
}

/** Key at linear slot `idx` (number): array → numeric index; table → insertion-order key. */
Value fuji_forof_key_at(Value v, Value idx_val) {
    if (!IS_OBJ(v) || !IS_NUMBER(idx_val)) return NIL_VAL;
    int i = (int)AS_NUMBER(idx_val);
    Obj* o = AS_OBJ(v);
    if (o->type == OBJ_ARRAY) {
        ObjArray* a = (ObjArray*)o;
        if (i < 0 || i >= a->count) return NIL_VAL;
        return NUMBER_VAL((double)i);
    }
    if (o->type == OBJ_TABLE) {
        ObjTable* t = (ObjTable*)o;
        if (i < 0 || i >= t->count) return NIL_VAL;
        return t->keys[i];
    }
    return NIL_VAL;
}

/** Value at linear slot `idx` (insertion order for tables). */
Value fuji_forof_value_at(Value v, Value idx_val) {
    if (!IS_OBJ(v) || !IS_NUMBER(idx_val)) return NIL_VAL;
    int i = (int)AS_NUMBER(idx_val);
    Obj* o = AS_OBJ(v);
    if (o->type == OBJ_ARRAY) {
        ObjArray* a = (ObjArray*)o;
        if (i < 0 || i >= a->count) return NIL_VAL;
        return a->elements[i];
    }
    if (o->type == OBJ_TABLE) {
        ObjTable* t = (ObjTable*)o;
        if (i < 0 || i >= t->count) return NIL_VAL;
        return t->values[i];
    }
    return NIL_VAL;
}

// Standard library functions
Value fuji_type(Value value) {
    if (IS_NIL(value)) {
        ObjString* str = allocate_string(3);
        memcpy(str->chars, "nil", 3);
        return OBJ_VAL((Obj*)str);
    } else if (IS_BOOL(value)) {
        ObjString* str = allocate_string(7);
        memcpy(str->chars, "boolean", 7);
        return OBJ_VAL((Obj*)str);
    } else if (IS_NUMBER(value)) {
        ObjString* str = allocate_string(6);
        memcpy(str->chars, "number", 6);
        return OBJ_VAL((Obj*)str);
    } else if (IS_OBJ(value)) {
        Obj* obj = AS_OBJ(value);
        const char* type_str = "object";
        switch (obj->type) {
            case OBJ_STRING: type_str = "string"; break;
            case OBJ_ARRAY: type_str = "array"; break;
            case OBJ_TABLE: type_str = "table"; break;
            case OBJ_FUNCTION: type_str = "function"; break;
            case OBJ_CLOSURE: type_str = "closure"; break;
            case OBJ_NATIVE: type_str = "native"; break;
            case OBJ_CELL: type_str = "cell"; break;
        }
        int len = strlen(type_str);
        ObjString* str = allocate_string(len);
        memcpy(str->chars, type_str, len);
        return OBJ_VAL((Obj*)str);
    }
    return NIL_VAL;
}

Value fuji_len(Value value) {
    if (IS_OBJ(value)) {
        Obj* obj = AS_OBJ(value);
        switch (obj->type) {
            case OBJ_STRING: {
                ObjString* str = (ObjString*)obj;
                return NUMBER_VAL(str->length);
            }
            case OBJ_ARRAY: {
                ObjArray* arr = (ObjArray*)obj;
                return NUMBER_VAL(arr->count);
            }
            case OBJ_TABLE: {
                ObjTable* table = (ObjTable*)obj;
                return NUMBER_VAL((double)table->count);
            }
            default:
                break;
        }
    }
    return NUMBER_VAL(0);
}

Value fuji_abs(Value value) {
    if (!IS_NUMBER(value)) return NIL_VAL;
    double num = AS_NUMBER(value);
    if (num < 0) num = -num;
    return NUMBER_VAL(num);
}

Value fuji_sqrt(Value value) {
    if (!IS_NUMBER(value)) return NIL_VAL;
    double num = AS_NUMBER(value);
    if (num < 0) return NIL_VAL;
    return NUMBER_VAL(sqrt(num));
}

Value fuji_time() {
    return NUMBER_VAL((double)time(NULL));
}

Value fuji_clock(void) {
    return NUMBER_VAL((double)clock() / (double)CLOCKS_PER_SEC);
}

Value fuji_wall_time(void) {
    return NUMBER_VAL((double)time(NULL));
}

void fuji_sleep(int64_t ms) {
    // Sleep implementation (platform-specific)
    #ifdef _WIN32
        Sleep((DWORD)ms);
    #else
        usleep(ms * 1000);
    #endif
}

Value fuji_number(Value value) {
    if (IS_NUMBER(value)) return value;
    if (IS_FALSE(value)) return NUMBER_VAL(0);
    if (IS_TRUE(value)) return NUMBER_VAL(1);
    if (IS_NIL(value)) return NUMBER_VAL(0);
    return NIL_VAL;
}

Value fuji_string(Value value) {
    if (IS_OBJ(value)) {
        Obj* obj = AS_OBJ(value);
        if (obj->type == OBJ_STRING) return value;
    }
    if (IS_NUMBER(value)) {
        char buffer[32];
        snprintf(buffer, 32, "%g", AS_NUMBER(value));
        int len = strlen(buffer);
        ObjString* str = allocate_string(len);
        memcpy(str->chars, buffer, len);
        return OBJ_VAL((Obj*)str);
    }
    if (IS_TRUE(value)) {
        ObjString* str = allocate_string(4);
        memcpy(str->chars, "true", 4);
        return OBJ_VAL((Obj*)str);
    }
    if (IS_FALSE(value)) {
        ObjString* str = allocate_string(5);
        memcpy(str->chars, "false", 5);
        return OBJ_VAL((Obj*)str);
    }
    if (IS_NIL(value)) {
        ObjString* str = allocate_string(3);
        memcpy(str->chars, "nil", 3);
        return OBJ_VAL((Obj*)str);
    }
    return NIL_VAL;
}

Value fuji_string_concat(Value a, Value b) {
    Value sa = fuji_string(a);
    Value sb = fuji_string(b);
    if (!IS_OBJ(sa) || AS_OBJ(sa)->type != OBJ_STRING) return sb;
    if (!IS_OBJ(sb) || AS_OBJ(sb)->type != OBJ_STRING) return sa;
    ObjString* A = (ObjString*)AS_OBJ(sa);
    ObjString* B = (ObjString*)AS_OBJ(sb);
    int len = A->length + B->length;
    ObjString* out = allocate_string(len);
    memcpy(out->chars, A->chars, A->length);
    memcpy(out->chars + A->length, B->chars, B->length);
    return OBJ_VAL((Obj*)out);
}

// Time functions for game development
Value fuji_delta_time(int arg_count, Value* args) {
    (void)arg_count;
    (void)args;
    fuji_frame_clock_ensure_init();
    double now = fuji_monotonic_seconds();
    double dt = now - fuji_frame_clock_last;
    fuji_frame_clock_last = now;
    if (dt <= 0.0 || dt > 0.25) {
        dt = 1.0 / 60.0;
    }
    return NUMBER_VAL(dt);
}

Value fuji_program_time(int arg_count, Value* args) {
    (void)arg_count;
    (void)args;
    fuji_frame_clock_ensure_init();
    double now = fuji_monotonic_seconds();
    return NUMBER_VAL(now - fuji_frame_clock_start);
}

Value fuji_timestamp(int arg_count, Value* args) {
    // Unix timestamp
    return NUMBER_VAL((double)time(NULL));
}

// Random functions for game development
Value fuji_random(int arg_count, Value* args) {
    switch (arg_count) {
    case 0:
        // random() → [0, 1)
        return NUMBER_VAL((double)rand() / RAND_MAX);
    case 1:
        // random(max) → [0, max)
        if (!IS_NUMBER(args[0])) return NIL_VAL;
        double max = AS_NUMBER(args[0]);
        return NUMBER_VAL((double)rand() / RAND_MAX * max);
    case 2:
        // random(min, max) → [min, max)
        if (!IS_NUMBER(args[0]) || !IS_NUMBER(args[1])) return NIL_VAL;
        double min = AS_NUMBER(args[0]);
        double max2 = AS_NUMBER(args[1]);
        return NUMBER_VAL(min + (double)rand() / RAND_MAX * (max2 - min));
    default:
        return NIL_VAL;
    }
}

Value fuji_randomInt(int arg_count, Value* args) {
    switch (arg_count) {
    case 1:
        // randomInt(max) → [0, max)
        if (!IS_NUMBER(args[0])) return NIL_VAL;
        double max = AS_NUMBER(args[0]);
        return NUMBER_VAL((double)(rand() % (int)max));
    case 2:
        // randomInt(min, max) → [min, max)
        if (!IS_NUMBER(args[0]) || !IS_NUMBER(args[1])) return NIL_VAL;
        double min = AS_NUMBER(args[0]);
        double max2 = AS_NUMBER(args[1]);
        return NUMBER_VAL(min + (double)(rand() % (int)(max2 - min)));
    default:
        return NIL_VAL;
    }
}

Value fuji_randomChoice(int arg_count, Value* args) {
    // For now, return nil (proper implementation would select random element from array)
    return NIL_VAL;
}

Value fuji_randomSeed(int arg_count, Value* args) {
    if (arg_count != 1 || !IS_NUMBER(args[0])) return NIL_VAL;
    unsigned int seed = (unsigned int)AS_NUMBER(args[0]);
    srand(seed);
    return NIL_VAL;
}

// Math functions for game development
Value fuji_lerp(int arg_count, Value* args) {
    if (arg_count != 3) return NIL_VAL;
    if (!IS_NUMBER(args[0]) || !IS_NUMBER(args[1]) || !IS_NUMBER(args[2])) return NIL_VAL;
    double a = AS_NUMBER(args[0]);
    double b = AS_NUMBER(args[1]);
    double t = AS_NUMBER(args[2]);
    return NUMBER_VAL(a + (b - a) * t);
}

Value fuji_clamp(int arg_count, Value* args) {
    if (arg_count != 3) return NIL_VAL;
    if (!IS_NUMBER(args[0]) || !IS_NUMBER(args[1]) || !IS_NUMBER(args[2])) return NIL_VAL;
    double val = AS_NUMBER(args[0]);
    double min = AS_NUMBER(args[1]);
    double max2 = AS_NUMBER(args[2]);
    
    if (val < min) return NUMBER_VAL(min);
    if (val > max2) return NUMBER_VAL(max2);
    return NUMBER_VAL(val);
}

Value fuji_distance(int arg_count, Value* args) {
    if (arg_count != 4) return NIL_VAL;
    if (!IS_NUMBER(args[0]) || !IS_NUMBER(args[1]) || !IS_NUMBER(args[2]) || !IS_NUMBER(args[3])) return NIL_VAL;
    double x1 = AS_NUMBER(args[0]);
    double y1 = AS_NUMBER(args[1]);
    double x2 = AS_NUMBER(args[2]);
    double y2 = AS_NUMBER(args[3]);
    
    double dx = x2 - x1;
    double dy = y2 - y1;
    return NUMBER_VAL(sqrt(dx * dx + dy * dy));
}

Value fuji_angleBetween(int arg_count, Value* args) {
    if (arg_count != 4) return NIL_VAL;
    if (!IS_NUMBER(args[0]) || !IS_NUMBER(args[1]) || !IS_NUMBER(args[2]) || !IS_NUMBER(args[3])) return NIL_VAL;
    double x1 = AS_NUMBER(args[0]);
    double y1 = AS_NUMBER(args[1]);
    double x2 = AS_NUMBER(args[2]);
    double y2 = AS_NUMBER(args[3]);
    
    return NUMBER_VAL(atan2(y2 - y1, x2 - x1));
}

Value fuji_map(int arg_count, Value* args) {
    if (arg_count != 5) return NIL_VAL;
    if (!IS_NUMBER(args[0]) || !IS_NUMBER(args[1]) || !IS_NUMBER(args[2]) || !IS_NUMBER(args[3]) || !IS_NUMBER(args[4])) return NIL_VAL;
    double val = AS_NUMBER(args[0]);
    double inMin = AS_NUMBER(args[1]);
    double inMax = AS_NUMBER(args[2]);
    double outMin = AS_NUMBER(args[3]);
    double outMax = AS_NUMBER(args[4]);
    
    double normalized = (val - inMin) / (inMax - inMin);
    return NUMBER_VAL(outMin + normalized * (outMax - outMin));
}

Value fuji_pi(int arg_count, Value* args) {
    return NUMBER_VAL(3.141592653589793);
}

Value fuji_e(int arg_count, Value* args) {
    return NUMBER_VAL(2.718281828459045);
}

Value fuji_sin(int arg_count, Value* args) {
    if (arg_count != 1 || !IS_NUMBER(args[0])) return NIL_VAL;
    return NUMBER_VAL(sin(AS_NUMBER(args[0])));
}

Value fuji_cos(int arg_count, Value* args) {
    if (arg_count != 1 || !IS_NUMBER(args[0])) return NIL_VAL;
    return NUMBER_VAL(cos(AS_NUMBER(args[0])));
}

Value fuji_tan(int arg_count, Value* args) {
    if (arg_count != 1 || !IS_NUMBER(args[0])) return NIL_VAL;
    return NUMBER_VAL(tan(AS_NUMBER(args[0])));
}

Value fuji_asin(int arg_count, Value* args) {
    if (arg_count != 1 || !IS_NUMBER(args[0])) return NIL_VAL;
    return NUMBER_VAL(asin(AS_NUMBER(args[0])));
}

Value fuji_acos(int arg_count, Value* args) {
    if (arg_count != 1 || !IS_NUMBER(args[0])) return NIL_VAL;
    return NUMBER_VAL(acos(AS_NUMBER(args[0])));
}

Value fuji_atan(int arg_count, Value* args) {
    if (arg_count != 1 || !IS_NUMBER(args[0])) return NIL_VAL;
    return NUMBER_VAL(atan(AS_NUMBER(args[0])));
}

Value fuji_atan2(int arg_count, Value* args) {
    if (arg_count != 2 || !IS_NUMBER(args[0]) || !IS_NUMBER(args[1])) return NIL_VAL;
    return NUMBER_VAL(atan2(AS_NUMBER(args[0]), AS_NUMBER(args[1])));
}

Value fuji_pow(int arg_count, Value* args) {
    if (arg_count != 2 || !IS_NUMBER(args[0]) || !IS_NUMBER(args[1])) return NIL_VAL;
    return NUMBER_VAL(pow(AS_NUMBER(args[0]), AS_NUMBER(args[1])));
}

Value fuji_exp(int arg_count, Value* args) {
    if (arg_count != 1 || !IS_NUMBER(args[0])) return NIL_VAL;
    return NUMBER_VAL(exp(AS_NUMBER(args[0])));
}

Value fuji_log(int arg_count, Value* args) {
    if (arg_count != 1 || !IS_NUMBER(args[0])) return NIL_VAL;
    return NUMBER_VAL(log(AS_NUMBER(args[0])));
}

Value fuji_log10(int arg_count, Value* args) {
    if (arg_count != 1 || !IS_NUMBER(args[0])) return NIL_VAL;
    return NUMBER_VAL(log10(AS_NUMBER(args[0])));
}

Value fuji_floor(int arg_count, Value* args) {
    if (arg_count != 1 || !IS_NUMBER(args[0])) return NIL_VAL;
    return NUMBER_VAL(floor(AS_NUMBER(args[0])));
}

Value fuji_ceil(int arg_count, Value* args) {
    if (arg_count != 1 || !IS_NUMBER(args[0])) return NIL_VAL;
    return NUMBER_VAL(ceil(AS_NUMBER(args[0])));
}

Value fuji_round(int arg_count, Value* args) {
    if (arg_count != 1 || !IS_NUMBER(args[0])) return NIL_VAL;
    return NUMBER_VAL(round(AS_NUMBER(args[0])));
}

Value fuji_trunc(int arg_count, Value* args) {
    if (arg_count != 1 || !IS_NUMBER(args[0])) return NIL_VAL;
    return NUMBER_VAL(trunc(AS_NUMBER(args[0])));
}

Value fuji_sign(int arg_count, Value* args) {
    if (arg_count != 1 || !IS_NUMBER(args[0])) return NIL_VAL;
    double val = AS_NUMBER(args[0]);
    if (val > 0) return NUMBER_VAL(1);
    if (val < 0) return NUMBER_VAL(-1);
    return NUMBER_VAL(0);
}

Value fuji_min(int arg_count, Value* args) {
    if (arg_count == 0) return NIL_VAL;
    double result = AS_NUMBER(args[0]);
    for (int i = 1; i < arg_count; i++) {
        if (!IS_NUMBER(args[i])) return NIL_VAL;
        double val = AS_NUMBER(args[i]);
        if (val < result) result = val;
    }
    return NUMBER_VAL(result);
}

Value fuji_max(int arg_count, Value* args) {
    if (arg_count == 0) return NIL_VAL;
    double result = AS_NUMBER(args[0]);
    for (int i = 1; i < arg_count; i++) {
        if (!IS_NUMBER(args[i])) return NIL_VAL;
        double val = AS_NUMBER(args[i]);
        if (val > result) result = val;
    }
    return NUMBER_VAL(result);
}

Value fuji_smoothstep(int arg_count, Value* args) {
    if (arg_count != 3) return NIL_VAL;
    if (!IS_NUMBER(args[0]) || !IS_NUMBER(args[1]) || !IS_NUMBER(args[2])) return NIL_VAL;
    double a = AS_NUMBER(args[0]);
    double b = AS_NUMBER(args[1]);
    double t = AS_NUMBER(args[2]);
    
    // Clamp t to [0, 1]
    if (t < 0) t = 0;
    if (t > 1) t = 1;
    
    // Smoothstep interpolation
    double result = t * t * (3 - 2 * t);
    return NUMBER_VAL(a + result * (b - a));
}

Value fuji_distanceSq(int arg_count, Value* args) {
    if (arg_count != 4) return NIL_VAL;
    if (!IS_NUMBER(args[0]) || !IS_NUMBER(args[1]) || !IS_NUMBER(args[2]) || !IS_NUMBER(args[3])) return NIL_VAL;
    double x1 = AS_NUMBER(args[0]);
    double y1 = AS_NUMBER(args[1]);
    double x2 = AS_NUMBER(args[2]);
    double y2 = AS_NUMBER(args[3]);
    
    double dx = x2 - x1;
    double dy = y2 - y1;
    return NUMBER_VAL(dx * dx + dy * dy);
}

Value fuji_normalize(int arg_count, Value* args) {
    if (arg_count != 2) return NIL_VAL;
    if (!IS_NUMBER(args[0]) || !IS_NUMBER(args[1])) return NIL_VAL;
    
    double x = AS_NUMBER(args[0]);
    double y = AS_NUMBER(args[1]);
    
    double len = sqrt(x * x + y * y);
    if (len == 0) {
        return NIL_VAL;  // Cannot normalize zero vector
    }
    
    // Return object {x: normalized_x, y: normalized_y}
    // For now, return nil (proper implementation would create a table with properties)
    return NIL_VAL;
}

Value fuji_hypot(int arg_count, Value* args) {
    if (arg_count != 2 || !IS_NUMBER(args[0]) || !IS_NUMBER(args[1])) return NIL_VAL;
    return NUMBER_VAL(hypot(AS_NUMBER(args[0]), AS_NUMBER(args[1])));
}

Value fuji_fmod(int arg_count, Value* args) {
    if (arg_count != 2 || !IS_NUMBER(args[0]) || !IS_NUMBER(args[1])) return NIL_VAL;
    return NUMBER_VAL(fmod(AS_NUMBER(args[0]), AS_NUMBER(args[1])));
}

Value fuji_degrees(int arg_count, Value* args) {
    if (arg_count != 1 || !IS_NUMBER(args[0])) return NIL_VAL;
    return NUMBER_VAL(AS_NUMBER(args[0]) * (180.0 / 3.141592653589793));
}

Value fuji_radians(int arg_count, Value* args) {
    if (arg_count != 1 || !IS_NUMBER(args[0])) return NIL_VAL;
    return NUMBER_VAL(AS_NUMBER(args[0]) * (3.141592653589793 / 180.0));
}

Value fuji_wrap(int arg_count, Value* args) {
    if (arg_count != 1 || !IS_NUMBER(args[0])) return NIL_VAL;
    double x = AS_NUMBER(args[0]);
    const double pi = 3.141592653589793;
    const double two_pi = 2.0 * pi;
    x = fmod(x + pi, two_pi);
    if (x < 0.0) x += two_pi;
    return NUMBER_VAL(x - pi);
}

Value fuji_approach(int arg_count, Value* args) {
    if (arg_count != 3 || !IS_NUMBER(args[0]) || !IS_NUMBER(args[1]) || !IS_NUMBER(args[2])) return NIL_VAL;
    double value = AS_NUMBER(args[0]);
    double target = AS_NUMBER(args[1]);
    double delta = AS_NUMBER(args[2]);
    if (delta < 0.0) delta = -delta;
    if (value < target) {
        value += delta;
        if (value > target) value = target;
    } else if (value > target) {
        value -= delta;
        if (value < target) value = target;
    }
    return NUMBER_VAL(value);
}

Value fuji_smoothdamp(int arg_count, Value* args) {
    if (arg_count != 4 || !IS_NUMBER(args[0]) || !IS_NUMBER(args[1]) || !IS_NUMBER(args[2]) || !IS_NUMBER(args[3])) return NIL_VAL;
    double current = AS_NUMBER(args[0]);
    double target = AS_NUMBER(args[1]);
    double velocity = AS_NUMBER(args[2]);
    double smooth_time = AS_NUMBER(args[3]);
    if (smooth_time <= 0.0001) smooth_time = 0.0001;
    double omega = 2.0 / smooth_time;
    double x = omega * (1.0 / 60.0);
    double exp_term = 1.0 / (1.0 + x + 0.48 * x * x + 0.235 * x * x * x);
    double change = current - target;
    double temp = (velocity + omega * change) * (1.0 / 60.0);
    double next_velocity = (velocity - omega * temp) * exp_term;
    double output = target + (change + temp) * exp_term;
    (void)next_velocity;
    return NUMBER_VAL(output);
}

// Type checking functions
Value fuji_isNumber(int arg_count, Value* args) {
    if (arg_count != 1) return NIL_VAL;
    return BOOL_VAL(IS_NUMBER(args[0]));
}

Value fuji_isString(int arg_count, Value* args) {
    if (arg_count != 1) return NIL_VAL;
    return BOOL_VAL(IS_OBJ(args[0]) && AS_OBJ(args[0])->type == OBJ_STRING);
}

Value fuji_isBool(int arg_count, Value* args) {
    if (arg_count != 1) return NIL_VAL;
    return BOOL_VAL(IS_BOOL(args[0]));
}

Value fuji_isNull(int arg_count, Value* args) {
    if (arg_count != 1) return NIL_VAL;
    return BOOL_VAL(IS_NIL(args[0]));
}

Value fuji_isArray(int arg_count, Value* args) {
    if (arg_count != 1) return NIL_VAL;
    return BOOL_VAL(IS_OBJ(args[0]) && AS_OBJ(args[0])->type == OBJ_ARRAY);
}

Value fuji_isObject(int arg_count, Value* args) {
    if (arg_count != 1) return NIL_VAL;
    return BOOL_VAL(IS_OBJ(args[0]) && AS_OBJ(args[0])->type == OBJ_TABLE);
}

Value fuji_isFunction(int arg_count, Value* args) {
    if (arg_count != 1) return NIL_VAL;
    return BOOL_VAL(IS_OBJ(args[0]) && AS_OBJ(args[0])->type == OBJ_CLOSURE);
}

// Conversion functions
Value fuji_bool(int arg_count, Value* args) {
    if (arg_count != 1) return NIL_VAL;
    
    Value val = args[0];
    if (IS_BOOL(val)) return val;
    if (IS_NUMBER(val)) return BOOL_VAL(AS_NUMBER(val) != 0);
    if (IS_NIL(val)) return BOOL_VAL(false);
    if (IS_OBJ(val)) return BOOL_VAL(true);
    
    return BOOL_VAL(false);
}

// Format function: args[0] is the format string; each "{}" is replaced by fuji_string(args[1]), args[2], ...
// If there are more placeholders than arguments, remaining "{}" are copied literally.
static int fuji_format_piece_len(Value v) {
    Value s = fuji_string(v);
    if (!IS_OBJ(s) || AS_OBJ(s)->type != OBJ_STRING) return 0;
    return ((ObjString*)AS_OBJ(s))->length;
}

Value fuji_format(int arg_count, Value* args) {
    if (arg_count < 1) return NIL_VAL;

    Value fmtv = fuji_string(args[0]);
    if (!IS_OBJ(fmtv) || AS_OBJ(fmtv)->type != OBJ_STRING) return NIL_VAL;
    ObjString* fmt = (ObjString*)AS_OBJ(fmtv);

    int total = 0;
    int i = 0;
    int ap = 1;
    while (i < fmt->length) {
        if (i + 1 < fmt->length && fmt->chars[i] == '{' && fmt->chars[i + 1] == '}') {
            if (ap < arg_count) {
                total += fuji_format_piece_len(args[ap]);
                ap++;
            } else {
                total += 2;
            }
            i += 2;
        } else {
            total++;
            i++;
        }
    }

    ObjString* out = allocate_string(total);
    int pos = 0;
    i = 0;
    ap = 1;
    while (i < fmt->length) {
        if (i + 1 < fmt->length && fmt->chars[i] == '{' && fmt->chars[i + 1] == '}') {
            if (ap < arg_count) {
                Value piece = fuji_string(args[ap]);
                ap++;
                if (IS_OBJ(piece) && AS_OBJ(piece)->type == OBJ_STRING) {
                    ObjString* p = (ObjString*)AS_OBJ(piece);
                    memcpy(out->chars + pos, p->chars, p->length);
                    pos += p->length;
                }
            } else {
                out->chars[pos++] = '{';
                out->chars[pos++] = '}';
            }
            i += 2;
        } else {
            out->chars[pos++] = fmt->chars[i];
            i++;
        }
    }
    return OBJ_VAL((Obj*)out);
}

// Integer range as half-open [from, to): fuji_range(from, to) -> array of numbers.
Value fuji_range(int arg_count, Value* args) {
    if (arg_count != 2) return NIL_VAL;
    if (!IS_NUMBER(args[0]) || !IS_NUMBER(args[1])) return NIL_VAL;
    int start = (int)AS_NUMBER(args[0]);
    int end = (int)AS_NUMBER(args[1]);
    if (end < start) {
        ObjArray* empty = allocate_array(1);
        return OBJ_VAL((Obj*)empty);
    }
    int n = end - start;
    ObjArray* out = allocate_array(n < 1 ? 1 : n);
    Value ov = OBJ_VAL((Obj*)out);
    for (int i = 0; i < n; i++) {
        fuji_array_push(ov, NUMBER_VAL((double)(start + i)));
    }
    return ov;
}

// Array methods for game development
//
// Higher-order helpers taking Fuji callbacks are lowered by the LLVM backend on `.map()` /
// `.filter()` etc.; argv variants remain unavailable from native Tier-1 glue.

Value fuji_array_map(int arg_count, Value* args) {
    (void)args;
    if (arg_count != 2) return NIL_VAL;
    return NIL_VAL;
}

Value fuji_array_filter(int arg_count, Value* args) {
    (void)args;
    if (arg_count != 2) return NIL_VAL;
    return NIL_VAL;
}

Value fuji_array_forEach(int arg_count, Value* args) {
    (void)args;
    if (arg_count != 2) return NIL_VAL;
    return NIL_VAL;
}

Value fuji_array_find(int arg_count, Value* args) {
    (void)args;
    if (arg_count != 2) return NIL_VAL;
    return NIL_VAL;
}

Value fuji_array_findIndex(int arg_count, Value* args) {
    (void)args;
    if (arg_count != 2) return NIL_VAL;
    return NUMBER_VAL(-1);
}

Value fuji_array_some(int arg_count, Value* args) {
    (void)args;
    if (arg_count != 2) return NIL_VAL;
    return BOOL_VAL(false);
}

Value fuji_array_every(int arg_count, Value* args) {
    (void)args;
    if (arg_count != 2) return NIL_VAL;
    return BOOL_VAL(false);
}

Value fuji_array_reduce(int arg_count, Value* args) {
    (void)args;
    if (arg_count != 2) return NIL_VAL;
    return NIL_VAL;
}

static int cmp_sort_values(Value a, Value b) {
    if (IS_NUMBER(a) && IS_NUMBER(b)) {
        double da = AS_NUMBER(a);
        double db = AS_NUMBER(b);
        return (da > db) - (da < db);
    }
    Value sa = fuji_string(a);
    Value sb = fuji_string(b);
    if (!IS_OBJ(sa) || AS_OBJ(sa)->type != OBJ_STRING) {
        sa = fuji_copy_string("", 0);
    }
    if (!IS_OBJ(sb) || AS_OBJ(sb)->type != OBJ_STRING) {
        sb = fuji_copy_string("", 0);
    }
    ObjString* A = (ObjString*)AS_OBJ(sa);
    ObjString* B = (ObjString*)AS_OBJ(sb);
    int minlen = A->length < B->length ? A->length : B->length;
    int c = memcmp(A->chars, B->chars, (size_t)minlen);
    if (c != 0) {
        return c > 0 ? 1 : -1;
    }
    return (A->length > B->length) - (A->length < B->length);
}

Value fuji_array_sort(int arg_count, Value* args) {
    if (arg_count != 1) return NIL_VAL;
    if (!IS_OBJ(args[0]) || AS_OBJ(args[0])->type != OBJ_ARRAY) return NIL_VAL;
    ObjArray* arr = (ObjArray*)AS_OBJ(args[0]);
    int n = arr->count;
    for (int i = 0; i < n - 1; i++) {
        for (int j = 0; j < n - 1 - i; j++) {
            if (cmp_sort_values(arr->elements[j], arr->elements[j + 1]) > 0) {
                Value tmp = arr->elements[j];
                arr->elements[j] = arr->elements[j + 1];
                arr->elements[j + 1] = tmp;
            }
        }
    }
    return args[0];
}

Value fuji_array_reverse(int arg_count, Value* args) {
    if (arg_count != 1) return NIL_VAL;
    if (!IS_OBJ(args[0]) || AS_OBJ(args[0])->type != OBJ_ARRAY) return NIL_VAL;
    ObjArray* src = (ObjArray*)AS_OBJ(args[0]);
    ObjArray* out = allocate_array(src->count < 1 ? 1 : src->count);
    Value ov = OBJ_VAL((Obj*)out);
    for (int i = src->count - 1; i >= 0; i--) {
        fuji_array_push(ov, src->elements[i]);
    }
    return ov;
}

Value fuji_array_indexOf(int arg_count, Value* args) {
    if (arg_count != 2) return NIL_VAL;
    if (!IS_OBJ(args[0]) || AS_OBJ(args[0])->type != OBJ_ARRAY) return NIL_VAL;
    ObjArray* arr = (ObjArray*)AS_OBJ(args[0]);
    Value needle = args[1];
    for (int i = 0; i < arr->count; i++) {
        if (values_equal(arr->elements[i], needle)) {
            return NUMBER_VAL((double)i);
        }
    }
    return NUMBER_VAL(-1);
}

Value fuji_array_includes(int arg_count, Value* args) {
    if (arg_count != 2) return NIL_VAL;
    if (!IS_OBJ(args[0]) || AS_OBJ(args[0])->type != OBJ_ARRAY) return NIL_VAL;
    ObjArray* arr = (ObjArray*)AS_OBJ(args[0]);
    Value needle = args[1];
    for (int i = 0; i < arr->count; i++) {
        if (values_equal(arr->elements[i], needle)) {
            return BOOL_VAL(true);
        }
    }
    return BOOL_VAL(false);
}

Value fuji_array_slice(int arg_count, Value* args) {
    if (arg_count != 3) return NIL_VAL;
    if (!IS_OBJ(args[0]) || AS_OBJ(args[0])->type != OBJ_ARRAY) return NIL_VAL;
    if (!IS_NUMBER(args[1]) || !IS_NUMBER(args[2])) return NIL_VAL;

    ObjArray* src = (ObjArray*)AS_OBJ(args[0]);
    int count = src->count;
    int start = (int)AS_NUMBER(args[1]);
    int end = (int)AS_NUMBER(args[2]);

    if (start < 0) start = count + start;
    if (end < 0) end = count + end;
    if (start < 0) start = 0;
    if (end > count) end = count;
    if (end < start) end = start;

    int len = end - start;
    ObjArray* out = allocate_array(len < 1 ? 1 : len);
    Value outv = OBJ_VAL((Obj*)out);
    for (int i = 0; i < len; i++) {
        fuji_array_push(outv, src->elements[start + i]);
    }
    return outv;
}

Value fuji_array_concat(int arg_count, Value* args) {
    if (arg_count < 1) return NIL_VAL;

    int total = 0;
    for (int a = 0; a < arg_count; a++) {
        if (!IS_OBJ(args[a]) || AS_OBJ(args[a])->type != OBJ_ARRAY) return NIL_VAL;
        ObjArray* src = (ObjArray*)AS_OBJ(args[a]);
        total += src->count;
    }

    ObjArray* out = allocate_array(total < 1 ? 1 : total);
    Value ov = OBJ_VAL((Obj*)out);
    for (int a = 0; a < arg_count; a++) {
        ObjArray* src = (ObjArray*)AS_OBJ(args[a]);
        for (int j = 0; j < src->count; j++) {
            fuji_array_push(ov, src->elements[j]);
        }
    }
    return ov;
}

Value fuji_array_join(int arg_count, Value* args) {
    if (arg_count != 2) return NIL_VAL;
    if (!IS_OBJ(args[0]) || AS_OBJ(args[0])->type != OBJ_ARRAY) return NIL_VAL;
    Value sepv = fuji_string(args[1]);
    if (!IS_OBJ(sepv) || AS_OBJ(sepv)->type != OBJ_STRING) return NIL_VAL;
    ObjArray* arr = (ObjArray*)AS_OBJ(args[0]);
    ObjString* sep = (ObjString*)AS_OBJ(sepv);
    Value acc = fuji_copy_string("", 0);
    for (int i = 0; i < arr->count; i++) {
        if (i > 0) {
            acc = fuji_string_concat(acc, sepv);
        }
        Value part = fuji_string(arr->elements[i]);
        acc = fuji_string_concat(acc, part);
    }
    return acc;
}

static void push_substring(Value arrv, const char* start, int len) {
    Value chunk = fuji_copy_string(start, len);
    fuji_array_push(arrv, chunk);
}

// String methods for game development
Value fuji_string_split(int arg_count, Value* args) {
    if (arg_count != 2) return NIL_VAL;
    if (!IS_OBJ(args[0]) || AS_OBJ(args[0])->type != OBJ_STRING) return NIL_VAL;
    if (!IS_OBJ(args[1]) || AS_OBJ(args[1])->type != OBJ_STRING) return NIL_VAL;

    ObjString* s = (ObjString*)AS_OBJ(args[0]);
    ObjString* delim = (ObjString*)AS_OBJ(args[1]);

    if (delim->length == 0) {
        Value ov = OBJ_VAL((Obj*)allocate_array(s->length < 1 ? 1 : s->length));
        for (int i = 0; i < s->length; i++) {
            push_substring(ov, &s->chars[i], 1);
        }
        return ov;
    }

    Value ov = OBJ_VAL((Obj*)allocate_array(8));
    int pos = 0;
    const char* hay = s->chars;
    int sl = s->length;
    int dl = delim->length;

    while (pos <= sl) {
        int found = -1;
        for (int i = pos; i <= sl - dl; i++) {
            if (memcmp(hay + i, delim->chars, (size_t)dl) == 0) {
                found = i;
                break;
            }
        }
        if (found < 0) {
            push_substring(ov, hay + pos, sl - pos);
            break;
        }
        push_substring(ov, hay + pos, found - pos);
        pos = found + dl;
    }
    return ov;
}

Value fuji_string_trim(int arg_count, Value* args) {
    if (arg_count != 1) return NIL_VAL;
    if (!IS_OBJ(args[0]) || AS_OBJ(args[0])->type != OBJ_STRING) return NIL_VAL;
    ObjString* s = (ObjString*)AS_OBJ(args[0]);
    int lo = 0;
    int hi = s->length;
    while (lo < hi && isspace((unsigned char)s->chars[lo])) {
        lo++;
    }
    while (hi > lo && isspace((unsigned char)s->chars[hi - 1])) {
        hi--;
    }
    return fuji_copy_string(s->chars + lo, hi - lo);
}

Value fuji_string_upper(int arg_count, Value* args) {
    if (arg_count != 1) return NIL_VAL;
    if (!IS_OBJ(args[0]) || AS_OBJ(args[0])->type != OBJ_STRING) return NIL_VAL;
    ObjString* s = (ObjString*)AS_OBJ(args[0]);
    ObjString* out = allocate_string(s->length);
    for (int i = 0; i < s->length; i++) {
        unsigned char c = (unsigned char)s->chars[i];
        out->chars[i] = (char)toupper((int)c);
    }
    return OBJ_VAL((Obj*)out);
}

Value fuji_string_lower(int arg_count, Value* args) {
    if (arg_count != 1) return NIL_VAL;
    if (!IS_OBJ(args[0]) || AS_OBJ(args[0])->type != OBJ_STRING) return NIL_VAL;
    ObjString* s = (ObjString*)AS_OBJ(args[0]);
    ObjString* out = allocate_string(s->length);
    for (int i = 0; i < s->length; i++) {
        unsigned char c = (unsigned char)s->chars[i];
        out->chars[i] = (char)tolower((int)c);
    }
    return OBJ_VAL((Obj*)out);
}

Value fuji_string_startsWith(int arg_count, Value* args) {
    if (arg_count != 2) return NIL_VAL;
    if (!IS_OBJ(args[0]) || AS_OBJ(args[0])->type != OBJ_STRING) return NIL_VAL;
    if (!IS_OBJ(args[1]) || AS_OBJ(args[1])->type != OBJ_STRING) return NIL_VAL;
    ObjString* s = (ObjString*)AS_OBJ(args[0]);
    ObjString* pre = (ObjString*)AS_OBJ(args[1]);
    if (pre->length > s->length) return BOOL_VAL(false);
    return BOOL_VAL(memcmp(s->chars, pre->chars, (size_t)pre->length) == 0);
}

Value fuji_string_endsWith(int arg_count, Value* args) {
    if (arg_count != 2) return NIL_VAL;
    if (!IS_OBJ(args[0]) || AS_OBJ(args[0])->type != OBJ_STRING) return NIL_VAL;
    if (!IS_OBJ(args[1]) || AS_OBJ(args[1])->type != OBJ_STRING) return NIL_VAL;
    ObjString* s = (ObjString*)AS_OBJ(args[0]);
    ObjString* suf = (ObjString*)AS_OBJ(args[1]);
    if (suf->length > s->length) return BOOL_VAL(false);
    int off = s->length - suf->length;
    return BOOL_VAL(memcmp(s->chars + off, suf->chars, (size_t)suf->length) == 0);
}

Value fuji_string_indexOf(int arg_count, Value* args) {
    if (arg_count != 2) return NIL_VAL;
    if (!IS_OBJ(args[0]) || AS_OBJ(args[0])->type != OBJ_STRING) return NIL_VAL;
    if (!IS_OBJ(args[1]) || AS_OBJ(args[1])->type != OBJ_STRING) return NIL_VAL;
    ObjString* s = (ObjString*)AS_OBJ(args[0]);
    ObjString* needle = (ObjString*)AS_OBJ(args[1]);
    if (needle->length == 0) return NUMBER_VAL(0);
    if (needle->length > s->length) return NUMBER_VAL(-1);
    for (int i = 0; i <= s->length - needle->length; i++) {
        if (memcmp(s->chars + i, needle->chars, (size_t)needle->length) == 0) {
            return NUMBER_VAL((double)i);
        }
    }
    return NUMBER_VAL(-1);
}

Value fuji_string_slice(int arg_count, Value* args) {
    if (arg_count != 3) return NIL_VAL;
    if (!IS_OBJ(args[0]) || AS_OBJ(args[0])->type != OBJ_STRING) return NIL_VAL;
    if (!IS_NUMBER(args[1]) || !IS_NUMBER(args[2])) return NIL_VAL;

    ObjString* s = (ObjString*)AS_OBJ(args[0]);
    int len = s->length;
    int start = (int)AS_NUMBER(args[1]);
    int end = (int)AS_NUMBER(args[2]);

    if (start < 0) start = len + start;
    if (end < 0) end = len + end;
    if (start < 0) start = 0;
    if (end > len) end = len;
    if (end < start) end = start;

    int span = end - start;
    return fuji_copy_string(s->chars + start, span);
}

Value fuji_string_replace(int arg_count, Value* args) {
    if (arg_count != 3) return NIL_VAL;
    if (!IS_OBJ(args[0]) || AS_OBJ(args[0])->type != OBJ_STRING) return NIL_VAL;
    if (!IS_OBJ(args[1]) || AS_OBJ(args[1])->type != OBJ_STRING) return NIL_VAL;
    if (!IS_OBJ(args[2]) || AS_OBJ(args[2])->type != OBJ_STRING) return NIL_VAL;

    ObjString* s = (ObjString*)AS_OBJ(args[0]);
    ObjString* from = (ObjString*)AS_OBJ(args[1]);
    ObjString* to = (ObjString*)AS_OBJ(args[2]);

    if (from->length == 0) return args[0];

    int idx = -1;
    for (int i = 0; i <= s->length - from->length; i++) {
        if (memcmp(s->chars + i, from->chars, (size_t)from->length) == 0) {
            idx = i;
            break;
        }
    }
    if (idx < 0) return args[0];

    int newLen = s->length - from->length + to->length;
    ObjString* out = allocate_string(newLen);
    memcpy(out->chars, s->chars, (size_t)idx);
    memcpy(out->chars + idx, to->chars, (size_t)to->length);
    memcpy(out->chars + idx + to->length, s->chars + idx + from->length,
           (size_t)(s->length - idx - from->length));
    return OBJ_VAL((Obj*)out);
}

Value fuji_string_replaceAll(int arg_count, Value* args) {
    if (arg_count != 3) return NIL_VAL;
    Value cur = args[0];
    if (!IS_OBJ(cur) || AS_OBJ(cur)->type != OBJ_STRING) return NIL_VAL;
    if (!IS_OBJ(args[1]) || AS_OBJ(args[1])->type != OBJ_STRING) return NIL_VAL;
    if (!IS_OBJ(args[2]) || AS_OBJ(args[2])->type != OBJ_STRING) return NIL_VAL;

    ObjString* from = (ObjString*)AS_OBJ(args[1]);
    if (from->length == 0) return cur;

    while (true) {
        Value args3[3] = { cur, args[1], args[2] };
        Value next = fuji_string_replace(3, args3);
        if (values_equal(next, cur)) {
            break;
        }
        cur = next;
    }
    return cur;
}

/** Literal substring match after string coercion (`matches(haystack, pattern)`). */
Value fuji_matches(int arg_count, Value* args) {
    if (arg_count != 2) return NIL_VAL;
    Value a = fuji_string(args[0]);
    Value b = fuji_string(args[1]);
    if (!IS_OBJ(a) || AS_OBJ(a)->type != OBJ_STRING) return NIL_VAL;
    if (!IS_OBJ(b) || AS_OBJ(b)->type != OBJ_STRING) return NIL_VAL;
    ObjString* hay = (ObjString*)AS_OBJ(a);
    ObjString* needle = (ObjString*)AS_OBJ(b);
    if (needle->length == 0) return TRUE_VAL;
    char* hit = strstr(hay->chars, needle->chars);
    return BOOL_VAL(hit != NULL);
}

// File I/O functions
Value fuji_readFile(int arg_count, Value* args) {
    if (arg_count < 1 || !IS_OBJ(args[0]) || AS_OBJ(args[0])->type != OBJ_STRING) {
        return fuji_err_str("readFile requires a string path");
    }
    ObjString* path = (ObjString*)AS_OBJ(args[0]);

    FILE* f = fopen(path->chars, "rb");
    if (!f) {
        char msg[512];
        snprintf(msg, sizeof(msg), "could not open '%s': %s", path->chars, strerror(errno));
        return fuji_err_str(msg);
    }
    if (fseek(f, 0, SEEK_END) != 0) {
        fclose(f);
        return fuji_err_str("could not seek file");
    }
    long size = ftell(f);
    if (size < 0) {
        fclose(f);
        return fuji_err_str("could not determine file size");
    }
    rewind(f);
    if (size == 0) {
        fclose(f);
        Value content = fuji_copy_string("", 0);
        Value inner[1] = { content };
        return fuji_ok(1, inner);
    }
    char* buf = (char*)malloc((size_t)size + 1u);
    if (!buf) {
        fclose(f);
        return fuji_err_str("out of memory reading file");
    }
    size_t got = fread(buf, 1, (size_t)size, f);
    fclose(f);
    if (got != (size_t)size) {
        free(buf);
        return fuji_err_str("could not read full file");
    }
    buf[size] = '\0';
    Value content = fuji_copy_string(buf, (int)size);
    free(buf);
    Value inner[1] = { content };
    return fuji_ok(1, inner);
}

Value fuji_writeFile(int arg_count, Value* args) {
    if (arg_count < 2 || !IS_OBJ(args[0]) || AS_OBJ(args[0])->type != OBJ_STRING) {
        return fuji_err_str("writeFile requires path string and content string");
    }
    if (!IS_OBJ(args[1]) || AS_OBJ(args[1])->type != OBJ_STRING) {
        return fuji_err_str("writeFile requires path string and content string");
    }
    ObjString* path = (ObjString*)AS_OBJ(args[0]);
    ObjString* body = (ObjString*)AS_OBJ(args[1]);

    FILE* f = fopen(path->chars, "wb");
    if (!f) {
        char msg[512];
        snprintf(msg, sizeof(msg), "could not open '%s' for write: %s", path->chars, strerror(errno));
        return fuji_err_str(msg);
    }
    if (body->length > 0) {
        size_t w = fwrite(body->chars, 1, (size_t)body->length, f);
        if (w != (size_t)body->length) {
            fclose(f);
            return fuji_err_str("writeFile: short write");
        }
    }
    if (fclose(f) != 0) {
        return fuji_err_str("writeFile: fclose failed");
    }
    Value inner[1] = { NIL_VAL };
    return fuji_ok(1, inner);
}

Value fuji_appendFile(int arg_count, Value* args) {
    if (arg_count != 2) return NIL_VAL;
    if (!IS_OBJ(args[0]) || AS_OBJ(args[0])->type != OBJ_STRING) return NIL_VAL;
    if (!IS_OBJ(args[1]) || AS_OBJ(args[1])->type != OBJ_STRING) return NIL_VAL;
    
    // For now, return false (proper implementation would append to file)
    return BOOL_VAL(false);
}

Value fuji_fileExists(int arg_count, Value* args) {
    if (arg_count != 1) return NIL_VAL;
    if (!IS_OBJ(args[0]) || AS_OBJ(args[0])->type != OBJ_STRING) return NIL_VAL;
    
    // For now, return false (proper implementation would check if file exists)
    return BOOL_VAL(false);
}

Value fuji_deleteFile(int arg_count, Value* args) {
    if (arg_count != 1) return NIL_VAL;
    if (!IS_OBJ(args[0]) || AS_OBJ(args[0])->type != OBJ_STRING) return NIL_VAL;
    
    // For now, return false (proper implementation would delete file)
    return BOOL_VAL(false);
}

// Debug utilities (legacy LLVM helper; uses truthy + panic)
void fuji_assert_llvm(Value cond, Value msg) {
    Value a[1] = { cond };
    Value b = fuji_bool(1, a);
    if (IS_TRUE(b)) {
        return;
    }
    if (IS_OBJ(msg) && AS_OBJ(msg)->type == OBJ_STRING) {
        Value pargs[1] = { msg };
        fuji_panic(1, pargs);
    }
    fuji_panic_str("assertion failed");
}

Value fuji_trace(int arg_count, Value* args) {
    for (int i = 0; i < arg_count; i++) {
        print_value(args[i]);
        if (i < arg_count - 1) printf(" ");
    }
    printf("\n");
    return NIL_VAL;
}

// JSON parsing (subset) — returns ok(value) or err(message)
typedef struct {
    const char* current;
    char err[256];
} FujiJsonParser;

static void fuji_json_skip_ws(FujiJsonParser* p) {
    while (*p->current == ' ' || *p->current == '\t' || *p->current == '\n' || *p->current == '\r') {
        p->current++;
    }
}

static void fuji_json_set_err(FujiJsonParser* p, const char* msg) {
    if (p->err[0] != '\0') {
        return;
    }
    snprintf(p->err, sizeof(p->err), "%s", msg);
}

static Value fuji_json_parse_value(FujiJsonParser* p);

static Value fuji_json_parse_string(FujiJsonParser* p) {
    if (*p->current != '"') {
        fuji_json_set_err(p, "expected '\"' for string");
        return NIL_VAL;
    }
    p->current++;
    const char* start = p->current;
    while (*p->current != '"' && *p->current != '\0') {
        p->current++;
    }
    if (*p->current != '"') {
        fuji_json_set_err(p, "unterminated string");
        return NIL_VAL;
    }
    int len = (int)(p->current - start);
    Value sv = fuji_copy_string(start, len);
    p->current++;
    return sv;
}

static Value fuji_json_parse_array(FujiJsonParser* p) {
    if (*p->current != '[') {
        fuji_json_set_err(p, "expected '['");
        return NIL_VAL;
    }
    p->current++;
    Value arrv = fuji_allocate_array(0);
    fuji_json_skip_ws(p);
    while (*p->current != ']' && *p->current != '\0') {
        Value elem = fuji_json_parse_value(p);
        if (p->err[0]) {
            return NIL_VAL;
        }
        fuji_array_push(arrv, elem);
        fuji_json_skip_ws(p);
        if (*p->current == ',') {
            p->current++;
            fuji_json_skip_ws(p);
        }
    }
    if (*p->current != ']') {
        fuji_json_set_err(p, "expected ']'");
        return NIL_VAL;
    }
    p->current++;
    return arrv;
}

static Value fuji_json_parse_object(FujiJsonParser* p) {
    if (*p->current != '{') {
        fuji_json_set_err(p, "expected '{'");
        return NIL_VAL;
    }
    p->current++;
    Value objv = fuji_allocate_object(64);
    fuji_json_skip_ws(p);
    while (*p->current != '}' && *p->current != '\0') {
        Value key = fuji_json_parse_string(p);
        if (p->err[0]) {
            return NIL_VAL;
        }
        fuji_json_skip_ws(p);
        if (*p->current != ':') {
            fuji_json_set_err(p, "expected ':' in object");
            return NIL_VAL;
        }
        p->current++;
        fuji_json_skip_ws(p);
        Value val = fuji_json_parse_value(p);
        if (p->err[0]) {
            return NIL_VAL;
        }
        fuji_object_set(objv, key, val);
        fuji_json_skip_ws(p);
        if (*p->current == ',') {
            p->current++;
            fuji_json_skip_ws(p);
        }
    }
    if (*p->current != '}') {
        fuji_json_set_err(p, "expected '}'");
        return NIL_VAL;
    }
    p->current++;
    return objv;
}

static Value fuji_json_parse_value(FujiJsonParser* p) {
    fuji_json_skip_ws(p);
    char c = *p->current;
    if (c == '\0') {
        fuji_json_set_err(p, "unexpected end of input");
        return NIL_VAL;
    }
    if (c == '"') {
        return fuji_json_parse_string(p);
    }
    if (c == '[') {
        return fuji_json_parse_array(p);
    }
    if (c == '{') {
        return fuji_json_parse_object(p);
    }
    if (c == 't') {
        if (strncmp(p->current, "true", 4) == 0) {
            p->current += 4;
            return TRUE_VAL;
        }
        fuji_json_set_err(p, "invalid literal");
        return NIL_VAL;
    }
    if (c == 'f') {
        if (strncmp(p->current, "false", 5) == 0) {
            p->current += 5;
            return FALSE_VAL;
        }
        fuji_json_set_err(p, "invalid literal");
        return NIL_VAL;
    }
    if (c == 'n') {
        if (strncmp(p->current, "null", 4) == 0) {
            p->current += 4;
            return NIL_VAL;
        }
        fuji_json_set_err(p, "invalid literal");
        return NIL_VAL;
    }
    if ((c >= '0' && c <= '9') || c == '-') {
        char* end = NULL;
        double d = strtod(p->current, &end);
        if (end == p->current) {
            fuji_json_set_err(p, "invalid number");
            return NIL_VAL;
        }
        p->current = end;
        return NUMBER_VAL(d);
    }
    fuji_json_set_err(p, "unexpected character");
    return NIL_VAL;
}

Value fuji_parseJSON(int arg_count, Value* args) {
    if (arg_count < 1 || !IS_OBJ(args[0]) || AS_OBJ(args[0])->type != OBJ_STRING) {
        return fuji_err_str("parseJSON requires a string");
    }
    ObjString* s = (ObjString*)AS_OBJ(args[0]);
    FujiJsonParser p;
    p.current = s->chars;
    p.err[0] = '\0';
    Value val = fuji_json_parse_value(&p);
    if (p.err[0]) {
        return fuji_err_str(p.err);
    }
    fuji_json_skip_ws(&p);
    if (*p.current != '\0') {
        return fuji_err_str("parseJSON: trailing data after value");
    }
    Value inner[1] = { val };
    return fuji_ok(1, inner);
}

Value fuji_toJSON(int arg_count, Value* args) {
    if (arg_count != 1) return NIL_VAL;
    
    // For now, return empty string (proper implementation would convert to JSON)
    return OBJ_VAL((Obj*)allocate_string(0));
}

// Graphics: not linked here (no Raylib). Use fuji build + FUJI_NATIVE_SOURCES / wrapgen for native games.
Value fuji_gfx_init_window(int arg_count, Value* args) { (void)arg_count; (void)args; return NIL_VAL; }
Value fuji_gfx_set_target_fps(int arg_count, Value* args) { (void)arg_count; (void)args; return NIL_VAL; }
Value fuji_gfx_window_should_close(int arg_count, Value* args) { (void)arg_count; (void)args; return NUMBER_VAL(1); }
Value fuji_gfx_close_window(int arg_count, Value* args) { (void)arg_count; (void)args; return NIL_VAL; }
Value fuji_gfx_begin_drawing(int arg_count, Value* args) { (void)arg_count; (void)args; return NIL_VAL; }
Value fuji_gfx_end_drawing(int arg_count, Value* args) { (void)arg_count; (void)args; return NIL_VAL; }
Value fuji_gfx_clear_background(int arg_count, Value* args) { (void)arg_count; (void)args; return NIL_VAL; }
Value fuji_gfx_draw_text(int arg_count, Value* args) { (void)arg_count; (void)args; return NIL_VAL; }
Value fuji_gfx_draw_rectangle(int arg_count, Value* args) { (void)arg_count; (void)args; return NIL_VAL; }
Value fuji_gfx_draw_circle(int arg_count, Value* args) { (void)arg_count; (void)args; return NIL_VAL; }
Value fuji_gfx_is_key_down(int arg_count, Value* args) { (void)arg_count; (void)args; return NUMBER_VAL(0); }
Value fuji_gfx_is_key_pressed(int arg_count, Value* args) { (void)arg_count; (void)args; return NUMBER_VAL(0); }
Value fuji_gfx_begin_mode3d(int arg_count, Value* args) { (void)arg_count; (void)args; return NIL_VAL; }
Value fuji_gfx_end_mode3d(int arg_count, Value* args) { (void)arg_count; (void)args; return NIL_VAL; }
Value fuji_gfx_draw_grid(int arg_count, Value* args) { (void)arg_count; (void)args; return NIL_VAL; }
Value fuji_gfx_draw_cube(int arg_count, Value* args) { (void)arg_count; (void)args; return NIL_VAL; }
Value fuji_gfx_draw_cube_wires(int arg_count, Value* args) { (void)arg_count; (void)args; return NIL_VAL; }
