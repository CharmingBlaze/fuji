#include "fuji.h"
#include <math.h>
#include <stdlib.h>
#include <time.h>
#ifdef _WIN32
#define WIN32_LEAN_AND_MEAN
#define NOGDI
#define NOUSER
#include <windows.h>
#else
#include <unistd.h>
#endif

#define STACK_MAX 1024
#define GRAY_STACK_MAX 1024
#define YOUNG_GEN_SIZE (1 * 1024 * 1024) // 1MB for demo
#define PROMOTION_AGE 3

typedef struct FujiTask {
    FujiValue closure;
    double target_time;
    double interval; // 0 for setTimeout, >0 for setInterval
    struct FujiTask* next;
} FujiTask;

typedef struct {
    FujiObj* objects;
    FujiValue stack[STACK_MAX];
    FujiValue* sp;
    
    FujiObj* gray_stack[GRAY_STACK_MAX];
    int gray_count;
    
    size_t bytes_allocated;
    size_t next_gc;
    
    double program_start;
    double last_frame;
    
    FujiTask* tasks;
} FujiHeap;

static double fuji_get_time_seconds() {
#ifdef _WIN32
    LARGE_INTEGER freq, counter;
    QueryPerformanceFrequency(&freq);
    QueryPerformanceCounter(&counter);
    return (double)counter.QuadPart / freq.QuadPart;
#else
    struct timespec ts;
    clock_gettime(CLOCK_MONOTONIC, &ts);
    return (double)ts.tv_sec + (double)ts.tv_nsec / 1e9;
#endif
}

static FujiHeap heap;

void fuji_init() {
    heap.objects = NULL;
    heap.sp = heap.stack;
    heap.bytes_allocated = 0;
    heap.next_gc = 10 * 1024 * 1024; // 10MB initial
    heap.gray_count = 0;
    
    heap.program_start = fuji_get_time_seconds();
    heap.last_frame = heap.program_start;
    heap.tasks = NULL;
    
    srand((unsigned int)time(NULL));
}

void fuji_runtime_init(void) {
    fuji_init();
}

void fuji_runtime_shutdown(void) {
    fuji_shutdown();
}

void fuji_push(FujiValue value) {
    if (heap.sp >= heap.stack + STACK_MAX) {
        fprintf(stderr, "Fuji: stack overflow\n");
        exit(1);
    }
    *heap.sp = value;
    heap.sp++;
}

FujiValue fuji_pop() {
    heap.sp--;
    return *heap.sp;
}

// Internal marking
static void gray_object(FujiObj* obj) {
    if (obj == NULL || obj->mark) return;
    obj->mark = 1;
    
    if (heap.gray_count >= GRAY_STACK_MAX) {
        // Fallback to recursive or just exit if stack is full
        // Production would reallocate gray stack
        return;
    }
    heap.gray_stack[heap.gray_count++] = obj;
}

static void blacken_object(FujiObj* obj) {
    switch (obj->type) {
        case OBJ_ARRAY:
        case OBJ_SET:
        case OBJ_TUPLE: {
            FujiArray* arr = (FujiArray*)obj;
            for (int i = 0; i < arr->count; i++) {
                if (IS_OBJ(arr->elements[i])) gray_object(AS_OBJ(arr->elements[i]));
            }
            break;
        }
        case OBJ_OBJECT:
        case OBJ_MAP: {
            FujiObject* ko = (FujiObject*)obj;
            for (int i = 0; i < ko->count; i++) {
                if (IS_OBJ(ko->entries[i].key)) gray_object(AS_OBJ(ko->entries[i].key));
                if (IS_OBJ(ko->entries[i].value)) gray_object(AS_OBJ(ko->entries[i].value));
            }
            break;
        }
        case OBJ_CLOSURE: {
            FujiClosure* closure = (FujiClosure*)obj;
            for (int i = 0; i < closure->upvalueCount; i++) {
                if (closure->upvalues[i]) gray_object((FujiObj*)closure->upvalues[i]);
            }
            break;
        }
        case OBJ_UPVALUE: {
            FujiUpvalue* upvalue = (FujiUpvalue*)obj;
            if (IS_OBJ(upvalue->closed)) gray_object(AS_OBJ(upvalue->closed));
            break;
        }
        case OBJ_BOUND_METHOD: {
            FujiBoundMethod* bound = (FujiBoundMethod*)obj;
            if (IS_OBJ(bound->receiver)) gray_object(AS_OBJ(bound->receiver));
            if (IS_OBJ(bound->method)) gray_object(AS_OBJ(bound->method));
            break;
        }
        default: break;
    }
}

void fuji_collect() {
    // Mark roots
    for (FujiValue* v = heap.stack; v < heap.sp; v++) {
        if (IS_OBJ(*v)) gray_object(AS_OBJ(*v));
    }
    
    // Mark tasks
    FujiTask* task = heap.tasks;
    while (task) {
        if (IS_OBJ(task->closure)) gray_object(AS_OBJ(task->closure));
        task = task->next;
    }
    
    // Trace references
    while (heap.gray_count > 0) {
        FujiObj* obj = heap.gray_stack[--heap.gray_count];
        blacken_object(obj);
    }
    
    // Sweep
    FujiObj** prev = &heap.objects;
    while (*prev) {
        FujiObj* obj = *prev;
        if (!obj->mark) {
            *prev = obj->next;
            // Free nested memory
            if (obj->type == OBJ_ARRAY) free(((FujiArray*)obj)->elements);
            if (obj->type == OBJ_OBJECT) free(((FujiObject*)obj)->entries);
            if (obj->type == OBJ_CLOSURE) free(((FujiClosure*)obj)->upvalues);
            free(obj);
        } else {
            obj->mark = 0; // Reset for next GC
            prev = &obj->next;
        }
    }
    
    heap.next_gc = heap.bytes_allocated * 2;
}

FujiObj* fuji_alloc(FujiObjType type, size_t size) {
    if (heap.bytes_allocated + size > heap.next_gc) {
        fuji_collect();
    }
    
    FujiObj* obj = (FujiObj*)malloc(size);
    obj->type = type;
    obj->mark = 0;
    obj->next = heap.objects;
    heap.objects = obj;
    
    heap.bytes_allocated += size;
    return obj;
}

FujiArray* fuji_new_array() {
    FujiArray* arr = (FujiArray*)fuji_alloc(OBJ_ARRAY, sizeof(FujiArray));
    arr->count = 0;
    arr->capacity = 0;
    arr->elements = NULL;
    return arr;
}

FujiObject* fuji_alloc_object() {
    FujiObject* obj = (FujiObject*)fuji_alloc(OBJ_OBJECT, sizeof(FujiObject));
    obj->count = 0;
    obj->capacity = 0;
    obj->entries = NULL;
    return obj;
}

FujiUpvalue* fuji_new_upvalue(FujiValue* slot) {
    FujiUpvalue* upvalue = (FujiUpvalue*)fuji_alloc(OBJ_UPVALUE, sizeof(FujiUpvalue));
    if (slot == NULL) {
        upvalue->location = &upvalue->closed;
    } else {
        upvalue->location = slot;
    }
    upvalue->closed = NULL_VAL;
    upvalue->next = NULL;
    return upvalue;
}

FujiClosure* fuji_new_closure(FujiFn fn, int upvalueCount) {
    size_t size = sizeof(FujiClosure) + sizeof(FujiUpvalue*) * upvalueCount;
    FujiClosure* closure = (FujiClosure*)fuji_alloc(OBJ_CLOSURE, size);
    closure->fn = fn;
    closure->upvalueCount = upvalueCount;
    closure->upvalues = (FujiUpvalue**)((char*)closure + sizeof(FujiClosure));
    for (int i = 0; i < upvalueCount; i++) closure->upvalues[i] = NULL;
    return closure;
}

void fuji_upvalue_set(FujiUpvalue* up, FujiValue val) {
    *up->location = val;
}

FujiValue fuji_upvalue_get(FujiUpvalue* up) {
    return *up->location;
}

void fuji_closure_set_upvalue(FujiValue closure, int index, FujiUpvalue* u) {
    AS_CLOSURE(closure)->upvalues[index] = u;
}

FujiValue fuji_obj_as_value(void* obj) {
    return OBJ_VAL((FujiObj*)obj);
}

void fuji_assert(FujiValue cond, FujiValue msg) {
    if (fuji_is_truthy(cond)) return;
    fprintf(stderr, "assertion failed");
    if (!IS_NULL(msg)) {
        fprintf(stderr, ": ");
        fuji_print_no_newline(msg);
    }
    fprintf(stderr, "\n");
    exit(1);
}
void fuji_abort_arg_error(int kind, int a, int b) {
    switch (kind) {
    case 0:
        fprintf(stderr, "expected at least %d arguments but got %d\n", a, b);
        break;
    case 1:
        fprintf(stderr, "expected at most %d arguments but got %d\n", a, b);
        break;
    case 2:
        fprintf(stderr, "expected %d arguments but got %d\n", a, b);
        break;
    case 3:
        fprintf(stderr, "missing default for parameter %d\n", a);
        break;
    default:
        fprintf(stderr, "internal: bad arg error kind %d\n", kind);
        break;
    }
    exit(1);
}

