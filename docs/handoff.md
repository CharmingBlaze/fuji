# Fuji implementation handoff

**Audience:** Engineers working on the **native LLVM** compiler (lexer → parser → sema → LLVM IR → clang → static binary).  
**Architecture overview:** [architecture.md](architecture.md).  
**Language surface:** [language.md](../language.md) (root) or [docs/language.md](language.md), [README.md](../README.md).  
**Shipped changes:** [CHANGELOG.md](../CHANGELOG.md).  
**Long-term roadmap (maintainers):** [MASTER_PLAN.md](MASTER_PLAN.md).

**Stack:** Go **1.22+**, module **`fuji`** ([go.mod](../go.mod)), LLVM IR via [llir/llvm](https://github.com/llir/llvm) v0.3.6, C11 runtime under **`runtime/src/`** (linked as **`runtime/libfuji_runtime.a`**).

There is **no bytecode VM** in the supported path: one pipeline end-to-end.

---

## 1. Product (CLI)

| Command | What it does | Toolchain |
|--------|----------------|-----------|
| **`fuji run`** | Load bundle → sema → LLVM → temp native exe → run | **llc** + **clang** + runtime archive (or embedded release binary) |
| **`fuji watch`** | Watch `.fuji` files under entry dir; rebuild + restart temp exe on change | Same as `run` |
| **`fuji build`** | Same lowering → linked **`*.exe`** / binary | Same |
| **`fuji check`** | Lexer + parser + **`sema.PrepareNativeBundle`** only (no LLVM) | Go only |
| **`fuji fmt`** | AST-based canonical format | Go only |
| **`fuji bundle`** | Package entry + assets for distribution | See **`cmd/fuji`** |

Release builds (`**-tags release**`) embed **llc** / **clang** (and **lld** on Windows) plus **`libfuji_runtime.a`** so end users need no LLVM install.

---

## 2. Repository map

| Path | Role |
|------|------|
| **`internal/lexer`** | Tokens; **`NewLexer(src, file)`** threads diagnostics paths |
| **`internal/parser`** | AST, **`LoadProgram`**, `#include` / import flattening, math prelude injection |
| **`internal/sema`** | **`PrepareNativeBundle`**, escape/shadow layout, arity, builtins prelude, **`Analyze`** |
| **`internal/codegen`** | LLVM **`ir.Module`** emission; **`internal/codegen/runtime.go`** declares **`fuji_*`** — must match **`runtime/src/fuji_runtime.c`** |
| **`internal/nativebuild`** | **llc** + **clang** invocation |
| **`internal/diagnostic`** | **`DiagnosticError`**, **`MultiError`**, snippet formatting |
| **`internal/formatter`** | **`fuji fmt`** |
| **`runtime/src`** | NaN-boxed **`Value`**, GC, objects, **`fuji_runtime_init`**, builtins |
| **`stdlib/*.fuji`** | Optional **`@math`**, **`@vec2`**, etc. |
| **`tests/*.fuji`**, **`examples/`** | Regression and demos |

---

## 3. Invariants (do not break)

1. **Sema blocks codegen** — **`PrepareNativeBundle`** errors must prevent emitting broken IR.
2. **Runtime / LLVM declare drift** — every **`fuji_*`** used from generated IR exists in C with the same ABI (**`internal/codegen/runtime.go`** ↔ **`fuji_runtime.c`**).
3. **Builtin names** — **`internal/sema/builtin_globals.go`** (sema prelude) and **`internal/codegen/builtin_register.go`** (codegen prelude) must agree for globals **`fuji build`** can link.
4. **Shadow stack** — every **`fuji_push_frame`** matched by **`fuji_pop_frame`** on **every** exit edge: each LLVM **`ret`** path must emit a pop (functions with multiple **`return`** statements have multiple pop sites). **`defer`** runs before pop on **`return`**.
5. **GC write barriers** — any new mutator that stores an **old → young** reference must use the same barrier pattern as **`fuji_object_set`** / **`fuji_array_set`**.

---

## 4. Session checklist

1. **`go test ./...`** (and **`go vet ./...`**) — same coverage as **`.github/workflows/ci.yml`** (Ubuntu, macOS, Windows).
2. **`go build -o fuji ./cmd/fuji`** then **`fuji run tests/hello.fuji`**
3. After runtime C changes: rebuild **`runtime/libfuji_runtime.a`** (see **`scripts/build-runtime.ps1`** / **`.sh`**)
4. **`CHANGELOG.md`** — update **`[Unreleased]`** for user-visible compiler or runtime changes; see **[releasing.md](releasing.md)** to cut **`v*`** tags.

---

## 5. Historical note

Older docs referred to a **bytecode VM** and **`kuji`** binaries; that dual path was removed so the tree matches **one** execution model. See **`docs/architecture.md`** § *What happened to the bytecode VM?* for context.
