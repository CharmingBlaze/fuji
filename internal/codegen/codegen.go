package codegen

import (
	"fmt"
	"strings"

	"fuji/internal/parser"
	"fuji/internal/sema"

	"github.com/llir/llvm/ir"
	"github.com/llir/llvm/ir/constant"
	"github.com/llir/llvm/ir/enum"
	"github.com/llir/llvm/ir/types"
	"github.com/llir/llvm/ir/value"
)

// PrepareNativeBundle prepares a bundle for native emission by performing
// capture analysis and local variable mapping.
func PrepareNativeBundle(bundle *parser.ProgramBundle) (*sema.NativeEmitContext, error) {
	return sema.PrepareNativeBundle(bundle)
}

// EmitLLVMIR generates LLVM IR for the bundle carried by ctx.
func EmitLLVMIR(ctx *sema.NativeEmitContext) (*ir.Module, error) {
	gen := NewGenerator(ctx)
	return gen.Generate(ctx.Bundle)
}

type Generator struct {
	mod                     *ir.Module
	block                   *ir.Block
	ctx                     *sema.NativeEmitContext
	funcs                   map[string]*ir.Func    // Function name → LLVM function
	locals                  map[string]value.Value // Variable name → stack slot
	globals                 map[string]*ir.Global  // Global variables
	currentFn               *ir.Func
	runtimeInit             *ir.Func
	runtimeShutdown         *ir.Func
	runtimeDeltaTime        *ir.Func
	runtimeProgramTime      *ir.Func
	runtimeTimestamp        *ir.Func
	runtimeTime             *ir.Func
	runtimeClock            *ir.Func
	runtimeSleep            *ir.Func
	runtimePrint            *ir.Func
	runtimePrintArgv        *ir.Func
	runtimePrintNewline     *ir.Func
	runtimeRandom           *ir.Func
	runtimeRandomInt        *ir.Func
	runtimeRandomChoice     *ir.Func
	runtimeRandomSeed       *ir.Func
	runtimeLerp             *ir.Func
	runtimeClamp            *ir.Func
	runtimeDistance         *ir.Func
	runtimeAngleBetween     *ir.Func
	runtimeMap              *ir.Func
	runtimePI               *ir.Func
	runtimeE                *ir.Func
	runtimeSin              *ir.Func
	runtimeCos              *ir.Func
	runtimeTan              *ir.Func
	runtimeAsin             *ir.Func
	runtimeAcos             *ir.Func
	runtimeAtan             *ir.Func
	runtimeAtan2            *ir.Func
	runtimePow              *ir.Func
	runtimeExp              *ir.Func
	runtimeLog              *ir.Func
	runtimeLog10            *ir.Func
	runtimeFloor            *ir.Func
	runtimeCeil             *ir.Func
	runtimeRound            *ir.Func
	runtimeTrunc            *ir.Func
	runtimeSign             *ir.Func
	runtimeMin              *ir.Func
	runtimeMax              *ir.Func
	runtimeSmoothstep       *ir.Func
	runtimeDistanceSq       *ir.Func
	runtimeNormalize        *ir.Func
	runtimeIsNumber         *ir.Func
	runtimeIsString         *ir.Func
	runtimeIsBool           *ir.Func
	runtimeIsNull           *ir.Func
	runtimeIsArray          *ir.Func
	runtimeIsObject         *ir.Func
	runtimeIsFunction       *ir.Func
	runtimeBool             *ir.Func
	runtimeFormat           *ir.Func
	runtimeArrayMap         *ir.Func
	runtimeArrayFilter      *ir.Func
	runtimeArrayForEach     *ir.Func
	runtimeArrayFind        *ir.Func
	runtimeArrayFindIndex   *ir.Func
	runtimeArraySome        *ir.Func
	runtimeArrayEvery       *ir.Func
	runtimeArrayReduce      *ir.Func
	runtimeArraySort        *ir.Func
	runtimeArrayReverse     *ir.Func
	runtimeArrayIndexOf     *ir.Func
	runtimeArrayIncludes    *ir.Func
	runtimeArraySlice       *ir.Func
	runtimeArrayConcat      *ir.Func
	runtimeStringSplit      *ir.Func
	runtimeStringTrim       *ir.Func
	runtimeStringUpper      *ir.Func
	runtimeStringLower      *ir.Func
	runtimeStringStartsWith *ir.Func
	runtimeStringEndsWith   *ir.Func
	runtimeStringIndexOf    *ir.Func
	runtimeStringSlice      *ir.Func
	runtimeStringReplace    *ir.Func
	runtimeStringReplaceAll *ir.Func
	runtimeReadFile         *ir.Func

	runtimeWriteFile    *ir.Func
	runtimeAppendFile   *ir.Func
	runtimeFileExists   *ir.Func
	runtimeDeleteFile   *ir.Func
	runtimeAssert       *ir.Func
	runtimeTrace        *ir.Func
	runtimeParseJSON    *ir.Func
	runtimeToJSON       *ir.Func
	runtimeAllocObj     *ir.Func
	runtimeObjGet       *ir.Func
	runtimeObjSet       *ir.Func
	runtimeUnboxNumber  *ir.Func
	runtimeBoxNumber    *ir.Func
	runtimeSet          *ir.Func
	runtimeAllocStr     *ir.Func
	runtimeAllocArray   *ir.Func
	runtimeArrayGet     *ir.Func
	runtimeArraySet     *ir.Func
	runtimeArrayPush    *ir.Func
	runtimeArrayPop     *ir.Func
	runtimeArrayLen     *ir.Func
	runtimeType         *ir.Func
	runtimeLen          *ir.Func
	runtimeAbs          *ir.Func
	runtimeSqrt         *ir.Func
	runtimeNumber       *ir.Func
	runtimeString       *ir.Func
	runtimeStringConcat *ir.Func
	runtimeRange        *ir.Func
	runtimeAllocCell    *ir.Func
	runtimeCellRead     *ir.Func
	runtimeCellWrite    *ir.Func
	runtimeGcCollect    *ir.Func
	runtimeGcFrameStep  *ir.Func
	runtimePushFrame    *ir.Func
	runtimePopFrame     *ir.Func
	runtimeOk           *ir.Func
	runtimeErr          *ir.Func
	runtimePanic        *ir.Func
	runtimePushCall     *ir.Func
	runtimePopCall      *ir.Func
	sourcePath          string
	localIsCell         map[string]bool
	loopStack           []loopContext
	tempN               int

	shadowLayout     *sema.ShadowLayout
	shadowFramePtr   value.Value
	shadowFrameArrTy *types.ArrayType
	shadowPushed     bool
	shadowTempNext   int
}

