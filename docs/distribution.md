# Fuji build, wrapper, and distribution guide

This guide explains the complete workflow for writing Fuji programs, generating wrappers for C/C++ libraries, building native executables, and packaging applications or games for distribution.

## 1. Run a Fuji program during development

`fuji run` compiles with the **LLVM native** pipeline and runs the resulting binary (same as `fuji build`, temp output):

```powershell
.\fuji.exe run .\game.fuji
```

Useful checks:

```powershell
.\fuji.exe check .\game.fuji
.\fuji.exe disasm .\game.fuji   # prints LLVM IR text
```

## 2. Build a native executable

`fuji build` emits LLVM IR, lowers it with **llc** to an object file, and links that object with **`runtime/libfuji_runtime.a`** (see `runtime/src/`) using **Clang**, plus any **`FUJI_NATIVE_SOURCES`** / **`FUJI_LINKFLAGS`** you set for third-party libraries.

```powershell
.\fuji.exe build .\game.fuji -o .\game.exe
```

Optional environment variables:

| Variable | Purpose |
|----------|---------|
| `FUJI_CLANG` | Full path to the Clang executable. |
| `CC` | Fallback compiler name if `FUJI_CLANG` is not set. |
| `FUJI_USE_LLD` | Set to `1` to request LLVM LLD linking. |
| `FUJI_PATH` | Extra Fuji source search paths. |
| `FUJI_WRAPPERS` | Extra generated wrapper search paths. |
| `FUJI_NATIVE_SOURCES` | C/C++ wrapper glue files to compile into the executable. |
| `FUJI_LINKFLAGS` | Native include/library/link flags passed to Clang. |

## 3. Generate wrappers for C/C++ libraries

Build **kujiwrap** from this repo (same module as `fuji`; sources live under `cmd/wrapgen`):

```powershell
go build -o kujiwrap.exe ./cmd/wrapgen
```

Run it on one or more headers (or use `kuji wrap …` if `kujiwrap.exe` sits next to `kuji.exe`):

```powershell
.\kujiwrap.exe -name mylib -headers .\native\mylib.h -out .\wrappers\mylib
```

Generated output:

| File | Purpose |
|------|---------|
| `mylib.fuji` | Readable Fuji source with `// fuji:extern` lines linking to C. |
| `wrapper.c` | C ABI glue that converts Fuji values to native C calls. |
| `README.md` | Human-readable usage guide. |
| `api_reference.md` | Function/type reference. |
| `examples.md` | Example usage (when docs are enabled). |
| `Makefile` / `CMakeLists.txt` | Optional (`-build`); off by default. |

Include the generated Fuji wrapper from your program:

```fuji
#include "wrappers/mylib/mylib.fuji"

let result = my_function(1, 2);
print(result);
```

## 4. Build with a wrapper and native library

Set the generated C glue and native link flags before building:

```powershell
$env:FUJI_NATIVE_SOURCES = '..\wrappers\mylib\wrapper.c ..\native\mylib.c'
$env:FUJI_LINKFLAGS = '-I..\native -L..\native\build -lmylib'
.\kuji.exe build .\app.fuji -o .\app.exe
```

For Raylib on Windows with the current local source tree:

```powershell
$env:FUJI_NATIVE_SOURCES = '..\wrappers\raylib_min\raylib_bridge.c'
$env:FUJI_LINKFLAGS = '-I..\temp_raylib\src -L..\temp_raylib\src -lraylib -lopengl32 -lgdi32 -lwinmm'
.\kuji.exe build .\raylib_brick_breaker.fuji -o .\raylib_brick_breaker.exe
```

## 5. Bundle an application or game for distribution

Use `fuji bundle` to create a clean folder that contains the executable, launcher, and metadata:

```powershell
.\kuji.exe bundle .\game.fuji -o .\dist\game
```

Bundle extra files such as assets, DLLs, licenses, or config files:

```powershell
$env:FUJI_BUNDLE_FILES = '.\raylib.dll .\LICENSE.txt .\assets\logo.png'
.\kuji.exe bundle .\game.fuji -o .\dist\game
```

The output folder contains:

| File | Purpose |
|------|---------|
| `game.exe` | The compiled application. |
| `run.bat` | Windows launcher. |
| `README.md` | User-facing run instructions. |
| `bundle-info.txt` | Build metadata and native link settings. |
| Extra files | Any files listed in `FUJI_BUNDLE_FILES`. |

## 6. Ship a wrapper package

A clean wrapper package should include:

```text
mylib-wrapper/
  mylib.fuji
  wrapper.c
  README.md
  api_reference.md
  examples.md
  native/
    mylib.dll or mylib.a if redistribution is allowed
    LICENSE.txt
```

Users can either include `mylib.fuji` directly or set:

```powershell
$env:FUJI_WRAPPERS = 'C:\path\to\mylib-wrapper'
```

## 7. Ship a complete Fuji toolchain folder

For distributing Fuji itself, build the CLI and include the docs and wrappers:

```powershell
go build -o kuji.exe .\cmd\kuji
go build -o cmd\wrapgen\wrapgen.exe .\cmd\wrapgen
```

Suggested release folder:

```text
kuji-release/
  kuji.exe
  wrapgen.exe
  README.md
  RELEASE.md
  DISTRIBUTION_GUIDE.md
  FUJI_PROGRAMMER_REFERENCE.md
  WRAPPERS.md
  wrappers/
  examples/
```

## 8. Professional release checklist

Before shipping an app or game:

- Build with `fuji build` or `fuji bundle`.
- Run the executable from the output folder, not from the source tree.
- Include native DLLs or data files required by the app.
- Include licenses for any third-party native libraries.
- Keep wrapper docs with the wrapper package.
- Avoid source-only temp folders in the final app bundle.
- Verify on a clean machine or clean terminal session.
