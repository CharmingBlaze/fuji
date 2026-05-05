# Distribution guide

See [DISTRIBUTION_GUIDE.md](../../DISTRIBUTION_GUIDE.md) for the full reference.

This page summarizes the key commands.

---

## Quick reference

| Goal | Command |
|------|---------|
| Run during development | `fuji run game.fuji` |
| Build native executable | `fuji build game.fuji -o game.exe` |
| Create distributable folder | `fuji bundle game.fuji -o dist/game` |
| Include extra files in bundle | Set `FUJI_BUNDLE_FILES` |
| Use a C library wrapper | Set `FUJI_NATIVE_SOURCES` and `FUJI_LINKFLAGS` |

---

## Bundle output

```
dist/mygame/
  mygame.exe          <- compiled application
  run.bat             <- Windows launcher
  README.md           <- user-facing instructions
  bundle-info.txt     <- build metadata
  (extra files)       <- DLLs, assets, licenses from FUJI_BUNDLE_FILES
```

---

## Environment variables

| Variable | Purpose |
|----------|---------|
| `FUJI_CLANG` | Path to Clang executable |
| `CC` | Fallback compiler |
| `FUJI_USE_LLD` | Set `1` to use LLD linker |
| `FUJI_PATH` | Extra Fuji source search paths |
| `FUJI_WRAPPERS` | Extra wrapper search paths |
| `FUJI_NATIVE_SOURCES` | C/C++ wrapper glue files |
| `FUJI_LINKFLAGS` | Flags passed to Clang linker |
| `FUJI_BUNDLE_FILES` | Extra files copied into bundle |

---

## Cross-compiling the CLI

```powershell
$env:GOOS = "linux"; $env:GOARCH = "amd64"
go build -o kuji-linux-amd64 ./cmd/kuji
```

Supported targets: `windows/amd64`, `linux/amd64`, `linux/arm64`, `darwin/amd64`, `darwin/arm64`.
