#!/usr/bin/env bash
# Native compile/run smoke used by GitHub Actions on Linux, macOS, and Windows (Git Bash).
# Requires: FUJI_CI_BIN (path to fuji or fuji.exe), CI_ARTIFACT_DIR (writable output directory).
set -euo pipefail

FUJI="${FUJI_CI_BIN:?FUJI_CI_BIN is required}"
OUT="${CI_ARTIFACT_DIR:?CI_ARTIFACT_DIR is required}"

EXE=""
case "$(uname -s)" in
  MINGW* | MSYS* | CYGWIN*) EXE=.exe ;;
esac

mkdir -p "$OUT"

run_out() {
  local lines="$1"
  local bin="$2"
  shift 2
  "$bin" "$@" | head -n "$lines"
}

echo "==> native smoke: FUJI=$FUJI OUT=$OUT EXE=$EXE"

"$FUJI" build tests/smoke_native.fuji -o "$OUT/smoke_native$EXE"
test -x "$OUT/smoke_native$EXE"
run_out 5 "$OUT/smoke_native$EXE"

"$FUJI" build tests/for_of_pairs.fuji -o "$OUT/for_of_pairs$EXE"
test -x "$OUT/for_of_pairs$EXE"
run_out 10 "$OUT/for_of_pairs$EXE"

"$FUJI" build tests/nullish_assign.fuji -o "$OUT/nullish_assign$EXE"
run_out 5 "$OUT/nullish_assign$EXE"

"$FUJI" build tests/for_of_dynamic_range.fuji -o "$OUT/for_of_dynamic_range$EXE"
run_out 5 "$OUT/for_of_dynamic_range$EXE"

"$FUJI" build tests/defer_test.fuji -o "$OUT/defer_test$EXE"
run_out 5 "$OUT/defer_test$EXE"

"$FUJI" build tests/assert_string_eq.fuji -o "$OUT/assert_string_eq$EXE"
run_out 5 "$OUT/assert_string_eq$EXE"

"$FUJI" build tests/unary_plus_fold.fuji -o "$OUT/unary_plus_fold$EXE"
run_out 5 "$OUT/unary_plus_fold$EXE"

"$FUJI" build tests/gc_shadow_multi_return_pop.fuji -o "$OUT/gc_shadow_multi_return_pop$EXE"
run_out 5 "$OUT/gc_shadow_multi_return_pop$EXE"

"$FUJI" build tests/gc_control_test.fuji -o "$OUT/gc_control_test$EXE"
run_out 5 "$OUT/gc_control_test$EXE"

"$FUJI" run tests/vec3_test.fuji
"$FUJI" run tests/vec3_math_lerp_direct_test.fuji
"$FUJI" run tests/math_lerp_member_test.fuji
"$FUJI" run tests/array_push_growth_test.fuji
"$FUJI" run tests/multi_while_flow_test.fuji

"$FUJI" build --debug tests/hello.fuji -o "$OUT/hello_debug$EXE"
test -x "$OUT/hello_debug$EXE"

echo "==> native smoke OK"