FujiValue fuji_argv_slice_to_array(FujiValue* argv, int start, int count) {
    FujiArray* arr = fuji_new_array();
    FujiValue wrapper = fuji_obj_as_value(arr);
    if (count <= 0 || argv == NULL) {
        return wrapper;
    }
    for (int i = 0; i < count; i++) {
        fuji_array_push(wrapper, argv[start + i]);
    }
    return wrapper;
}

FujiValue fuji_new_bound_method(FujiValue receiver, FujiValue method) {
    FujiBoundMethod* bound = (FujiBoundMethod*)fuji_alloc(OBJ_BOUND_METHOD, sizeof(FujiBoundMethod));
    bound->receiver = receiver;
    bound->method = method;
    return OBJ_VAL(bound);
}

FujiValue fuji_new_native(FujiNativeFn fn) {
    FujiNative* native = (FujiNative*)fuji_alloc(OBJ_NATIVE, sizeof(FujiNative));
    native->fn = fn;
    return OBJ_VAL(native);
}

void fuji_array_push(FujiValue array, FujiValue value) {
    if (!IS_OBJ(array)) return;
    FujiObjType t = AS_OBJ(array)->type;
    if (t != OBJ_ARRAY && t != OBJ_SET) return;
    FujiArray* arr = (FujiArray*)AS_OBJ(array);
    if (arr->count >= arr->capacity) {
        arr->capacity = arr->capacity < 8 ? 8 : arr->capacity * 2;
        arr->elements = (FujiValue*)realloc(arr->elements, sizeof(FujiValue) * arr->capacity);
    }
    arr->elements[arr->count++] = value;
    FUJI_WRITE_BARRIER(arr, value);
}

// --- Native Methods ---

FujiValue fuji_method_string_upper(int argCount, FujiValue* args) {
    FujiValue receiver = args[0]; // Bound method passes receiver as first arg
    if (!IS_OBJ(receiver) || AS_OBJ(receiver)->type != OBJ_STRING) return NULL_VAL;
    FujiString* s = (FujiString*)AS_OBJ(receiver);
    FujiString* res = fuji_copy_string(s->chars, s->length);
    for (int i = 0; i < res->length; i++) {
        if (res->chars[i] >= 'a' && res->chars[i] <= 'z') res->chars[i] -= 32;
    }
    return OBJ_VAL(res);
}

FujiValue fuji_method_array_push(int argCount, FujiValue* args) {
    FujiValue receiver = args[0];
    if (!IS_OBJ(receiver) || AS_OBJ(receiver)->type != OBJ_ARRAY) return NULL_VAL;
    FujiValue val = args[1];
    fuji_array_push(receiver, val);
    return receiver;
}

FujiValue fuji_method_array_pop(int argCount, FujiValue* args) {
    FujiValue receiver = args[0];
    if (!IS_OBJ(receiver) || AS_OBJ(receiver)->type != OBJ_ARRAY) return NULL_VAL;
    FujiArray* arr = (FujiArray*)AS_OBJ(receiver);
    if (arr->count <= 0) return NULL_VAL;
    arr->count--;
    return arr->elements[arr->count];
}

FujiValue fuji_method_array_length(int argCount, FujiValue* args) {
    FujiValue receiver = args[0];
    if (!IS_OBJ(receiver) || AS_OBJ(receiver)->type != OBJ_ARRAY) return NULL_VAL;
    FujiArray* arr = (FujiArray*)AS_OBJ(receiver);
    return NUMBER_VAL((double)arr->count);
}

void fuji_object_set(FujiValue obj, FujiValue key, FujiValue value) {
    if (!IS_OBJ(obj)) return;
    FujiObjType t = AS_OBJ(obj)->type;
    if (t != OBJ_OBJECT && t != OBJ_MAP) return;
    FujiObject* o = (FujiObject*)AS_OBJ(obj);
    
    // Simple linear search for now, can be optimized later
    for (int i = 0; i < o->count; i++) {
        if (AS_BOOL(fuji_eq(o->entries[i].key, key))) {
            o->entries[i].value = value;
            return;
        }
    }
    
    if (o->count == o->capacity) {
        o->capacity = o->capacity == 0 ? 8 : o->capacity * 2;
        o->entries = (FujiEntry*)realloc(o->entries, sizeof(FujiEntry) * o->capacity);
    }
    
    o->entries[o->count].key = key;
    o->entries[o->count].value = value;
    o->count++;
}

FujiValue fuji_object_delete(FujiValue obj, FujiValue key) {
    if (!IS_OBJ(obj)) return BOOL_VAL(false);
    FujiObjType t = AS_OBJ(obj)->type;
    if (t != OBJ_OBJECT && t != OBJ_MAP) return BOOL_VAL(false);
    FujiObject* o = (FujiObject*)AS_OBJ(obj);
    for (int i = 0; i < o->count; i++) {
        if (AS_BOOL(fuji_eq(o->entries[i].key, key))) {
            o->entries[i] = o->entries[o->count - 1];
            o->count--;
            return BOOL_VAL(true);
        }
    }
    return BOOL_VAL(false);
}

FujiValue fuji_map_method_set(int argCount, FujiValue* args) {
    if (argCount < 3) return NULL_VAL;
    if (!IS_OBJ(args[0]) || AS_OBJ(args[0])->type != OBJ_MAP) return NULL_VAL;
    fuji_object_set(args[0], args[1], args[2]);
    return NULL_VAL;
}

FujiValue fuji_map_method_get(int argCount, FujiValue* args) {
    if (argCount < 2) return NULL_VAL;
    if (!IS_OBJ(args[0]) || AS_OBJ(args[0])->type != OBJ_MAP) return NULL_VAL;
    FujiObject* m = (FujiObject*)AS_OBJ(args[0]);
    FujiValue key = args[1];
    for (int i = 0; i < m->count; i++) {
        if (AS_BOOL(fuji_eq(m->entries[i].key, key))) return m->entries[i].value;
    }
    return NULL_VAL;
}

FujiValue fuji_map_method_has(int argCount, FujiValue* args) {
    if (argCount < 2) return BOOL_VAL(false);
    if (!IS_OBJ(args[0]) || AS_OBJ(args[0])->type != OBJ_MAP) return BOOL_VAL(false);
    FujiObject* m = (FujiObject*)AS_OBJ(args[0]);
    FujiValue key = args[1];
    for (int i = 0; i < m->count; i++) {
        if (AS_BOOL(fuji_eq(m->entries[i].key, key))) return BOOL_VAL(true);
    }
    return BOOL_VAL(false);
}

FujiValue fuji_map_method_delete(int argCount, FujiValue* args) {
    if (argCount < 2) return BOOL_VAL(false);
    return fuji_object_delete(args[0], args[1]);
}

FujiValue fuji_map_method_size(int argCount, FujiValue* args) {
    if (argCount < 1) return NUMBER_VAL(0);
    if (!IS_OBJ(args[0]) || AS_OBJ(args[0])->type != OBJ_MAP) return NUMBER_VAL(0);
    FujiObject* m = (FujiObject*)AS_OBJ(args[0]);
    return NUMBER_VAL((double)m->count);
}

static int fuji_set_index_of(FujiArray* s, FujiValue v) {
    for (int i = 0; i < s->count; i++) {
        if (AS_BOOL(fuji_eq(s->elements[i], v))) return i;
    }
    return -1;
}

FujiValue fuji_set_method_add(int argCount, FujiValue* args) {
    if (argCount < 2) return NULL_VAL;
    if (!IS_OBJ(args[0]) || AS_OBJ(args[0])->type != OBJ_SET) return NULL_VAL;
    FujiArray* s = (FujiArray*)AS_OBJ(args[0]);
    if (fuji_set_index_of(s, args[1]) >= 0) return NULL_VAL;
    fuji_array_push(args[0], args[1]);
    return NULL_VAL;
}

FujiValue fuji_set_method_has(int argCount, FujiValue* args) {
    if (argCount < 2) return BOOL_VAL(false);
    if (!IS_OBJ(args[0]) || AS_OBJ(args[0])->type != OBJ_SET) return BOOL_VAL(false);
    return BOOL_VAL(fuji_set_index_of((FujiArray*)AS_OBJ(args[0]), args[1]) >= 0);
}

FujiValue fuji_set_method_remove(int argCount, FujiValue* args) {
    if (argCount < 2) return BOOL_VAL(false);
    if (!IS_OBJ(args[0]) || AS_OBJ(args[0])->type != OBJ_SET) return BOOL_VAL(false);
    FujiArray* s = (FujiArray*)AS_OBJ(args[0]);
    int i = fuji_set_index_of(s, args[1]);
    if (i < 0) return BOOL_VAL(false);
    s->elements[i] = s->elements[s->count - 1];
    s->count--;
    return BOOL_VAL(true);
}

