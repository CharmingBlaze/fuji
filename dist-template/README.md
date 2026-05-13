# dist-template

Small **user-facing doc and example snippets** used in some packaging flows. The **full** GitHub Release SDK zip is assembled from the **repository root** (`stdlib/`, `docs/`, `wrappers/`, `examples/`, root `*.md`, plus vendored raylib under `third_party/`) by **`scripts/package-release-sdk.sh`** (Linux CI) or **`scripts/assemble-offline-sdk.ps1`** (Windows maintainers). See **`docs/distribution.md`** §7.

Do not put compiler internals, test files, or build scripts in this folder.
Users will see everything here.
