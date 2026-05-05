#include "gc.h"
#include "value.h"
#include "object.h"
#include "shadow_stack.h"
#include <stdbool.h>
#include <stdint.h>
#include <stdio.h>
#include <stdlib.h>
#include <string.h>

/* Stack base for conservative root scan; set in fuji_runtime_init before any allocation. */
extern void* gc_stack_base;

extern Value fuji_globals[];
extern int fuji_globals_count;

void fuji_mark_module_cache(void);
void fuji_mark_open_upvalues(void);

#define REMEMBERED_SET_MAX 4096
#define NURSERY_SIZE (256u * 1024u)

typedef struct {
    Obj* objects;
    size_t bytes_allocated;
    size_t next_gc;
    bool gc_disabled;
    bool collecting;
    bool use_shadow_stack;
    GCStats stats;

    uint8_t nursery_buf[NURSERY_SIZE];
    uint8_t* nursery_top;
    uint8_t* nursery_end;
    size_t nursery_live_bytes;

    Obj* remembered[REMEMBERED_SET_MAX];
    int remembered_count;
} GCState;

static GCState gc_state = {
    .objects = NULL,
    .bytes_allocated = 0,
    .next_gc = 1024u * 1024u,
    .gc_disabled = false,
    .collecting = false,
    .use_shadow_stack = false,
    .nursery_top = NULL,
    .nursery_end = NULL,
    .nursery_live_bytes = 0,
    .remembered_count = 0,
};

static size_t gc_obj_total_bytes(Obj* obj);
static bool gc_is_heap_object_exact(Obj* candidate);
static Obj* gc_find_containing_obj(void* p);
static void maybe_mark_stack_word(uint64_t word);
static void maybe_mark_stack_addr(uintptr_t addr);
static void gc_unmark_all(void);

void gc_collect_minor(void);

static bool obj_header_in_nursery_range(const Obj* obj) {
    const uint8_t* p = (const uint8_t*)obj;
    return p >= gc_state.nursery_buf && p < gc_state.nursery_end;
}

void gc_register_object(Obj* obj) {
    if (obj_header_in_nursery_range(obj)) {
        obj->generation = GEN_NURSERY;
    } else {
        obj->generation = GEN_OLD;
    }
    obj->next = gc_state.objects;
    gc_state.objects = obj;
}

void gc_free(void* ptr, size_t size) {
    if (ptr == NULL) {
        return;
    }
    const uint8_t* p = (const uint8_t*)ptr;
    if (p >= gc_state.nursery_buf && p < gc_state.nursery_end) {
        (void)size;
        return;
    }
    if (size <= gc_state.bytes_allocated) {
        gc_state.bytes_allocated -= size;
    } else {
        gc_state.bytes_allocated = 0;
    }
    free(ptr);
}

static void* nursery_alloc(size_t size) {
    size = (size + 7u) & ~7u;
    if (gc_state.nursery_top == NULL || gc_state.nursery_end == NULL) {
        return NULL;
    }
    if (gc_state.nursery_top + size > gc_state.nursery_end) {
        gc_collect_minor();
        if (gc_state.nursery_top + size > gc_state.nursery_end) {
            return NULL;
        }
    }
    void* ptr = (void*)gc_state.nursery_top;
    gc_state.nursery_top += size;
    gc_state.nursery_live_bytes += size;
    return ptr;
}

void* gc_alloc(size_t size) {
    if (size == 0) {
        size = 1;
    }
    if (gc_state.gc_disabled) {
        void* ptr = malloc(size);
        if (ptr == NULL) {
            fprintf(stderr, "fuji: out of memory\n");
            exit(1);
        }
        return ptr;
    }
    if (size <= 4096u && !gc_state.collecting) {
        void* n = nursery_alloc(size);
        if (n != NULL) {
            return n;
        }
    }
    gc_state.bytes_allocated += size;
    if (!gc_state.collecting && gc_state.bytes_allocated > gc_state.next_gc) {
        gc_collect();
    }
    void* ptr = malloc(size);
    if (ptr == NULL) {
        fprintf(stderr, "fuji: out of memory\n");
        exit(1);
    }
    return ptr;
}

