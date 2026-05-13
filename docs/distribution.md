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

Build **fujiwrap** from this repo (same module as `fuji`; sources live under `cmd/wrapgen`):

```powershell
go build -o fujiwrap.exe ./cmd/wrapgen
```

Run it on one or more headers (or use **`fuji wrap …`**, which discovers **`fujiwrap.exe`** next to **`fuji.exe`**):

```powershell
.\fujiwrap.exe -name mylib -headers .\native\mylib.h -out .\wrappers\mylib
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
.\fuji.exe build .\app.fuji -o .\app.exe
```

For Raylib on Windows with the current local source tree:

```powershell
$env:FUJI_NATIVE_SOURCES = '..\wrappers\raylib_min\raylib_bridge.c'
$env:FUJI_LINKFLAGS = '-I..\temp_raylib\src -L..\temp_raylib\src -lraylib -lopengl32 -lgdi32 -lwinmm'
.\fuji.exe build .\raylib_brick_breaker.fuji -o .\raylib_brick_breaker.exe
```

## 5. Bundle an application or game for distribution

Use `fuji bundle` to create a clean folder that contains the executable, launcher, and metadata:

```powershell
.\fuji.exe bundle .\game.fuji -o .\dist\game
```

Bundle extra files or directories such as assets, DLLs, licenses, or config files:

```powershell
# Windows path-list style (;)
$env:FUJI_BUNDLE_FILES = '.\raylib.dll;.\LICENSE.txt;.\assets'
.\fuji.exe bundle .\game.fuji -o .\dist\game
```

```bash
# Linux/macOS path-list style (:)
FUJI_BUNDLE_FILES="./libfoo.so:./assets:./LICENSE.txt" fuji bundle game.fuji -o dist/game
```

The output folder contains:

| File | Purpose |
|------|---------|
| `game.exe` | The compiled application. |
| `run.bat` | Windows launcher. |
| `README.md` | User-facing run instructions. |
| `bundle-info.txt` | Build metadata and native link settings. |
| Extra files | Any files or directories listed in `FUJI_BUNDLE_FILES`. |

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

## 7. Official Releases: Windows, Linux, macOS SDK zips

GitHub **Releases** (tags `v*`) attach **offline SDK archives** built by CI for **every mainstream OS**:

| OS | Artifact |
|----|----------|
| **Windows** x64 | `fuji-<tag>-sdk-windows-amd64.zip` |
| **Linux** x64 | `fuji-<tag>-sdk-linux-amd64.zip` |
| **Linux** ARM64 | `fuji-<tag>-sdk-linux-arm64.zip` |
| **macOS** Intel | `fuji-<tag>-sdk-darwin-amd64.zip` |
| **macOS** Apple Silicon | `fuji-<tag>-sdk-darwin-arm64.zip` |

Each zip unpacks to a single folder containing **`fuji`** (or **`fuji.exe`**), **`fujiwrap`**, **`stdlib/`**, the full **`docs/`** tree, **every repo-root `*.md`**, **`wrappers/`** (including the **full raylib** binding + `wrapper.c` + reference docs), **`third_party/raylib_static/stage/`** with **raylib 5.0** headers and libraries (and **`raylib.dll`** on Windows; Linux ARM64 may ship a **README** instead of prebuilt libs — see that file), **`examples/`**, and **`SDK_README.txt`**. Keep **`stdlib/`** next to the compiler so `@` imports resolve with no extra download. Raylib workflow: [guides/raylib.md](guides/raylib.md).

Maintainers reproduce the archives by pushing a version tag; the workflow is **[`.github/workflows/release.yml`](.github/workflows/release.yml)** and **`scripts/package-release-sdk.sh`**.

### 7.1 Local Windows SDK (same layout as the release zip)

After you have release **`fuji.exe`** / **`fujiwrap.exe`** (embedded Clang + llc + lld + runtime), assemble the full distributable folder:

```powershell
powershell -File scripts/assemble-offline-sdk.ps1 -FujiExe .\fuji-release.exe -FujiwrapExe .\fujiwrap.exe -Zip
```

Or build the compiler then package in one step:

```powershell
powershell -File scripts/build-release.ps1 -PackageSdk -PackageSdkZip
```

Options: **`-SkipRaylib`** (no download; add **`third_party/raylib_static/stage`** yourself), **`-RaylibZipPath`** (use a downloaded **`raylib-5.0_win64_mingw-w64.zip`** offline). The result matches what users get from **Releases** for Windows amd64: static **`libraylib.a`** for linking plus **`raylib.dll`** from the official prebuild when you need it beside shipped **`.exe`** files.

### Manual “toolchain folder” layout (maintainers)

If you build from source locally, mirror the same layout next to **`fuji`**:

```powershell
go build -o fuji.exe .\cmd\fuji
go build -o fujiwrap.exe .\cmd\wrapgen
```

```text
install-root/
  fuji.exe           # or `fuji` on Linux/macOS
  fujiwrap.exe
  stdlib/
  docs/
  *.md               # all repo-root markdown
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
