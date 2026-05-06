# Fuji

![Fuji Logo](assets/fuji-logo.png)

**Fuji** is a **JavaScript-style** language that compiles **`.fuji`** programs to **native executables** (LLVM). You write familiar syntax—objects, functions, control flow—and ship a **single binary** (plus any assets you bundle). Players do **not** install a separate Fuji runtime.

The **[language.md](language.md)** page is the compact catalog of syntax, operators, builtins, and stdlib surface.

---

## Start here (important)

- Use the prebuilt **release binaries** from [GitHub Releases](https://github.com/CharmingBlaze/fuji/releases).
- Do **not** compile Fuji from source for normal usage. This repo includes maintainer/developer internals and is not the beginner install path.
- If your goal is to make games/apps with Fuji, use `fuji` + `fujiwrap` binaries and the guides below.

---

## Get the compiler (`fuji`)

**Recommended:** download a **release build** from **[GitHub Releases](https://github.com/CharmingBlaze/fuji/releases)** (tags **`v*`**). Those binaries embed **Clang**, **`libfuji_runtime.a`**, and on Windows **lld**, so you can **`fuji build`** / **`fuji run`** without installing LLVM yourself.

| Platform | Compiler binary | Typical download name |
|----------|-----------------|-------------------------|
| Windows (x64) | `fuji` | **`fuji-windows-amd64.exe`** |
| Linux (x64) | `fuji` | **`fuji-linux-amd64`** |
| Linux (ARM64) | `fuji` | **`fuji-linux-arm64`** |
| macOS Intel | `fuji` | **`fuji-darwin-amd64`** |
| macOS Apple Silicon | `fuji` | **`fuji-darwin-arm64`** |

Put the binary on your **`PATH`**, or run it by full path. Then:

```bash
fuji version
```

On **Linux / macOS**, mark the file executable after download: `chmod +x fuji-linux-amd64` (example).

---

## Get the C header wrapper (`fujiwrap`)

**`fujiwrap`** turns C/C++ headers into **`.fuji`** bindings plus a **`wrapper.c`** you link with **`FUJI_NATIVE_SOURCES`**. It ships **next to `fuji`** on the same **[Releases](https://github.com/CharmingBlaze/fuji/releases)** page when published for your platform (for example **`fujiwrap-windows-amd64.exe`**, **`fujiwrap-linux-amd64`**, **`fujiwrap-darwin-arm64`**).

Run **`fuji wrap …`** from the CLI: **`fuji`** looks for **`fujiwrap`** beside itself, then **`wrapgen`**, then **`kujiwrap`**, then **`PATH`**.

```bash
fuji wrap --help
```

Full workflow: **[docs/wrappers.md](docs/wrappers.md)** and **`fuji help`** (environment variables **`FUJI_NATIVE_SOURCES`**, **`FUJI_LINKFLAGS`**, **`FUJI_BUNDLE_FILES`**, etc.).

---

## Use `fuji` on your project

From a directory that contains your entry **`.fuji`** file (and any **`#include`** / **`@module`** dependencies the loader can resolve):

```bash
# Run (builds a temp exe, runs it, deletes it)
fuji run examples/hello.fuji

# Same pipeline, optional flags
fuji run --no-opt path/to/big_program.fuji

# Native executable you keep
fuji build mygame.fuji -o mygame.exe

# Debug symbols + unoptimized (easier stepping in a debugger)
fuji build --debug mygame.fuji -o mygame_debug.exe

# Rebuild + rerun when .fuji files change under the entry file’s folder
fuji watch src/main.fuji

# Folder you zip or ship (exe + README + extras)
fuji bundle mygame.fuji -o dist/MyGame

# Check syntax/semantics only (no LLVM)
fuji check mygame.fuji

# Canonical formatting
fuji fmt mygame.fuji
fuji fmt --check .
```

See **`fuji help`** for the full command list (`disasm`, `paths`, `doctor`, `bundle`, …).

---

## Documentation (using Fuji only)

| Document | Contents |
|----------|----------|
| [docs/guides/getting-started.md](docs/guides/getting-started.md) | Beginner onboarding (release binary path only) |
| [language.md](language.md) | Keywords, operators, statements, builtins |
| [docs/user_guide.md](docs/user_guide.md) | Longer guide to writing Fuji |
| [docs/reference.md](docs/reference.md) | Stdlib and builtins |
| [docs/language.md](docs/language.md) | Language specification |
| [docs/distribution.md](docs/distribution.md) | Shipping games and bundles |
| [docs/wrappers.md](docs/wrappers.md) | C/C++ FFI, **`fujiwrap`**, Raylib-style workflows |
| [docs/MASTER_PLAN.md](docs/MASTER_PLAN.md) | Maintainer roadmap (language, GC, tooling) — **not** required to *use* Fuji |

---

## Optional examples

- **`examples/`** — small programs and games (**[examples/games/README.md](examples/games/README.md)**).
- **Brick Breaker (Windows, Raylib path):** `powershell -ExecutionPolicy Bypass -File .\scripts\play-brick-breaker.ps1` (from a full checkout with scripts).

---

## Changelog

**[CHANGELOG.md](CHANGELOG.md)** lists shipped features and fixes by version.