FujiValue fuji_set_method_size(int argCount, FujiValue* args) {
    if (argCount < 1) return NUMBER_VAL(0);
    if (!IS_OBJ(args[0]) || AS_OBJ(args[0])->type != OBJ_SET) return NUMBER_VAL(0);
    return NUMBER_VAL((double)((FujiArray*)AS_OBJ(args[0]))->count);
}

FujiValue fuji_map_new(int argCount, FujiValue* args) {
    (void)args;
    if (argCount != 0) return NULL_VAL;
    FujiObject* m = (FujiObject*)fuji_alloc(OBJ_MAP, sizeof(FujiObject));
    m->count = 0;
    m->capacity = 0;
    m->entries = NULL;
    return OBJ_VAL((FujiObj*)m);
}

FujiValue fuji_tuple_new(int argCount, FujiValue* args) {
    if (argCount < 2) return NULL_VAL;
    FujiArray* arr = (FujiArray*)fuji_alloc(OBJ_TUPLE, sizeof(FujiArray));
    arr->capacity = argCount;
    arr->count = argCount;
    arr->elements = (FujiValue*)malloc(sizeof(FujiValue) * (size_t)argCount);
    if (!arr->elements) return NULL_VAL;
    memcpy(arr->elements, args, sizeof(FujiValue) * (size_t)argCount);
    return OBJ_VAL((FujiObj*)arr);
}

FujiValue fuji_input(int argCount, FujiValue* args) {
    if (argCount >= 1 && IS_OBJ(args[0]) && AS_OBJ(args[0])->type == OBJ_STRING) {
        FujiString* p = (FujiString*)AS_OBJ(args[0]);
        printf("%.*s", p->length, p->chars);
        fflush(stdout);
    }
    char buf[4096];
    if (fgets(buf, sizeof buf, stdin) == NULL) {
        return OBJ_VAL(fuji_copy_string("", 0));
    }
    size_t n = strlen(buf);
    while (n > 0 && (buf[n - 1] == '\n' || buf[n - 1] == '\r')) {
        buf[--n] = '\0';
    }
    return OBJ_VAL(fuji_copy_string(buf, (int)n));
}

FujiValue fuji_set_new(int argCount, FujiValue* args) {
    if (argCount > 1) return NULL_VAL;
    FujiArray* s = (FujiArray*)fuji_alloc(OBJ_SET, sizeof(FujiArray));
    s->count = 0;
    s->capacity = 0;
    s->elements = NULL;
    FujiValue sv = OBJ_VAL((FujiObj*)s);
    if (argCount >= 1) {
        if (!IS_OBJ(args[0]) || AS_OBJ(args[0])->type != OBJ_ARRAY) return NULL_VAL;
        FujiArray* a = (FujiArray*)AS_OBJ(args[0]);
        for (int i = 0; i < a->count; i++) {
            FujiValue v = a->elements[i];
            if (fuji_set_index_of(s, v) < 0) fuji_array_push(sv, v);
        }
    }
    return sv;
}

FujiValue fuji_get_index(FujiValue obj, FujiValue index) {
    if (!IS_OBJ(obj)) return NULL_VAL;
    FujiObj* o = AS_OBJ(obj);
    FujiValue recv = OBJ_VAL(o);

    if (o->type == OBJ_TUPLE) {
        if (IS_NUMBER(index)) {
            FujiArray* arr = (FujiArray*)o;
            int idx = (int)AS_NUMBER(index);
            if (idx < 0) idx = arr->count + idx;
            if (idx >= 0 && idx < arr->count) return arr->elements[idx];
            return NULL_VAL;
        }
        return NULL_VAL;
    }

    if (o->type == OBJ_ARRAY) {
        if (IS_NUMBER(index)) {
            FujiArray* arr = (FujiArray*)o;
            int idx = (int)AS_NUMBER(index);
            if (idx < 0) idx = arr->count + idx;
            if (idx >= 0 && idx < arr->count) return arr->elements[idx];
            return NULL_VAL;
        }
    } else if (o->type == OBJ_STRING) {
        if (IS_NUMBER(index)) {
            FujiString* s = (FujiString*)o;
            int idx = (int)AS_NUMBER(index);
            if (idx < 0) idx = s->length + idx;
            if (idx >= 0 && idx < s->length) {
                return OBJ_VAL(fuji_copy_string(s->chars + idx, 1));
            }
            return NULL_VAL;
        }
    } else if (o->type == OBJ_OBJECT || o->type == OBJ_MAP) {
        FujiObject* mobj = (FujiObject*)o;
        for (int i = 0; i < mobj->count; i++) {
            if (AS_BOOL(fuji_eq(mobj->entries[i].key, index))) {
                FujiValue val = mobj->entries[i].value;
                if (IS_OBJ(val) && (AS_OBJ(val)->type == OBJ_CLOSURE || AS_OBJ(val)->type == OBJ_NATIVE)) {
                    return fuji_new_bound_method(obj, val);
                }
                return val;
            }
        }
        if (o->type == OBJ_MAP && IS_OBJ(index) && AS_OBJ(index)->type == OBJ_STRING) {
            FujiString* s = (FujiString*)AS_OBJ(index);
            if (strcmp(s->chars, "set") == 0) return fuji_new_bound_method(recv, fuji_new_native(fuji_map_method_set));
            if (strcmp(s->chars, "get") == 0) return fuji_new_bound_method(recv, fuji_new_native(fuji_map_method_get));
            if (strcmp(s->chars, "has") == 0) return fuji_new_bound_method(recv, fuji_new_native(fuji_map_method_has));
            if (strcmp(s->chars, "delete") == 0) return fuji_new_bound_method(recv, fuji_new_native(fuji_map_method_delete));
            if (strcmp(s->chars, "size") == 0) return fuji_new_bound_method(recv, fuji_new_native(fuji_map_method_size));
        }
    } else if (o->type == OBJ_SET && IS_OBJ(index) && AS_OBJ(index)->type == OBJ_STRING) {
        FujiString* s = (FujiString*)AS_OBJ(index);
        if (strcmp(s->chars, "add") == 0) return fuji_new_bound_method(recv, fuji_new_native(fuji_set_method_add));
        if (strcmp(s->chars, "has") == 0) return fuji_new_bound_method(recv, fuji_new_native(fuji_set_method_has));
        if (strcmp(s->chars, "remove") == 0) return fuji_new_bound_method(recv, fuji_new_native(fuji_set_method_remove));
        if (strcmp(s->chars, "size") == 0) return fuji_new_bound_method(recv, fuji_new_native(fuji_set_method_size));
    }

    if (o->type == OBJ_STRING) {
        if (IS_OBJ(index) && AS_OBJ(index)->type == OBJ_STRING) {
            FujiString* strk = (FujiString*)AS_OBJ(index);
            if (strcmp(strk->chars, "upper") == 0) {
                return fuji_new_bound_method(recv, fuji_new_native(fuji_method_string_upper));
            }
        }
    } else if (o->type == OBJ_ARRAY) {
        if (IS_OBJ(index) && AS_OBJ(index)->type == OBJ_STRING) {
            FujiString* strk = (FujiString*)AS_OBJ(index);
            if (strcmp(strk->chars, "length") == 0) {
                return fuji_new_bound_method(recv, fuji_new_native(fuji_method_array_length));
            }
            if (strcmp(strk->chars, "push") == 0) {
                return fuji_new_bound_method(recv, fuji_new_native(fuji_method_array_push));
            }
            if (strcmp(strk->chars, "pop") == 0) {
                return fuji_new_bound_method(recv, fuji_new_native(fuji_method_array_pop));
            }
        }
    }

    return NULL_VAL;
}

FujiValue fuji_call(FujiValue callee, int argCount, FujiValue* args) {
    if (!IS_OBJ(callee)) {
        fprintf(stderr, "Fuji: callee is not an object\n");
        return NULL_VAL;
    }
    FujiObj* obj = AS_OBJ(callee);
    if (obj->type == OBJ_CLOSURE) {
        FujiClosure* closure = (FujiClosure*)obj;
        return closure->fn(NULL_VAL, argCount, args, closure->upvalues);
    } else if (obj->type == OBJ_NATIVE) {
        return ((FujiNative*)obj)->fn(argCount, args);
    } else if (obj->type == OBJ_BOUND_METHOD) {
        FujiBoundMethod* bound = (FujiBoundMethod*)obj;
        if (IS_OBJ(bound->method) && AS_OBJ(bound->method)->type == OBJ_CLOSURE) {
            FujiClosure* closure = (FujiClosure*)AS_OBJ(bound->method);
            return closure->fn(bound->receiver, argCount, args, closure->upvalues);
        } else if (IS_OBJ(bound->method) && AS_OBJ(bound->method)->type == OBJ_NATIVE) {
            FujiNative* native = (FujiNative*)AS_OBJ(bound->method);
            // Native functions still expect receiver as first arg
            FujiValue* newArgs = (FujiValue*)malloc(sizeof(FujiValue) * (argCount + 1));
            newArgs[0] = bound->receiver;
            memcpy(newArgs + 1, args, sizeof(FujiValue) * argCount);
            FujiValue res = native->fn(argCount + 1, newArgs);
            free(newArgs);
            return res;
        }
    }
    
    fprintf(stderr, "Fuji: callee type %d is not callable\n", obj->type);
    return NULL_VAL;
}