static void remembered_set_add(Obj* obj) {
    if (obj == NULL) {
        return;
    }
    for (int i = 0; i < gc_state.remembered_count; i++) {
        if (gc_state.remembered[i] == obj) {
            return;
        }
    }
    if (gc_state.remembered_count < REMEMBERED_SET_MAX) {
        gc_state.remembered[gc_state.remembered_count++] = obj;
        return;
    }
    gc_collect();
    gc_state.remembered_count = 0;
}

static void gc_reset_nursery(void) {
    gc_state.nursery_top = gc_state.nursery_buf;
    gc_state.nursery_live_bytes = 0;
}

static void gc_unlink_gen_dead(void) {
    Obj** previous = &gc_state.objects;
    Obj* obj = gc_state.objects;
    while (obj != NULL) {
        if (obj->generation == GEN_DEAD) {
            *previous = obj->next;
            free_object(obj);
            obj = *previous;
        } else {
            previous = &obj->next;
            obj = obj->next;
        }
    }
}

static bool nursery_has_live_object(void) {
    for (Obj* o = gc_state.objects; o != NULL; o = o->next) {
        const uint8_t* p = (const uint8_t*)o;
        if (p >= gc_state.nursery_buf && p < gc_state.nursery_end && o->generation != GEN_DEAD) {
            return true;
        }
    }
    return false;
}

void gc_collect_minor(void) {
    if (gc_state.collecting) {
        return;
    }
    gc_state.collecting = true;

    gc_unmark_all();
    if (gc_state.use_shadow_stack) {
        gc_mark_shadow_stack();
    } else {
        gc_mark_stack_conservative();
    }
    if (fuji_globals_count > 0) {
        for (int i = 0; i < fuji_globals_count; i++) {
            gc_mark_value(fuji_globals[i]);
        }
    }
    fuji_mark_module_cache();
    fuji_mark_open_upvalues();

    for (int i = 0; i < gc_state.remembered_count; i++) {
        if (gc_state.remembered[i] != NULL) {
            gc_mark_object(gc_state.remembered[i]);
        }
    }

    for (Obj* obj = gc_state.objects; obj != NULL; obj = obj->next) {
        if (obj->generation == GEN_NURSERY) {
            if (obj->is_marked) {
                obj->generation = GEN_YOUNG;
            } else {
                obj->generation = GEN_DEAD;
            }
            obj->is_marked = false;
        }
    }

    gc_unlink_gen_dead();
    gc_state.remembered_count = 0;

    if (!nursery_has_live_object()) {
        gc_reset_nursery();
    }

    gc_unmark_all();
    gc_state.stats.collections++;
    gc_state.collecting = false;
}

size_t gc_nursery_used_bytes(void) {
    if (gc_state.nursery_top == NULL) {
        return 0;
    }
    return (size_t)(gc_state.nursery_top - gc_state.nursery_buf);
}

size_t gc_nursery_capacity_bytes(void) {
    return (size_t)sizeof(gc_state.nursery_buf);
}

static size_t gc_obj_total_bytes(Obj* obj) {
    switch (obj->type) {
        case OBJ_STRING:
            return sizeof(ObjString) + (size_t)((ObjString*)obj)->length + 1u;
        case OBJ_ARRAY: {
            ObjArray* a = (ObjArray*)obj;
            return sizeof(ObjArray) + (size_t)a->capacity * sizeof(Value);
        }
        case OBJ_TABLE: {
            ObjTable* t = (ObjTable*)obj;
            return sizeof(ObjTable) + (size_t)t->capacity * 2u * sizeof(Value);
        }
        case OBJ_FUNCTION: {
            ObjFunction* f = (ObjFunction*)obj;
            return sizeof(ObjFunction);
        }
        case OBJ_CLOSURE: {
            ObjClosure* c = (ObjClosure*)obj;
            return sizeof(ObjClosure) + (size_t)c->upvalue_count * sizeof(Value);
        }
        case OBJ_NATIVE:
            return sizeof(ObjNative);
        case OBJ_CELL:
            return sizeof(ObjCell);
        default:
            return sizeof(Obj);
    }
}

