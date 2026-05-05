#!/usr/bin/env bash
# Nightly / manual: download pinned raylib.h, run wrapgen, compare SHA256 + binding count to testdata/wrapgen_golden/manifest.json
#
# Update golden after intentional generator changes or re-pinning the header:
#   UPDATE_GOLDEN=1 bash scripts/wrapgen-golden-audit.sh
# then tighten min_functions / max_functions in manifest.json if the defaults are too wide.
set -euo pipefail

ROOT="$(cd "$(dirname "$0")/.." && pwd)"
cd "$ROOT"

MANIFEST="$ROOT/testdata/wrapgen_golden/manifest.json"
if [[ ! -f "$MANIFEST" ]]; then
  echo "missing $MANIFEST"
  exit 1
fi

sha256_file() {
  if command -v sha256sum >/dev/null 2>&1; then
    sha256sum "$1" | awk '{print $1}' | tr '[:upper:]' '[:lower:]'
  else
    shasum -a 256 "$1" | awk '{print $1}' | tr '[:upper:]' '[:lower:]'
  fi
}

read_manifest() {
  python3 <<'PY'
import json, pathlib, sys
p = pathlib.Path("testdata/wrapgen_golden/manifest.json")
m = json.loads(p.read_text(encoding="utf-8"))
r = m["raylib_pinned"]
keys = ["header_url", "header_sha256", "raylib_fuji_sha256", "wrapper_c_sha256", "min_functions", "max_functions"]
for k in keys:
    print(r[k])
PY
}

mapfile -t GOLD < <(read_manifest)
HDR_URL="${GOLD[0]}"
expect_header_sha="${GOLD[1]}"
expect_fuji_sha="${GOLD[2]}"
expect_wrap_sha="${GOLD[3]}"
min_fn="${GOLD[4]}"
max_fn="${GOLD[5]}"

WORKDIR="${WRAPGEN_GOLDEN_WORKDIR:-$(mktemp -d)}"
cleanup() { rm -rf "$WORKDIR"; }
if [[ -z "${WRAPGEN_GOLDEN_WORKDIR:-}" ]]; then
  trap cleanup EXIT
fi

HDR="$WORKDIR/raylib_pinned.h"
echo "==> Download $HDR_URL"
curl -fsSL -o "$HDR" "$HDR_URL"

got_header_sha="$(sha256_file "$HDR")"
echo "    header sha256: $got_header_sha"
if [[ "$got_header_sha" != "$expect_header_sha" ]]; then
  echo "Pinned raylib.h SHA256 mismatch (tag content changed on GitHub?)."
  echo "  expected: $expect_header_sha"
  echo "  got:      $got_header_sha"
  echo "Fix: UPDATE_GOLDEN=1 then review manifest.json"
  exit 1
fi

OUT="$WORKDIR/out"
rm -rf "$OUT"
mkdir -p "$OUT"
echo "==> wrapgen"
go run ./cmd/wrapgen -name raylib -headers "$HDR" -out "$OUT" -docs=false -build=false

fn_line="$(grep -E -c '// (fuji|kuji):extern' "$OUT/raylib.fuji" || true)"
echo "    // fuji:extern count: $fn_line"
if [[ "$fn_line" -lt "$min_fn" ]] || [[ "$fn_line" -gt "$max_fn" ]]; then
  echo "Binding count outside [$min_fn, $max_fn] — generator or API drift."
  exit 1
fi

got_kuji="$(sha256_file "$OUT/raylib.fuji")"
got_wrap="$(sha256_file "$OUT/wrapper.c")"
echo "    raylib.fuji sha256: $got_kuji"
echo "    wrapper.c   sha256: $got_wrap"

if [[ "${UPDATE_GOLDEN:-}" == "1" ]]; then
  echo "==> UPDATE_GOLDEN: writing manifest.json"
  python3 <<PY
import json, pathlib
fn = int("$fn_line")
lo = max(1, fn - 40)
hi = fn + 40
data = {
    "raylib_pinned": {
        "tag": "6.0",
        "header_url": "$HDR_URL",
        "header_sha256": "$got_header_sha",
        "raylib_fuji_sha256": "$got_kuji",
        "wrapper_c_sha256": "$got_wrap",
        "min_functions": lo,
        "max_functions": hi,
    }
}
path = pathlib.Path("testdata/wrapgen_golden/manifest.json")
path.write_text(json.dumps(data, indent=2) + "\n", encoding="utf-8")
print("Wrote", path, "(tighten min_functions/max_functions manually if desired)")
PY
  exit 0
fi

if [[ "$got_kuji" != "$expect_fuji_sha" ]]; then
  echo "raylib.fuji SHA256 mismatch."
  echo "  expected: $expect_fuji_sha"
  echo "  got:      $got_kuji"
  exit 1
fi
if [[ "$got_wrap" != "$expect_wrap_sha" ]]; then
  echo "wrapper.c SHA256 mismatch."
  echo "  expected: $expect_wrap_sha"
  echo "  got:      $got_wrap"
  exit 1
fi

echo "==> Golden audit OK"
