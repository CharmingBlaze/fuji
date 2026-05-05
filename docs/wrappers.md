# Distributing and using C/C++ library wrappers (Raylib, etc.)

Fuji stays on a **single native pipeline**: the Go frontend lowers programs to **LLVM IR**, then **`fuji build`** runs **llc** on that IR and links the resulting object with **`runtime/libfuji_runtime.a`** (C sources under `runtime/src/`) using **Clang** as the driver, plus optional **`FUJI_NATIVE_SOURCES`** / **`FUJI_LINKFLAGS`** for your C glue. Wrappers are **ordinary `.fuji` modules** (plus your own glue) that you ship alongside—or inside—a **wrapper root directory** so `import "raylib"` (or `@raylib/core`) resolves without vendoring copies into every app.

## 1. Resolver: where imports look

For each `import "name"` / `#include "file"` string, the loader tries, in order:

1. Next to the importing file  
2. Each directory in **`FUJI_PATH`** (same semantics as a normal search path)  
3. Each directory in **`FUJI_WRAPPERS`** (path list, same separator as `PATH` on your OS)

Under each root it looks for `name`, `name.fuji`, or `name/index.fuji` (and the same for `@segment/...` style modules). Set **`FUJI_WRAPPERS`** to the folder you distribute (e.g. unzip `kuji-wrappers/` next to `kuji.exe` and point users at it once).

Example (PowerShell):

```powershell
$env:FUJI_WRAPPERS = "C:\tools\kuji-wrappers"
fuji run .\game.fuji
```

Example (Unix):

```bash
export FUJI_WRAPPERS="$HOME/kuji-wrappers"
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

- **`cmd/wrapgen`** (build the binary as **`kujiwrap`**) — header-driven generator in the **same Go module** as `fuji`. Run `go build -o kujiwrap ./cmd/wrapgen`. It emits readable `.fuji`, `wrapper.c`, and Markdown — no Python, no extra runtime.  
- **`kuji wrap …`** — convenience: forwards to `kujiwrap` when it is installed next to `fuji`.  
- **`cmd/kuji-single`** / **`cmd/dist`** — optional tools to embed or extract a **`wrappers/`** tree and set **`FUJI_WRAPPERS`** for turnkey bundles.

Your distribution can be:

| Artifact | Role |
|----------|------|
| `fuji` (or `kuji.exe`) | Compiler + VM + `fuji build` |
| `wrappers/` directory (zip) | Pre-built or generated `*.fuji` trees per library |
| README in `wrappers/` | Per-library **FUJI_LINKFLAGS** hints |

## 4. Raylib specifically

1. Ship `wrappers/raylib/` (or similar) containing the generated **`raylib.fuji`** (or `raylib/index.fuji`).  
2. In the app: `import "raylib"` (or the path you documented).  
3. Build: set **`FUJI_WRAPPERS`** to the parent of `raylib/` and **`FUJI_LINKFLAGS`** so clang finds **libraylib** and required frameworks/libs on that OS.

Until every symbol is lowered automatically, some declarations may still need to match whatever ABI your wrapper emits (same as any FFI layer).

## 5. Checklist for “easy for users”

- [ ] Ship **`wrappers.zip`** with one top-level folder; user sets **`FUJI_WRAPPERS`** to that folder’s absolute path.  
- [ ] Document **`FUJI_LINKFLAGS`** per OS for each heavy library you support.  
- [ ] Optionally ship a **small launcher script** that sets both env vars and runs `fuji run` / `fuji build`.

This keeps **one** language and **one** LLVM link step—wrappers are data + Fuji source on the search path, not a second compiler fork.