FujiString* fuji_copy_string(const char* chars, int length) {
    FujiString* str = (FujiString*)fuji_alloc(OBJ_STRING, sizeof(FujiString) + length + 1);
    str->length = length;
    memcpy(str->chars, chars, length);
    str->chars[length] = '\0';
    return str;
}

FujiValue fuji_allocate_string(int length, const char* chars) {
    return OBJ_VAL(fuji_copy_string(chars, length));
}

void fuji_write_barrier(FujiObj* obj, FujiValue value) {
    (void)obj; (void)value;
    // No-op for Mark-Sweep GC
}

// --- Operations ---

FujiValue fuji_add(FujiValue a, FujiValue b) {
    if (IS_NUMBER(a) && IS_NUMBER(b)) {
        return NUMBER_VAL(AS_NUMBER(a) + AS_NUMBER(b));
    }
    
    if (IS_OBJ(a) && AS_OBJ(a)->type == OBJ_STRING) {
        FujiString* sa = (FujiString*)AS_OBJ(a);
        char buf[64];
        const char* sb;
        int blen;
        
        if (IS_NUMBER(b)) {
            blen = sprintf(buf, "%g", AS_NUMBER(b));
            sb = buf;
        } else if (IS_OBJ(b) && AS_OBJ(b)->type == OBJ_STRING) {
            FujiString* s_obj = (FujiString*)AS_OBJ(b);
            sb = s_obj->chars;
            blen = s_obj->length;
        } else {
            return NULL_VAL;
        }
        
        int total_len = sa->length + blen;
        FujiString* res = (FujiString*)fuji_alloc(OBJ_STRING, sizeof(FujiString) + total_len + 1);
        res->length = total_len;
        memcpy(res->chars, sa->chars, sa->length);
        memcpy(res->chars + sa->length, sb, blen);
        res->chars[total_len] = '\0';
        return OBJ_VAL(res);
    }
    
    return NULL_VAL;
}

FujiValue fuji_to_string_val(FujiValue v) {
    char buf[256];
    int n;
    if (IS_NUMBER(v)) {
        n = snprintf(buf, sizeof(buf), "%g", AS_NUMBER(v));
        if (n < 0) n = 0;
        return OBJ_VAL(fuji_copy_string(buf, n));
    }
    if (IS_BOOL(v)) {
        if (AS_BOOL(v)) {
            return OBJ_VAL(fuji_copy_string("true", 4));
        }
        return OBJ_VAL(fuji_copy_string("false", 5));
    }
    if (IS_NULL(v)) {
        return OBJ_VAL(fuji_copy_string("null", 4));
    }
    if (IS_OBJ(v) && AS_OBJ(v)->type == OBJ_STRING) {
        return v;
    }
    n = snprintf(buf, sizeof(buf), "(value)");
    return OBJ_VAL(fuji_copy_string(buf, n));
}

FujiValue fuji_parse_number(FujiValue v) {
    if (IS_NUMBER(v)) {
        return v;
    }
    if (IS_OBJ(v) && AS_OBJ(v)->type == OBJ_STRING) {
        FujiString* s = (FujiString*)AS_OBJ(v);
        char* end = NULL;
        double d = strtod(s->chars, &end);
        if (end == s->chars) {
            return NULL_VAL;
        }
        while (*end == ' ' || *end == '\t' || *end == '\n' || *end == '\r') {
            end++;
        }
        if (*end != '\0') {
            return NULL_VAL;
        }
        return NUMBER_VAL(d);
    }
    return NULL_VAL;
}

FujiValue fuji_range(FujiValue from, FujiValue to) {
    if (!IS_NUMBER(from) || !IS_NUMBER(to)) {
        return NULL_VAL;
    }
    long long lo = (long long)trunc(AS_NUMBER(from));
    long long hi = (long long)trunc(AS_NUMBER(to));
    FujiValue arr = OBJ_VAL(fuji_new_array());
    if (lo <= hi) {
        for (long long i = lo; i <= hi; i++) {
            fuji_array_push(arr, NUMBER_VAL((double)i));
        }
    } else {
        for (long long i = lo; i >= hi; i--) {
            fuji_array_push(arr, NUMBER_VAL((double)i));
        }
    }
    return arr;
}

FujiValue fuji_is_truthy(FujiValue v) {
    if (IS_NULL(v)) return BOOL_VAL(false);
    if (IS_BOOL(v)) return v;
    if (IS_NUMBER(v)) return BOOL_VAL(AS_NUMBER(v) != 0);
    return BOOL_VAL(true);
}

void fuji_print_no_newline(FujiValue val) {
    if (IS_NUMBER(val)) {
        printf("%g", AS_NUMBER(val));
    } else if (IS_BOOL(val)) {
        printf("%s", AS_BOOL(val) ? "true" : "false");
    } else if (IS_NULL(val)) {
        printf("null");
    } else if (IS_OBJ(val)) {
        FujiObj* obj = AS_OBJ(val);
        if (obj->type == OBJ_STRING) {
            printf("%s", ((FujiString*)obj)->chars);
        } else if (obj->type == OBJ_ARRAY) {
            FujiArray* arr = (FujiArray*)obj;
            printf("[");
            for (int i = 0; i < arr->count; i++) {
                fuji_print_no_newline(arr->elements[i]);
                if (i < arr->count - 1) printf(", ");
            }
            printf("]");
        } else {
            printf("<obj %p type %d>", obj, obj->type);
        }
    }
}

void fuji_print_space(void) { printf(" "); }

void fuji_print_newline(void) { printf("\n"); }

FujiValue fuji_print(FujiValue val) {
    fuji_print_no_newline(val);
    fuji_print_newline();
    return NULL_VAL;
}

// Rest of arithmetic...
FujiValue fuji_sub(FujiValue a, FujiValue b) { return IS_NUMBER(a) && IS_NUMBER(b) ? NUMBER_VAL(AS_NUMBER(a) - AS_NUMBER(b)) : NULL_VAL; }
FujiValue fuji_mul(FujiValue a, FujiValue b) { return IS_NUMBER(a) && IS_NUMBER(b) ? NUMBER_VAL(AS_NUMBER(a) * AS_NUMBER(b)) : NULL_VAL; }
FujiValue fuji_div(FujiValue a, FujiValue b) { return IS_NUMBER(a) && IS_NUMBER(b) ? NUMBER_VAL(AS_NUMBER(a) / AS_NUMBER(b)) : NULL_VAL; }
FujiValue fuji_gt(FujiValue a, FujiValue b) { return IS_NUMBER(a) && IS_NUMBER(b) ? BOOL_VAL(AS_NUMBER(a) > AS_NUMBER(b)) : BOOL_VAL(false); }
FujiValue fuji_lt(FujiValue a, FujiValue b) { return IS_NUMBER(a) && IS_NUMBER(b) ? BOOL_VAL(AS_NUMBER(a) < AS_NUMBER(b)) : BOOL_VAL(false); }
FujiValue fuji_eq(FujiValue a, FujiValue b) {
    if (a == b) return BOOL_VAL(true);
    if (IS_OBJ(a) && IS_OBJ(b)) {
        FujiObj* oa = AS_OBJ(a);
        FujiObj* ob = AS_OBJ(b);
        if (oa->type == OBJ_STRING && ob->type == OBJ_STRING) {
            FujiString* sa = (FujiString*)oa;
            FujiString* sb = (FujiString*)ob;
            return BOOL_VAL(sa->length == sb->length && memcmp(sa->chars, sb->chars, sa->length) == 0);
        }
    }
    return BOOL_VAL(false);
}

FujiValue fuji_neq(FujiValue a, FujiValue b) { return BOOL_VAL(!AS_BOOL(fuji_eq(a, b))); }
FujiValue fuji_le(FujiValue a, FujiValue b) { return IS_NUMBER(a) && IS_NUMBER(b) ? BOOL_VAL(AS_NUMBER(a) <= AS_NUMBER(b)) : BOOL_VAL(false); }
FujiValue fuji_ge(FujiValue a, FujiValue b) { return IS_NUMBER(a) && IS_NUMBER(b) ? BOOL_VAL(AS_NUMBER(a) >= AS_NUMBER(b)) : BOOL_VAL(false); }
FujiValue fuji_mod(FujiValue a, FujiValue b) { return IS_NUMBER(a) && IS_NUMBER(b) ? NUMBER_VAL(fmod(AS_NUMBER(a), AS_NUMBER(b))) : NULL_VAL; }
FujiValue fuji_pow(FujiValue a, FujiValue b) { return IS_NUMBER(a) && IS_NUMBER(b) ? NUMBER_VAL(pow(AS_NUMBER(a), AS_NUMBER(b))) : NULL_VAL; }
FujiValue fuji_negate(FujiValue a) { return IS_NUMBER(a) ? NUMBER_VAL(-AS_NUMBER(a)) : NULL_VAL; }