// loopContext tracks information about the current loop for break/continue.
type loopContext struct {
	condBlock  *ir.Block
	incBlock   *ir.Block
	afterBlock *ir.Block
}

// NewGenerator creates a new LLVM IR generator.
func NewGenerator(ctx *sema.NativeEmitContext) *Generator {
	mod := ir.NewModule()

	// Declare runtime functions
	runtimeFuncs := declareRuntimeFunctions(mod)

	gen := &Generator{
		mod:                     mod,
		ctx:                     ctx,
		funcs:                   make(map[string]*ir.Func),
		locals:                  make(map[string]value.Value),
		globals:                 make(map[string]*ir.Global),
		runtimeInit:             runtimeFuncs["FUJI_runtime_init"],
		runtimeShutdown:         runtimeFuncs["FUJI_runtime_shutdown"],
		runtimeDeltaTime:        runtimeFuncs["FUJI_delta_time"],
		runtimeProgramTime:      runtimeFuncs["FUJI_program_time"],
		runtimeTimestamp:        runtimeFuncs["FUJI_timestamp"],
		runtimeTime:             runtimeFuncs["FUJI_time"],
		runtimeClock:            runtimeFuncs["FUJI_clock"],
		runtimeSleep:            runtimeFuncs["FUJI_sleep"],
		runtimePrint:            runtimeFuncs["FUJI_print"],
		runtimePrintArgv:        runtimeFuncs["FUJI_print_argv"],
		runtimePrintNewline:     runtimeFuncs["FUJI_print_newline"],
		runtimeRandom:           runtimeFuncs["FUJI_random"],
		runtimeRandomInt:        runtimeFuncs["FUJI_randomInt"],
		runtimeRandomChoice:     runtimeFuncs["FUJI_randomChoice"],
		runtimeRandomSeed:       runtimeFuncs["FUJI_randomSeed"],
		runtimeLerp:             runtimeFuncs["FUJI_lerp"],
		runtimeClamp:            runtimeFuncs["FUJI_clamp"],
		runtimeDistance:         runtimeFuncs["FUJI_distance"],
		runtimeAngleBetween:     runtimeFuncs["FUJI_angleBetween"],
		runtimeMap:              runtimeFuncs["FUJI_map"],
		runtimePI:               runtimeFuncs["FUJI_pi"],
		runtimeE:                runtimeFuncs["FUJI_e"],
		runtimeSin:              runtimeFuncs["FUJI_sin"],
		runtimeCos:              runtimeFuncs["FUJI_cos"],
		runtimeTan:              runtimeFuncs["FUJI_tan"],
		runtimeAsin:             runtimeFuncs["FUJI_asin"],
		runtimeAcos:             runtimeFuncs["FUJI_acos"],
		runtimeAtan:             runtimeFuncs["FUJI_atan"],
		runtimeAtan2:            runtimeFuncs["FUJI_atan2"],
		runtimePow:              runtimeFuncs["FUJI_pow"],
		runtimeExp:              runtimeFuncs["FUJI_exp"],
		runtimeLog:              runtimeFuncs["FUJI_log"],
		runtimeLog10:            runtimeFuncs["FUJI_log10"],
		runtimeFloor:            runtimeFuncs["FUJI_floor"],
		runtimeCeil:             runtimeFuncs["FUJI_ceil"],
		runtimeRound:            runtimeFuncs["FUJI_round"],
		runtimeTrunc:            runtimeFuncs["FUJI_trunc"],
		runtimeSign:             runtimeFuncs["FUJI_sign"],
		runtimeMin:              runtimeFuncs["FUJI_min"],
		runtimeMax:              runtimeFuncs["FUJI_max"],
		runtimeSmoothstep:       runtimeFuncs["FUJI_smoothstep"],
		runtimeDistanceSq:       runtimeFuncs["FUJI_distanceSq"],
		runtimeNormalize:        runtimeFuncs["FUJI_normalize"],
		runtimeIsNumber:         runtimeFuncs["FUJI_isNumber"],
		runtimeIsString:         runtimeFuncs["FUJI_isString"],
		runtimeIsBool:           runtimeFuncs["FUJI_isBool"],
		runtimeIsNull:           runtimeFuncs["FUJI_isNull"],
		runtimeIsArray:          runtimeFuncs["FUJI_isArray"],
		runtimeIsObject:         runtimeFuncs["FUJI_isObject"],
		runtimeIsFunction:       runtimeFuncs["FUJI_isFunction"],
		runtimeBool:             runtimeFuncs["FUJI_bool"],
		runtimeFormat:           runtimeFuncs["FUJI_format"],
		runtimeArrayMap:         runtimeFuncs["FUJI_array_map"],
		runtimeArrayFilter:      runtimeFuncs["FUJI_array_filter"],
		runtimeArrayForEach:     runtimeFuncs["FUJI_array_forEach"],
		runtimeArrayFind:        runtimeFuncs["FUJI_array_find"],
		runtimeArrayFindIndex:   runtimeFuncs["FUJI_array_findIndex"],
		runtimeArraySome:        runtimeFuncs["FUJI_array_some"],
		runtimeArrayEvery:       runtimeFuncs["FUJI_array_every"],
		runtimeArrayReduce:      runtimeFuncs["FUJI_array_reduce"],
		runtimeArraySort:        runtimeFuncs["FUJI_array_sort"],
		runtimeArrayReverse:     runtimeFuncs["FUJI_array_reverse"],
		runtimeArrayIndexOf:     runtimeFuncs["FUJI_array_indexOf"],
		runtimeArrayIncludes:    runtimeFuncs["FUJI_array_includes"],
		runtimeArraySlice:       runtimeFuncs["FUJI_array_slice"],
		runtimeArrayConcat:      runtimeFuncs["FUJI_array_concat"],
		runtimeStringSplit:      runtimeFuncs["FUJI_string_split"],
		runtimeStringTrim:       runtimeFuncs["FUJI_string_trim"],
		runtimeStringUpper:      runtimeFuncs["FUJI_string_upper"],
		runtimeStringLower:      runtimeFuncs["FUJI_string_lower"],
		runtimeStringStartsWith: runtimeFuncs["FUJI_string_startsWith"],
		runtimeStringEndsWith:   runtimeFuncs["FUJI_string_endsWith"],
		runtimeStringIndexOf:    runtimeFuncs["FUJI_string_indexOf"],
		runtimeStringSlice:      runtimeFuncs["FUJI_string_slice"],
		runtimeStringReplace:    runtimeFuncs["FUJI_string_replace"],
		runtimeStringReplaceAll: runtimeFuncs["FUJI_string_replaceAll"],
		runtimeReadFile:         runtimeFuncs["FUJI_readFile"],
		runtimeWriteFile:        runtimeFuncs["FUJI_writeFile"],
		runtimeAppendFile:       runtimeFuncs["FUJI_appendFile"],
		runtimeFileExists:       runtimeFuncs["FUJI_fileExists"],
		runtimeDeleteFile:       runtimeFuncs["FUJI_deleteFile"],
		runtimeAssert:           runtimeFuncs["FUJI_assert"],
		runtimeTrace:            runtimeFuncs["FUJI_trace"],
		runtimeParseJSON:        runtimeFuncs["FUJI_parseJSON"],
		runtimeToJSON:           runtimeFuncs["FUJI_toJSON"],
		runtimeAllocObj:         runtimeFuncs["FUJI_allocate_object"],
		runtimeObjGet:           runtimeFuncs["FUJI_object_get"],
		runtimeObjSet:           runtimeFuncs["FUJI_object_set"],
		runtimeUnboxNumber:      runtimeFuncs["FUJI_unbox_number"],
		runtimeBoxNumber:        runtimeFuncs["FUJI_box_number"],
		runtimeSet:              runtimeFuncs["FUJI_set"],
		runtimeAllocStr:         runtimeFuncs["FUJI_allocate_string"],
		runtimeAllocArray:       runtimeFuncs["FUJI_allocate_array"],
		runtimeArrayGet:         runtimeFuncs["FUJI_array_get"],
		runtimeArraySet:         runtimeFuncs["FUJI_array_set"],
		runtimeArrayPush:        runtimeFuncs["FUJI_array_push"],
		runtimeArrayPop:         runtimeFuncs["FUJI_array_pop"],
		runtimeArrayLen:         runtimeFuncs["FUJI_array_length"], // distinct symbol from FUJI_len
		runtimeType:             runtimeFuncs["FUJI_type"],
		runtimeLen:              runtimeFuncs["FUJI_len"],
		runtimeAbs:              runtimeFuncs["FUJI_abs"],
		runtimeSqrt:             runtimeFuncs["FUJI_sqrt"],
		runtimeNumber:           runtimeFuncs["FUJI_number"],
		runtimeString:           runtimeFuncs["FUJI_string"],
		runtimeStringConcat:     runtimeFuncs["FUJI_string_concat"],
		runtimeRange:            runtimeFuncs["FUJI_range"],
		runtimeAllocCell:        runtimeFuncs["FUJI_alloc_cell"],
		runtimeCellRead:         runtimeFuncs["FUJI_cell_read"],
		runtimeCellWrite:        runtimeFuncs["FUJI_cell_write"],
		runtimeGcCollect:        runtimeFuncs["FUJI_gc_collect"],
		runtimeGcFrameStep:      runtimeFuncs["FUJI_gc_frame_step"],
		runtimePushFrame:        runtimeFuncs["FUJI_push_frame"],
		runtimePopFrame:         runtimeFuncs["FUJI_pop_frame"],
		runtimeOk:               runtimeFuncs["FUJI_ok"],
		runtimeErr:              runtimeFuncs["FUJI_err"],
		runtimePanic:            runtimeFuncs["FUJI_panic"],
		runtimePushCall:         runtimeFuncs["FUJI_push_call"],
		runtimePopCall:          runtimeFuncs["FUJI_pop_call"],
		tempN:                   0,
	}

	return gen
}

