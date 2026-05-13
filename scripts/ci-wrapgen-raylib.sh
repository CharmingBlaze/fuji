#!/usr/bin/env bash
# PR CI: Raylib header + wrapgen + fuji link smoke.
# Linux: apt installs lib when needed; macOS: Homebrew. Vendored stdlib/sys-include/raylib.h
# supplies the header, but the native lib is still required to link the smoke binary.
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
  echo "==> Using vendored raylib header: $HDR"
else
  echo "==> No stdlib/sys-include/raylib.h — installing system Raylib (Linux only path here)"
  case "$(uname -s)" in
    Linux)
      sudo apt-get update -qq
      sudo apt-get install -y libraylib-dev pkg-config
      for c in /usr/include/raylib.h /usr/local/include/raylib.h; do
        if [[ -f "$c" ]]; then HDR="$c"; break; fi
      done
      ;;
    Darwin)
      brew list raylib &>/dev/null || brew install raylib
      _rb="$(brew --prefix raylib 2>/dev/null || true)"
      if [[ -n "$_rb" && -f "$_rb/include/raylib.h" ]]; then
        HDR="$_rb/include/raylib.h"
      fi
      ;;
  esac
  if [[ -z "$HDR" ]]; then
    echo "raylib.h not found (add stdlib/sys-include/raylib.h or install Raylib dev)"
    exit 1
  fi
  echo "    header: $HDR"
fi

ensure_raylib_for_link() {
  if pkg-config --exists raylib 2>/dev/null; then
    echo "==> raylib pkg-config ok"
    return 0
  fi
  case "$(uname -s)" in
    Linux)
      echo "==> Installing libraylib-dev for link (Linux)"
      sudo apt-get update -qq
      sudo apt-get install -y libraylib-dev pkg-config
      ;;
    Darwin)
      echo "==> Installing raylib via Homebrew for link (macOS)"
      brew list raylib &>/dev/null || brew install raylib
      _rb="$(brew --prefix raylib)"
      _pc="$_rb/lib/pkgconfig"
      if [[ -d "$_pc" ]]; then
        export PKG_CONFIG_PATH="${_pc}${PKG_CONFIG_PATH:+:}${PKG_CONFIG_PATH:-}"
      fi
      ;;
    *)
      echo "Unsupported OS: $(uname -s)" >&2
      exit 1
      ;;
  esac
}

ensure_raylib_for_link

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
  case "$(uname -s)" in
    Darwin)
      _rb="$(brew --prefix raylib)"
      export FUJI_LINKFLAGS="-I${_rb}/include -L${_rb}/lib -lraylib"
      ;;
    *)
      export FUJI_LINKFLAGS="-I/usr/include -lraylib"
      ;;
  esac
fi

FUJI_BIN="${FUJI_BIN:-}"
if [[ -z "$FUJI_BIN" || ! -x "$FUJI_BIN" ]]; then
  FUJI_BIN="$ROOT/.ci_fuji"
  go build -trimpath -o "$FUJI_BIN" ./cmd/fuji
fi

"$FUJI_BIN" build "$ROOT/tests/raylib_wrapgen_smoke.fuji" -o /tmp/raylib_wrapgen_smoke
test -x /tmp/raylib_wrapgen_smoke
echo "==> OK: /tmp/raylib_wrapgen_smoke linked successfully"
