# test

Fuji bindings for the **test** library.

---

## Files in this folder

| File | Description |
|------|-------------|
| `test.fuji` | Include this in your Fuji program. |
| `wrapper.c` | Compiled automatically by `fuji build` and `fuji bundle`. You do not need to touch this file. |
| `api_reference.md` | Full reference for every function, struct, and constant. |
| `examples.md` | Ready-to-run code examples. |

---

## Library summary

- **3** functions

---

## Quick start

**Step 1.** Include the bindings at the top of your Fuji program:

```fuji
#include "test.fuji"
```

**Step 2.** Call functions directly by name:

```fuji
let result = add(a, b);
print(result);
```

**Step 3.** Build or bundle:

```powershell
set FUJI_NATIVE_SOURCES=test\wrapper.c
set FUJI_LINKFLAGS=-I<include-dir> -L<lib-dir> -ltest

fuji build  mygame.fuji -o mygame.exe
fuji bundle mygame.fuji -o dist\mygame
```

---

## Troubleshooting

**Undefined symbol**  
Make sure `FUJI_NATIVE_SOURCES` points to `wrapper.c`.

**Missing header or library**  
Add `-I<dir>` for headers and `-L<dir> -ltest` for the library in `FUJI_LINKFLAGS`.

**Unexpected return values**  
Check the type conversions in `wrapper.c`. Pointer and struct types may need manual adjustment for complex cases.

---

## See also

- [API Reference](api_reference.md)
- [Examples](examples.md)
