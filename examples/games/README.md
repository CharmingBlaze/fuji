# Example games

## Text-only (no Raylib)

Runs anywhere **`fuji`** runs:

```bash
fuji run examples/games/lunar_lander_text.fuji
```

## Raylib (graphics)

Raylib samples need generated bindings plus linker flags. See **`tests/raylib_wrapgen_smoke.fuji`** and **`scripts/ci-wrapgen-raylib.sh`** for the usual **`FUJI_WRAPPERS`**, **`FUJI_NATIVE_SOURCES`**, and **`FUJI_LINKFLAGS`** setup.

After wrappers are built:

```bash
export FUJI_WRAPPERS=/path/to/wrappers/raylib
export FUJI_NATIVE_SOURCES=/path/to/wrappers/raylib/wrapper.c
export FUJI_LINKFLAGS="$(pkg-config --libs --cflags raylib)"   # platform-specific
fuji run demos/demo_3d.fuji    # or another raylib-backed demo under demos/
```
