# Embedded LLVM toolchain + runtime (release builds only)

Populated **before** `go build -tags release` (CI or `scripts/build-release.*`). Normal dev builds use `-tags` without `release` and never compile these assets.

Expected layout after population:

- `windows/amd64/clang.exe`, `llc.exe`, `lld.exe`, `libfuji_runtime.a`
- `linux/amd64/clang`, `llc`, `libfuji_runtime.a`
- `linux/arm64/clang`, `llc`, `libfuji_runtime.a`
- `darwin/amd64/clang`, `llc`, `libfuji_runtime.a`
- `darwin/arm64/clang`, `llc`, `libfuji_runtime.a`

`llc` is bundled on every platform so `fuji build` never depends on host LLVM tools.
`lld.exe` is bundled on Windows so Clang can use `-fuse-ld=lld` without relying on MSVC.

Large binaries are gitignored; small `README.txt` placeholders keep each OS directory present for `go:embed`.
