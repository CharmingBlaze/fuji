#!/usr/bin/env bash
# Build self-contained Fuji SDK archives for GitHub Releases.
# Requires: unzip's companion `zip`, bash, repo checkout at stdlib/, docs/, language.md, README.md.
#
# Usage: package-release-sdk.sh <VERSION> <ARTIFACTS_DIR> [OUT_DIR]
# Example (from repo root): bash scripts/package-release-sdk.sh "$GITHUB_REF_NAME" artifacts ./sdk-zips

set -euo pipefail

VERSION="${1:?first arg: version/tag e.g. v1.2.3}"
ART_ROOT="${2:?second arg: artifacts root (contains windows/, linux-amd64/, …)}"
OUT_DIR="${3:-./sdk-zips}"

if ! command -v zip >/dev/null 2>&1; then
  echo "zip(1) is required" >&2
  exit 1
fi

REPO_ROOT="$(cd "$(dirname "$0")/.." && pwd)"
cd "$REPO_ROOT"

require_file() {
  local f="$1"
  if [[ ! -f "$f" ]]; then
    echo "missing required file: $f" >&2
    exit 1
  fi
}

zip_sdk() {
  local slug="$1"          # windows-amd64
  local fj_src="$2"        # source path to compiler binary
  local fw_src="$3"        # source path to fujiwrap binary
  local fj_out="$4"        # basename inside archive (fuji or fuji.exe)
  local fw_out="$5"        # fujiwrap or fujiwrap.exe

  require_file "$fj_src"
  require_file "$fw_src"
  require_dir "$REPO_ROOT/stdlib"
  require_dir "$REPO_ROOT/docs"

  local root_name="fuji-${VERSION}-${slug}"
  local stage
  stage="$(mktemp -d)"

  mkdir -p "$stage/$root_name"
  cp -a "$fj_src" "$stage/$root_name/$fj_out"
  cp -a "$fw_src" "$stage/$root_name/$fw_out"
  cp -a "$REPO_ROOT/stdlib" "$stage/$root_name/stdlib"
  cp -a "$REPO_ROOT/docs" "$stage/$root_name/docs"

  for f in language.md README.md; do
    if [[ -f "$REPO_ROOT/$f" ]]; then
      cp -a "$REPO_ROOT/$f" "$stage/$root_name/"
    fi
  done

  # Small runnable corpus (optional; helps offline users).
  if [[ -d "$REPO_ROOT/examples" ]]; then
    cp -a "$REPO_ROOT/examples" "$stage/$root_name/examples"
  fi

  cat >"$stage/$root_name/SDK_README.txt" <<EOF
Fuji SDK ${VERSION} (${slug})
==============================

Platforms: Windows, Linux, and macOS each have dedicated release zips (this archive is ${slug}).
Everything offline: compiler, fujiwrap, stdlib, docs, examples — no Go or LLVM install required to compile .fuji.

Layout
------
  $(printf '%s' "$fj_out")     — Fuji compiler (${slug})
  $(printf '%s' "$fw_out")     — fujiwrap (C header → .fuji + wrapper.c)
  stdlib/         — shipped .fuji modules (@ imports, #includes)
  docs/           — guides (commands, wrappers, distribution, …)
  language.md     — language reference (repo root copy)
  README.md       — project overview & links

Use
---
  Keep this folder together so stdlib/ sits next to fuji (or run from this directory):

  • Windows:
      .\fuji.exe version
      .\fuji.exe run examples\hello.fuji

  • Linux / macOS:
      chmod +x fuji fujiwrap
      ./fuji version
      ./fuji run examples/hello.fuji

Wrapper tool
------------
  ./fuji wrap --help          (same as ./fujiwrap … when both sit here)

Embedded toolchain (release builds): Clang + runtime are bundled inside $(printf '%s' "$fj_out") — no separate LLVM download for Fuji itself.
Linking third-party native libs still uses your own headers/libs (see docs/wrappers.md).

EOF

  mkdir -p "$OUT_DIR"
  local out_zip
  out_zip="$(cd "$OUT_DIR" && pwd)/fuji-${VERSION}-sdk-${slug}.zip"
  rm -f "$out_zip"
  if [[ "$fj_out" != *.exe ]]; then
    chmod +x "$stage/$root_name/$fj_out" "$stage/$root_name/$fw_out"
  fi
  (cd "$stage" && zip -rq "$out_zip" "$root_name")
  rm -rf "$stage"
  echo "wrote $out_zip"
}

require_dir() {
  local d="$1"
  if [[ ! -d "$d" ]]; then
    echo "missing required directory: $d" >&2
    exit 1
  fi
}

ART_ROOT_ABS="$(cd "$ART_ROOT" && pwd)"

zip_sdk "windows-amd64" \
  "$ART_ROOT_ABS/windows/fuji-windows-amd64.exe" \
  "$ART_ROOT_ABS/windows/fujiwrap-windows-amd64.exe" \
  "fuji.exe" \
  "fujiwrap.exe"

zip_sdk "linux-amd64" \
  "$ART_ROOT_ABS/linux-amd64/fuji-linux-amd64" \
  "$ART_ROOT_ABS/linux-amd64/fujiwrap-linux-amd64" \
  "fuji" \
  "fujiwrap"

zip_sdk "linux-arm64" \
  "$ART_ROOT_ABS/linux-arm64/fuji-linux-arm64" \
  "$ART_ROOT_ABS/linux-arm64/fujiwrap-linux-arm64" \
  "fuji" \
  "fujiwrap"

zip_sdk "darwin-amd64" \
  "$ART_ROOT_ABS/macos/fuji-darwin-amd64" \
  "$ART_ROOT_ABS/macos/fujiwrap-darwin-amd64" \
  "fuji" \
  "fujiwrap"

zip_sdk "darwin-arm64" \
  "$ART_ROOT_ABS/macos/fuji-darwin-arm64" \
  "$ART_ROOT_ABS/macos/fujiwrap-darwin-arm64" \
  "fuji" \
  "fujiwrap"

echo "All SDK zips OK under $OUT_DIR"
