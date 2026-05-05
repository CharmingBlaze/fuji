# Fuji

**Fuji** is a **JavaScript-syntax** scripting language that compiles to **native binaries** via LLVM. That combination is **implemented today**, not only planned: **simple JavaScript**—objects, functions, familiar control flow—and **simple C**—**`fuji build`** produces a native executable (LLVM IR → object + **`libfuji_runtime.a`**), explicit statements, predictable performance, and **no separate runtime** for players to install.

The **[language.md](language.md)** catalog lists the full surface (including **string/array methods**, **template `` `${}` `` literals**, **`math.*`** prelude, **`typeof`**, **`delete`**, **spread**, **destructuring**, **`?.`**, **`??`**, and **`matches()`** for substring checks).

A small language for **games and apps**: **plain `.fuji` files** you can read, diff, and ship. The compiler and CLI are **`fuji`**. Use **`fuji fmt`** for canonical formatting (`CONTRIBUTING.md`).

| Who | What they need |
|-----|----------------|
| **Players / Customers** | Only your **bundle** folder (`.exe` + assets). No Go, Python, or C++. |
| **You (Author)** | **Go** to build **`fuji`** once. Either a **self-contained release binary** (no LLVM install) or **Clang** + **llc** on `PATH`, plus a C compiler to build **`runtime/libfuji_runtime.a`** when not using a prebuilt archive. Optional **`fujiwrap`** (same sources as **`wrapgen`**) for C header bindings. |

See **[CHANGELOG.md](CHANGELOG.md)** for **v0.1.0** (Result errors, GC, embedded toolchain, first tagged release). Contributors: **[CONTRIBUTING.md](CONTRIBUTING.md)** (tests, **`fuji fmt --check`**, runtime build).

### Verify v0.1.0 locally (Windows)

Use **Developer Command Prompt for VS 2022** (or any shell where **`clang`**, **`llc`**, and a C compiler are on **`PATH`**).

```bat
cd path\to\fuji-main
.\scripts\build-runtime.ps1
go build -o fuji.exe ./cmd/fuji
fuji.exe build tests\hello.fuji -o hello.exe
hello.exe
fuji.exe build tests\error_handling_test.fuji -o error_test.exe
error_test.exe
fuji.exe build tests\phase1_surface.fuji -o surface.exe
surface.exe
```

**Expected (roughly):** **`hello.fuji`** prints the Fuji runtime init line plus type lines for number, string, boolean, and null. **`error_handling_test.fuji`** ends with **`PASS: error handling works`**. **`phase1_surface.fuji`** follows whatever that test prints today (see the file under **`tests/`**).

From a **GitHub Release** artifact (**`fuji-windows-amd64.exe`**), the same commands work **without** installing LLVM, because the binary embeds **`llc`**, **`lld`**, and **`libfuji_runtime.a`** (**`go build -tags release`**).

---

## Quick Start

Fuji compiles **only** through the LLVM native pipeline (`fuji run` builds a temp binary and runs it — same lowering as `fuji build`).

**Build the toolchain:**
```bash
make build    # ./bin/fuji and ./bin/fujiwrap (see Makefile; `make wrapgen` also builds legacy bin/wrapgen)
```

**Run an example:**
```bash
./bin/fuji run examples/hello.fuji
```

**Build a native app** (requires **`llc`**, **`clang`**, and **`runtime/libfuji_runtime.a`** unless using a **release-tagged** `fuji` binary):
```bash
./bin/fuji build examples/hello.fuji -o hello.exe
```

**Bundle for distribution:**
```bash
./bin/fuji bundle examples/games/production_brick_breaker.fuji -o dist/BrickBreaker
```

---

## Project Structure

- **`bin/`**: Compiled toolchain binaries.
- **`cmd/`**: CLI source code (compiler and wrapper generator).
- **`docs/`**: Extensive documentation, guides, and references.
- **`examples/`**: Sample Fuji programs and games.
- **`internal/`**: Core compiler and runtime implementation.
- **`stdlib/`**: The Fuji standard library.
- **`wrappers/`**: C/C++ library wrappers.
- **`dist/`**: Output directory for generated application bundles.

---

## Documentation

| Document | Contents |
|----------|----------|
| [language.md](language.md) | **Single-page catalog** — keywords, operators, statements, builtins |
| [docs/user_guide.md](docs/user_guide.md) | Comprehensive guide to Fuji |
| [docs/reference.md](docs/reference.md) | Standard library and built-in functions |
| [docs/language.md](docs/language.md) | Syntax and language specification |
| [docs/distribution.md](docs/distribution.md) | How to ship your games |
| [docs/wrappers.md](docs/wrappers.md) | C/C++ FFI and Raylib integration |
| [docs/architecture.md](docs/architecture.md) | Compiler pipeline (LLVM native only) |
| [docs/windows-native-toolchain.md](docs/windows-native-toolchain.md) | Windows 11: LLVM 14, MSVC, Make, CGo for go-llvm |

For the included Brick Breaker on Windows (raylib shim path), run:
```powershell
powershell -ExecutionPolicy Bypass -File .\scripts\play-brick-breaker.ps1
```

**Release downloads** (GitHub **Releases**, tag `v*`): self-contained **`fuji`** binaries with embedded **`llc`**, **`lld`**, and **`libfuji_runtime.a`** are built by **`.github/workflows/release.yml`** (`go build -tags release`). No LLVM install is required to use those artifacts for **`fuji build`** / **`fuji run`**.

---

## Requirements

- **Authors:** Go **1.22+**, **Clang** and **llc** on `PATH` for native builds, plus a C toolchain to build **`runtime/libfuji_runtime.a`** (see [CONTRIBUTING.md](CONTRIBUTING.md)).
- **Players:** Nothing required — just your shipped executable and data files.
