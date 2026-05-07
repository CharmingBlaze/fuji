# Fuji — full language reference (how to use everything)

This document is the **complete user-facing language reference**: what you can write, with **runnable-style examples** beside each feature. For CLI usage (`fuji run`, `fuji build`, …), see **`docs/commands.md`**. For a gentler walkthrough, see **`docs/using-the-language.md`**. For C/C++ interop, see **`docs/wrappers.md`**.

**Run a fragment:** save as `main.fuji`, then `fuji run main.fuji` (or define **`func main()`** as below). When examples show only declarations, wrap them in **`func main() { … }`** if your runner expects an entrypoint.

Identifiers and keywords are **case-insensitive** (normalized internally). Builtin names are **case-insensitive** at call sites.

If anything here disagrees with the compiler, trust **`fuji check`**, **`fuji run`**, and the sources **`internal/parser`**, **`internal/lexer/token.go`**, and **`internal/codegen/builtin_register.go`**.

---

## Table of contents

1. [Program structure](#1-program-structure)
2. [Comments and statements](#2-comments-and-statements)
3. [Identifiers and reserved words](#3-identifiers-and-reserved-words)
4. [Literals](#4-literals)
5. [Variables and destructuring](#5-variables-and-destructuring)
6. [Operators](#6-operators)
7. [Control flow statements](#7-control-flow-statements)
8. [Expression forms (`if`, `switch`, chaining)](#8-expression-forms-if-switch-chaining)
9. [Functions and `this`](#9-functions-and-this)
10. [Struct types](#10-struct-types)
11. [Enum types](#11-enum-types)
12. [Objects, arrays, and iteration](#12-objects-arrays-and-iteration)
13. [`#include` and `import`](#13-include-and-import)
14. [Native FFI hint](#14-native-ffi-hint)
15. [Global builtins](#15-global-builtins)
16. [String methods](#16-string-methods)
17. [Array methods](#17-array-methods)
18. [The `math` object](#18-the-math-object)
19. [Truthy and precedence](#19-truthy-and-precedence)
20. [Keyword index](#20-keyword-index)
21. [Canonical builtin names](#21-canonical-builtin-names)

---

## 1. Program structure

A program is a sequence of top-level declarations and statements. The usual entry for **`fuji run`** / **`fuji build`** is **`func main()`** (name **`main`** is matched case-insensitively). Top-level expression statements are also allowed (script style).

```fuji
func main() {
    print("hello");
}
```

```fuji
// script style: statements at top level
print("hi");
```

---

## 2. Comments and statements

```fuji
// line comment to end of line

/*
   block comment
*/

let x = 1;
print(x);
```

Every statement ends with **`;`**. Whitespace separates tokens.

---

## 3. Identifiers and reserved words

Use **`let`** for variables. **`var`** is reserved and rejected — use **`let`** instead. **`this`** is a keyword (object methods).

```fuji
let PlayerScore = 100;   // same as playerscore after normalization
let _tmp = 0;
```

---

## 4. Literals

```fuji
let n = 42;
let pi = 3.14;
let sci = 1e2;
let hex = 0xff;
let bits = 0b1010;

let s = "hello\n\t\"quoted\"";
let u = "pi: \u{03c0}";

let flag = true;
let empty = null;

let arr = [1, 2, 3];
let spread = [0, ...arr, 4];

let obj = {
    x: 1,
    y: 2,
    greet() {
        print("hi from this");
        print(this.x);
    }
};

let tpl = `count = ${ len(arr) }`;
```

Arrays: **`...expr`** inside **`[ ... ]`** splices another **array** into the literal. Objects support **shorthand methods** with **`this`**. Template strings use backticks and **`${ expr }`**.

---

## 5. Variables and destructuring

```fuji
let a = 1;
let b;          // null

a = a + 1;
a += 2;

let point = { x: 10, y: 20 };
let { x, y } = point;
```

Object destructuring binds properties by name from the right-hand value.

---

## 6. Operators

**Arithmetic:** `+` `-` `*` `/` `%` `**` (power), unary `+` `-`.

**Bitwise:** `&` `|` `^` `~` `<<` `>>` `>>>` and compounds `&=` `|=` `^=` `<<=` `>>=`.

**Compare:** `<` `<=` `>` `>=` `==` `!=` `===` `!==`.

**Logic:** `&&` `||` `!`.

**Assignment:** `=` `+=` `-=` `*=` `/=` `%=` and the bitwise compounds above. Also **`??=`** (nullish assign) where the grammar allows.

**Increment:** `++` `--` prefix and postfix where permitted.

**Nullish coalescing:** **`a ?? b`** — **`b`** is used only when **`a`** is **`null`** (not when **`a`** is **`0`** or **`""`**).

**Optional chaining:** **`obj?.field`**, **`obj?.[index]`** — if the receiver is **`null`**, the result is **`null`** without touching the property.

**`typeof`:** unary **`typeof x`** — same kind of answer as **`type(x)`** (string name of the type).

```fuji
let n = null;
let x = n ?? 10;

let o = null;
let v = o?.x;

let t = typeof 3;
```

**Range (for iteration):** **`low..high`** in contexts like **`for (let i of lo..hi)`**.

**Spread / rest:** **`...name`** in function parameters; **`...arr`** in array literals (see §4, §9).

---

## 7. Control flow statements

Braces are **required** on **`if`**, **`else`**, loops, and **`switch`** bodies (no single-line braceless branches).

**`if` / `else`**

```fuji
if (x > 0) {
    print("positive");
} else {
    print("non-positive");
}
```

**`while` / `do`–`while`**

```fuji
while (n > 0) {
    n = n - 1;
}

do {
    n = n + 1;
} while (n < 3);
```

**`for` (C-style)**

```fuji
for (let i = 0; i < 5; i += 1) {
    print(i);
}

for (;;) {
    break;
}
```

**`for`–`in` (keys)**

```fuji
let obj = { a: 1, b: 2 };
for (let k in obj) {
    print(k);
}
```

**`for`–`of` (values)**

```fuji
let arr = [10, 20];
for (let v of arr) {
    print(v);
}
```

**`for`–`of` with range**

```fuji
let lo = 0;
let hi = 3;
for (let i of lo..hi) {
    print(i);
}
```

**`for`–`of` with pairs**

```fuji
let obj = { a: 1, b: 2 };
for (let [k, v] of obj) {
    print(k, v);
}

let xs = ["x", "y"];
for (let [ik, iv] of xs) {
    print(ik, iv);
}
```

**`switch` (statement)** — cases fall through unless you **`break`**.

```fuji
switch (x) {
    case 1:
        print("one");
        break;
    case 2:
        print("two");
        break;
    default:
        print("other");
}
```

**`return`**, **`break`**, **`continue`**

```fuji
func f() {
    while (true) {
        break;
    }
    return 0;
}
```

**`defer`** — runs when the enclosing **`func`** exits, **last deferred runs first** (LIFO).

```fuji
func work() {
    defer print("third");
    defer print("second");
    defer print("first");
    print("body");
}
```

**`delete`** — removes an **own** property from a table; target must be a property access (not optional chaining).

```fuji
let t = { a: 1, b: 2 };
delete t.a;
```

---

## 8. Expression forms (`if`, `switch`, chaining)

**`if` expression**

```fuji
let sign = if (x > 0) { 1 } else { -1 };
```

**`switch` expression** — arms use **`=>`**

```fuji
let label = switch (kind) {
    case 1 => "a"
    default => "z"
};
```

**Calls, members, index**

```fuji
let f = print;
f("hi");

let o = { x: 1 };
let a = o.x;
let b = o["x"];
let c = arr[0];
```

**`import` as an expression** — **`import("path")`** (see loader / project layout; many projects prefer **`#include`** for local files).

---

## 9. Functions and `this`

**Declaration and function value**

```fuji
func add(a, b) {
    return a + b;
}

let mul = func(a, b) {
    return a * b;
};
```

**Default parameters and rest**

```fuji
func greet(name = "world") {
    print("Hello, " + name);
}

func sum(...xs) {
    let n = xs.length();
    return n;
}
```

**Closures** capture **`let`** bindings from outer scopes.

**`this`** in a method shorthand is the receiver for **`obj.method()`** calls.

```fuji
let player = {
    x: 0,
    move(dx) {
        this.x = this.x + dx;
    }
};
player.move(1);
```

---

## 10. Struct types

Structs give **ordered fields** and **`TypeName { field: value, … }`** construction. Fields are readable and assignable like object properties.

```fuji
struct Point {
    x,
    y
}

let p = Point { x: 3, y: 4 };

assert(p.x == 3, "struct field x");
assert(p.y == 4, "struct field y");

p.y = 10;
assert(p.y == 10, "struct field assign");
```

See also **`tests/struct_test.fuji`**.

---

## 11. Enum types

Enums declare **named members**. Each member is a number **`0`**, **`1`**, **`2`**, … in source order. Use **`EnumName.Member`** anywhere a value is expected (including **`switch`** **`case`**).

```fuji
enum Dir {
    Up,
    Down,
    Left,
    Right
}

let u = Dir.Up;
let r = Dir.Right;

assert(u == 0, "enum Up");
assert(r == 3, "enum Right");

switch (u) {
    case Dir.Up:
        print("case up");
    default:
        print("default arm");
}
```

See also **`tests/enum_test.fuji`**.

---

## 12. Objects, arrays, and iteration

**Length and indexing**

```fuji
let a = [10, 20, 30];
print(len(a));
print(a[1]);
a[1] = 99;
```

**Dynamic array ops** (methods; case-insensitive)

```fuji
let items = [];
items.push(1);
items.pop();
print(items.length());
```

**`len`** applies to strings (length), arrays (element count), and plain objects (**entry count**).

**Iteration** — see §7 (`for–in`, `for–of`, pairs, ranges).

Higher-order **`map`**, **`filter`**, **`find`**, **`reduce`** on arrays are supported as methods (§17).

---

## 13. `#include` and `import`

Textual inclusion merges another `.fuji` file into the program bundle:

```fuji
#include "lib/helpers.fuji"
#include "../stdlib/vec3.fuji"
```

Paths are relative to the including file unless your tool resolves **`@`** modules (install / **`FUJI_PATH`**). Expression form **`import("path")`** exists for modular loading; shipped games often standardize on **`#include`** plus **`@`** resolution.

---

## 14. Native FFI hint

A special comment declares a binding to a native symbol for the **`let`** / **`func`** declaration that follows:

```fuji
// fuji: extern my_sin sin 1

let my_sin;   // arity and symbol wired per comment
```

Details follow the parser rules for **`NativeDirective`**.

---

## 15. Global builtins

Grouped by role. Names are lowercase here; **`DeltaTime`**, **`readfile`**, etc. work the same.

**Errors / results / control**

```fuji
let r = ok(value);
let e = err("msg");
panic("fatal");
assert(x == 1, "expected one");
```

**Printing and formatting**

```fuji
print("a", 2, true);
trace("debug");
let s = format("n={}", 42);
```

**Lengths and types**

```fuji
let n = len(arr);
print(type(x), typeof x);
print(isNumber(x), isArray(x));
```

**Coercion and JSON**

```fuji
let x = number("3.14");
let s = string(42);
let b = bool(1);

let parsed = parseJSON("{}");
let out = toJSON(obj);
```

**Files**

```fuji
let text = readFile("in.txt");
writeFile("out.txt", text);
appendFile("log.txt", "\n");
if (fileExists("x.bin")) {
    deleteFile("x.bin");
}
```

**Time**

```fuji
let t = time();
let c = clock();
let ts = timestamp();
let p = programTime();
sleep(0);

let dt = deltaTime();   // framedelta-style hook when wired
```

**Random**

```fuji
randomSeed(12345);
let u = random();
let i = randomInt(5);           // single arg = [0, max)
let j = randomInt(0, 100);       // two args = [min, max)
```

`randomChoice` is registered; check your runtime build for full array support.

**Math as globals**

```fuji
let y = sin(0);
let z = clamp(x, 0, 1);
let d = distance(0, 0, bx, by);
let m = lerp(0, 100, 0.5);
let hypotenuse = hypot(3, 4);
let w = wrap(angle, low, high);
```

**Substring helper (not full regex)**

```fuji
if (matches(body, needle)) {
    print("found");
}
```

**Garbage collection (native / games)**

```fuji
gc();
gcCollect();
gcDisable();
gcEnable();
gcFrameStep(budgetMicrosOrSimilar);
let st = gcStats();
```

Exact runtime shapes for **`gcStats`** and **`ok`/`err`** follow the argv / native ABI.

---

## 16. String methods

Call as **`receiver.method(...)`**. Names are **case-insensitive** (`toUpper` / `toupper`).

| Method | Example shape |
|--------|----------------|
| **`split(delimiter)`** | `"a,b".split(",")` — empty delim → one char per segment |
| **`trim()`** | `"  hi  ".trim()` |
| **`toUpper()`** / **`toLower()`** | `"Ab".toUpper()` |
| **`replace(a, b)`** | single replacement |
| **`replaceAll(a, b)`** | replace every occurrence |
| **`indexOf(needle)`** | index or **`−1`** |
| **`includes(needle)`** | substring test |
| **`slice(start, end)`** | half-open range |
| **`startsWith`**, **`endsWith`** | prefix / suffix |

```fuji
let s = "  Hello ";
print(s.trim().toLower());
print(s.includes("ell"));
```

**`slice`**, **`indexOf`**, **`includes`** are **overloaded** with arrays—the compiler picks based on receiver type.

---

## 17. Array methods

| Method | Notes |
|--------|--------|
| **`concat(...vals)`** | new array appended |
| **`push(x)`**, **`pop()`** | mutate tail |
| **`length()`** | like **`len(arr)`** |
| **`join(sep)`** | string |
| **`map(cb)`**, **`filter(cb)`**, **`find(cb)`**, **`reduce(cb [, initial])`** | callback forms per runtime |
| **`slice(start, end)`** | copy range |
| **`sort()`**, **`reverse()`** | in place |
| **`indexOf(item)`**, **`includes(item)`** | search |

```fuji
let a = [1, 2, 3];
let b = a.map(func(x) { return x * 2; });

let doubled = [];
for (let item of a) {
    doubled.push(item * 2);
}
```

Higher-order methods are lowered to LLVM loops with indirect calls in the native compiler.

---

## 18. The `math` object

If you do not shadow the name **`math`**, the compiler injects **`let math = { … };`** so namespace calls work (**`Math`** normalizes to **`math`**):

```fuji
let x = math.floor(3.9);
let y = math.lerp(0, 100, 0.5);
let z = math.sqrt(2);
```

Many members route through **`math.*`** argv fast paths (**`floor`**, **`sin`**, **`lerp`**, **`hypot`**, **`wrap`**, …). Others are duplicated onto **`math`** from globals by the prelude (e.g. **`sqrt`**, **`abs`**, **`pi`**, **`e`**, **`random`**—see **`internal/parser/prelude.go`**). Anything not on the **`math`** value remains available as a **global** call (for example **`distance(...)`**, **`smoothstep(...)`**).

---

## 19. Truthy and precedence

**Truthy / falsy (conceptually):** **`false`**, **`null`**, **`0`**, **`""`** are falsy for conditions; most other values are truthy.

**Precedence (high → low, approximate):** call, member, index → unary (`+`, `-`, `!`, `typeof`) → `* / %` → `+ -` → comparisons → **`==`** family → **`&&`** → **`||`** → **`??`** → assignment. Use parentheses when in doubt.

---

## 20. Keyword index

Spelling lowercase; any case accepted. **`var`** is reserved (**use `let`**).

```
break case continue default defer delete do else enum false for func if import in
let null of return struct switch this true typeof while
```

Directive: **`#include "..."`**.

---

## 21. Canonical builtin names

Registered in **`builtin_register.go`** (and **`builtin_globals.go`** for the resolver). Alphabetical:

`abs`, `acos`, `appendFile`, `approach`, `asin`, `assert`, `atan`, `atan2`, `bool`, `ceil`, `clamp`, `clock`, `cos`, `deltaTime`, `deleteFile`, `distance`, `distanceSq`, `e`, `err`, `exp`, `fileExists`, `floor`, `format`, `fmod`, `gc`, `gcCollect`, `gcDisable`, `gcEnable`, `gcFrameStep`, `gcStats`, `hypot`, `isArray`, `isBool`, `isFunction`, `isNull`, `isNumber`, `isObject`, `isString`, `len`, `lerp`, `log`, `log10`, `map`, `matches`, `max`, `min`, `normalize`, `number`, `ok`, `panic`, `parseJSON`, `pi`, `pow`, `print`, `programTime`, `random`, `randomChoice`, `randomInt`, `randomSeed`, `readFile`, `round`, `sign`, `sin`, `sleep`, `smoothdamp`, `smoothstep`, `sqrt`, `string`, `tan`, `timestamp`, `time`, `toJSON`, `trace`, `trunc`, `type`, `typeof`, `wrap`, `writeFile`

---

## See also

| Where | Purpose |
|-------|---------|
| **`docs/commands.md`** | CLI: `run`, `build`, `fmt`, … |
| **`docs/using-the-language.md`** | Beginner-oriented narrative |
| **`docs/wrappers.md`** | **`fuji wrap`** |
| **`docs/language/syntax.md`** | Compact syntax sheet |
| **`tests/*.fuji`** | Executable behavior examples |