// Generate emits LLVM IR for a program bundle.
func (g *Generator) Generate(bundle *parser.ProgramBundle) (*ir.Module, error) {
	entry := bundle.Entry
	if entry == nil {
		return nil, fmt.Errorf("no entry program")
	}

	g.sourcePath = "<entry>"
	if ep, err := parser.BundleEntryPath(bundle); err == nil {
		g.sourcePath = ep
	}
	g.registerBuiltinFuncs()

	// Create user_main function (implicit entry point for all top-level code)
	userMain := g.mod.NewFunc("user_main", types.I64, ir.NewParam("this", types.I64))
	g.funcs["user_main"] = userMain
	g.currentFn = userMain

	// Start emitting into user_main's entry block
	g.block = userMain.NewBlock("entry")
	g.locals = make(map[string]value.Value)
	g.localIsCell = make(map[string]bool)

	// Allocate 'this' slot
	thisSlot := g.block.NewAlloca(types.I64)
	g.locals["this"] = thisSlot
	g.block.NewStore(constant.NewInt(types.I64, 0), thisSlot)

	prevShadowLayout := g.shadowLayout
	prevShadowFramePtr := g.shadowFramePtr
	prevShadowFrameArrTy := g.shadowFrameArrTy
	prevShadowPushed := g.shadowPushed
	prevShadowTempNext := g.shadowTempNext

	g.shadowLayout = g.ctx.ShadowEntry
	var emptyParams []string
	g.beginShadowFrame(g.ctx.ShadowEntry, thisSlot, emptyParams)
	g.emitCallTracePush("user_main", g.sourcePath, 0)

	// Emit all top-level declarations and statements
	for _, decl := range entry.Declarations {
		if err := g.emitDecl(decl); err != nil {
			return nil, err
		}
	}

	// User-defined `function main()` is emitted as LLVM symbol `fuji_user_main`; invoke after top-level runs.
	if mainFn := g.funcs["main"]; mainFn != nil {
		g.block.NewCall(mainFn, constant.NewInt(types.I64, 0))
	}

	// Terminate user_main
	if g.block.Term == nil {
		g.emitCallTracePop()
		g.emitShadowPop()
		g.block.NewRet(constant.NewInt(types.I64, 0))
	}
	g.shadowLayout = prevShadowLayout
	g.shadowFramePtr = prevShadowFramePtr
	g.shadowFrameArrTy = prevShadowFrameArrTy
	g.shadowPushed = prevShadowPushed
	g.shadowTempNext = prevShadowTempNext
	g.currentFn = nil

	// Create main function that calls the entry point
	mainFn := g.mod.NewFunc("main", types.I32)
	entryBlock := mainFn.NewBlock("entry")
	g.block = entryBlock
	g.currentFn = mainFn

	// Call runtime init
	g.block.NewCall(g.runtimeInit)

	// Call the user's main function if it exists
	if um := g.funcs["user_main"]; um != nil {
		g.block.NewCall(um, constant.NewInt(types.I64, 0))
	}

	g.block.NewRet(constant.NewInt(types.I32, 0))
	g.currentFn = nil

	return g.mod, nil
}

