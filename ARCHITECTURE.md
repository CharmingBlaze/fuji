# Fuji — Architecture Overview

For the full compiler pipeline and invariants, see [`docs/architecture.md`](docs/architecture.md) and [`docs/handoff.md`](docs/handoff.md).

## Repository layout

| Path | What it is |
|---|---|
| `cmd/fuji/` | CLI entry point — `fuji build`, `fuji run`, `fuji check`, etc. |
| `cmd/wrapgen/` | `fujiwrap` — C header → Fuji bindings generator |
| `internal/lexer/` | Tokenizer |
| `internal/parser/` | AST types + recursive-descent parser |
| `internal/sema/` | Semantic analysis, escape analysis, shadow layout |
| `internal/codegen/` | LLVM IR emission (Go, using llir/llvm) |
| `internal/nativebuild/` | Invokes llc + clang + linker |
| `internal/formatter/` | `fuji fmt` |
| `internal/diagnostic/` | Error reporting with Rust-style source snippets |
| `internal/fujihome/` | Toolchain discovery and embedded binary extraction |
| `internal/embed/` | Release-build: embedded llc, lld, libfuji_runtime.a |
| `runtime/src/` | C runtime: GC, NaN-boxing, objects, shadow stack |
| `stdlib/` | Standard library as `.fuji` files |
| `wrappers/` | Pre-generated Raylib and other C library bindings |
| `api/` | Go API for embedding the Fuji compiler |
| `tests/` | `.fuji` test programs |
| `examples/` | Sample programs and games |
| `scripts/` | Build scripts for runtime and release packages |
| `docs/` | All documentation |
| `dist-template/` | Template for what goes into the release SDK zip |
| `_legacy/` | Old artifacts kept for reference — not part of the build |

## Pipeline

```
source.fuji
  → internal/lexer       (tokens)
  → internal/parser      (AST)
  → internal/sema        (analysis, escape analysis)
  → internal/codegen     (LLVM IR)
  → llc                  (object file)
  → clang + libfuji_runtime.a  (native binary)
```

## Key invariant

Every symbol in `internal/codegen/runtime.go` must have a matching implementation in `runtime/src/fuji_runtime.c` with the exact same C calling convention. If these drift, the linker produces a broken binary silently.

## Running tests

```bash
go test ./...
bash scripts/build-runtime.sh
./fuji run tests/hello.fuji
```
