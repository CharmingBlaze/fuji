# Fuji CLI — every command

Use the **`fuji`** binary from [GitHub Releases](https://github.com/CharmingBlaze/fuji/releases). For full inline help:

```bash
fuji help
```

---

## `fuji run` / `fuji native`

Compile the entry **`.fuji`** file to a native executable (temporary), run it, then delete the temp binary.

```bash
fuji run [--no-opt] <file.fuji>
fuji native [--no-opt] <file.fuji>   # same as run
```

- **`--no-opt`** — skips LLVM IR optimisation and uses a less aggressive native compile. Use if a very large program hits a flaky Clang optimiser on your machine.

**Example**

```bash
fuji run src/main.fuji
fuji run --no-opt src/main.fuji
```

---

## `fuji watch`

Rebuild and rerun whenever **`.fuji`** files under the entry file’s directory change.

```bash
fuji watch [--no-opt] <file.fuji>
```

**Example**

```bash
fuji watch src/main.fuji
```

---

## `fuji check`

Parse, resolve imports, and run semantic analysis only — **no** native compile.

```bash
fuji check <file.fuji>
```

Prints **`OK`** if the program is valid.

**Example**

```bash
fuji check src/main.fuji
```

---

## `fuji fmt`

Format **`.fuji`** sources with the canonical formatter.

```bash
fuji fmt [--check] <file.fuji> [more files...]
fuji fmt [--check] ./...
```

- **`--check`** — do not write files; exit with an error if any file would change (for CI).

**Examples**

```bash
fuji fmt src/main.fuji
fuji fmt ./...
fuji fmt --check .
```

---

## `fuji disasm`

Print **LLVM IR** for the program (after parse, sema, and codegen). Useful for debugging the compiler pipeline, not for everyday game dev.

```bash
fuji disasm <file.fuji>
```

---

## `fuji build`

Produce a **native executable** you keep.

```bash
fuji build [--no-opt] [--debug] <file.fuji> [-o <exe>]
```

- **`--no-opt`** — same meaning as for `run`.
- **`--debug`** — emit debug symbols and favour debuggable builds (implies **`--no-opt`** for the native step).
- **`-o`** — output path; if omitted, a default name next to the source is used.

**Examples**

```bash
fuji build game.fuji -o dist/game.exe
fuji build --debug game.fuji -o dist/game_debug.exe
```

---

## `fuji bundle`

Build the program and write a **folder** you can zip or ship (executable, helper scripts, README, etc.).

```bash
fuji bundle <file.fuji> [-o <dir>]
```

Default output directory is **`dist`** if you omit **`-o`**.

**Example**

```bash
fuji bundle game.fuji -o dist/MyGame
```

Copy extra files (DLLs, assets) with the **`FUJI_BUNDLE_FILES`** environment variable (space-separated paths). See **`fuji help`** for the exact format on your OS.

---

## `fuji wrap`

Forward arguments to **`fujiwrap`** (C/C++ header → **`.fuji`** bindings + **`wrapper.c`**). **`fuji`** looks for **`fujiwrap`** next to itself, then **`wrapgen`**, **`kujiwrap`**, then **`PATH`**.

```bash
fuji wrap --help
fuji wrap -name mylib -headers ./include/mylib.h -out ./wrappers/mylib
```

Full workflow: **[wrappers.md](wrappers.md)**.

---

## `fuji paths`

Print **machine-readable** toolchain paths (for scripts and CI).

```bash
fuji paths
```

---

## `fuji doctor`

Human-readable **health check**: clang resolution, runtime library, stdlib, install directory, etc.

```bash
fuji doctor
```

---

## `fuji version`

Print version and platform.

```bash
fuji version
fuji --version
```

---

## `fuji help`

Print the full help screen (commands, environment variables, examples).

```bash
fuji help
fuji --help
```

---

## Environment variables (quick reference)

| Variable | Purpose |
|----------|---------|
| **`FUJI_CLANG`** / **`CC`** | Clang driver for native builds (if not using embedded release toolchain) |
| **`FUJI_PATH`** | Extra **`@module`** search directories |
| **`FUJI_WRAPPERS`** | Pre-built wrapper **`.fuji`** trees |
| **`FUJI_NATIVE_SOURCES`** | C/C++ sources linked into your app (e.g. **`wrapper.c`**) |
| **`FUJI_LINKFLAGS`** | Extra linker flags (**`-l`**, **`-L`**, frameworks, …) |
| **`FUJI_BUNDLE_FILES`** | Extra files copied into **`fuji bundle`** output |

See **`fuji help`** for the complete list and notes.