func (g *Generator) emitDecl(decl parser.Decl) error {
	switch d := decl.(type) {
	case *parser.FuncDecl:
		return g.emitFuncDecl(d)
	case *parser.LetDecl:
		return g.emitLetDecl(d)
	case *parser.IncludeDecl:
		return g.emitIncludeDecl(d)
	case parser.Stmt:
		return g.emitStmt(d)
	}
	return nil
}

func (g *Generator) emitIncludeDecl(_ *parser.IncludeDecl) error {
	// Module includes are handled at the parser level by loading and parsing the included file
	// The codegen just needs to skip the include declaration since the included
	// declarations are already merged into the AST before code generation
	return nil
}

func (g *Generator) emitLetDecl(d *parser.LetDecl) error {
	name := d.Name.Lexeme

	if d.Native != nil {
		return g.emitNativeExternLet(d)
	}

	// If we're not in a function (g.block is nil), create a global variable
	if g.block == nil {
		// Create a global variable with pointer type
		global := g.mod.NewGlobalDef(name, constant.NewNull(types.NewPointer(types.I64)))
		g.globals[name] = global
		g.locals[name] = global // Store in locals for easy access

		// Initialize if there's an initial value
		if d.Init != nil {
			// For global variables, we need to initialize them in a constructor
			// For now, just store the initial value in the global
			initVal, err := g.emitExpr(d.Init)
			if err != nil {
				return err
			}
			// Store the initial value into the global
			// This needs to be done in the main function, not here
			// For now, we'll skip this and handle it in main
			_ = initVal
		}
		return nil
	}

	var storageSlot value.Value
	var useStack bool
	switch {
	case g.ctx.StackDecls[d]:
		useStack = true
		storageSlot = g.block.NewAlloca(types.I64)
		g.localIsCell[name] = false
	case g.ctx.EscapingDecls[d]:
		useStack = false
		storageSlot = g.block.NewCall(g.runtimeAllocCell)
		g.localIsCell[name] = true
	default:
		/* Conservative: heap cell if escape status unknown (must not keep stack slot past return). */
		useStack = false
		storageSlot = g.block.NewCall(g.runtimeAllocCell)
		g.localIsCell[name] = true
	}
	g.locals[name] = storageSlot
	g.shadowStoreLet(d, storageSlot)

	if d.Init != nil {
		initVal, err := g.emitExpr(d.Init)
		if err != nil {
			return err
		}
		boxed := g.emitAsFujiI64(initVal)
		if useStack {
			g.block.NewStore(boxed, storageSlot)
		} else {
			g.block.NewCall(g.runtimeCellWrite, storageSlot, boxed)
		}
	}

	return nil
}

