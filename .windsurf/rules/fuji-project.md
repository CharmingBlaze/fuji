---
description: Fuji compiler — workspace facts for accurate edits
globs: **/*
---

## Module and entrypoints

- Go module: `fuji` (root `go.mod`, Go 1.22+). Main CLI: `cmd/fuji`.
- Compiler implementation: **`internal/parser`**, **`internal/lexer`**, **`internal/sema`**, **`internal/codegen`**, **`internal/nativebuild`**, **`internal/runtime`**, **`internal/fujihome`**. C runtime sources: **`runtime/src/`** (`fuji_runtime.c`, etc.) and optional embed under **`internal/runtime/data/`** (`fuji.c`, `fuji.h`).
- Wrapper generator: **`cmd/wrapgen`** — build as **`wrapgen`** (`go build -o wrapgen ./cmd/wrapgen`); same Go module as `fuji`. **`fuji wrap`** forwards to **`wrapgen`** (or legacy **`kujiwrap`**) when placed next to **`fuji`**.

## Commands to verify changes

From repo root (Windows PowerShell):

- `go test ./...`
- `go build -o fuji.exe ./cmd/fuji` then `.\fuji.exe run .\tests\smoke_native.fuji`
- Native gate: `.\fuji.exe build .\tests\smoke_native.fuji -o .\tests\smoke_native.exe` then run the exe (requires **`runtime/libfuji_runtime.a`** and LLVM tools unless using a **release**-tagged binary).

## Architecture (do not contradict docs in repo)

- **`fuji run`** / **`fuji build`** / **`fuji bundle`**: LLVM native pipeline (no separate bytecode VM).
- Raylib is **not** vendored as part of the language; local `temp_raylib/` is optional and gitignored. Wrappers live under `wrappers/`.

## When refactoring

Prefer small PR-sized steps; keep `go test ./...` and native smoke tests passing.
