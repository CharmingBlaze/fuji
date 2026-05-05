package fujiruntime

import _ "embed"

// This package embeds a large alternate C runtime (data/fuji.c, data/fuji.h).
//
// The default `fuji build` / [internal/nativebuild.Build] pipeline does **not** link these
// files today: it links the C runtime built from [runtime/src] into
// runtime/libfuji_runtime.a. Keep LLVM `declare` names in [internal/codegen/runtime.go]
// aligned with that tree. Treat this embed as experimental / reserved unless a future
// backend switches to shipping monolithic fuji.c again.

//go:embed data/fuji.c
var FujiC []byte

//go:embed data/fuji.h
var FujiH []byte
