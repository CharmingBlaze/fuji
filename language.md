# Fuji — complete language catalog

**Shape of the language (implemented):** the **JavaScript-like** surface—syntax, objects, arrays, functions, templates, common methods—is lowered by **`internal/parser`** and **`internal/codegen`**. The **C-like** part is the **native** artifact and runtime: LLVM emission, **`runtime/src/`** (NaN-boxed values, GC), **`fuji build`** output—no bytecode VM in the shipping path.

Single-page overview of **everything the Fuji surface syntax and native builtins expose today**. Deep tutorials live in **[docs/user_guide.md](docs/user_guide.md)**; the shorter **[docs/language.md](docs/language.md)** overlaps partially — prefer **tests under `tests/`** when behavior is ambiguous.

**Implementation:** lexer **`internal/lexer/`**, grammar **`internal/parser/`**, LLVM lowering **`internal/codegen/`**, C runtime **`runtime/src/`**.

---

## Table of contents

1. [Program structure](#1-program-structure)  
2. [Comments and whitespace](#2-comments-and-whitespace)  
3. [Identifiers and case](#3-identifiers-and-case)  
4. [Literals](#4-literals)  
5. [Keywords](#5-keywords)  
6. [Operators and punctuation](#6-operators-and-punctuation)  
7. [Types (logical)](#7-types-logical)  
8. [Declarations](#8-declarations)  
9. [Statements](#9-statements)  
10. [Expressions](#10-expressions)  
11. [Functions](#11-functions)  
12. [Objects and arrays](#12-objects-and-arrays)  
13. [Modules and includes](#13-modules-and-includes)  
14. [Native builtins (global calls)](#14-native-builtins-global-calls)  
15. [Truthy / falsy](#15-truthy--falsy)  
16. [Precedence (summary)](#16-precedence-summary)  
17. [String methods](#17-string-methods)  
18. [Array methods](#18-array-methods)  
19. [Templates, math namespace, pattern matching, and operators](#19-templates-math-namespace-pattern-matching-and-operators)  

---

## 1. Program structure

- A **program** is a sequence of **declarations** and **statements** at top level (`parser.Parse`).
- **`func main()`** is the usual entry point for **`fuji run`** / **`fuji build`** when present (entry resolution is case-insensitive for **`main`**).
- Statements at top level are allowed (expression statements, control flow, etc.) depending on loader rules.

---

## 2. Comments and whitespace

| Form | Syntax |
|------|--------|
| Line comment | `// …` |
| Block comment | `/* … */` |

Whitespace separates tokens; **semicolons `;` terminate statements** (required between statements).

---

## 3. Identifiers and case

- **Reserved words** and **ASCII identifiers** are matched **case-insensitively**; the AST stores normalized spellings (typically lowercase).
- **`this`** is a keyword.
- **`var`** is **reserved** — use **`let`** instead.

---

## 4. Literals

| Kind | Examples / notes |
|------|------------------|
| Number | Decimal (`42`, `3.14`), scientific (`1e2`), **hex** (`0xff`), **binary** (`0b1010`). |
| String | `"..."` with escapes (`\n`, `\t`, `\u{...}`, etc.). |
| Boolean | `true`, `false`. |
| Null | `null`. |
| Array | `[e1, e2, …]` — elements may use **spread** **`...expr`** to splice another array’s elements (see [§12](#12-objects-and-arrays)). |
| Object | `{ key: expr, … }` — supports **computed keys** **`[expr]: value`** where grammar allows; **method shorthand** **`name() { … }`** with **`this`**. |
| Template string | `` `text ${ expr } more` `` — backtick literals with **`${ … }`** holes; expressions are evaluated and coerced like normal concatenation. |

---

## 5. Keywords

From **`internal/lexer/token.go`** (spelling shown lowercase; any case accepted):

`break`, `case`, `continue`, `default`, `delete`, `do`, `else`, `false`, `for`, `func`, `if`, `import`, `in`, `let`, `null`, `of`, `return`, `switch`, `this`, `true`, `typeof`, `while`, `var` (reserved error).

Directive-style tokens:

- **`#include`** — include path string.

---

## 6. Operators and punctuation

**Arithmetic / bitwise**

`+` `-` `*` `/` `%` `**` (power), unary `+` `-`,  
`&` `|` `^` `~` `<<` `>>` `>>>`  

**Comparison**

`<` `<=` `>` `>=` `==` `!=` `===` `!==`

**Logical**

`&&` `||` `!`

**Nullish**

`??` — nullish coalescing: yields the right operand only when the left is **`null`** (not for other falsy values).

**Optional chaining**

`?.` — member (`obj?.prop`) or computed (`obj?.[expr]`); if the receiver is **`null`**, the whole expression is **`null`** without evaluating the property/index.

**Assignment**

`=` `+=` `-=` `*=` `/=` `%=` `&=` `|=` `^=` `<<=` `>>=`

**Increment**

`++` `--` (prefix and postfix where grammar permits)

**Other punctuation**

`( ) { } [ ] , ; . : ?`  
`=>` — switch/if **expression** arms, lambda-style branches  
`..` — range expressions  
`...` — **rest** parameters in function headers (**`...name`**) **and** **spread** inside array literals (**`[...a, b]`**).

**Unary**

`typeof` — unary operator (**`typeof x`**, same spelling as JS); behaviour matches the **`type(...)`** builtin (type name string).

---

## 7. Types (logical)

Dynamic types at runtime; LLVM uses NaN-boxed **`i64`** for tagged values.

| Kind | Notes |
|------|--------|
| Number | IEEE-754 double in C runtime. |
| String | Byte-oriented UTF-8 storage. |
| Bool | |
| Null | |
| Array | Dense vector with `len`. |
| Object / table | String-like keys + values; **`len(obj)`** is **entry count**. |
| Function | Named or **`func(...) { … }`** expressions; closures supported. |

---

## 8. Declarations

```fuji
let name = expr;
let name;              // null

let { x, y } = point; // object destructuring: binds fields `x` and `y` from an object expression

func name(params) {
    statements
}

#include "relative.fuji"
```

**Native FFI hint** (special `//` comment parsed by **`parser`**):

```fuji
// fuji: extern bindingName symbol_name [arity]
```

Used with **`let`** / **`func`** declarations that follow (see **`parser_stmt.go`**).

---

## 9. Statements

All **`if`**, **`else`**, loops, and **`switch`** bodies require **`{ … }`** braces (no braceless single-statement branches).

| Statement | Forms |
|-----------|--------|
| Expression statement | `expr;` |
| Block | `{ … }` |
| **`if`** | `if (cond) stmt [else stmt]` — **`if (…) { … } else { … }`** |
| **`switch`** | `switch (subject) { case value: stmts … default: stmts }` — lowering **chains cases** and **falls through** between **`case`** bodies unless **`break`** ends a branch early (see **`emitSwitchStmt`**). |
| **`while`** | `while (cond) stmt` |
| **`do`** **`while`** | `do stmt while (cond);` |
| **`for` (C-style)** | `for ([let binds [, …]] ; cond ; incr [, …]) stmt` — init may be empty (`for (;;)`) ; **`let`** bindings may repeat with commas ; increments are comma-separated expressions. |
| **`for`** **`in`** | `for (let key in iterable) stmt` — **keys**: arrays yield **numeric indices** as values ; tables yield keys in **slot order**. |
| **`for`** **`of`** | `for (let x of iterable) stmt` — **values** in slot order. |
| **`for`** **`of`** destructuring | `for (let [k, v] of iterable) stmt` — arrays: **`k`** numeric index, **`v`** element ; tables: insertion-order **key** / **value** pairs (`fuji_forof_*` runtime). |
| **`return`** | `return;` / `return expr;` |
| **`break`** / **`continue`** | Loop / switch (**`switch`**:** **`break`** exits **`switch`** where emitted). |
| **`delete`** | **`delete target;`** where **`target`** is a property access (**`obj.field`** or **`obj["key"]`** — dot form parses as bracket access). Removes an **own** key from a table object via **`fuji_object_remove`**. Prefer non-optional access here. |

---

## 10. Expressions

- **Assignment** (`=` and compounds), **logical** (`&&` `||`), **binary**, **unary** (including **`typeof expr`**), **postfix** (`++` `--` calls indexing).
- **Nullish coalescing** **`left ?? right`** — see [§6](#6-operators-and-punctuation).
- **Optional chaining** on **`?.`** member/index reads — see [§6](#6-operators-and-punctuation).
- **Call** `f(a,b)` — builtins resolve **case-insensitively**.
- **Member** `obj.field`
- **Index** `obj[key]` — arrays, strings, tables.
- **`import`** expression — **`import("path")`** (see tests / loader).
- **Array / object / tuple literals**, **range** **`from..to`** where the grammar defines **`RangeExpr`**.
- **`if`** **expression**: `let x = if (c) { a } else { b };`
- **`switch`** **expression**: `switch (x) { case 1 => "a" default => "z" }`

---

## 11. Functions

```fuji
func f(a, b = defaultExpr, ...rest) {
    return expr;
}

let g = func(x) { return x + 1; };
```

- **Default parameters** (literal defaults where supported).
- **Rest parameter**: **`...name`** collects remaining arguments.
- **Closures** capture **`let`** bindings (**LLVM** path + GC roots).

---

## 12. Objects and arrays

- **Object literals**, **dot** / **bracket** access, **`this`** on **`obj.method()`**.
- **Array literals**: **`[a, ...rest, b]`** — **`...expr`** concatenates **`expr`** into the literal using **`concat`** semantics; **`expr`** should be an **array** (other runtime kinds are unsupported for spread).
- **Arrays**: **`push`**, **`pop`**, **`length`** as callable-style methods on values (see **`docs/language.md`** / tests).
- **`len(x)`** — string length, array count, **object entry count**.
- **Instance methods** on strings and arrays are spelled like JS (**`split`**, **`map`**, …) but identifiers are **case-insensitive**, so **`toUpper`** and **`toupper`** both resolve — see [§17](#17-string-methods) and [§18](#18-array-methods).

Iteration summary:

| Loop | Arrays | Tables / objects |
|------|--------|-------------------|
| **`for (let k in x)`** | index keys (`0`… as numbers) | key per slot |
| **`for (let v of x)`** | elements | values per slot |
| **`for (let [k,v] of x)`** | index + element | key + value |

Order is **runtime slot order** (insertion order for tables).

---

## 13. Modules and includes

| Feature | Role |
|---------|------|
| **`#include "path.fuji"`** | Textual include ; participates in **`parser` bundle** loading. |
| **`import("path")`** | Expression module load (see **`parser`** / **`loader`**). |

---

## 14. Native builtins (global calls)

Registered in **`internal/codegen/builtin_register.go`** (names here lowercase; source may use any case):

**Errors / control**

| Name | Role |
|------|------|
| `ok` | Result success (`argv` convention). |
| `err` | Result failure. |
| `panic` | Fatal error / stderr / trace / exit. |
| `assert` | Assertion helper. |

**I/O / strings**

| Name | Role |
|------|------|
| `print` | Print values (newline-terminated batch). |
| `format` | Format string helper. |
| `readFile`, `writeFile`, `appendFile`, `fileExists`, `deleteFile` | Filesystem (**Results** where applicable). |
| `trace` | Debug trace. |

**Values**

| Name | Role |
|------|------|
| `len` | Length / count. |
| `type`, `typeof` | Type name string (**`typeof x`** unary is equivalent to **`type(x)`**). |
| `number`, `string` | Coercion / parsing. |
| `bool` | Bool coercion. |
| `parseJSON`, `toJSON` | JSON helpers. |

**Time**

| Name | Role |
|------|------|
| `clock`, `time`, `timestamp`, `programTime`, `sleep` | Wall / program timing (`deltaTime` also registered). |

**Random**

| Name | Role |
|------|------|
| `random`, `randomInt`, `randomChoice`, `randomSeed` | RNG helpers. |

**Math (scalar / argv)**  

These remain **globals** for backward compatibility:  
`abs`, `sqrt`, `lerp`, `clamp`, `distance`, `angleBetween`, `map`, `pi`, `e`,  
`sin`, `cos`, `tan`, `asin`, `acos`, `atan`, `atan2`,  
`pow`, `exp`, `log`, `log10`,  
`floor`, `ceil`, `round`, `trunc`, `sign`, `min`, `max`,  
`smoothstep`, `distanceSq`, `normalize`

The compiler **prepends** a prelude **`let math = { floor: floor, … };`** (names folded to lowercase **`math`**) when your entry file does not already declare **`math`**, so **`math.floor(x)`** works like a namespace (see [`internal/parser/prelude.go`](internal/parser/prelude.go)).

**Pattern-style substring check**

| Name | Role |
|------|------|
| `matches` | **`matches(haystack, pattern)`** — both operands coerced to strings; returns whether **`pattern`** occurs as a **substring** (implemented with **`strstr`**, not full regex). |


**Type predicates**

`isNumber`, `isString`, `isBool`, `isNull`, `isArray`, `isObject`, `isFunction`

**GC**

| Name | Role |
|------|------|
| `gc` | **`fuji_gc_collect`** |
| `gcFrameStep` | Budgeted incremental step |

Symbols exist in **`runtime.go`** declarations without a Fuji **`registerBuiltinFuncs`** entry until wired (example: **`fuji_wall_time`**). Treat **`builtin_register.go`** as the **user-visible native surface**.

---

## 15. Truthy / falsy

Runtime **`FUJI_is_truthy`** (conceptually): **`false`**, **`null`**, **`0`**, **`""`** are falsy ; most other values truthy.

---

## 16. Precedence (summary)

Highest → lowest (approximate): **call, member, index** → unary (`typeof`, …) → `* / %` → `+ -` → comparisons → **`==` `!=` `===` `!==`** → **`&&`** → **`||`** → **`??`** → assignment.

Exact rules: **`internal/parser/parser_expr.go`** (Pratt parsing).

---

## 17. String methods

Called as **`receiver.method(...)`** (names **case-insensitive**). Implemented via **`fuji_string_*`** argv helpers unless noted.

| Method | Arguments | Notes |
|--------|-----------|--------|
| **`split`** | `(delimiter)` | Empty delimiter splits into single-character strings. |
| **`trim`** | `()` | |
| **`toUpper`** / **`toupper`** | `()` | |
| **`toLower`** / **`tolower`** | `()` | |
| **`replace`** | `(search, replacement)` | Single replace (runtime semantics). |
| **`replaceAll`** | `(search, replacement)` | |
| **`indexOf`** | `(needle)` | Index or **`-1`**; also defined on arrays (disambiguated at codegen). |
| **`includes`** | `(needle)` | Substring / element test (arrays vs strings disambiguated at codegen). |
| **`slice`** | `(start, end)` | Half-open range; **string** vs **array** receiver disambiguated at codegen. |
| **`startsWith`** | `(prefix)` | |
| **`endsWith`** | `(suffix)` | |

---

## 18. Array methods

| Method | Arguments | Notes |
|--------|-----------|--------|
| **`concat`** | `(…values)` | Variadic **`concat`**. |
| **`join`** | `(separator)` | **`fuji_array_join`**. |
| **`map`** | `(callback)` | Callback **`(element)`**; new array. |
| **`filter`** | `(callback)` | Truthy predicate **`(element)`**. |
| **`reduce`** | `(callback)` or **`(callback, initial)`** | Callback **`(accumulator, element)`**; without **`initial`**, first element seeds the accumulator (empty array yields **`null`** without **`initial`**). |
| **`find`** | `(callback)` | First element satisfying predicate, else **`null`**. |
| **`slice`** | `(start, end)` | Copy range (half-open). |
| **`sort`** | `()` | In-place sort (runtime). |
| **`reverse`** | `()` | In-place reverse. |
| **`indexOf`** | `(item)` | |
| **`includes`** | `(item)` | |

Higher-order methods (**`map`**, **`filter`**, **`find`**, **`reduce`**) are lowered in **`internal/codegen/methods.go`** (LLVM loops + indirect calls).

---

## 19. Templates, math namespace, pattern matching, and operators

Summary of surface syntax beyond classic Fuji:

| Feature | Syntax / notes |
|---------|----------------|
| Template literals | `` `Hello, ${ name }!` `` |
| **`math`** namespace | Prelude **`math`** object (**§14**) shadows **`Math`** spelling because identifiers fold to lowercase. |
| **Substring “regex”** | Global **`matches(haystack, pattern)`** — literal substring match, not a regex engine (**§14**). |
| **`typeof`** | Unary **`typeof value`** — same idea as **`type(value)`**. |
| **`delete`** | **`delete obj.key;`** (**§9**). |
| Spread | **`[...a, item]`** (**§12**). |
| Destructuring | **`let { x, y } = expr;`** (**§8**). |
| Optional chaining | **`obj?.prop`**, **`obj?.[i]`** (**§6**, **§10**). |
| Nullish coalescing | **`a ?? b`** — only substitutes when **`a`** is **`null`** (**§6**, **§10**). |

---

## See also

- **[docs/user_guide.md](docs/user_guide.md)** — narrative guide.  
- **[docs/language/syntax.md](docs/language/syntax.md)** — syntax cheat sheet.  
- **[docs/reference.md](docs/reference.md)** — stdlib / API-oriented notes if present.  
- **`tests/*.fuji`** — executable examples.

---

*Generated as a consolidated catalog; when docs disagree with **`internal/parser`** or **`builtin_register.go`**, trust the implementation.*
