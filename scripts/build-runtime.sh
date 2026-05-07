#!/usr/bin/env bash
set -euo pipefail
root="$(cd "$(dirname "$0")/.." && pwd)"

cc="${CC:-clang}"
if ! command -v "$cc" >/dev/null 2>&1; then
  cc="gcc"
fi
ar_bin="${AR:-ar}"

obj="$root/runtime/obj"
mkdir -p "$obj"

"$cc" -c "$root/runtime/src/value.c"        -O2 -std=c11 -I"$root/runtime/src" -o "$obj/value.o"
"$cc" -c "$root/runtime/src/object.c"       -O2 -std=c11 -I"$root/runtime/src" -o "$obj/object.o"
"$cc" -c "$root/runtime/src/gc.c"           -O2 -std=c11 -I"$root/runtime/src" -o "$obj/gc.o"
"$cc" -c "$root/runtime/src/fuji_runtime.c" -O2 -std=c11 -I"$root/runtime/src" -o "$obj/fuji_runtime.o"
"$ar_bin" rcs "$root/runtime/libfuji_runtime.a" "$obj/value.o" "$obj/object.o" "$obj/gc.o" "$obj/fuji_runtime.o"
