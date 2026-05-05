#!/usr/bin/env bash
# PR CI: install Raylib dev headers, regenerate bindings with wrapgen, then fuji build a tiny sample.
# Binding drift is covered implicitly by this link smoke (raylib.fuji + wrapper.c must exist and link).
# Usage: from repo root, with sudo-capable apt (e.g. GitHub Actions ubuntu-latest).
set -euo pipefail

ROOT="$(cd "$(dirname "$0")/.." && pwd)"
cd "$ROOT"

if [[ ! -f "runtime/libfuji_runtime.a" ]]; then
  echo "runtime/libfuji_runtime.a missing; run scripts/build-runtime.sh first"
  exit 1
fi

HDR=""
if [[ -f "$ROOT/stdlib/sys-include/raylib.h" ]]; then
  HDR="$ROOT/stdlib/sys-include/raylib.h"
  echo "==> Using vendored raylib 6 header: $HDR"
else
  echo "==> Installing Raylib dev headers + pkg-config (Debian/Ubuntu); no stdlib/sys-include/raylib.h"
  sudo apt-get update -qq
  sudo apt-get install -y libraylib-dev pkg-config
  for c in /usr/include/raylib.h /usr/local/include/raylib.h; do
    if [[ -f "$c" ]]; then HDR="$c"; break; fi
  done
  if [[ -z "$HDR" ]]; then
    echo "raylib.h not found (add stdlib/sys-include/raylib.h or install libraylib-dev)"
    exit 1
  fi
  echo "    header: $HDR"
fi

OUT="${WRAP_OUT:-$ROOT/wrappers/raylib_ci_generated}"
rm -rf "$OUT"
mkdir -p "$OUT"

echo "==> fujiwrap / wrapgen (this may take a minute)"
go run ./cmd/wrapgen -name raylib -headers "$HDR" -out "$OUT" -docs=false -build=false -v

if [[ ! -s "$OUT/raylib.fuji" ]] || [[ ! -s "$OUT/wrapper.c" ]]; then
  echo "wrapgen did not produce raylib.fuji / wrapper.c"
  exit 1
fi

echo "==> Link smoke: fuji build tests/raylib_wrapgen_smoke.fuji"
export FUJI_WRAPPERS="$OUT"
export FUJI_NATIVE_SOURCES="$OUT/wrapper.c"
if pkg-config --exists raylib 2>/dev/null; then
  export FUJI_LINKFLAGS
  FUJI_LINKFLAGS="$(pkg-config --libs --cflags raylib)"
else
  export FUJI_LINKFLAGS="-I/usr/include -lraylib"
fi

FUJI_BIN="${FUJI_BIN:-}"
if [[ -z "$FUJI_BIN" || ! -x "$FUJI_BIN" ]]; then
  FUJI_BIN="$ROOT/.ci_fuji"
  go build -trimpath -o "$FUJI_BIN" ./cmd/fuji
fi

"$FUJI_BIN" build "$ROOT/tests/raylib_wrapgen_smoke.fuji" -o /tmp/raylib_wrapgen_smoke
test -x /tmp/raylib_wrapgen_smoke
echo "==> OK: /tmp/raylib_wrapgen_smoke linked successfully"
