# Getting started with Fuji (beginner path)

This guide is only about **using Fuji**.

## 1) Install Fuji the easy way

On [GitHub Releases](https://github.com/CharmingBlaze/fuji/releases), download **one** file: the **SDK zip** for your OS (for example **`fuji-vX.Y.Z-sdk-windows-amd64.zip`** on Windows 64-bit). Unzip it; inside you get **`fuji`** (or **`fuji.exe`**), **`fujiwrap`**, **`stdlib/`**, **`docs/`**, **`wrappers/`**, **`examples/`**, and (on most platforms) **`third_party/raylib_static/stage/`** with Raylib already vendored.

That is the full, **offline** toolchain: **`fuji` does not download LLVM, Raylib, or anything else** when you compile. Embedded Clang/llc live inside the release binary and are written to a **local temp folder** the first time you build or run — still no network.

Then:

```bash
fuji version
```

(Loose **`fuji-*`** binaries without the zip are also on the release page if you only want the compiler and will arrange **`stdlib/`** yourself.)

Do not compile Fuji from source for normal usage. The source tree is maintainer-focused and not the intended beginner install route.

---

## 2) Your first program

Create `hello.fuji`:

```fuji
print("Hello, Fuji!");
```

Run it:

```bash
fuji run hello.fuji
```

---

## 3) Commands you will use every day

```bash
# Run (temporary executable)
fuji run game.fuji

# Build a native executable
fuji build game.fuji -o game.exe

# Debug-friendly build
fuji build --debug game.fuji -o game_debug.exe

# Check parse + semantic errors
fuji check game.fuji

# Format source
fuji fmt game.fuji
fuji fmt --check .

# Rebuild/rerun when files change
fuji watch game.fuji

# Package for sharing
fuji bundle game.fuji -o dist/MyGame
```

For **every** command, flags, and copy-paste examples, see **[docs/commands.md](../commands.md)** (or run **`fuji help`**).

---

## 4) Using the wrapper tool (`fujiwrap`)

`fujiwrap` generates `.fuji` bindings + `wrapper.c` from C/C++ headers.

```bash
fuji wrap --help
```

Typical flow:

1. Generate bindings from a header.
2. Import generated `.fuji` module in your game.
3. Build/run with native glue via `FUJI_NATIVE_SOURCES` and linker flags via `FUJI_LINKFLAGS`.

Full details: `docs/wrappers.md`.

---

## 5) Learn the whole language

- `docs/using-the-language.md` — **start here**: how to use the language end-to-end (syntax, types, control flow, modules, builtins, stdlib)
- `language.md` — complete language catalog (operators, statements, builtins, methods)
- `docs/user_guide.md` — practical beginner/intermediate walkthrough
- `docs/reference.md` — builtins and runtime-facing APIs
- `docs/guides/game-dev.md` — game-focused usage patterns
- `docs/distribution.md` — shipping and bundles

If docs and behavior ever differ, run a tiny `.fuji` example and trust the CLI result.
