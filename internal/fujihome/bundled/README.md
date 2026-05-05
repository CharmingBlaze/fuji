# Legacy: bundled LLVM layout

Release pipeline used to populate **`internal/fujihome/bundled/`** with `llc`, `lld`, and `libfuji_runtime.a`.

That layout is **replaced** by **`internal/embed/<GOOS>/<GOARCH>/`** (Clang + runtime, and `lld.exe` on Windows). See **`internal/embed/README.md`**.

This directory remains for older notes and tooling references only.
