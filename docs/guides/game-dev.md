# Game development with Fuji and Raylib

Fuji compiles to native code via LLVM. External C libraries like Raylib connect through generated wrapper glue. This guide shows the full path from Raylib source to a running game.

---

## 1. Get Raylib

Download or build Raylib. On Windows with MSYS2/MinGW:

```powershell
# Or build from source (place in temp_raylib/ at project root)
# https://github.com/raysan5/raylib/releases
```

You need `libraylib.a` (or `.lib`) and the `raylib.h` header.

---

## 2. Generate the Fuji wrapper

```powershell
.\bin\fujiwrap.exe `
  -name raylib `
  -headers .\path\to\raylib.h `
  -out .\wrappers\raylib
```

This produces:

```
wrappers/raylib/
  raylib.fuji       <- Fuji declarations
  wrapper.c         <- C ABI glue
  README.md
  api_reference.md
  examples.md
```

---

## 3. Write your game

```fuji
#include "wrappers/raylib/raylib.fuji"

let WIDTH = 800;
let HEIGHT = 600;

InitWindow(WIDTH, HEIGHT, "My Game");
SetTargetFPS(60);

while (!WindowShouldClose()) {
    BeginDrawing();
    ClearBackground(20, 20, 20, 255);
    DrawText("Hello, Fuji!", 320, 280, 20, 255, 255, 255, 255);
    EndDrawing();
}

CloseWindow();
```

---

## 4. Build the native executable

Set environment variables so `fuji build` can find the wrapper glue and Raylib:

```powershell
$env:FUJI_NATIVE_SOURCES = '.\wrappers\raylib\wrapper.c'
$env:FUJI_LINKFLAGS = '-I.\path\to\raylib\src -L.\path\to\raylib\src -lraylib -lopengl32 -lgdi32 -lwinmm'
.\kuji.exe build .\game.fuji -o .\game.exe
```

---

## 5. Bundle for distribution

```powershell
$env:FUJI_NATIVE_SOURCES = '.\wrappers\raylib\wrapper.c'
$env:FUJI_LINKFLAGS = '-I.\path\to\raylib\src -L.\path\to\raylib\src -lraylib -lopengl32 -lgdi32 -lwinmm'
$env:FUJI_BUNDLE_FILES = '.\raylib.dll'
.\kuji.exe bundle .\game.fuji -o .\dist\mygame
```

The `dist\mygame` folder contains everything a player needs.

---

## Example: Raylib brick breaker

The canonical example is at `examples/games/raylib_brick_breaker.fuji`. Build it:

```powershell
$env:FUJI_NATIVE_SOURCES = '.\wrappers\raylib_min\raylib_bridge.c'
$env:FUJI_LINKFLAGS = '-I.\temp_raylib\src -L.\temp_raylib\src -lraylib -lopengl32 -lgdi32 -lwinmm'
.\kuji.exe bundle .\examples\games\raylib_brick_breaker.fuji -o .\dist\brick_breaker
```

---

## Raylib wrapper ABI notes

Each Raylib function is wrapped by a C function with the Fuji native ABI:

```c
FujiValue FUJI_wrap_raylib_InitWindow(int argCount, FujiValue* args);
```

The `.fuji` file maps Fuji names to these symbols via `// fuji:extern` directives. The LLVM emitter produces direct native calls — no interpreter overhead.

---

## Performance notes

- Fuji compiles to LLVM IR; Clang generates fully optimized native code.
- The game loop runs at native speed with no GC pauses in the hot path.
- Value boxing/unboxing (NaN-boxing) is the main overhead; it is within 10% of hand-written C for typical game loops.
