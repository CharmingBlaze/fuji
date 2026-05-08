# Contributing to Fuji

## ⚠️ This is for contributors only

If you want to **use** Fuji to make games or applications, **do not build from source**. Download a release binary from [GitHub Releases](https://github.com/CharmingBlaze/fuji/releases) instead. The release binary is self-contained and does not require any of the toolchain described in this document.

Building from source requires Go 1.22+, Clang + llc on PATH, and a C toolchain. The embedded-toolchain release build (`-tags release`) additionally requires having the LLVM binaries available to embed. This is a non-trivial setup that is only worth doing if you are modifying the compiler itself.

## Build

From the repo root (Go **1.22+**):

```bash
go build -o bin/fuji ./cmd/fuji
```

Build the C runtime archive when linking native binaries (needed for **`go test`** / **`fuji build`**):

- Linux / macOS: `bash scripts/build-runtime.sh`
- Windows: `powershell -File scripts/build-runtime.ps1`

## Tests

```bash
go vet ./...
go test ./... -count=1
```

## Releases

See **[docs/releasing.md](docs/releasing.md)** for version bumps, changelog, and tagging **`v*`** so **`.github/workflows/release.yml`** runs.

## Git: commit and push to `main`

Step-by-step: **[docs/git-workflow.md](docs/git-workflow.md)**.

## Formatting

Canonical `.fuji` style:

```bash
go build -o bin/fuji ./cmd/fuji
./bin/fuji fmt ./...
./bin/fuji fmt --check ./...
```

CI runs **`fuji fmt --check ./...`** after building **`fuji`**.

## Compiler diagnostics

Lexer and parser failures include the source line and a caret when reporting errors from **`parser.LoadProgram`** (used by **`fuji run`**, **`fuji build`**, etc.). **`fuji check`** runs the same load step plus **`sema.PrepareNativeBundle`** (semantic errors, no LLVM).

## Native toolchain

LLVM **`clang`** / **`llc`** (and optionally **`lld`**) are required for **`fuji build`** / **`fuji run`** unless you use a **release** build that embeds the toolchain. See **`fuji doctor`** and **`fuji help`**.

## Example programs

- **`tests/`** — language and smoke tests.
- **`examples/games/`** — small standalone demos (**`examples/games/README.md`**).
- **`demos/`** — larger samples (some need Raylib + wrappers).