void gc_mark_object(Obj* obj) {
    if (obj == NULL || obj->is_marked) {
        return;
    }
    if (obj->generation == GEN_DEAD) {
        return;
    }
    obj->is_marked = true;

    switch (obj->type) {
        case OBJ_ARRAY: {
            ObjArray* array = (ObjArray*)obj;
            for (int i = 0; i < array->count; i++) {
                gc_mark_value(array->elements[i]);
            }
            break;
        }
        case OBJ_TABLE: {
            ObjTable* table = (ObjTable*)obj;
            for (int i = 0; i < table->count; i++) {
                gc_mark_value(table->keys[i]);
                gc_mark_value(table->values[i]);
            }
            break;
        }
        case OBJ_CLOSURE: {
            ObjClosure* closure = (ObjClosure*)obj;
            gc_mark_object((Obj*)closure->function);
            for (int i = 0; i < closure->upvalue_count; i++) {
                gc_mark_value(closure->upvalues[i]);
            }
            break;
        }
        case OBJ_CELL: {
            ObjCell* cell = (ObjCell*)obj;
            gc_mark_value(cell->value);
            break;
        }
        default:
            break;
    }
}

void gc_mark_value(Value v) {
    if (IS_OBJ(v)) {
        gc_mark_object(AS_OBJ(v));
    }
}

void gc_write_barrier(Obj* parent, Value new_val) {
    if (gc_state.collecting || parent == NULL) {
        return;
    }
    if (!IS_OBJ(new_val)) {
        return;
    }
    Obj* child = AS_OBJ(new_val);
    if (child == NULL) {
        return;
    }
    if (parent->generation == GEN_OLD && child->generation != GEN_OLD) {
        remembered_set_add(parent);
    }
}

static bool gc_is_heap_object_exact(Obj* candidate) {
    if (candidate == NULL) {
        return false;
    }
    for (Obj* obj = gc_state.objects; obj != NULL; obj = obj->next) {
        if (obj->generation == GEN_DEAD) {
            continue;
        }
        if (obj == candidate) {
            return true;
        }
    }
    return false;
}

static Obj* gc_find_containing_obj(void* p) {
    if (p == NULL) {
        return NULL;
    }
    uintptr_t c = (uintptr_t)p;
    for (Obj* obj = gc_state.objects; obj != NULL; obj = obj->next) {
        if (obj->generation == GEN_DEAD) {
            continue;
        }
        uintptr_t base = (uintptr_t)obj;
        size_t sz = gc_obj_total_bytes(obj);
        uintptr_t end = base + sz;
        if (c >= base && c < end) {
            return obj;
        }
    }
    return NULL;
}

void gc_mark_stack_conservative(void) {
    if (gc_stack_base == NULL) {
        return;
    }
    void* stack_top;
    GC_GET_SP(stack_top);
    uintptr_t a = (uintptr_t)stack_top;
    uintptr_t b = (uintptr_t)gc_stack_base;
    uintptr_t lo = a < b ? a : b;
    uintptr_t hi = a < b ? b : a;
    lo = (lo + 7u) & ~(uintptr_t)7u;
    for (uintptr_t p = lo; p + sizeof(uint64_t) <= hi; p += sizeof(uint64_t)) {
        uint64_t word = *(uint64_t*)p;
        maybe_mark_stack_word(word);
        maybe_mark_stack_addr((uintptr_t)word);
    }
}

static void maybe_mark_stack_word(uint64_t word) {
    Value v = word;
    if (IS_OBJ(v)) {
        Obj* o = AS_OBJ(v);
        if (gc_is_heap_object_exact(o)) {
            gc_mark_object(o);
        } else {
            Obj* c = gc_find_containing_obj(o);
            if (c != NULL) {
                gc_mark_object(c);
            }
        }
    }
}