FujiValue fuji_bit_and(FujiValue a, FujiValue b) { return IS_NUMBER(a) && IS_NUMBER(b) ? NUMBER_VAL((double)((int64_t)AS_NUMBER(a) & (int64_t)AS_NUMBER(b))) : NULL_VAL; }
FujiValue fuji_bit_or(FujiValue a, FujiValue b) { return IS_NUMBER(a) && IS_NUMBER(b) ? NUMBER_VAL((double)((int64_t)AS_NUMBER(a) | (int64_t)AS_NUMBER(b))) : NULL_VAL; }
FujiValue fuji_bit_xor(FujiValue a, FujiValue b) { return IS_NUMBER(a) && IS_NUMBER(b) ? NUMBER_VAL((double)((int64_t)AS_NUMBER(a) ^ (int64_t)AS_NUMBER(b))) : NULL_VAL; }
FujiValue fuji_bit_not(FujiValue a) { return IS_NUMBER(a) ? NUMBER_VAL((double)(~(int64_t)AS_NUMBER(a))) : NULL_VAL; }
FujiValue fuji_shl(FujiValue a, FujiValue b) { return IS_NUMBER(a) && IS_NUMBER(b) ? NUMBER_VAL((double)((int64_t)AS_NUMBER(a) << (int64_t)AS_NUMBER(b))) : NULL_VAL; }
FujiValue fuji_shr(FujiValue a, FujiValue b) { return IS_NUMBER(a) && IS_NUMBER(b) ? NUMBER_VAL((double)((int64_t)AS_NUMBER(a) >> (int64_t)AS_NUMBER(b))) : NULL_VAL; }
FujiValue fuji_ushr(FujiValue a, FujiValue b) { return IS_NUMBER(a) && IS_NUMBER(b) ? NUMBER_VAL((double)((uint64_t)(int64_t)AS_NUMBER(a) >> (int64_t)AS_NUMBER(b))) : NULL_VAL; }

void fuji_set_index(FujiValue obj, FujiValue index, FujiValue value) {
    if (!IS_OBJ(obj)) return;
    FujiObj* o = AS_OBJ(obj);
    if (o->type == OBJ_ARRAY && IS_NUMBER(index)) {
        FujiArray* arr = (FujiArray*)o;
        int idx = (int)AS_NUMBER(index);
        if (idx < 0) idx = arr->count + idx;
        if (idx >= 0 && idx < arr->count) {
            arr->elements[idx] = value;
            FUJI_WRITE_BARRIER(arr, value);
        }
    } else if (o->type == OBJ_OBJECT || o->type == OBJ_MAP) {
        fuji_object_set(obj, index, value);
    }
}

int fuji_for_in_len(FujiValue iterable) {
    if (!IS_OBJ(iterable)) return 0;
    FujiObj* o = AS_OBJ(iterable);
    if (o->type == OBJ_ARRAY) return ((FujiArray*)o)->count;
    if (o->type == OBJ_OBJECT) return ((FujiObject*)o)->count;
    if (o->type == OBJ_STRING) return ((FujiString*)o)->length;
    return 0;
}

FujiValue fuji_for_in_get(FujiValue iterable, int index) {
    if (!IS_OBJ(iterable)) return NULL_VAL;
    FujiObj* o = AS_OBJ(iterable);
    if (o->type == OBJ_ARRAY) {
        FujiArray* arr = (FujiArray*)o;
        if (index >= 0 && index < arr->count) return NUMBER_VAL((double)index);
        return NULL_VAL;
    }
    if (o->type == OBJ_OBJECT) {
        FujiObject* obj = (FujiObject*)o;
        if (index >= 0 && index < obj->count) return obj->entries[index].key;
        return NULL_VAL;
    }
    if (o->type == OBJ_STRING) {
        FujiString* s = (FujiString*)o;
        if (index >= 0 && index < s->length) return NUMBER_VAL((double)index);
        return NULL_VAL;
    }
    return NULL_VAL;
}

FujiValue fuji_slice(FujiValue obj, FujiValue start_val, FujiValue end_val) {
    if (!IS_OBJ(obj)) return NULL_VAL;
    FujiObj* o = AS_OBJ(obj);
    int length = 0;
    if (o->type == OBJ_ARRAY) length = ((FujiArray*)o)->count;
    else if (o->type == OBJ_STRING) length = ((FujiString*)o)->length;
    else return NULL_VAL;

    int start = 0;
    if (IS_NULL(start_val)) start = 0;
    else {
        start = (int)AS_NUMBER(start_val);
        if (start < 0) start = length + start;
    }

    int end = length;
    if (IS_NULL(end_val)) end = length;
    else {
        end = (int)AS_NUMBER(end_val);
        if (end < 0) end = length + end;
    }

    if (start < 0) start = 0;
    if (end > length) end = length;
    if (start > end) start = end;

    if (o->type == OBJ_ARRAY) {
        FujiArray* src = (FujiArray*)o;
        FujiArray* res = fuji_new_array();
        int count = end - start;
        for (int i = 0; i < count; i++) {
            fuji_array_push(OBJ_VAL(res), src->elements[start + i]);
        }
        return OBJ_VAL(res);
    } else {
        FujiString* src = (FujiString*)o;
        return OBJ_VAL(fuji_copy_string(src->chars + start, end - start));
    }
}

void fuji_shutdown() {
    /* Standalone Fuji executables are process-lifetime runtimes.
       The OS reclaims heap pages on exit; doing a manual sweep here risks
       double-freeing objects already handled by runtime/host finalizers. */
    heap.objects = NULL;
    heap.sp = heap.stack;
}

FujiValue fuji_len(FujiValue val) {
    if (!IS_OBJ(val)) return NUMBER_VAL(0);
    FujiObj* obj = AS_OBJ(val);
    switch (obj->type) {
        case OBJ_STRING: return NUMBER_VAL(((FujiString*)obj)->length);
        case OBJ_ARRAY: return NUMBER_VAL(((FujiArray*)obj)->count);
        case OBJ_OBJECT:
        case OBJ_MAP: return NUMBER_VAL(((FujiObject*)obj)->count);
        case OBJ_SET:
        case OBJ_TUPLE: return NUMBER_VAL(((FujiArray*)obj)->count);
        default: return NUMBER_VAL(0);
    }
}

FujiValue fuji_type(FujiValue val) {
    const char* type_str = "unknown";
    if (IS_NUMBER(val)) type_str = "number";
    else if (IS_BOOL(val)) type_str = "bool";
    else if (IS_NULL(val)) type_str = "null";
    else if (IS_OBJ(val)) {
        switch (AS_OBJ(val)->type) {
            case OBJ_STRING: type_str = "string"; break;
            case OBJ_ARRAY: type_str = "array"; break;
            case OBJ_OBJECT: type_str = "object"; break;
            case OBJ_MAP: type_str = "map"; break;
            case OBJ_SET: type_str = "set"; break;
            case OBJ_TUPLE: type_str = "tuple"; break;
            case OBJ_CLOSURE: type_str = "function"; break;
            case OBJ_NATIVE: type_str = "function"; break;
            case OBJ_BOUND_METHOD: type_str = "function"; break;
            default: type_str = "object"; break;
        }
    }
    return OBJ_VAL(fuji_copy_string(type_str, strlen(type_str)));
}

FujiValue fuji_clock() {
    return NUMBER_VAL((double)clock() / CLOCKS_PER_SEC);
}

FujiValue fuji_wall_time(void) {
    return NUMBER_VAL((double)time(NULL));
}

void fuji_sleep_ms(FujiValue ms_val) {
    if (!IS_NUMBER(ms_val)) return;
    double ms = AS_NUMBER(ms_val);
    if (ms <= 0) return;
#ifdef _WIN32
    Sleep((DWORD)(ms + 0.5));
#else
    struct timespec ts;
    ts.tv_sec = (time_t)(ms / 1000.0);
    ts.tv_nsec = (long)fmod(ms, 1000.0) * 1000000L;
    if (ts.tv_nsec < 0) ts.tv_nsec = 0;
    nanosleep(&ts, NULL);
#endif
}

static int fuji_rand_seeded;

FujiValue fuji_random_unit(void) {
    if (!fuji_rand_seeded) {
        srand((unsigned int)time(NULL));
        fuji_rand_seeded = 1;
    }
    return NUMBER_VAL((double)rand() / (double)RAND_MAX);
}

