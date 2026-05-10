# Fuji

Fuji is a **simple C and JavaScript-like language** for making games and applications.

Write familiar syntax — C-style control flow, JavaScript-style objects and closures — and compile to a **single native binary**. No VM. No interpreter. Players run your binary without installing anything.

---

## If you want to make games or apps with Fuji

**Do not compile this repository from source.** That path is for maintainers only and expects **Go** and **LLVM** on the machine that builds Fuji.

**To use Fuji you only need the release downloads:** **`fuji`** and **`fujiwrap`** from [GitHub Releases](https://github.com/CharmingBlaze/fuji/releases), plus the **`stdlib/`** folder from the same **SDK zip** (keep it next to the executables). Those binaries are self-contained — they embed the native compiler pieces they need. **You do not install Go, LLVM, or any other compiler** on the computer where you only write and build `.fuji` programs.

Then read:
- [`language.md`](language.md) — everything you can write in Fuji, with examples
- [`docs/guides/getting-started.md`](docs/guides/getting-started.md) — your first program
- [`docs/commands.md`](docs/commands.md) — every `fuji` CLI command
- [`docs/using-the-language.md`](docs/using-the-language.md) — full language tutorial

---

## What the language looks like

```c
struct Player {
    x, y, speed, health
}

enum State {
    Idle, Running, Dead
}

func update(player, dt) {
    player.x = player.x + player.speed * dt;
    if (player.health <= 0) {
        return State.Dead;
    }
    return State.Running;
}

let p = Player { x: 0, y: 0, speed: 200, health: 100 };
let state = update(p, 0.016);
print(state);
```

---

## Quick start (after downloading the release binary)

```bash
# Run a program directly
fuji run mygame.fuji

# Compile to a binary
fuji build mygame.fuji -o mygame

# Check for errors without compiling
fuji check mygame.fuji

# Format your code
fuji fmt mygame.fuji

# Watch for changes and rebuild automatically
fuji watch mygame.fuji

# Bundle for distribution (binary + assets)
fuji bundle mygame.fuji -o dist/MyGame
```

---

## SDK downloads

Each release includes self-contained SDK zips. **You only add `fuji`, `fujiwrap`, and the shipped folders (e.g. `stdlib/`) — no Go install and no LLVM install** on the machine where you write Fuji code.

| Platform | Download |
|---|---|
| Windows (x64) | `fuji-vX.Y.Z-sdk-windows-amd64.zip` |
| Linux (x64) | `fuji-vX.Y.Z-sdk-linux-amd64.zip` |
| Linux (ARM64) | `fuji-vX.Y.Z-sdk-linux-arm64.zip` |
| macOS Intel | `fuji-vX.Y.Z-sdk-darwin-amd64.zip` |
| macOS Apple Silicon | `fuji-vX.Y.Z-sdk-darwin-arm64.zip` |

Unpack into any folder and run `fuji` from that folder so `stdlib/` sits next to the executable.

---

## Documentation

| File | Contents |
|---|---|
| [`language.md`](language.md) | Compact reference — all syntax, operators, builtins |
| [`docs/guides/getting-started.md`](docs/guides/getting-started.md) | First steps |
| [`docs/using-the-language.md`](docs/using-the-language.md) | Full language tutorial |
| [`docs/commands.md`](docs/commands.md) | Every CLI command |
| [`docs/reference.md`](docs/reference.md) | Stdlib and builtins reference |
| [`docs/distribution.md`](docs/distribution.md) | Shipping games and apps |
| [`docs/wrappers.md`](docs/wrappers.md) | Wrapping C libraries with `fujiwrap` |
| [`CHANGELOG.md`](CHANGELOG.md) | What changed in each version |

---

## For maintainers and contributors

See [`CONTRIBUTING.md`](CONTRIBUTING.md) for how to build the compiler from source, run tests, and contribute changes.

The build requires Go 1.22+, Clang, llc, and a C toolchain. See [`docs/handoff.md`](docs/handoff.md) for the internal architecture.

**Do not tell users to build from source.** Release binaries are the correct install path.

---

## Examples

- [`examples/`](examples/) — sample programs and games
- Run Brick Breaker on Windows: `powershell -ExecutionPolicy Bypass -File scripts/play-brick-breaker.ps1`