// emitNativeExternLet binds a Fuji name to an LLVM declaration for a C symbol
// implementing FujiValue (*)(int argCount, FujiValue* args) (i64, i32, i64* in IR).
func (g *Generator) emitNativeExternLet(d *parser.LetDecl) error {
	sym := strings.TrimSpace(d.Native.Symbol)
	if sym == "" {
		return fmt.Errorf("native extern missing C symbol for %q", d.Name.Lexeme)
	}
	fn := g.ensureNativeExternFunc(sym)
	name := d.Name.Lexeme
	g.funcs[name] = fn
	if kn := strings.TrimSpace(d.Native.BindingName); kn != "" && kn != name {
		g.funcs[kn] = fn
	}
	// Native bindings are declarations only; ignore initializer (typically `0`).
	return nil
}

func (g *Generator) ensureNativeExternFunc(symbol string) *ir.Func {
	for _, f := range g.mod.Funcs {
		if f.Name() == symbol {
			return f
		}
	}
	return g.mod.NewFunc(symbol, types.I64,
		ir.NewParam("arg_count", types.I32),
		ir.NewParam("args", types.NewPointer(types.I64)))
}

func (g *Generator) emitExpr(expr parser.Expr) (value.Value, error) {
	switch e := expr.(type) {
	case *parser.LiteralExpr:
		return g.emitLiteral(e)
	case *parser.IdentifierExpr:
		return g.emitIdentifier(e)
	case *parser.PrefixExpr:
		return g.emitPrefix(e)
	case *parser.InfixExpr:
		return g.emitInfix(e)
	case *parser.CallExpr:
		return g.emitCall(e)
	case *parser.ObjectExpr:
		return g.emitObject(e)
	case *parser.ArrayExpr:
		return g.emitArray(e)
	case *parser.IndexExpr:
		return g.emitIndex(e)
	case *parser.SliceExpr:
		return g.emitSlice(e)
	case *parser.AssignExpr:
		return g.emitAssign(e)
	case *parser.FuncExpr:
		return g.emitFuncExpr(e)
	case *parser.ThisExpr:
		return g.emitThis(e)
	case *parser.TemplateExpr:
		return g.emitTemplate(e)
	case *parser.RangeExpr:
		return g.emitRange(e)
	case *parser.TupleExpr:
		return g.emitTuple(e)
	case *parser.ImportExpr:
		return g.emitImport(e)
	case *parser.SwitchExpr:
		return g.emitSwitchExpr(e)
	case *parser.UpdateExpr:
		return g.emitUpdate(e)
	default:
		return nil, fmt.Errorf("unsupported expression type: %T", expr)
	}
}