FujiValue fuji_io_read_file(int argCount, FujiValue* args) {
    if (argCount < 1 || !IS_OBJ(args[0]) || AS_OBJ(args[0])->type != OBJ_STRING) return NULL_VAL;
    const char* path = ((FujiString*)AS_OBJ(args[0]))->chars;
    
    FILE* file = fopen(path, "rb");
    if (file == NULL) return NULL_VAL;
    
    fseek(file, 0L, SEEK_END);
    size_t fileSize = ftell(file);
    rewind(file);
    
    char* buffer = (char*)malloc(fileSize + 1);
    if (buffer == NULL) {
        fclose(file);
        return NULL_VAL;
    }
    
    size_t bytesRead = fread(buffer, sizeof(char), fileSize, file);
    buffer[bytesRead] = '\0';
    fclose(file);
    
    FujiValue res = OBJ_VAL(fuji_copy_string(buffer, bytesRead));
    free(buffer);
    return res;
}

FujiValue fuji_io_write_file(int argCount, FujiValue* args) {
    if (argCount < 2 || !IS_OBJ(args[0]) || AS_OBJ(args[0])->type != OBJ_STRING) return NULL_VAL;
    if (!IS_OBJ(args[1]) || AS_OBJ(args[1])->type != OBJ_STRING) return NULL_VAL;
    
    const char* path = ((FujiString*)AS_OBJ(args[0]))->chars;
    FujiString* content = (FujiString*)AS_OBJ(args[1]);
    
    FILE* file = fopen(path, "wb");
    if (file == NULL) return BOOL_VAL(false);
    
    fwrite(content->chars, sizeof(char), content->length, file);
    fclose(file);
    return BOOL_VAL(true);
}

FujiValue fuji_module_init_io() {
    static FujiValue module_exports = 0;
    if (module_exports != 0) return module_exports;
    module_exports = OBJ_VAL(fuji_alloc_object());
    fuji_object_set(module_exports, OBJ_VAL(fuji_copy_string("read", 4)), fuji_new_native(fuji_io_read_file));
    fuji_object_set(module_exports, OBJ_VAL(fuji_copy_string("write", 5)), fuji_new_native(fuji_io_write_file));
    return module_exports;
}

// --- JSON Module ---

typedef struct {
    char* buffer;
    int capacity;
    int length;
} FujiStringBuilder;

static void sb_init(FujiStringBuilder* sb) {
    sb->capacity = 128;
    sb->length = 0;
    sb->buffer = (char*)malloc(sb->capacity);
}

static void sb_append(FujiStringBuilder* sb, const char* str, int len) {
    if (sb->length + len >= sb->capacity) {
        while (sb->length + len >= sb->capacity) sb->capacity *= 2;
        sb->buffer = (char*)realloc(sb->buffer, sb->capacity);
    }
    memcpy(sb->buffer + sb->length, str, len);
    sb->length += len;
    sb->buffer[sb->length] = '\0';
}

static void fuji_json_stringify_recursive(FujiStringBuilder* sb, FujiValue val) {
    if (IS_NUMBER(val)) {
        char buf[32];
        int len = sprintf(buf, "%g", AS_NUMBER(val));
        sb_append(sb, buf, len);
    } else if (IS_BOOL(val)) {
        if (AS_BOOL(val)) sb_append(sb, "true", 4);
        else sb_append(sb, "false", 5);
    } else if (IS_NULL(val)) {
        sb_append(sb, "null", 4);
    } else if (IS_OBJ(val)) {
        FujiObj* obj = AS_OBJ(val);
        if (obj->type == OBJ_STRING) {
            FujiString* s = (FujiString*)obj;
            sb_append(sb, "\"", 1);
            sb_append(sb, s->chars, s->length);
            sb_append(sb, "\"", 1);
        } else if (obj->type == OBJ_ARRAY) {
            FujiArray* arr = (FujiArray*)obj;
            sb_append(sb, "[", 1);
            for (int i = 0; i < arr->count; i++) {
                fuji_json_stringify_recursive(sb, arr->elements[i]);
                if (i < arr->count - 1) sb_append(sb, ",", 1);
            }
            sb_append(sb, "]", 1);
        } else if (obj->type == OBJ_OBJECT) {
            FujiObject* o = (FujiObject*)obj;
            sb_append(sb, "{", 1);
            for (int i = 0; i < o->count; i++) {
                FujiString* key = (FujiString*)AS_OBJ(o->entries[i].key);
                sb_append(sb, "\"", 1);
                sb_append(sb, key->chars, key->length);
                sb_append(sb, "\":", 2);
                fuji_json_stringify_recursive(sb, o->entries[i].value);
                if (i < o->count - 1) sb_append(sb, ",", 1);
            }
            sb_append(sb, "}", 1);
        }
    }
}

FujiValue fuji_json_stringify(int argCount, FujiValue* args) {
    if (argCount < 1) return NULL_VAL;
    FujiStringBuilder sb;
    sb_init(&sb);
    fuji_json_stringify_recursive(&sb, args[0]);
    FujiValue res = OBJ_VAL(fuji_copy_string(sb.buffer, sb.length));
    free(sb.buffer);
    return res;
}

typedef struct {
    const char* current;
} FujiJsonParser;

static void skip_ws(FujiJsonParser* p) {
    while (*p->current == ' ' || *p->current == '\t' || *p->current == '\n' || *p->current == '\r') p->current++;
}

static FujiValue fuji_json_parse_recursive(FujiJsonParser* p);

static FujiValue fuji_json_parse_string(FujiJsonParser* p) {
    p->current++; // skip "
    const char* start = p->current;
    while (*p->current != '"' && *p->current != '\0') p->current++;
    int len = p->current - start;
    FujiString* s = fuji_copy_string(start, len);
    if (*p->current == '"') p->current++;
    return OBJ_VAL(s);
}

static FujiValue fuji_json_parse_array(FujiJsonParser* p) {
    p->current++; // skip [
    FujiArray* arr = fuji_new_array();
    skip_ws(p);
    while (*p->current != ']' && *p->current != '\0') {
        fuji_array_push(OBJ_VAL(arr), fuji_json_parse_recursive(p));
        skip_ws(p);
        if (*p->current == ',') { p->current++; skip_ws(p); }
    }
    if (*p->current == ']') p->current++;
    return OBJ_VAL(arr);
}

static FujiValue fuji_json_parse_object(FujiJsonParser* p) {
    p->current++; // skip {
    FujiObject* obj = fuji_alloc_object();
    skip_ws(p);
    while (*p->current != '}' && *p->current != '\0') {
        FujiValue key = fuji_json_parse_string(p);
        skip_ws(p);
        if (*p->current == ':') p->current++;
        skip_ws(p);
        fuji_object_set(OBJ_VAL(obj), key, fuji_json_parse_recursive(p));
        skip_ws(p);
        if (*p->current == ',') { p->current++; skip_ws(p); }
    }
    if (*p->current == '}') p->current++;
    return OBJ_VAL(obj);
}

static FujiValue fuji_json_parse_recursive(FujiJsonParser* p) {
    skip_ws(p);
    char c = *p->current;
    if (c == '"') return fuji_json_parse_string(p);
    if (c == '[') return fuji_json_parse_array(p);
    if (c == '{') return fuji_json_parse_object(p);
    if (c == 't') { p->current += 4; return BOOL_VAL(true); }
    if (c == 'f') { p->current += 5; return BOOL_VAL(false); }
    if (c == 'n') { p->current += 4; return NULL_VAL; }
    if ((c >= '0' && c <= '9') || c == '-') {
        char* end;
        double d = strtod(p->current, &end);
        p->current = end;
        return NUMBER_VAL(d);
    }
    return NULL_VAL;
}

FujiValue fuji_json_parse(int argCount, FujiValue* args) {
    if (argCount < 1 || !IS_OBJ(args[0]) || AS_OBJ(args[0])->type != OBJ_STRING) return NULL_VAL;
    FujiJsonParser p;
    p.current = ((FujiString*)AS_OBJ(args[0]))->chars;
    return fuji_json_parse_recursive(&p);
}

FujiValue fuji_module_init_json() {
    static FujiValue module_exports = 0;
    if (module_exports != 0) return module_exports;
    module_exports = OBJ_VAL(fuji_alloc_object());
    fuji_object_set(module_exports, OBJ_VAL(fuji_copy_string("parse", 5)), fuji_new_native(fuji_json_parse));
    fuji_object_set(module_exports, OBJ_VAL(fuji_copy_string("stringify", 9)), fuji_new_native(fuji_json_stringify));
    return module_exports;
}

