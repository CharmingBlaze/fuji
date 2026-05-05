# Fuji language reference

Fuji is a dynamically-typed scripting language with JavaScript-like syntax. It compiles to native machine code through LLVM IR and the C runtime (`fuji build`, `fuji run`).

---

## Variables

```fuji
let x = 42;
let name = "Player";
let active = true;
let nothing = null;
let score;           // declared, value is null
```

Assignment:

```fuji
x = x + 1;
x += 1;
x -= 1;
x *= 2;
x /= 2;
```

---

## Types

| Type | Examples |
|------|---------|
| Number | `0`, `42`, `3.14`, `-10`, `0xFF` |
| String | `"hello"`, `"line\n"` |
| Boolean | `true`, `false` |
| Null | `null` |
| Array | `[1, 2, 3]` |
| Object | `{x: 10, y: 20}` |
| Function | `func(a, b) { return a + b; }` |

---

## Operators

```fuji
// Arithmetic
a + b    a - b    a * b    a / b    a % b

// Comparison
a == b   a != b   a < b   a > b   a <= b   a >= b

// Logical
a && b   a || b   !a

// Prefix / postfix update
++x   --x   x++   x--

// Compound assignment
x += 1   x -= 1   x *= 2   x /= 2
```

---

## Functions

```fuji
func add(a, b) {
    return a + b;
}

let result = add(3, 4);   // 7
```

Anonymous function expression:

```fuji
let square = func(n) {
    return n * n;
};
```

Closures:

```fuji
func makeCounter() {
    let count = 0;
    return func() {
        count += 1;
        return count;
    };
}

let c = makeCounter();
print(c());   // 1
print(c());   // 2
```

---

## Control flow

```fuji
if (x > 0) {
    print("positive");
} else {
    print("non-positive");
}

while (running) {
    update();
}

for (let i = 0; i < 10; i += 1) {
    print(i);
}
```

---

## For-of loops

Iterate array values:

```fuji
let items = ["sword", "shield", "potion"];

for (let item of items) {
    print(item);
}
```

Iterate with index:

```fuji
for (let i, item of items) {
    print(i, ":", item);
}
```

Iterate object keys and values:

```fuji
let player = {name: "Hero", hp: 100};

for (let key, value of player) {
    print(key, "=", value);
}
```

---

## Arrays

```fuji
let arr = [10, 20, 30];

print(arr[0]);          // 10
print(len(arr));        // 3

arr[1] = 99;
```

---

## Objects

```fuji
let player = {
    name: "Hero",
    hp: 100,
    x: 0,
    y: 0
};

print(player.name);     // Hero
player.hp = 75;
```

---

## Switch

```fuji
switch (state) {
    case "menu":
        renderMenu();
        break;
    case "playing":
        updateGame();
        break;
    default:
        break;
}
```

---

## Modules

```fuji
#include "math.fuji"
#include "wrappers/raylib/raylib.fuji"
```

Resolution order: local path, `FUJI_PATH` directories, `FUJI_WRAPPERS` directories.

---

## Standard library

| Function | Description |
|----------|-------------|
| `print(...)` | Print values, space-separated |
| `type(v)` | Return type name as string |
| `number(v)` | Convert to number |
| `string(v)` | Convert to string |
| `len(v)` | Length of array or string |
| `time()` | Current time in seconds (float) |
| `sleep(ms)` | Sleep for milliseconds |
| `abs(n)` | Absolute value |
| `sqrt(n)` | Square root |
| `random()` | Random float in [0, 1) |
