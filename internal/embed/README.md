# Embedded Clang + runtime (release builds only)

Populated **before** `go build -tags release` (CI or `scripts/build-release.*`). Normal dev builds use `-tags` without `release` and never compile these assets.

Expected layout after population:

- `windows/amd64/clang.exe`, `lld.exe`, `libfuji_runtime.a`
- `linux/amd64/clang`, `libfuji_runtime.a`
- `linux/arm64/clang`, `libfuji_runtime.a`
- `darwin/amd64/clang`, `libfuji_runtime.a`
- `darwin/arm64/clang`, `libfuji_runtime.a`

`lld.exe` is bundled on Windows so Clang can use `-fuse-ld=lld` without relying on MSVC.

Large binaries are gitignored; small `README.txt` placeholders keep each OS directory present for `go:embed`.