func (g *Generator) emitPrefix(e *parser.PrefixExpr) (value.Value, error) {
	// Emit the right-hand side expression
	right, err := g.emitExpr(e.Right)
	if err != nil {
		return nil, err
	}

	// Handle different prefix operators
	switch e.Operator {
	case "!":
		truthy := g.emitTruthy(right)
		notTruthy := g.block.NewXor(truthy, constant.NewBool(true))
		return g.block.NewZExt(notTruthy, types.I64), nil
	case "-":
		// Unary negation on numeric values (NaN-boxed)
		ri := g.emitAsFujiI64(right)
		rd := g.block.NewCall(g.runtimeUnboxNumber, ri)
		neg := g.block.NewFNeg(rd)
		return g.block.NewCall(g.runtimeBoxNumber, neg), nil
	default:
		return nil, fmt.Errorf("unsupported prefix operator: %s", e.Operator)
	}
}

// emitBoxBoolNaN maps i1 (predicate result) to Fuji true/false NaN-boxed values.
func (g *Generator) emitBoxBoolNaN(cmp value.Value) value.Value {
	falseVal := constant.NewInt(types.I64, 0x7ffc000000000002)
	trueVal := constant.NewInt(types.I64, 0x7ffc000000000003)
	return g.block.NewSelect(cmp, trueVal, falseVal)
}

