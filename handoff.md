# Fuji Project Handoff

## Project Overview
Fuji is a high-performance programming language designed for game development and systems programming. We are currently in the process of shipping **Fuji v1.0**, which transitions the project from a VM-based interpreter to a production-grade compiler with an LLVM backend.

## Architecture (v1.0)
The compiler follows a traditional multi-pass architecture:
1. **Frontend (Go)**:
   - `internal/lexer`: High-performance byte-based lexer.
   - `internal/ast`: Clean, decoupled AST definitions.
   - `internal/parser`: Recursive descent parser (Pratt parsing for expressions).
   - `internal/sema`: Semantic analysis and type checking.
2. **Middle-end (LLVM IR)**:
   - `internal/codegen`: AST to LLVM IR lowering.
3. **Runtime (C)**:
   - `runtime/`: Tri-generational garbage collector and standard library primitives.

## Current Status (Work in Progress)
- [x] **Lexer**: Fully implemented and tested. Optimized for speed using byte-slice scanning.
- [x] **AST**: Redefined in `internal/ast` to be compiler-friendly and decoupled from runtime.
- [/] **Parser**: Currently being rewritten to use the new Lexer and AST.
- [ ] **Codegen**: Pending (Targeting LLVM 17).
- [ ] **Runtime**: Pending (C implementation with tri-generational GC).

## Development Guide
### Prerequisites
- Go 1.21+
- LLVM 17
- C compiler (gcc/clang)

### Running Tests
```bash
go test ./internal/...
```

### Building the Compiler
```bash
go build -o kuji ./cmd/kuji
```

## Key Design Decisions
- **Byte-based Lexer**: We avoid `rune` conversions for performance. UTF-8 is supported but ASCII is the fast path.
- **Integer-based Tokens**: `TokenType` is a `uint8` for fast switch-based dispatch in the parser.
- **Decoupled AST**: The AST does not depend on the VM or Codegen, allowing for better testability and modularity.

## Contact & links

- **Mission:** ship Fuji v1.0 production release.
- **Upstream:** [github.com/CharmingBlaze/fuji](https://github.com/CharmingBlaze/fuji).