FujiValue fuji_abs_num(FujiValue v) { if (!IS_NUMBER(v)) return NULL_VAL; return NUMBER_VAL(fabs(AS_NUMBER(v))); }
FujiValue fuji_sqrt_num(FujiValue v) { if (!IS_NUMBER(v)) return NULL_VAL; return NUMBER_VAL(sqrt(AS_NUMBER(v))); }
FujiValue fuji_cbrt_num(FujiValue v) { if (!IS_NUMBER(v)) return NULL_VAL; return NUMBER_VAL(cbrt(AS_NUMBER(v))); }
FujiValue fuji_sin_num(FujiValue v) { if (!IS_NUMBER(v)) return NULL_VAL; return NUMBER_VAL(sin(AS_NUMBER(v))); }
FujiValue fuji_cos_num(FujiValue v) { if (!IS_NUMBER(v)) return NULL_VAL; return NUMBER_VAL(cos(AS_NUMBER(v))); }
FujiValue fuji_tan_num(FujiValue v) { if (!IS_NUMBER(v)) return NULL_VAL; return NUMBER_VAL(tan(AS_NUMBER(v))); }
FujiValue fuji_asin_num(FujiValue v) { if (!IS_NUMBER(v)) return NULL_VAL; return NUMBER_VAL(asin(AS_NUMBER(v))); }
FujiValue fuji_acos_num(FujiValue v) { if (!IS_NUMBER(v)) return NULL_VAL; return NUMBER_VAL(acos(AS_NUMBER(v))); }
FujiValue fuji_atan_num(FujiValue v) { if (!IS_NUMBER(v)) return NULL_VAL; return NUMBER_VAL(atan(AS_NUMBER(v))); }
FujiValue fuji_atan2_num(FujiValue y, FujiValue x) { if (!IS_NUMBER(y) || !IS_NUMBER(x)) return NULL_VAL; return NUMBER_VAL(atan2(AS_NUMBER(y), AS_NUMBER(x))); }
FujiValue fuji_log_num(FujiValue v) { if (!IS_NUMBER(v)) return NULL_VAL; return NUMBER_VAL(log(AS_NUMBER(v))); }
FujiValue fuji_log2_num(FujiValue v) { if (!IS_NUMBER(v)) return NULL_VAL; return NUMBER_VAL(log2(AS_NUMBER(v))); }
FujiValue fuji_log10_num(FujiValue v) { if (!IS_NUMBER(v)) return NULL_VAL; return NUMBER_VAL(log10(AS_NUMBER(v))); }
FujiValue fuji_exp_num(FujiValue v) { if (!IS_NUMBER(v)) return NULL_VAL; return NUMBER_VAL(exp(AS_NUMBER(v))); }
FujiValue fuji_floor_num(FujiValue v) { if (!IS_NUMBER(v)) return NULL_VAL; return NUMBER_VAL(floor(AS_NUMBER(v))); }
FujiValue fuji_ceil_num(FujiValue v) { if (!IS_NUMBER(v)) return NULL_VAL; return NUMBER_VAL(ceil(AS_NUMBER(v))); }
FujiValue fuji_round_num(FujiValue v) { if (!IS_NUMBER(v)) return NULL_VAL; return NUMBER_VAL(round(AS_NUMBER(v))); }
FujiValue fuji_trunc_num(FujiValue v) { if (!IS_NUMBER(v)) return NULL_VAL; return NUMBER_VAL(trunc(AS_NUMBER(v))); }

FujiValue fuji_min_num(int argCount, FujiValue* args) {
    if (argCount == 0) return NULL_VAL;
    double res = AS_NUMBER(args[0]);
    for (int i = 1; i < argCount; i++) {
        double v = AS_NUMBER(args[i]);
        if (v < res) res = v;
    }
    return NUMBER_VAL(res);
}

FujiValue fuji_max_num(int argCount, FujiValue* args) {
    if (argCount == 0) return NULL_VAL;
    double res = AS_NUMBER(args[0]);
    for (int i = 1; i < argCount; i++) {
        double v = AS_NUMBER(args[i]);
        if (v > res) res = v;
    }
    return NUMBER_VAL(res);
}

FujiValue fuji_clamp_num(FujiValue v, FujiValue min, FujiValue max) {
    if (!IS_NUMBER(v) || !IS_NUMBER(min) || !IS_NUMBER(max)) return NULL_VAL;
    double dv = AS_NUMBER(v);
    double dmin = AS_NUMBER(min);
    double dmax = AS_NUMBER(max);
    if (dv < dmin) return min;
    if (dv > dmax) return max;
    return v;
}

FujiValue fuji_lerp_num(FujiValue a, FujiValue b, FujiValue t) {
    if (!IS_NUMBER(a) || !IS_NUMBER(b) || !IS_NUMBER(t)) return NULL_VAL;
    double da = AS_NUMBER(a);
    double db = AS_NUMBER(b);
    double dt = AS_NUMBER(t);
    return NUMBER_VAL(da + (db - da) * dt);
}

FujiValue fuji_smoothstep_num(FujiValue a, FujiValue b, FujiValue t) {
    if (!IS_NUMBER(a) || !IS_NUMBER(b) || !IS_NUMBER(t)) return NULL_VAL;
    double da = AS_NUMBER(a);
    double db = AS_NUMBER(b);
    double tv = AS_NUMBER(t);
    double dt = (tv - da) / (db - da);
    if (dt < 0) dt = 0;
    if (dt > 1) dt = 1;
    return NUMBER_VAL(dt * dt * (3 - 2 * dt));
}

FujiValue fuji_map_num(FujiValue v, FujiValue inMin, FujiValue inMax, FujiValue outMin, FujiValue outMax) {
    if (!IS_NUMBER(v) || !IS_NUMBER(inMin) || !IS_NUMBER(inMax) || !IS_NUMBER(outMin) || !IS_NUMBER(outMax)) return NULL_VAL;
    double dv = AS_NUMBER(v);
    double dim = AS_NUMBER(inMin);
    double dix = AS_NUMBER(inMax);
    double dom = AS_NUMBER(outMin);
    double dox = AS_NUMBER(outMax);
    return NUMBER_VAL(dom + (dox - dom) * ((dv - dim) / (dix - dim)));
}

FujiValue fuji_sign_num(FujiValue v) {
    if (!IS_NUMBER(v)) return NULL_VAL;
    double dv = AS_NUMBER(v);
    if (dv > 0) return NUMBER_VAL(1);
    if (dv < 0) return NUMBER_VAL(-1);
    return NUMBER_VAL(0);
}

FujiValue fuji_hypot_num(FujiValue a, FujiValue b) {
    if (!IS_NUMBER(a) || !IS_NUMBER(b)) return NULL_VAL;
    return NUMBER_VAL(hypot(AS_NUMBER(a), AS_NUMBER(b)));
}

FujiValue fuji_distance_num(FujiValue x1, FujiValue y1, FujiValue x2, FujiValue y2) {
    if (!IS_NUMBER(x1) || !IS_NUMBER(y1) || !IS_NUMBER(x2) || !IS_NUMBER(y2)) return NULL_VAL;
    double dx = AS_NUMBER(x2) - AS_NUMBER(x1);
    double dy = AS_NUMBER(y2) - AS_NUMBER(y1);
    return NUMBER_VAL(sqrt(dx * dx + dy * dy));
}

FujiValue fuji_distance_sq_num(FujiValue x1, FujiValue y1, FujiValue x2, FujiValue y2) {
    if (!IS_NUMBER(x1) || !IS_NUMBER(y1) || !IS_NUMBER(x2) || !IS_NUMBER(y2)) return NULL_VAL;
    double dx = AS_NUMBER(x2) - AS_NUMBER(x1);
    double dy = AS_NUMBER(y2) - AS_NUMBER(y1);
    return NUMBER_VAL(dx * dx + dy * dy);
}

FujiValue fuji_angle_between_num(FujiValue x1, FujiValue y1, FujiValue x2, FujiValue y2) {
    if (!IS_NUMBER(x1) || !IS_NUMBER(y1) || !IS_NUMBER(x2) || !IS_NUMBER(y2)) return NULL_VAL;
    return NUMBER_VAL(atan2(AS_NUMBER(y2) - AS_NUMBER(y1), AS_NUMBER(x2) - AS_NUMBER(x1)));
}

FujiValue fuji_normalize_num(FujiValue x, FujiValue y) {
    if (!IS_NUMBER(x) || !IS_NUMBER(y)) return NULL_VAL;
    double dx = AS_NUMBER(x);
    double dy = AS_NUMBER(y);
    double len = sqrt(dx * dx + dy * dy);
    double nx = 0, ny = 0;
    if (len > 0) {
        nx = dx / len;
        ny = dy / len;
    }
    FujiValue objVal = OBJ_VAL(fuji_alloc_object());
    fuji_object_set(objVal, OBJ_VAL(fuji_copy_string("x", 1)), NUMBER_VAL(nx));
    fuji_object_set(objVal, OBJ_VAL(fuji_copy_string("y", 1)), NUMBER_VAL(ny));
    return objVal;
}

FujiValue fuji_delta_time(void) {
    double now = fuji_get_time_seconds();
    double dt = now - heap.last_frame;
    heap.last_frame = now;
    if (dt > 0.1) dt = 0.1; 
    return NUMBER_VAL(dt);
}

FujiValue fuji_time(void) {
    return NUMBER_VAL(fuji_get_time_seconds() - heap.program_start);
}

FujiValue fuji_timestamp(void) {
    return NUMBER_VAL((double)time(NULL));
}

FujiValue fuji_random(int argCount, FujiValue* args) {
    double r = (double)rand() / (double)RAND_MAX;
    if (argCount == 0) return NUMBER_VAL(r);
    if (argCount == 1) return NUMBER_VAL(r * AS_NUMBER(args[0]));
    double min = AS_NUMBER(args[0]);
    double max = AS_NUMBER(args[1]);
    return NUMBER_VAL(min + r * (max - min));
}

