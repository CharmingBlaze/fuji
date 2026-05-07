# Using the Fuji language (complete beginner guide)

This page explains **how to write and run Fuji programs** and what the language offers end-to-end. For a compact **syntax and builtin catalog**, keep **[language.md](../language.md)** open beside this guide. For **CLI commands** (`fuji run`, `fuji build`, …), see **[commands.md](commands.md)**.

---

## 1. How you run Fuji code

You always work with **`.fuji`** text files and the **`fuji`** command-line tool (from [Releases](https://github.com/CharmingBlaze/fuji/releases)):

| Goal | Command |
|------|---------|
| Run once | `fuji run main.fuji` |
| Keep an `.exe` / binary | `fuji build main.fuji -o game.exe` |
| Check mistakes without compiling | `fuji check main.fuji` |
| Format code | `fuji fmt main.fuji` or `fuji fmt ./...` |
| Live reload while editing | `fuji watch main.fuji` |
| Ship a folder | `fuji bundle main.fuji -o dist/MyGame` |

If your program defines **`func main()`**, that function is the usual entry point. You can also use **top-level statements** (code not inside any `func`) like a script.

---

## 2. Comments and program shape

```fuji
// Line comment

/* Block comment */

print("hi");
```

Statements end with **`;`**. Identifiers and keywords are **case-insensitive** (Fuji normalizes spelling).

---

## 3. Values and literals

| Kind | Examples |
|------|------------|
| Number | `42`, `3.14`, `1e2`, `0xff`, `0b1010` |
| String | `"hello\n"` with escapes |
| Boolean | `true`, `false` |
| Null | `null` |
| Array | `[1, 2, ...a, 3]` — **`...expr`** spreads an array into the literal |
| Object | `{ x: 1, y: 2 }` — keys are usually strings; **method shorthand** `draw() { … }` uses **`this`** |

**Template strings** (backticks) embed expressions:

```fuji
let name = "Ada";
print(`Hello, ${name}!`);
```

---

## 4. Variables and assignment

```fuji
let x = 10;
let y;        // starts as null
x = x + 1;

let { a, b } = point;   // object destructuring (binds fields a, b)
```

Compound assignments work: **`+=`**, **`-=`**, **`*=`**, etc.  
Increment/decrement: **`++`** and **`--`** (prefix and postfix where grammar allows).

---

## 5. Operators (what you can write)

- **Arithmetic:** `+` `-` `*` `/` `%` `**` (power), unary `+` `-`
- **Bitwise:** `&` `|` `^` `~` `<<` `>>` `>>>` and compound forms
- **Compare:** `<` `<=` `>` `>=` `==` `!=` `===` `!==`
- **Logic:** `&&` `||` `!`
- **Nullish:** `a ?? b` — right side only when **`a`** is **`null`** (not for `0` or `""`)
- **Optional chaining:** `obj?.field`, `obj?.[i]` — short-circuits to **`null`** if receiver is **`null`**
- **Delete:** `delete obj.key;` removes an own property from a table object

Precedence is “normal JS-like”; when in doubt, use parentheses.

---

## 6. Control flow

All **`if`**, **`else`**, loops, and **`switch`** bodies use **braces** `{ … }` (no single-line braceless branches).

| Construct | Meaning |
|-----------|---------|
| `if (cond) { … } else { … }` | Conditional |
| `while (cond) { … }` | Loop |
| `do { … } while (cond);` | Loop runs at least once |
| `for (init; cond; step) { … }` | Classic C-style `for` |
| `for (let k in arr) { … }` | Keys: array indices, object keys in order |
| `for (let v of arr) { … }` | Values in order |
| `for (let [k, v] of arr) { … }` | Pairs (arrays: index+element; tables: key+value) |
| `switch (x) { case 1: … default: … }` | Switch; use **`break`** to exit a branch |
| `return expr;` | From a function |
| `break` / `continue` | Loops / `switch` |

**Expression forms:**

```fuji
let n = if (x > 0) { 1 } else { -1 };
let label = switch (kind) { case 1 => "a" default => "z" };
```

**`defer`** — run a call when the enclosing **`func`** exits (LIFO if multiple defers).

---

## 7. Functions and closures

```fuji
func add(a, b) {
    return a + b;
}

let mul = func(a, b) { return a * b; };

func greet(name = "world") {
    print("Hello, " + name);
}
```

- **Rest:** `func sum(...xs) { … }` collects extra arguments into an array.
- **Closures** capture outer `let` bindings.

**`this`** — inside an object method, **`this`** is the receiver for **`obj.method()`** calls.

---

## 8. Arrays and objects in practice

**Length and indexing**

```fuji
let a = [10, 20, 30];
print(len(a));      // 3
print(a[1]);        // 20
a[1] = 99;
```

**Growing arrays**

```fuji
let items = [];
items.push(1);
items.push(2);
print(items.length());   // method style
```

**Objects**

```fuji
let player = {
    x: 0,
    y: 0,
    move(dx, dy) {
        this.x = this.x + dx;
        this.y = this.y + dy;
    }
};
player.move(1, 2);
```

---

## 9. String and array methods (high level)

Names are **case-insensitive** (e.g. **`toUpper`** / **`toupper`**).

**Strings (examples):** `split`, `trim`, `toUpper`, `toLower`, `replace`, `replaceAll`, `indexOf`, `slice`, `startsWith`, `endsWith`, …

**Arrays (examples):** `push`, `pop`, `concat`, `join`, `map`, `filter`, `reduce`, `find`, `slice`, `sort`, `reverse`, `indexOf`, `includes`, …

**Ambiguous receivers:** some names (e.g. **`slice`**, **`indexOf`**, **`includes`**) work on both strings and arrays; Fuji picks based on the value type.

For a **full method table**, see **[language.md](../language.md)** sections on string and array methods.

---

## 10. Math: globals and `math.*`

Scalar math exists as **globals** (`sin`, `cos`, `lerp`, `clamp`, `floor`, …).

The compiler also provides a **`math`** object (unless you shadow it) so you can write:

```fuji
let y = math.floor(3.9);
let z = math.lerp(0, 100, 0.5);
```

Use **[language.md](../language.md)** for the full list of math-related globals and namespace members.

---

## 11. Built-in library (what ships in `fuji`)

These are **ordinary Fuji files** under the repo’s **`stdlib/`** folder. In your project, pull them in with **`#include`** (relative path) or **`@module`** resolution if your install ships **`stdlib`** next to **`fuji`**.

Examples of modules you may use:

| Module | Role |
|--------|------|
| `stdlib/math.fuji` | Re-exports / aliases around the `math` prelude |
| `stdlib/vec2.fuji`, `stdlib/vec3.fuji` | Small vector helpers |
| `stdlib/timer.fuji` | Timing helpers |
| `stdlib/json.fuji`, `stdlib/io.fuji`, … | JSON, I/O, strings, arrays |

Typical pattern:

```fuji
#include "../stdlib/vec3.fuji"

func main() {
    let v = create(1, 2, 3);
    print(length(v));
}
```

Adjust the **`#include`** path so it resolves from your file’s directory.

---

## 12. Modules: `#include`, `@`, and `import()`

| Mechanism | Use |
|-----------|-----|
| **`#include "path.fuji"`** | Textual include; good for local libraries and stdlib snippets |
| **`import("@math")`** / `@array` | Resolved via **`FUJI_PATH`**, install **`stdlib`**, or project layout |
| **`import("relative.fuji")`** | Expression import (see **[language.md](../language.md)** and tests for edge cases) |

For C/C++ libraries, use **`fuji wrap`** and **`FUJI_NATIVE_SOURCES`** — **[wrappers.md](wrappers.md)**.

---

## 13. Built-in functions you call every day

Groupings (names are case-insensitive when called):

- **Output / debug:** `print`, `trace`, `assert`, `panic`
- **Values:** `len`, `type`, `typeof`, `number`, `string`, `bool`, `parseJSON`, `toJSON`
- **Files:** `readFile`, `writeFile`, `appendFile`, `fileExists`, `deleteFile`
- **Time / random:** `time`, `clock`, `sleep`, `random`, `randomInt`, …
- **GC (games):** `gc`, `gcCollect`, `gcDisable`, `gcEnable`, `gcFrameStep`, `gcStats`, …

The authoritative alphabetical list is in **[language.md](../language.md)** under native builtins.

---

## 14. FFI hint (advanced)

You can mark native symbols with a special comment (see **[language.md](../language.md)** “Native FFI hint”). Most games use **`fuji wrap`** instead.

---

## 15. Where to read next

| Document | Use it for |
|----------|------------|
| **[language.md](../language.md)** (repo root) | Every keyword, operator, builtin, method — single long catalog |
| **[language.md](language.md)** (this folder) | Formal syntax + semantics reference |
| **[docs/commands.md](commands.md)** | Every `fuji` CLI command |
| **[docs/wrappers.md](wrappers.md)** | Headers → `.fuji` + `wrapper.c` |
| **[docs/distribution.md](distribution.md)** | Shipping bundles and extras |
| **`tests/*.fuji`** | Small runnable examples when behavior is unclear |

If two docs disagree, **`fuji check`** and **`fuji run`** on a tiny program are the final word.