func (g *Generator) emitTruthy(v value.Value) value.Value {
	if v.Type().Equal(types.I1) {
		return v
	}
	if !v.Type().Equal(types.I64) {
		v = g.emitAsFujiI64(v)
	}
	isZero := g.block.NewICmp(enum.IPredEQ, v, constant.NewInt(types.I64, 0))
	isNil := g.block.NewICmp(enum.IPredEQ, v, constant.NewInt(types.I64, 0x7ffc000000000001))
	isFalse := g.block.NewICmp(enum.IPredEQ, v, constant.NewInt(types.I64, 0x7ffc000000000002))
	zeroOrNil := g.block.NewOr(isZero, isNil)
	falsey := g.block.NewOr(zeroOrNil, isFalse)
	return g.block.NewXor(falsey, constant.NewBool(true))
}

// isNativeArgvCallee reports whether fn uses the embedded runtime convention
// FujiValue fn(int argCount, FujiValue* args) (FujiValue as i64).
func isNativeArgvCallee(fn *ir.Func) bool {
	if fn == nil || len(fn.Params) != 2 {
		return false
	}
	if fn.Params[0].Typ != types.I32 {
		return false
	}
	return fn.Params[1].Typ.Equal(types.NewPointer(types.I64))
}

func (g *Generator) emitAsFujiI64(v value.Value) value.Value {
	if fn, ok := v.(*ir.Func); ok {
		return g.block.NewPtrToInt(fn, types.I64)
	}
	switch v.Type() {
	case types.I64:
		return v
	case types.I32:
		return g.block.NewSExt(v, types.I64)
	case types.I1:
		return g.block.NewZExt(v, types.I64)
	case types.Double:
		return g.block.NewBitCast(v, types.I64)
	default:
		return g.block.NewBitCast(v, types.I64)
	}
}

// emitArgvRuntime calls a runtime function declared as i64 name(i32 arg_count, i64* argv).
func (g *Generator) emitArgvRuntime(fn *ir.Func, args []value.Value) value.Value {
	n := len(args)
	zero := constant.NewInt(types.I32, 0)
	if n == 0 {
		return g.block.NewCall(fn, constant.NewInt(types.I32, 0), constant.NewNull(types.NewPointer(types.I64)))
	}
	arrTy := types.NewArray(uint64(n), types.I64)
	slot := g.block.NewAlloca(arrTy)
	for i, arg := range args {
		argI64 := g.emitAsFujiI64(arg)
		elemPtr := g.block.NewGetElementPtr(arrTy, slot, zero, constant.NewInt(types.I32, int64(i)))
		g.block.NewStore(argI64, elemPtr)
	}
	argvPtr := g.block.NewGetElementPtr(arrTy, slot, zero, zero)
	return g.block.NewCall(fn, constant.NewInt(types.I32, int64(n)), argvPtr)
}

// indirectFujiFuncPtrType is i64 (i64 this, i64 × argCount)* — matches every user/closure function in this backend.
func indirectFujiFuncPtrType(argCount int) *types.PointerType {
	paramTys := make([]types.Type, 0, 1+argCount)
	for i := 0; i < 1+argCount; i++ {
		paramTys = append(paramTys, types.I64)
	}
	return types.NewPointer(types.NewFunc(types.I64, paramTys...))
}