static void maybe_mark_stack_addr(uintptr_t addr) {
    if (addr == 0) {
        return;
    }
    Obj* o = gc_find_containing_obj((void*)addr);
    if (o != NULL) {
        gc_mark_object(o);
    }
}

void gc_mark_shadow_stack(void) {
    for (int i = 0; i < fuji_shadow_depth; i++) {
        Value** ptrs = fuji_shadow_stack[i].slot_ptrs;
        int n = fuji_shadow_stack[i].count;
        if (ptrs == NULL || n <= 0) {
            continue;
        }
        for (int j = 0; j < n; j++) {
            Value* p = ptrs[j];
            if (p != NULL) {
                gc_mark_value(*p);
            }
        }
    }
}

void gc_set_use_shadow_stack(bool enabled) {
    gc_state.use_shadow_stack = enabled;
}

void gc_mark_roots(void) {
    if (gc_state.use_shadow_stack) {
        gc_mark_shadow_stack();
    } else {
        gc_mark_stack_conservative();
    }
    if (fuji_globals_count > 0) {
        for (int i = 0; i < fuji_globals_count; i++) {
            gc_mark_value(fuji_globals[i]);
        }
    }
    fuji_mark_module_cache();
    fuji_mark_open_upvalues();
}

static void gc_unmark_all(void) {
    for (Obj* obj = gc_state.objects; obj != NULL; obj = obj->next) {
        obj->is_marked = false;
    }
}

void gc_sweep(void) {
    Obj** previous = &gc_state.objects;
    Obj* obj = gc_state.objects;
    while (obj != NULL) {
        if (!obj->is_marked) {
            *previous = obj->next;
            free_object(obj);
            obj = *previous;
        } else {
            previous = &obj->next;
            obj = obj->next;
        }
    }
}

void gc_collect(void) {
    if (gc_state.collecting) {
        return;
    }
    gc_state.collecting = true;
    size_t before = gc_state.bytes_allocated;

    gc_unmark_all();
    gc_mark_roots();
    gc_sweep();

    for (Obj* obj = gc_state.objects; obj != NULL; obj = obj->next) {
        obj->generation = GEN_OLD;
    }
    gc_state.remembered_count = 0;

    gc_state.next_gc = gc_state.bytes_allocated * 2;
    if (gc_state.next_gc < 1024u * 1024u) {
        gc_state.next_gc = 1024u * 1024u;
    }

    size_t freed = 0;
    if (before > gc_state.bytes_allocated) {
        freed = before - gc_state.bytes_allocated;
    }
    gc_state.stats.bytes_freed += freed;
    gc_state.stats.collections++;
    gc_state.stats.bytes_allocated = gc_state.bytes_allocated;

    if (!nursery_has_live_object()) {
        gc_reset_nursery();
    }

    gc_state.collecting = false;
}

void gc_init(void) {
    gc_state.objects = NULL;
    gc_state.bytes_allocated = 0;
    gc_state.next_gc = 1024u * 1024u;
    gc_state.gc_disabled = false;
    gc_state.collecting = false;
    gc_state.use_shadow_stack = false;
    gc_state.nursery_top = gc_state.nursery_buf;
    gc_state.nursery_end = gc_state.nursery_buf + sizeof(gc_state.nursery_buf);
    gc_state.nursery_live_bytes = 0;
    gc_state.remembered_count = 0;
    memset(&gc_state.stats, 0, sizeof(gc_state.stats));
}

GCStats gc_get_stats(void) {
    gc_state.stats.bytes_allocated = gc_state.bytes_allocated;
    return gc_state.stats;
}

void gc_set_disabled(bool disabled) {
    gc_state.gc_disabled = disabled;
}

bool gc_is_disabled(void) {
    return gc_state.gc_disabled;
}

void gc_set_next_threshold(size_t bytes) {
    gc_state.next_gc = bytes;
}
