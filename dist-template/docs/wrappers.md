# Distributing and using C/C++ library wrappers (Raylib, etc.)

Fuji stays on a **single native pipeline**: the Go frontend lowers programs to **LLVM IR**, then **`fuji build`** runs **llc** on that IR and links the resulting object with **`runtime/libfuji_runtime.a`** (C sources under `runtime/src/`) using **Clang** as the driver, plus optional **`FUJI_NATIVE_SOURCES`** / **`FUJI_LINKFLAGS`** for your C glue. Wrappers are **ordinary `.fuji` modules** (plus your own glue) that you ship alongside—or inside—a **wrapper root directory** so `import "raylib"` (or `@raylib/core`) resolves without vendoring copies into every app.

## 1. Resolver: where imports look

For each `import "name"` / `#include "file"` string, the loader tries, in order:

1. Next to the importing file  
2. Each directory in **`FUJI_PATH`** (same semantics as a normal search path)  
3. Each directory in **`FUJI_WRAPPERS`** (path list, same separator as `PATH` on your OS)

Under each root it looks for `name`, `name.fuji`, or `name/index.fuji` (and the same for `@segment/...` style modules). Set **`FUJI_WRAPPERS`** to the folder you distribute (e.g. unzip a `wrappers/` bundle next to **`fuji.exe`** and point users at it once).

Example (PowerShell):

```powershell
$env:FUJI_WRAPPERS = "C:\tools\fuji-wrappers"
fuji run .\game.fuji
```

Example (Unix):

```bash
export FUJI_WRAPPERS="$HOME/fuji-wrappers"
fuji run ./game.fuji
```

## 2. Native link: actually link the C/C++ library

Generated or hand-written Fuji bindings call into C symbols. **`fuji build`** must pass the right flags to **clang** so the final executable links Raylib (or SQLite, etc.).

Set **`FUJI_LINKFLAGS`** to extra tokens (split on whitespace) inserted **before** `-o`:

```bash
export FUJI_LINKFLAGS="-L/opt/homebrew/lib -lraylib -framework CoreVideo -framework IOKit -framework Cocoa -framework GLUT -framework OpenGL"
fuji build game.fuji -o game
```

On Windows you might use `-L` to a MinGW lib folder and `-lraylib` (exact flags depend on how Raylib was built). Adjust per library and platform.

## 3. Producing wrappers

- **`cmd/wrapgen`** — build the binary as **`fujiwrap`** (`go build -o fujiwrap ./cmd/wrapgen`). Header-driven generator in the **same Go module** as `fuji`. It emits readable `.fuji`, `wrapper.c`, and Markdown — no Python, no extra runtime. The legacy name **`wrapgen`** is the same program.  
- **`fuji wrap …`** — forwards to **`fujiwrap`** (or **`wrapgen`** / **`kujiwrap`**) next to **`fuji`** or on **`PATH`**.  
- Optional release tooling may embed or ship a **`wrappers/`** tree and document **`FUJI_WRAPPERS`** for turnkey bundles.

Your distribution can be:

| Artifact | Role |
|----------|------|
| `fuji` | Compiler + `fuji build` / `fuji run` |
| `wrappers/` directory (zip) | Pre-built or generated `*.fuji` trees per library |
| README in `wrappers/` | Per-library **FUJI_LINKFLAGS** hints |

## 4. Raylib specifically

Official SDK zips ship **`wrappers/raylib/`** (full `.fuji` + `wrapper.c` + Markdown/HTML) and **raylib 5.0 prebuilds** under **`third_party/raylib_static/stage/`** (headers + `libraylib.a` + Windows `raylib.dll`). The compiler picks up that tree automatically when linking (see **`FUJI_USE_VENDORED_RAYLIB`** and **`FUJI_RAYLIB_STAGE`** in `fuji -help`).

1. In source: `#include "raylib/raylib.fuji"` (resolved against **`wrappers/`** next to **`fuji`** — see [guides/raylib.md](guides/raylib.md)).  
2. Set **`FUJI_NATIVE_SOURCES`** to **`wrappers/raylib/wrapper.c`**.  
3. On Windows, copy **`raylib.dll`** next to the `.exe` when using a dynamic link or bundle it with **`FUJI_BUNDLE_FILES`**.  
4. If you do not use the vendored stage (e.g. Linux ARM64 without an upstream binary), set **`FUJI_LINKFLAGS`** to your system `-I` / `-L` / `-lraylib` (and frameworks on macOS).

Until every symbol is lowered automatically, some declarations may still need to match whatever ABI your wrapper emits (same as any FFI layer).

## 5. Checklist for “easy for users”

- [ ] Ship **`wrappers.zip`** with one top-level folder; user sets **`FUJI_WRAPPERS`** to that folder’s absolute path.  
- [ ] Document **`FUJI_LINKFLAGS`** per OS for each heavy library you support.  
- [ ] Optionally ship a **small launcher script** that sets both env vars and runs `fuji run` / `fuji build`.

This keeps **one** language and **one** LLVM link step—wrappers are data + Fuji source on the search path, not a second compiler fork.
