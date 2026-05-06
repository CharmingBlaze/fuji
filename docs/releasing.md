# Cutting a Fuji release

**Audience:** maintainers tagging **`v*`** on the default branch so **[`.github/workflows/release.yml`](../.github/workflows/release.yml)** publishes artifacts.

## Checklist

1. **`CHANGELOG.md`** — Move notes from **`[Unreleased]`** into a dated **`[X.Y.Z]`** section (see **`[0.2.0]`** as a template).
2. **Version strings** — Align **`cmd/fuji/main.go`** `version` and **`cmd/wrapgen/wrapgen_version.go`** `WrapgenVersion` with the tag (release builds may also pass `-ldflags "-X main.version=..."` if you use that flow).
3. **CI green** — **`go vet ./...`**, **`go test ./...`**, and **`fuji fmt --check`** (see **[CONTRIBUTING.md](../CONTRIBUTING.md)**).
4. **Tag** — `git tag v0.2.0 && git push origin v0.2.0` (example). The workflow builds **`-tags release`** `fuji` binaries with embedded Clang + **`libfuji_runtime.a`** (and **lld** on Windows).
5. **Post-release** — Open **`[Unreleased]`** again; bump dev versions (e.g. **`0.3.0-dev`**) if you ship nightlies from `main`.

## Notes

- **Linux CI** (**`ci.yml`**) does not embed the toolchain; **`release.yml`** repopulates **`internal/embed/...`** per job from the runner’s Clang/LLVM before **`go build -tags release`**.
- Smoke tests in release jobs build **`tests/hello.fuji`** and **`tests/gc_shadow_multi_return_pop.fuji`** after each platform binary is produced.
