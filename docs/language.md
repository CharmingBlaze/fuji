# Fuji Language Reference

Fuji is a small, JavaScript-like scripting language. This document is the **authoritative language reference** for programmers. The reference implementation uses a **Go** frontend (lexer, parser, AST, analyses) and an **LLVM** backend (`fuji build` / `fuji run`) linking the C runtime — **one execution path**, native binaries only.

## Compilation pipeline (LLVM)

```
Source (.fuji)
  → Lexer → Parser → AST
  → Native bundle prep (captures, symbols)
  → LLVM IR (llir)
  → clang → object → executable
```

---

## 1. Variable declarations

### Syntax

```javascript
let identifier = expression;
let identifier;  // initializes to null
```

### Rules

- Block scoped (like JavaScript `let`).
- Mutable; no hoisting—must be declared before use in its scope.

---

## 2. Data types

### Primitives

- **Number** — IEEE-754 double. Decimal, scientific notation, and **hexadecimal** (`0xFF`) literals are supported.
- **String** — UTF-8 source; runtime strings are byte-oriented with escapes (`\n`, `\t`, `\u{...}`, etc.).
- **Boolean** — `true`, `false`.
- **Null** — `null`.

### Composite types

- **Array** — `[1, 2, 3]`, nested arrays, mixed elements.
- **Object** — `{ key: value, ... }`; **method shorthand** `name() { ... }` with `this` is supported.

---

## 3. Operators

### Arithmetic

`+`, `-`, `*`, `/`, `%`, unary `-`, unary `+`, `**` (power).

### Bitwise (extension vs. minimal JS subset)

`&`, `|`, `^`, `~`, `<<`, `>>`, `>>>` and compound forms (`&=`, …).

### Comparison and equality

`<`, `<=`, `>`, `>=`, `==`, `!=`, `===`, `!==`.

### Logical

`&&`, `||`, unary `!` (implemented with truthiness helpers in the runtime).

### Assignment

`=`, `+=`, `-=`, `*=`, `/=`, and other compound assignment tokens the lexer defines.

### Increment / decrement

`++` and `--` as prefix and postfix on variables.

---

## 4. Control flow

### `if` / `else if` / `else`

As in C/JavaScript; `else` binds to the nearest `if`.

### `switch` / `case` / `default`

Switch on a value; `break` exits the switch. Fall-through is possible if you omit `break`.

### `while` and `do` … `while`

Standard forms.

### `for` (C-style)

```javascript
for (let i = 0; i < 10; i = i + 1) { ... }
for (; condition; ) { ... }  // empty initializer / increment allowed
```

### `for` … `in` / `for` … `of`

The parser accepts **`in` or `of`** after the loop variable (they are equivalent in Fuji today):

```javascript
for (let item of items) { print(item); }
for (let ch in "hello") { print(ch); }  // iterates code units as one-char strings (LLVM runtime)
```

**Single binding only** — `for (let index, item of arr)` (two variables) is *not* implemented yet.

**Semantics**

- **Array** — loop variable is each **element** (tree runtime and LLVM/C runtime agree).
- **Object** — loop variable is each **key** (string). For values, index the object or use a helper pattern.
- **String** — under the LLVM runtime, iteration yields each character as a length-1 string.

### `break` / `continue`

Work in `while`, `do`/`while`, `for`, and `for`-`in`/`of` loops.

---

## 5. Functions

### Declaration

```javascript
func name(a, b) {
  return a + b;
}
```

Default arguments (compile-time literals) and rest parameters are supported in the compiler; the LLVM path lowers them in the callee prologue.

### Calls

`f(a, b)`; first-class closures.

### Function expressions

```javascript
let add = func(a, b) { return a + b; };
```

Closures capture `let` bindings with upvalues (LLVM path).

### `return`

`return;` returns null; `return expr;` returns a value.

---

## 6. Closures

Lexical closures over `let` bindings are supported (see tests under `tests/`).

---

## 7. Objects

Literals, dot and bracket access, dynamic `obj[prop] = v`, and method shorthand with **`this`** bound on `obj.method()` calls (bound method objects in the C runtime).

---

## 8. Arrays

Literals, `arr[i]` read/write, `len(arr)`.

