#ifndef FUJI_RUNTIME_H
#define FUJI_RUNTIME_H

#include "value.h"
#include "object.h"
#include "gc.h"
#include <stddef.h>
#include <stdbool.h>

extern void* gc_stack_base;
extern Value* fuji_globals;
extern int fuji_globals_count;
extern int fuji_globals_capacity;
extern Value** fuji_global_slots;
extern int fuji_global_slots_count;
extern int fuji_global_slots_capacity;

// Runtime initialization
void fuji_runtime_set_stack_base(void* base);
void fuji_runtime_init_ex(void* stack_base);
void fuji_runtime_init(void);

// Runtime cleanup
void fuji_runtime_shutdown(void);
void fuji_globals_init(void);
void fuji_globals_free(void);
void fuji_register_global(Value v);
void fuji_register_global_slot(Value* slot);

void fuji_gc_set_threshold(size_t bytes);
void fuji_gc_disable(void);
void fuji_gc_enable(void);
void fuji_gc_collect(void);
void fuji_gc_use_shadow_stack(bool enable);
void fuji_gc_frame_step(double budget_ms);

#include "shadow_stack.h"

Value* fuji_alloc_cell(void);
Value fuji_cell_read(Value* cell);
void fuji_cell_write(Value* cell, Value v);

// Native function declarations
Value fuji_print_argv(int arg_count, Value* args);
Value fuji_print_val(Value v);
void fuji_print_newline(void);
Value fuji_typeof(int arg_count, Value* args);
Value fuji_get_index(Value obj, Value index);
void fuji_assert_llvm(Value cond, Value msg);

/** Result helpers and panic (see language docs). */
Value fuji_ok(int argc, Value* argv);
Value fuji_err(int argc, Value* argv);
void fuji_panic(int argc, Value* argv);
Value fuji_assert(int argc, Value* argv);
Value fuji_err_str(const char* msg);
void fuji_panic_str(const char* msg);
/** After a full mark pass, drop intern slots for unmarked strings (weak refs); call before gc_sweep. */
void fuji_sweep_intern_table(void);

/** Optional call stack for panic stack traces (native codegen). */
void fuji_push_call(const char* fn_name, const char* file_name, int line);
void fuji_pop_call(void);
void fuji_print_stack_trace(void);
Value fuji_clock(void);
Value fuji_wall_time(void);
Value fuji_allocate_string(int length, char* chars);
/** Copy UTF-8 bytes into a string object (wrapgen / FFI helpers). */
Value fuji_copy_string(const char* chars, int length);
Value fuji_allocate_object(int property_count);
Value fuji_object_get(Value obj, Value key);
Value fuji_object_set(Value obj, Value key, Value value);
Value fuji_object_remove(Value obj, Value key);

/** Native lowering for `for-in` / `for-of` over arrays and tables (slot order). */
Value fuji_forof_length(Value v);
Value fuji_forof_key_at(Value v, Value idx_val);
Value fuji_forof_value_at(Value v, Value idx_val);

/** NaN-box helpers for LLVM codegen — never bitcast i64 to double without these. */
double fuji_unbox_number(Value v);
Value fuji_box_number(double d);

/** Unified index read / write (arrays, strings, tables). */
Value fuji_get(Value obj, Value key);
Value fuji_set(Value obj, Value key, Value value);
int fuji_get_shadow_depth(void);
/** High-water shadow stack depth since runtime init (for diagnostics). */
int fuji_shadow_stack_high_water(void);

// Bool value helper (not in value.h, needed for wrapper)
static inline Value BOOL_VAL(bool b) {
    return b ? TRUE_VAL : FALSE_VAL;
}

#endif // FUJI_RUNTIME_H
