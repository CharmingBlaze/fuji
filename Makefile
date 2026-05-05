# Fuji — build from repo root with GNU Make (MSYS2/Unix/Git Bash).
# On Windows, use GNU Make (e.g. mingw32-make from MSYS2/MinGW), not a shadowed `make.bat`.
# Requires sh.exe on PATH (Git for Windows / MSYS) for mkdir -p / rm -rf in recipes.
SHELL := sh.exe
.SHELLFLAGS := -ec
#
#   make              — runtime static lib, then fuji + fujiwrap, then go test
#   make runtime-lib  — only runtime/libfuji_runtime.a
#   make raylib-lib   — CMake Raylib into third_party/raylib_static/stage/ (needs raylib/ + cmake)

.PHONY: all build fuji fujiwrap wrapgen runtime-lib raylib-lib raylib-clean test soak-test fmt clean

all: build test

build: fuji fujiwrap

fuji:
	@mkdir -p bin
	go build -trimpath -ldflags "-s -w" -o bin/fuji ./cmd/fuji

fujiwrap:
	@mkdir -p bin
	go build -trimpath -ldflags "-s -w" -o bin/fujiwrap ./cmd/wrapgen

wrapgen: fujiwrap
	@mkdir -p bin
	go build -trimpath -ldflags "-s -w" -o bin/wrapgen ./cmd/wrapgen

# Static library required for native `fuji build` and for codegen tests.
runtime-lib:
	$(MAKE) -C runtime

raylib-lib:
	$(MAKE) -C third_party/raylib_static

raylib-clean:
	$(MAKE) -C third_party/raylib_static clean

test: runtime-lib
	go test ./... -count=1

soak-test: fuji runtime-lib
	./bin/fuji run tests/gc_pressure_expr.fuji
	./bin/fuji run tests/globals_perf.fuji
	./bin/fuji run tests/gc_soak.fuji

fmt:
	gofmt -w .

clean:
	rm -rf bin .FUJI_build
	go clean -cache
	$(MAKE) -C runtime clean