### Instance methods (value-oriented)

These resolve through the runtime (`arr["push"]` style in lowered form; in source you write method calls):

| Method     | Meaning |
|-----------|---------|
| `push(x)` | Append `x` |
| `pop()`   | Remove and return last element, or null |
| `length()`| Return element count as a number |

Example:

```javascript
let a = [1, 2];
a.push(3);
print(a.length());  // 3
print(a.pop());     // 3
```

---

## 9. Modules

- **`#include "relative.fuji"`** — textual include, single load; used in the bundle loader.
- **`import("path")`** — expression form for module object / init (see implementation and tests).

---

## 10. Comments

`//` line comments and `/* ... */` block comments.

---

## 11. Standard library (global functions)

| Function | Notes |
|----------|--------|
| `print(a, b, …)` | stdout, space-separated, newline at end |
| `type(x)` | `"number"`, `"string"`, `"bool"`, `"null"`, `"array"`, `"object"`, `"function"` (LLVM matches these strings) |
| `number(x)` | Parse string to number; pass-through number; else null |
| `string(x)` | Coerce to string (LLVM: `FUJI_to_string_val`; richer formatting may evolve) |
| `len(x)` | Array, string, or object entry count |
| `abs`, `sqrt`, `random` | Math helpers |
| `clock`, `time`, `sleep(ms)` | time / sleep |
| `keys(obj)` | Object key array (used by bytecode lowering of some loops) |

Additional natives (`gc`, type predicates, JSON/math/io modules) exist for the VM and embedding—see `internal/kuji/native.go` and `internal/runtime/data/kuji.c`.

---

## 12. Reserved words (non-exhaustive)

`let`, `func`, `if`, `else`, `for`, `while`, `do`, `switch`, `case`, `default`, `break`, `continue`, `return`, `true`, `false`, `null`, `this`, `in`, `of`, `import`, …

---

## 13. Precedence (summary)

Highest: call, member, index → unary → `* / %` → `+ -` → comparisons → `==` `!=` `===` `!==` → `&&` → `||` → assignment (see parser sources for exact grammar).

---

## 14. Truthy / falsy

Aligned with the C runtime’s `FUJI_is_truthy`: **`false`**, **`null`**, **`0`**, and **`""`** are falsy; other numbers, non-empty strings, arrays, objects, and functions are truthy.

---

## 15. Grammar (outline)

The parser packages (`parser.go`, `parser_stmt.go`, `parser_expr.go`) are the ground truth. At a high level:

- `program = declaration* EOF`
- `declaration = let | func | include | statement`
- `statement = exprStmt | if | switch | while | doWhile | for | forIn | return | break | continue | block`
- `expression` chains assignment, logical, binary, unary, postfix, call, primary—including `func` expressions, arrays, objects, templates.

---

## 16. Example programs

See `tests/hello.fuji`, `tests/loops.fuji`, `tests/closure.fuji`, `tests/control.fuji`, and others.

---

## 17. Type reference (logical)

| Kind    | Example | Notes |
|---------|---------|--------|
| Number  | `42`, `0x2a` | 64-bit NaN-boxed in LLVM |
| String  | `"hi"` | |
| Bool    | `true` | |
| Null    | `null` | |
| Array   | `[]` | |
| Object  | `{}` | |
| Function| `func () {}` | |

---

## 18. Implementation notes (read this when debugging)

1. **LLVM vs VM** — Some programs are only exercised on one tier; when extending the language, update **lexer/parser**, **native capture pass**, **LLVM emitter**, and **C runtime** together where applicable.
2. **`for`-`in` iterator decl** — Capture analysis stores a stable `*LetDecl` per `ForInStmt` (`NativeEmitContext.forInDecl`) so LLVM stack slots line up with `locals` from the capture walk.
3. **Strict equality** — `===` / `!==` are tokenized; lowering uses runtime helpers consistent with `==` for numbers where applicable—see compiler and runtime.
4. **`cmd/dist`** — May not build in all checkouts; prefer `go test ./internal/kuji/... ./cmd/kuji/...` for CI-style verification.

---

*This reference is maintained with the compiler sources. When in doubt, prefer behavior demonstrated by tests under `tests/` and `internal/kuji/*_test.go`.*
