# Getting started with Fuji

---

## Prerequisites

- Go 1.22 or later
- **Clang** and **llc** (LLVM) for **`fuji run`**, **`fuji build`**, and **`fuji bundle`** (all use the native pipeline)
- A C compiler (`gcc`/`clang`) and `ar` to build **`runtime/libfuji_runtime.a`** once (see [CONTRIBUTING.md](../../CONTRIBUTING.md))

On Windows, install LLVM from https://releases.llvm.org or via `winget install LLVM.LLVM`, and use MinGW-w64 for `gcc`/`ar` when running `scripts/build-runtime.ps1`.

---

## Build the CLI (and optional wrapgen)

From the repository root, **`fuji`** and **wrapgen** (C-header → `.fuji` generator) share the same `go.mod`.

```powershell
git clone <repo>
cd fuji-main

# One-shot (Windows): see scripts under repo root, or:
make build

# Or by hand
go build -o fuji.exe ./cmd/fuji
go build -o wrapgen.exe ./cmd/wrapgen
```

On macOS/Linux you can use `make build` to populate `./bin/`.

---

## Run a script

**`fuji run`** compiles with LLVM and executes the binary (same pipeline as **`fuji build`** — you need Clang + llc + built **`runtime/libfuji_runtime.a`**).

```powershell
.\fuji.exe run tests\hello.fuji
```

---

## Check and inspect

```powershell
.\fuji.exe check .\myscript.fuji        # parse and validate only
.\fuji.exe disasm .\myscript.fuji       # print LLVM IR (after sema + codegen)
```

---

## Build a native executable

```powershell
.\fuji.exe build .\myscript.fuji -o .\myscript.exe
```

Requires Clang on `PATH` (or set `FUJI_CLANG`).

---

## Bundle for distribution

```powershell
.\fuji.exe bundle .\mygame.fuji -o .\dist\mygame
```

This creates:

```
dist/mygame/
  mygame.exe
  run.bat
  README.md
  bundle-info.txt
```

To include extra files (DLLs, assets, licenses):

```powershell
$env:FUJI_BUNDLE_FILES = ".\raylib.dll .\assets\logo.png .\LICENSE.txt"
.\fuji.exe bundle .\mygame.fuji -o .\dist\mygame
```

---

## Hello World

```fuji
print("Hello, World!");
```

Run it:

```powershell
.\fuji.exe run .\hello.fuji
```

---

## Next steps

- [Language reference](../language/syntax.md) — full syntax documentation
- [Game development guide](game-dev.md) — using Raylib with Fuji
- [Distribution guide](distribution.md) — shipping apps and games
- [Wrapper guide](../../WRAPPERS.md) — integrating C/C++ libraries
