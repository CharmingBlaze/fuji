# Game Development with Fuji

This is a quick-start overview. For a complete, beginner-friendly guide to Raylib functions, colors, key codes, and a full working game, see **[raylib.md](raylib.md)**.

---

## The 3 things you need

1. **A `.fuji` source file** — your game code
2. **The wrapper** — `wrappers/raylib_shim/raylib.fuji` (already in the repo)
3. **The C bridge** — `wrappers/raylib_min/raylib_bridge.c` (already in the repo)

---

## Quickstart

**Set the bridge once per terminal session:**

```powershell
# Windows
$env:FUJI_NATIVE_SOURCES = "wrappers\raylib_min\raylib_bridge.c"
```
```bash
# Linux / macOS
export FUJI_NATIVE_SOURCES=wrappers/raylib_min/raylib_bridge.c
```

**Write your game (`game.fuji`):**

```fuji
#include "wrappers/raylib_shim/raylib.fuji"

func main() {
    InitWindow(800, 600, "My Game");
    SetTargetFPS(60);

    while (!WindowShouldClose()) {
        BeginDrawing();
        ClearBackground(0x181818FF);
        DrawText("Hello, Fuji!", 300, 280, 30, 0xFFFFFFFF);
        EndDrawing();
    }

    CloseWindow();
}
```

**Run it:**

```
fuji run game.fuji
```

**Build a binary:**

```
fuji build game.fuji -o game.exe
```

**Ship it:**

```powershell
$env:FUJI_BUNDLE_FILES = "raylib.dll"
fuji bundle game.fuji -o dist\mygame
```

---

## Game loop structure

Every Fuji + Raylib game follows the same pattern:

```fuji
#include "wrappers/raylib_shim/raylib.fuji"

func main() {
    // 1. Initialize
    InitWindow(800, 600, "Game");
    SetTargetFPS(60);

    // 2. Game loop
    while (!WindowShouldClose()) {

        // 3. Update (input, physics, logic)
        // ...

        // 4. Draw
        BeginDrawing();
        ClearBackground(0x000000FF);
        // ... draw calls ...
        EndDrawing();
    }

    // 5. Cleanup
    CloseWindow();
}
```

---

## Where to go next

| Document | What it covers |
|----------|----------------|
| **[raylib.md](raylib.md)** | All functions, colors, key codes, complete Pong game example |
| **`examples/games/brick_breaker.fuji`** | Full brick breaker game you can run and study |
| **[../wrappers.md](../wrappers.md)** | How the wrapper system works, extending it with more C functions |
| **[../commands.md](../commands.md)** | All `fuji` CLI commands |
| **[../../language.md](../../language.md)** | Full language reference |
