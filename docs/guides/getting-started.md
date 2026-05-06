# Getting started with Fuji (beginner path)

This guide is only about **using Fuji**.

## 1) Install Fuji the easy way

Download the latest `fuji` binary from [GitHub Releases](https://github.com/CharmingBlaze/fuji/releases), then run:

```bash
fuji version
```

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

- `language.md` — complete language catalog (operators, statements, builtins, methods)
- `docs/user_guide.md` — practical beginner/intermediate walkthrough
- `docs/reference.md` — builtins and runtime-facing APIs
- `docs/guides/game-dev.md` — game-focused usage patterns
- `docs/distribution.md` — shipping and bundles

If docs and behavior ever differ, run a tiny `.fuji` example and trust the CLI result.
