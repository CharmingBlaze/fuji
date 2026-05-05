# Contributing to Fuji

## Prerequisites

- **Go** 1.22 or later
- For **native** `fuji build` / `bundle` / `native`: **Clang** and **llc** on `PATH` (or set `FUJI_CLANG` / `FUJI_LLC`), plus a C toolchain able to build the runtime archive (`gcc` and `ar`, or MSVC equivalents on Windows).

## Build the C runtime (`libfuji_runtime.a`)

Native programs link against `runtime/libfuji_runtime.a`, built from `runtime/src/*.c`.

- **Linux / macOS (GNU Make):** from the repo root:

  ```bash
  make -C runtime
  ```

  Or: `bash scripts/build-runtime.sh`

- **Windows (PowerShell):** from the repo root (requires MinGW-w64 `gcc` + `ar` on PATH, or set `CC` / `AR`):

  ```powershell
  .\scripts\build-runtime.ps1
  ```

  For a full **LLVM + MSVC + go-llvm** layout (PATH, CGo, Developer Prompt), see **[docs/windows-native-toolchain.md](docs/windows-native-toolchain.md)**. With only the [official LLVM installer](https://github.com/llvm/llvm-project/releases/tag/llvmorg-14.0.6), use e.g. `CC=clang` and `AR=llvm-ar` when invoking `make -C runtime` or `build-runtime.ps1`.

**Release binaries (embedded LLVM + `libfuji_runtime.a`):** pushing a **`v*`** tag runs **`.github/workflows/release.yml`**, which populates **`internal/fujihome/bundled/<GOOS>/<GOARCH>/`** on the runner, then **`go build -tags release`** so **`llc`**, **`lld`**, and the static archive ship inside **`fuji`**. See **`internal/fujihome/bundled/README.md`**.

**Optional: LLVM IR optimisation in the compiler (`OptimiseIR`):** CI and release workflows set **`CGO_ENABLED=0`**, so **`internal/codegen/optimise_stub.go`** is used (no-op). To compile the real **`optimise_cgo.go`** implementation you need **CGo**, LLVM **development** libraries matching the **`llvm14`** build tag in **tinygo.org/x/go-llvm**, and:

```bash
go build -tags llvm14 ./cmd/fuji/...
```

(Linux example: `llvm-config-14` on `CPPFLAGS`/`LDFLAGS`; Windows/macOS paths are in the tinygo.org module’s `llvm_config_llvm14.go`.)

The root `Makefile` target `make test` runs `make -C runtime` first so `go test` (including codegen link-symbol checks) has the archive present.

## Build the compilers

```bash
go build -o bin/kuji ./cmd/kuji
go build -o bin/kujiwrap ./cmd/wrapgen
```

On Windows you can use `.\scripts\build.ps1`, which builds both binaries into `./bin`.

## Run tests and vet

```bash
go vet ./...
go test ./... -count=1
```

## Native smoke check

After `runtime/libfuji_runtime.a` exists:

```bash
./bin/fuji build tests/smoke_native.fuji -o /tmp/smoke_native
/tmp/smoke_native
```

Optional third-party glue (e.g. Raylib) uses `FUJI_NATIVE_SOURCES` and `FUJI_LINKFLAGS`; omit them for the smoke program.

## Release version string

`cmd/kuji` embeds `var version = "0.2.0-dev"`. Release builds should pass:

```bash
go build -ldflags "-X main.version=1.2.3" -o bin/kuji ./cmd/kuji
```

## LLVM ↔ C runtime contract

When you add or rename a runtime entry point:

1. Update `declareRuntimeFunctions` in [`internal/codegen/runtime.go`](internal/codegen/runtime.go) (LLVM symbol name and signature).
2. Implement the same symbol in [`runtime/src/fuji_runtime.c`](runtime/src/fuji_runtime.c) (and declare it in [`fuji_runtime.h`](runtime/src/fuji_runtime.h) before use from codegen).
3. Rebuild `runtime/libfuji_runtime.a` and run `go test ./internal/codegen/...`.

### C symbol naming: `fuji_snake_case`

The linkable C API uses **`fuji_snake_case`** (`fuji_unbox_number`, `fuji_get`, `fuji_alloc_cell`, …). That is the only convention for **exported runtime functions**; it matches the existing `fuji_*` symbols in `fuji_runtime.c` and what LLVM `declare` / `call` use.

In Go, `runtime.go` may use map keys like `FUJI_print` or `FUJI_get` as **organising labels** — the string passed to `mod.NewFunc` must still be the real C symbol (e.g. `fuji_print_val`, `fuji_get`). Do not introduce a second `FUJI_*` naming layer in C unless the whole runtime is being renamed.

**`FUJI_` in docs** usually means **environment variables** (`FUJI_CLANG`, `FUJI_PATH`, `FUJI_LINKFLAGS`, …) or project branding, not C function names. External briefs that say `FUJI_snake_case` for runtime entry points should be read as **`fuji_snake_case`** for this repository.

The embed under `internal/runtime/data/` is **not** the default link target for native `fuji build`; fix native bugs in `runtime/src/` unless you are explicitly working on that alternate tree.
