#!/usr/bin/env bash
# Placeholder: full Jolt Physics wrapgen in CI requires vendored Jolt headers + a stable parse surface.
# Enable by setting JOLT_HEADERS to a comma-separated list of .h paths, then run this script from repo root.
set -euo pipefail

ROOT="$(cd "$(dirname "$0")/.." && pwd)"
cd "$ROOT"

if [[ -z "${JOLT_HEADERS:-}" ]]; then
  echo "JOLT_HEADERS unset — skipping Jolt wrapgen CI (set to e.g. Jolt/Jolt.h,... to enable)."
  exit 0
fi

OUT="${JOLT_WRAP_OUT:-$ROOT/wrappers/jolt_ci_generated}"
rm -rf "$OUT"
mkdir -p "$OUT"

go run ./cmd/wrapgen -name jolt -headers "$JOLT_HEADERS" -out "$OUT" -docs=false -v
echo "Jolt wrapgen wrote to $OUT"
