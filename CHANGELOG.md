# Changelog

All notable changes to this project are documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.1.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

### Added

- **Language surface (native / LLVM)** — template literals with **`${}`**; unary **`typeof`**; **`delete obj.key`**; array literal spread **`[...a]`**; **`??`** and **`?.`**; **`let { x, y } = obj`** destructuring; **`matches(haystack, pattern)`** (substring / **`strstr`**, not full regex); string methods (**`split`**, **`join`** on arrays, **`trim`**, **`replace`** / **`replaceAll`**, **`indexOf`**, **`toUpper`** / **`tolower`**, **`slice`**, **`startsWith`** / **`endsWith`**); array methods (**`map`**, **`filter`**, **`reduce`**, **`find`**, **`slice`**, **`sort`**, **`reverse`**, **`includes`**, **`concat`**, **`join`**); prelude **`math`** object (**`math.floor`**, …) alongside existing math globals; runtime symbols **`fuji_array_join`**, **`fuji_object_remove`**, **`fuji_matches`** (rebuild **`runtime/libfuji_runtime.a`** after pulling).
- **`fuji fmt`** — canonical AST-based formatting (`internal/formatter`): 4-space indent, spacing around operators and after commas, `if (` / `while (` / `for (` style, `} else {` kept on one line when both branches are blocks, `import "…"` expressions, `fuji fmt --check`, directory and **`./...`** expansion (skips `.git`, `.FUJI_build`, `bin`, `node_modules`, `vendor`). Top-level **consecutive** expression statements (e.g. back-to-back `print`) stay compact; other top-level declarations stay separated by a blank line.
- **Classic `for (init; cond; step)`** — parsed and lowered alongside `while` / `for-in` / `for-of` (`internal/parser`, `internal/codegen`).
- **`for (let [k, v] of iterable)`** — destructuring **`for-of`** for **arrays** (numeric index + element) and **objects/tables** (insertion-order key + value); runtime helpers **`fuji_forof_length`**, **`fuji_forof_key_at`**, **`fuji_forof_value_at`** (`runtime/src/fuji_runtime.c`).
- **Release binaries (`-tags release`)** — embedded Clang + runtime under **`internal/embed/`** (Windows also bundles **`lld.exe`** for linking); **`scripts/build-release.ps1`** / **`scripts/build-release.sh`**; **`release.yml`** publishes **`fujiwrap`** where applicable (`CHANGELOG` superseded notes under 0.1.0 remain historical).
- **Runtime hardening (GC/debug):** dynamic global-slot root registration (`fuji_register_global_slot` + generated top-level slot wiring), opt-in **`FUJI_GC_DEBUG`** stats/checks (remembered-set overflow counter, shadow stack depth high-water mark, global slot stats, object-header sanity checks), and new stress scenarios (`tests/gc_pressure_expr.fuji`, `tests/globals_perf.fuji`, `tests/gc_soak.fuji`).

### Changed

- **`len(table)`** — returns **entry count** for object/table values (`fuji_len`).
- **`for-of`** / **`for-in`** lowering — uses slot-ordered runtime iteration for arrays and tables (**`for-in`** binds **keys**; **`for-of`** binds **values**; **`tests/phase1_surface.fuji`** uses **`for-of`** where element sums are intended).
- **Documentation** — loop-style guide (`docs/user_guide.md`), **`for-of`** / **`syntax.md`** examples, single-line vs multi-line braced bodies.
- **Language (case):** reserved words and ASCII identifiers are **case-insensitive** (identifiers and object property keys are normalized to lowercase in the AST). **`@` module** specifiers and **`#include`** are matched case-insensitively; **`// fuji:…`** lines treat the `fuji:` prefix and the **`extern`** keyword case-insensitively, while the **C symbol** on extern lines stays spelled exactly for the linker. **`main`** / **`Main`** / etc. still map to the native entry (`fuji_user_main`) case-insensitively. Whitespace between tokens remains free-form; **semicolons between statements are still required.**
- **Branding and paths:** sources use the **`.fuji`** extension; C runtime and LLVM symbols use the **`fuji_`** prefix; static archive **`runtime/libfuji_runtime.a`**. The desktop IDE module is **`fuji-ide`** (Wails app **`fuji-ide`**). **Wrapper generator:** canonical binary name **`fujiwrap`** (`go build -o fujiwrap ./cmd/wrapgen`); generated files credit **`fujiwrap`**. **`fuji wrap`** prefers **`fujiwrap`** next to **`fuji`**, then **`wrapgen`**, then legacy **`kujiwrap`**.
- **CI hardening:** Linux CI now runs a serialized, time-bounded GC soak pass (`gc_pressure_expr`, `globals_perf`, `gc_soak`) to catch long-run/rooting regressions.

