# Bundled LLVM + runtime (release builds only)

This directory is populated **before** `go build -tags release` (for example in CI). It is **not** used for normal development builds.

Expected layout (after population, before `go build -tags release`):

- `windows/amd64/llc.exe`, `lld.exe`, `libfuji_runtime.a`
- `linux/amd64/llc`, `lld`, `libfuji_runtime.a`
- `linux/arm64/llc`, `lld`, `libfuji_runtime.a`
- `darwin/amd64/llc`, `lld`, `libfuji_runtime.a`
- `darwin/arm64/llc`, `lld`, `libfuji_runtime.a`

Large binaries are listed in `.gitignore`; keep this `README.md` so the `bundled/` tree exists in git for documentation.

## Building `fuji` with embedded tools

From the repo root, after copying `llc`, `lld`, and `libfuji_runtime.a` into `bundled/<GOOS>/<GOARCH>/`:

```bash
go build -trimpath -tags release -ldflags="-s -w" -o fuji ./cmd/fuji
```

Without `-tags release`, `fuji` uses the normal resolution path ([ClangWithSource], [LLCWithSource], `runtime/libfuji_runtime.a` next to the project).
