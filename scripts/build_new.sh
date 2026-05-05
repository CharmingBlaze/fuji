#!/bin/bash
# Build script for Fuji using the new codegen and runtime

FUJI_FILE=$1
OUTPUT=${2:-"output"}

if [ -z "$FUJI_FILE" ]; then
    echo "Usage: $0 <file.fuji> [output]"
    exit 1
fi

echo "Building $FUJI_FILE..."

# Parse and generate LLVM IR
go run cmd/kuji/main.go check "$FUJI_FILE" || exit 1

# Generate LLVM IR
go run cmd/kuji/main.go disasm "$FUJI_FILE" > output.ll || exit 1

# Compile LLVM IR to object
llc -filetype=obj output.ll -o output.o || exit 1

# Link with new runtime library
gcc -static -O3 -s output.o runtime/libfuji_runtime.a -lm -o "$OUTPUT" || exit 1

# Clean up
rm -f output.ll output.o

echo "Built: $OUTPUT"