func (g *Generator) emitCall(e *parser.CallExpr) (value.Value, error) {
	// Check if this is a method call (e.g., obj.method())
	// Set this value if the function is called on an object
	var objForThis value.Value
	var memberExpr *parser.IndexExpr
	if ix, ok := e.Function.(*parser.IndexExpr); ok {
		memberExpr = ix
		obj, err := g.emitExpr(ix.Object)
		if err != nil {
			return nil, err
		}
		objForThis = obj
	}

	// Native: arr.concat(a, b, ...) -> fuji_array_concat([this, ...args])
	if memberExpr != nil && objForThis != nil {
		if lit, ok := memberExpr.Index.(*parser.LiteralExpr); ok {
			if name, ok := lit.Value.(string); ok && name == "concat" {
				callArgs := []value.Value{objForThis}
				for _, arg := range e.Arguments {
					val, err := g.emitExpr(arg)
					if err != nil {
						return nil, err
					}
					callArgs = append(callArgs, val)
				}
				return g.emitArgvRuntime(g.runtimeArrayConcat, callArgs), nil
			}
		}
	}

	// Emit the function
	fnVal, err := g.emitExpr(e.Function)
	if err != nil {
		return nil, err
	}

	// Set this for this call (JavaScript-like behavior)
	var thisVal value.Value
	if objForThis != nil {
		thisVal = objForThis
	} else {
		// In standalone function calls, this is undefined/null
		thisVal = constant.NewInt(types.I64, 0)
	}

	var args []value.Value
	for _, arg := range e.Arguments {
		val, err := g.emitExpr(arg)
		if err != nil {
			return nil, err
		}
		args = append(args, val)
	}

	if fn, ok := fnVal.(*ir.Func); ok && fn == g.runtimePrint {
		zero := constant.NewInt(types.I64, 0)
		if len(args) == 0 {
			g.block.NewCall(g.runtimePrintNewline)
			return zero, nil
		}
		if len(args) == 1 {
			return g.block.NewCall(g.runtimePrint, g.emitAsFujiI64(args[0])), nil
		}
		return g.emitArgvRuntime(g.runtimePrintArgv, args), nil
	}

	if fn, ok := fnVal.(*ir.Func); ok && fn == g.runtimeGcFrameStep {
		var budget value.Value
		if len(args) == 0 {
			budget = constant.NewFloat(types.Double, 0)
		} else {
			budget = g.block.NewCall(g.runtimeUnboxNumber, g.emitAsFujiI64(args[0]))
		}
		g.block.NewCall(g.runtimeGcFrameStep, budget)
		return constant.NewInt(types.I64, 0), nil
	}

	if fn, ok := fnVal.(*ir.Func); ok && isNativeArgvCallee(fn) {
		argCount := len(args)
		zero := constant.NewInt(types.I32, 0)
		var argvPtr value.Value
		if argCount == 0 {
			argvPtr = constant.NewNull(types.NewPointer(types.I64))
		} else {
			arrTy := types.NewArray(uint64(argCount), types.I64)
			slot := g.block.NewAlloca(arrTy)
			for i, arg := range args {
				argI64 := g.emitAsFujiI64(arg)
				elemPtr := g.block.NewGetElementPtr(arrTy, slot, zero, constant.NewInt(types.I32, int64(i)))
				g.block.NewStore(argI64, elemPtr)
			}
			argvPtr = g.block.NewGetElementPtr(arrTy, slot, zero, zero)
		}
		call := g.block.NewCall(fn, constant.NewInt(types.I32, int64(argCount)), argvPtr)
		if types.Equal(fn.Sig.RetType, types.Void) {
			_ = call
			return constant.NewInt(types.I64, 0), nil
		}
		return call, nil
	}

	if fn, ok := fnVal.(*ir.Func); ok && len(fn.Params) == len(args)+1 {
		finalArgs := []value.Value{thisVal}
		finalArgs = append(finalArgs, args...)
		call := g.block.NewCall(fn, finalArgs...)
		if types.Equal(fn.Sig.RetType, types.Void) {
			_ = call
			return constant.NewInt(types.I64, 0), nil
		}
		return call, nil
	}

	// Function value stored as i64 (ptrtoint of the LLVM function); recover pointer and call with this + args.
	if fnVal.Type().Equal(types.I64) {
		fnPtrTy := indirectFujiFuncPtrType(len(args))
		fnPtr := g.block.NewIntToPtr(g.emitAsFujiI64(fnVal), fnPtrTy)
		finalArgs := append([]value.Value{thisVal}, args...)
		call := g.block.NewCall(fnPtr, finalArgs...)
		return call, nil
	}

	call := g.block.NewCall(fnVal, args...)
	if fn, ok := fnVal.(*ir.Func); ok && types.Equal(fn.Sig.RetType, types.Void) {
		_ = call
		return constant.NewInt(types.I64, 0), nil
	}
	return call, nil
}
