package codegen

import (
	"github.com/llir/llvm/ir"
	"github.com/llir/llvm/ir/types"
)

// argvI64 declares Value name(int arg_count, Value* args) as i64(i32, i64*).
func argvI64(mod *ir.Module, cName string) *ir.Func {
	return mod.NewFunc(cName, types.I64,
		ir.NewParam("arg_count", types.I32),
		ir.NewParam("args", types.NewPointer(types.I64)))
}

// declareRuntimeFunctions declares all runtime functions used by the generated code.
// LLVM symbol names and calling conventions must match runtime/src/fuji_runtime.c.
func declareRuntimeFunctions(mod *ir.Module) map[string]*ir.Func {
	functions := make(map[string]*ir.Func)

	functions["FUJI_runtime_init"] = mod.NewFunc("fuji_runtime_init", types.Void)
	functions["FUJI_runtime_shutdown"] = mod.NewFunc("fuji_runtime_shutdown", types.Void)

	// print: fuji_print_val for one value; fuji_print_argv for multiple (space-separated, one newline).
	functions["FUJI_print"] = mod.NewFunc("fuji_print_val", types.I64,
		ir.NewParam("val", types.I64))
	functions["FUJI_print_argv"] = argvI64(mod, "fuji_print_argv")
	functions["FUJI_print_newline"] = mod.NewFunc("fuji_print_newline", types.Void)
	typeFn := mod.NewFunc("fuji_type", types.I64, ir.NewParam("val", types.I64))
	functions["FUJI_type"] = typeFn
	functions["FUJI_typeof"] = typeFn

	// Time (C uses argv for delta/timestamp/program; scalars for clock/wall/time/sleep).
	functions["FUJI_delta_time"] = argvI64(mod, "fuji_delta_time")
	functions["FUJI_clock"] = mod.NewFunc("fuji_clock", types.I64)
	functions["FUJI_timestamp"] = argvI64(mod, "fuji_timestamp")
	fujiTimeFn := mod.NewFunc("fuji_time", types.I64)
	functions["FUJI_time"] = fujiTimeFn
	functions["FUJI_program_time"] = argvI64(mod, "fuji_program_time")
	functions["FUJI_wall_time"] = mod.NewFunc("fuji_wall_time", types.I64)
	functions["FUJI_sleep"] = mod.NewFunc("fuji_sleep", types.Void,
		ir.NewParam("ms", types.I64))

	// Random
	functions["FUJI_random"] = argvI64(mod, "fuji_random")
	functions["FUJI_randomInt"] = argvI64(mod, "fuji_randomInt")
	functions["FUJI_randomChoice"] = argvI64(mod, "fuji_randomChoice")
	functions["FUJI_randomSeed"] = argvI64(mod, "fuji_randomSeed")

	// Math (C implementations are argv-style)
	functions["FUJI_lerp"] = argvI64(mod, "fuji_lerp")
	functions["FUJI_clamp"] = argvI64(mod, "fuji_clamp")
	functions["FUJI_distance"] = argvI64(mod, "fuji_distance")
	functions["FUJI_angleBetween"] = argvI64(mod, "fuji_angleBetween")
	functions["FUJI_map"] = argvI64(mod, "fuji_map")

	functions["FUJI_pi"] = argvI64(mod, "fuji_pi")
	functions["FUJI_e"] = argvI64(mod, "fuji_e")

	functions["FUJI_sin"] = argvI64(mod, "fuji_sin")
	functions["FUJI_cos"] = argvI64(mod, "fuji_cos")
	functions["FUJI_tan"] = argvI64(mod, "fuji_tan")
	functions["FUJI_asin"] = argvI64(mod, "fuji_asin")
	functions["FUJI_acos"] = argvI64(mod, "fuji_acos")
	functions["FUJI_atan"] = argvI64(mod, "fuji_atan")
	functions["FUJI_atan2"] = argvI64(mod, "fuji_atan2")

	functions["FUJI_pow"] = argvI64(mod, "fuji_pow")
	functions["FUJI_exp"] = argvI64(mod, "fuji_exp")
	functions["FUJI_log"] = argvI64(mod, "fuji_log")
	functions["FUJI_log10"] = argvI64(mod, "fuji_log10")

	functions["FUJI_floor"] = argvI64(mod, "fuji_floor")
	functions["FUJI_ceil"] = argvI64(mod, "fuji_ceil")
	functions["FUJI_round"] = argvI64(mod, "fuji_round")
	functions["FUJI_trunc"] = argvI64(mod, "fuji_trunc")

	functions["FUJI_sign"] = argvI64(mod, "fuji_sign")
	functions["FUJI_min"] = argvI64(mod, "fuji_min")
	functions["FUJI_max"] = argvI64(mod, "fuji_max")
	functions["FUJI_smoothstep"] = argvI64(mod, "fuji_smoothstep")
	functions["FUJI_distanceSq"] = argvI64(mod, "fuji_distanceSq")
	functions["FUJI_normalize"] = argvI64(mod, "fuji_normalize")

	// Type checks (argv in C)
	functions["FUJI_isNumber"] = argvI64(mod, "fuji_isNumber")
	functions["FUJI_isString"] = argvI64(mod, "fuji_isString")
	functions["FUJI_isBool"] = argvI64(mod, "fuji_isBool")
	functions["FUJI_isNull"] = argvI64(mod, "fuji_isNull")
	functions["FUJI_isArray"] = argvI64(mod, "fuji_isArray")
	functions["FUJI_isObject"] = argvI64(mod, "fuji_isObject")
	functions["FUJI_isFunction"] = argvI64(mod, "fuji_isFunction")

	functions["FUJI_bool"] = argvI64(mod, "fuji_bool")

	functions["FUJI_format"] = argvI64(mod, "fuji_format")

	functions["FUJI_array_map"] = argvI64(mod, "fuji_array_map")
	functions["FUJI_array_filter"] = argvI64(mod, "fuji_array_filter")
	functions["FUJI_array_forEach"] = argvI64(mod, "fuji_array_forEach")
	functions["FUJI_array_find"] = argvI64(mod, "fuji_array_find")
	functions["FUJI_array_findIndex"] = argvI64(mod, "fuji_array_findIndex")
	functions["FUJI_array_some"] = argvI64(mod, "fuji_array_some")
	functions["FUJI_array_every"] = argvI64(mod, "fuji_array_every")
	functions["FUJI_array_reduce"] = argvI64(mod, "fuji_array_reduce")
	functions["FUJI_array_sort"] = argvI64(mod, "fuji_array_sort")
	functions["FUJI_array_reverse"] = argvI64(mod, "fuji_array_reverse")
	functions["FUJI_array_indexOf"] = argvI64(mod, "fuji_array_indexOf")
	functions["FUJI_array_includes"] = argvI64(mod, "fuji_array_includes")
	functions["FUJI_array_slice"] = argvI64(mod, "fuji_array_slice")
	functions["FUJI_array_concat"] = argvI64(mod, "fuji_array_concat")

	functions["FUJI_string_split"] = argvI64(mod, "fuji_string_split")
	functions["FUJI_string_trim"] = argvI64(mod, "fuji_string_trim")
	functions["FUJI_string_upper"] = argvI64(mod, "fuji_string_upper")
	functions["FUJI_string_lower"] = argvI64(mod, "fuji_string_lower")
	functions["FUJI_string_startsWith"] = argvI64(mod, "fuji_string_startsWith")
	functions["FUJI_string_endsWith"] = argvI64(mod, "fuji_string_endsWith")
	functions["FUJI_string_indexOf"] = argvI64(mod, "fuji_string_indexOf")
	functions["FUJI_string_slice"] = argvI64(mod, "fuji_string_slice")
	functions["FUJI_string_replace"] = argvI64(mod, "fuji_string_replace")
	functions["FUJI_string_replaceAll"] = argvI64(mod, "fuji_string_replaceAll")

	functions["FUJI_readFile"] = argvI64(mod, "fuji_readFile")
	functions["FUJI_writeFile"] = argvI64(mod, "fuji_writeFile")
	functions["FUJI_appendFile"] = argvI64(mod, "fuji_appendFile")
	functions["FUJI_fileExists"] = argvI64(mod, "fuji_fileExists")
	functions["FUJI_deleteFile"] = argvI64(mod, "fuji_deleteFile")

	functions["FUJI_ok"] = argvI64(mod, "fuji_ok")
	functions["FUJI_err"] = argvI64(mod, "fuji_err")
	functions["FUJI_values_equal"] = mod.NewFunc("fuji_values_equal", types.I64,
		ir.NewParam("a", types.I64), ir.NewParam("b", types.I64))
	functions["FUJI_panic"] = mod.NewFunc("fuji_panic", types.Void,
		ir.NewParam("arg_count", types.I32),
		ir.NewParam("args", types.NewPointer(types.I64)))
	functions["FUJI_assert"] = argvI64(mod, "fuji_assert")
	functions["FUJI_assert_llvm"] = mod.NewFunc("fuji_assert_llvm", types.Void,
		ir.NewParam("cond", types.I64),
		ir.NewParam("msg", types.I64))
	functions["FUJI_trace"] = argvI64(mod, "fuji_trace")

	functions["FUJI_parseJSON"] = argvI64(mod, "fuji_parseJSON")
	functions["FUJI_toJSON"] = argvI64(mod, "fuji_toJSON")

	functions["FUJI_allocate_object"] = mod.NewFunc("fuji_allocate_object", types.I64,
		ir.NewParam("property_count", types.I32))

	functions["FUJI_unbox_number"] = mod.NewFunc("fuji_unbox_number", types.Double,
		ir.NewParam("v", types.I64))
	functions["FUJI_box_number"] = mod.NewFunc("fuji_box_number", types.I64,
		ir.NewParam("d", types.Double))

	getFn := mod.NewFunc("fuji_get", types.I64,
		ir.NewParam("obj", types.I64),
		ir.NewParam("key", types.I64))
	functions["FUJI_get"] = getFn
	functions["FUJI_object_get"] = getFn
	functions["FUJI_array_get"] = getFn

	functions["FUJI_set"] = mod.NewFunc("fuji_set", types.I64,
		ir.NewParam("obj", types.I64),
		ir.NewParam("key", types.I64),
		ir.NewParam("val", types.I64))

	functions["FUJI_object_set"] = mod.NewFunc("fuji_object_set", types.I64,
		ir.NewParam("obj", types.I64),
		ir.NewParam("key", types.I64),
		ir.NewParam("value", types.I64))

	functions["FUJI_allocate_string"] = mod.NewFunc("fuji_allocate_string", types.I64,
		ir.NewParam("length", types.I32),
		ir.NewParam("chars", types.NewPointer(types.I8)))

	functions["FUJI_allocate_array"] = mod.NewFunc("fuji_allocate_array", types.I64,
		ir.NewParam("length", types.I32))
	functions["FUJI_array_set"] = mod.NewFunc("fuji_array_set", types.Void,
		ir.NewParam("arr", types.I64),
		ir.NewParam("index", types.I64),
		ir.NewParam("value", types.I64))
	functions["FUJI_array_push"] = mod.NewFunc("fuji_array_push", types.Void,
		ir.NewParam("arr", types.I64),
		ir.NewParam("value", types.I64))
	functions["FUJI_array_pop"] = mod.NewFunc("fuji_array_pop", types.I64,
		ir.NewParam("arr", types.I64))
	lenFn := mod.NewFunc("fuji_len", types.I64,
		ir.NewParam("value", types.I64))
	functions["FUJI_len"] = lenFn
	arrayLenFn := mod.NewFunc("fuji_array_length", types.I64,
		ir.NewParam("value", types.I64))
	functions["FUJI_array_length"] = arrayLenFn

	functions["FUJI_abs"] = mod.NewFunc("fuji_abs", types.I64,
		ir.NewParam("value", types.I64))
	functions["FUJI_sqrt"] = mod.NewFunc("fuji_sqrt", types.I64,
		ir.NewParam("value", types.I64))
	functions["FUJI_number"] = mod.NewFunc("fuji_number", types.I64,
		ir.NewParam("value", types.I64))
	functions["FUJI_string"] = mod.NewFunc("fuji_string", types.I64,
		ir.NewParam("value", types.I64))
	functions["FUJI_string_concat"] = mod.NewFunc("fuji_string_concat", types.I64,
		ir.NewParam("a", types.I64),
		ir.NewParam("b", types.I64))

	functions["FUJI_range"] = argvI64(mod, "fuji_range")

	cellPtr := types.NewPointer(types.I64)
	functions["FUJI_alloc_cell"] = mod.NewFunc("fuji_alloc_cell", cellPtr)
	functions["FUJI_cell_read"] = mod.NewFunc("fuji_cell_read", types.I64,
		ir.NewParam("cell", cellPtr))
	functions["FUJI_cell_write"] = mod.NewFunc("fuji_cell_write", types.Void,
		ir.NewParam("cell", cellPtr),
		ir.NewParam("val", types.I64))
	functions["FUJI_gc_collect"] = mod.NewFunc("fuji_gc_collect", types.Void)
	functions["FUJI_gc_frame_step"] = mod.NewFunc("fuji_gc_frame_step", types.Void,
		ir.NewParam("budget_ms", types.Double))

	ptrPtr := types.NewPointer(types.NewPointer(types.I64))
	functions["FUJI_push_frame"] = mod.NewFunc("fuji_push_frame", types.Void,
		ir.NewParam("slots", ptrPtr),
		ir.NewParam("count", types.I32))
	functions["FUJI_pop_frame"] = mod.NewFunc("fuji_pop_frame", types.Void)

	functions["FUJI_push_call"] = mod.NewFunc("fuji_push_call", types.Void,
		ir.NewParam("fn_name", types.NewPointer(types.I8)),
		ir.NewParam("file_name", types.NewPointer(types.I8)),
		ir.NewParam("line", types.I32))
	functions["FUJI_pop_call"] = mod.NewFunc("fuji_pop_call", types.Void)

	// C malloc for packed closure environments (linker resolves from CRT).
	functions["malloc"] = mod.NewFunc("malloc", types.NewPointer(types.I8), ir.NewParam("size", types.I64))

	return functions
}