FujiValue fuji_random_int(int argCount, FujiValue* args) {
    if (argCount == 1) return NUMBER_VAL(rand() % (int)AS_NUMBER(args[0]));
    int min = (int)AS_NUMBER(args[0]);
    int max = (int)AS_NUMBER(args[1]);
    if (max <= min) return NUMBER_VAL(min);
    return NUMBER_VAL(min + (rand() % (max - min)));
}

FujiValue fuji_random_choice(FujiValue array) {
    if (!IS_OBJ(array) || AS_OBJ(array)->type != OBJ_ARRAY) return NULL_VAL;
    FujiArray* arr = (FujiArray*)AS_OBJ(array);
    if (arr->count == 0) return NULL_VAL;
    return arr->elements[rand() % arr->count];
}

void fuji_random_seed(FujiValue seed) {
    if (IS_NUMBER(seed)) srand((unsigned int)AS_NUMBER(seed));
}

FujiValue fuji_set_timeout(FujiValue closure, FujiValue ms) {
    if (!IS_NUMBER(ms)) return NULL_VAL;
    FujiTask* task = (FujiTask*)malloc(sizeof(FujiTask));
    task->closure = closure;
    task->target_time = fuji_get_time_seconds() + (AS_NUMBER(ms) / 1000.0);
    task->interval = 0;
    task->next = heap.tasks;
    heap.tasks = task;
    return NULL_VAL;
}

FujiValue fuji_set_interval(FujiValue closure, FujiValue ms) {
    if (!IS_NUMBER(ms)) return NULL_VAL;
    double interval = AS_NUMBER(ms) / 1000.0;
    FujiTask* task = (FujiTask*)malloc(sizeof(FujiTask));
    task->closure = closure;
    task->target_time = fuji_get_time_seconds() + interval;
    task->interval = interval;
    task->next = heap.tasks;
    heap.tasks = task;
    return NULL_VAL;
}

void fuji_poll_tasks() {
    double now = fuji_get_time_seconds();
    FujiTask** prev = &heap.tasks;
    while (*prev) {
        FujiTask* task = *prev;
        if (now >= task->target_time) {
            // Call closure
            fuji_call(task->closure, 0, NULL);
            
            if (task->interval > 0) {
                task->target_time = now + task->interval;
                prev = &task->next;
            } else {
                *prev = task->next;
                free(task);
            }
        } else {
            prev = &task->next;
        }
    }
}

FujiValue fuji_math_sin(int argCount, FujiValue* args) { if (argCount < 1 || !IS_NUMBER(args[0])) return NULL_VAL; return NUMBER_VAL(sin(AS_NUMBER(args[0]))); }
FujiValue fuji_math_cos(int argCount, FujiValue* args) { if (argCount < 1 || !IS_NUMBER(args[0])) return NULL_VAL; return NUMBER_VAL(cos(AS_NUMBER(args[0]))); }
FujiValue fuji_math_tan(int argCount, FujiValue* args) { if (argCount < 1 || !IS_NUMBER(args[0])) return NULL_VAL; return NUMBER_VAL(tan(AS_NUMBER(args[0]))); }
FujiValue fuji_math_sqrt(int argCount, FujiValue* args) { if (argCount < 1 || !IS_NUMBER(args[0])) return NULL_VAL; return NUMBER_VAL(sqrt(AS_NUMBER(args[0]))); }
FujiValue fuji_math_abs(int argCount, FujiValue* args) { if (argCount < 1 || !IS_NUMBER(args[0])) return NULL_VAL; return NUMBER_VAL(fabs(AS_NUMBER(args[0]))); }
FujiValue fuji_math_floor(int argCount, FujiValue* args) { if (argCount < 1 || !IS_NUMBER(args[0])) return NULL_VAL; return NUMBER_VAL(floor(AS_NUMBER(args[0]))); }
FujiValue fuji_math_ceil(int argCount, FujiValue* args) { if (argCount < 1 || !IS_NUMBER(args[0])) return NULL_VAL; return NUMBER_VAL(ceil(AS_NUMBER(args[0]))); }
FujiValue fuji_math_clamp(int argCount, FujiValue* args) { if (argCount < 3) return NULL_VAL; return fuji_clamp_num(args[0], args[1], args[2]); }
FujiValue fuji_math_lerp(int argCount, FujiValue* args) { if (argCount < 3) return NULL_VAL; return fuji_lerp_num(args[0], args[1], args[2]); }
FujiValue fuji_math_sign(int argCount, FujiValue* args) { if (argCount < 1) return NULL_VAL; return fuji_sign_num(args[0]); }
FujiValue fuji_math_random(int argCount, FujiValue* args) { return fuji_random(argCount, args); }
FujiValue fuji_math_randomInt(int argCount, FujiValue* args) { return fuji_random_int(argCount, args); }
FujiValue fuji_math_randomChoice(int argCount, FujiValue* args) { if (argCount < 1) return NULL_VAL; return fuji_random_choice(args[0]); }
FujiValue fuji_math_randomSeed(int argCount, FujiValue* args) { if (argCount >= 1) fuji_random_seed(args[0]); return NULL_VAL; }
FujiValue fuji_math_deltaTime(int argCount, FujiValue* args) { (void)argCount; (void)args; return fuji_delta_time(); }
FujiValue fuji_math_time(int argCount, FujiValue* args) { (void)argCount; (void)args; return fuji_time(); }
FujiValue fuji_math_timestamp(int argCount, FujiValue* args) { (void)argCount; (void)args; return fuji_timestamp(); }
FujiValue fuji_math_sleep(int argCount, FujiValue* args) { if (argCount >= 1) fuji_sleep_ms(args[0]); return NULL_VAL; }

FujiValue fuji_module_init_math() {
    static FujiValue module_exports = 0;
    if (module_exports != 0) return module_exports;
    module_exports = OBJ_VAL(fuji_alloc_object());
    fuji_object_set(module_exports, OBJ_VAL(fuji_copy_string("sin", 3)), fuji_new_native(fuji_math_sin));
    fuji_object_set(module_exports, OBJ_VAL(fuji_copy_string("cos", 3)), fuji_new_native(fuji_math_cos));
    fuji_object_set(module_exports, OBJ_VAL(fuji_copy_string("tan", 3)), fuji_new_native(fuji_math_tan));
    fuji_object_set(module_exports, OBJ_VAL(fuji_copy_string("sqrt", 4)), fuji_new_native(fuji_math_sqrt));
    fuji_object_set(module_exports, OBJ_VAL(fuji_copy_string("abs", 3)), fuji_new_native(fuji_math_abs));
    fuji_object_set(module_exports, OBJ_VAL(fuji_copy_string("floor", 5)), fuji_new_native(fuji_math_floor));
    fuji_object_set(module_exports, OBJ_VAL(fuji_copy_string("ceil", 4)), fuji_new_native(fuji_math_ceil));
    fuji_object_set(module_exports, OBJ_VAL(fuji_copy_string("clamp", 5)), fuji_new_native(fuji_math_clamp));
    fuji_object_set(module_exports, OBJ_VAL(fuji_copy_string("lerp", 4)), fuji_new_native(fuji_math_lerp));
    fuji_object_set(module_exports, OBJ_VAL(fuji_copy_string("sign", 4)), fuji_new_native(fuji_math_sign));
    fuji_object_set(module_exports, OBJ_VAL(fuji_copy_string("random", 6)), fuji_new_native(fuji_math_random));
    fuji_object_set(module_exports, OBJ_VAL(fuji_copy_string("randomInt", 9)), fuji_new_native(fuji_math_randomInt));
    fuji_object_set(module_exports, OBJ_VAL(fuji_copy_string("randomChoice", 12)), fuji_new_native(fuji_math_randomChoice));
    fuji_object_set(module_exports, OBJ_VAL(fuji_copy_string("randomSeed", 10)), fuji_new_native(fuji_math_randomSeed));
    fuji_object_set(module_exports, OBJ_VAL(fuji_copy_string("deltaTime", 9)), fuji_new_native(fuji_math_deltaTime));
    fuji_object_set(module_exports, OBJ_VAL(fuji_copy_string("time", 4)), fuji_new_native(fuji_math_time));
    fuji_object_set(module_exports, OBJ_VAL(fuji_copy_string("timestamp", 9)), fuji_new_native(fuji_math_timestamp));
    fuji_object_set(module_exports, OBJ_VAL(fuji_copy_string("sleep", 5)), fuji_new_native(fuji_math_sleep));
    fuji_object_set(module_exports, OBJ_VAL(fuji_copy_string("pi", 2)), NUMBER_VAL(3.14159265358979323846));
    fuji_object_set(module_exports, OBJ_VAL(fuji_copy_string("e", 1)), NUMBER_VAL(2.71828182845904523536));
    return module_exports;
}

