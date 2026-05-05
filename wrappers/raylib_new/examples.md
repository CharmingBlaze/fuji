# raylib - Examples

Copy any snippet into your Fuji program and adjust the arguments.

---

## Include the library

```fuji
#include "raylib.fuji"
```

---

## Functions

### void

```fuji
let result = void();
print(result);
```

### bool

```fuji
let result = bool();
print(result);
```

### bool

```fuji
let result = bool();
print(result);
```

### InitWindow

```fuji
let result = InitWindow(width, height, *title);
print(result);
```

### CloseWindow

```fuji
let result = CloseWindow();
print(result);
```

### WindowShouldClose

```fuji
let result = WindowShouldClose();
print(result);
```

### IsWindowReady

```fuji
let result = IsWindowReady();
print(result);
```

### IsWindowFullscreen

```fuji
let result = IsWindowFullscreen();
print(result);
```

*See [api_reference.md](api_reference.md) for all 554 functions.*

---

## Structs

Structs are passed as Fuji objects with matching field names:

```fuji
let obj = { x: 0,  y: 0,  component: 0 };
```

---

## Enum values

```fuji
let bool_false = 0;
let bool_true = 1;
```