### Fixed

- **`fujihome`**: removed duplicate **`FUJI_CLANG`** / **`FUJI_LLC`** env branches in **`ClangWithSource`** / **`LLCWithSource`**.

## [0.1.0] - 2026-05-05

First self-contained **Fuji** release: native LLVM pipeline, embedded toolchain option, GC stack, and Result-based errors.

### Added

- **`ok(value)`** / **`err(message)`** — Result objects `{ ok, value, error }` for recoverable errors (plain tables; no exceptions).
- **`panic(message)`** — unrecoverable error: message to stderr, optional call-stack trace, **`exit(1)`** (not catchable).
- **`assert(condition, message?)`** — debug assertions using truthy rules; failures use the panic path.
- **`readFile`** / **`writeFile`** — return **`ok(...)`** or **`err(...)`** (real I/O; errno messages on failure).
- **`parseJSON`** — returns **`ok(parsed)`** or **`err(message)`**; JSON subset (objects, arrays, strings, numbers, booleans, null).
- **Call stack tracking** — codegen emits call-frame push/pop into the C runtime so panics can print a Fuji stack trace.
- **Self-contained release binary** (Windows, Linux, macOS): **`go build -tags release`** embeds LLVM **`llc`** / **`lld`** and **`libfuji_runtime.a`**; **`fujihome.FindToolchain`** extracts them on first use. GitHub Actions **`.github/workflows/release.yml`** publishes **`fuji-windows-amd64.exe`**, **`fuji-linux-amd64`**, **`fuji-linux-arm64`**, **`fuji-darwin-amd64`**, **`fuji-darwin-arm64`**.
- **`fuji build --no-opt`** — disables IR optimisation and forces **`llc -O0`**; **`nativebuild.BuildOptions`** and related API surface the same.
- **`codegen.OptimiseIR`** — stub when **`CGO_ENABLED=0`** or without **`-tags llvm14`**; optional **tinygo.org/x/go-llvm** **`default<O2>`** pipeline when CGo + LLVM dev libs are available.
- Builtin **`gcFrameStep(ms)`** — maps to **`fuji_gc_frame_step`** for game-loop GC integration (optional budget parameter).
- **Nursery GC** and **shadow stack** — bump nursery, precise roots via shadow frames, write barrier (see runtime and sema packages).
- **Builtin registration** in codegen — global stdlib names (**`readFile`**, **`sin`**, **`ok`**, **`assert`**, etc.) resolve consistently for native emission.
- Lexer/parser: **`var`** is reserved; use **`let`** for declarations.
- GitHub Actions **`.github/workflows/ci.yml`** — **`go vet`**, **`go test`**, C runtime build, native smoke.
- **`scripts/build-runtime.sh`** and **`scripts/build-runtime.ps1`** — produce **`runtime/libfuji_runtime.a`**.
- **`tests/smoke_native.fuji`**, **`tests/error_handling_test.fuji`** — CI and local checks.
- **`internal/codegen`** tests that LLVM runtime symbols exist in the static archive.
- **`cmd/wrapgen`** parse regression test (**`testdata/minimal.h`**).

### Fixed

- **`null`** literal codegen now uses the correct NaN-boxed **`NIL_VAL`** pattern (e.g. **`type(null)`** in **`tests/hello.fuji`**).
- Array literal **`[1, 2, 3]`** and index expressions **`arr[i]`** (including chained indexing) parse correctly (Pratt handlers for **`[`** and indexing).

### Changed

- Docs and scripts consistently reference **`runtime/libfuji_runtime.a`** and **`runtime/src/fuji_runtime.c`**.
- **`runtime/Makefile`** writes **`runtime/libfuji_runtime.a`** at the path **`fuji build`** expects.
- Root **`make test`** depends on **`make -C runtime`** so the archive exists before link tests.
- Documentation describes the native pipeline: LLVM IR → **llc** → object → **Clang** + **`runtime/libfuji_runtime.a`** (or embedded toolchain when built with **`release`** tag).
- LLVM **`declare`** names aligned with **`fuji_runtime.c`** (including **`fuji_print_val`**, **`fuji_get`**, argv-style builtins).

### Removed

- Bytecode VM package (**`internal/vm/`**) — execution is the **LLVM native** pipeline only (**`fuji run`**, **`fuji build`**, **`fuji bundle`**, **`api.Run`**).
- Standalone interpreter / **`Value = any`** Go runtime — values are NaN-boxed in C and **i64** in IR.
- Separate **`fuji native`** as a distinct pipeline — **`fuji run`** is the primary path (legacy alias may remain).

### Deprecated

- The large embed under **`internal/runtime/data/`** is not used by default **`fuji build`** (see **`internal/runtime/embed.go`**).
