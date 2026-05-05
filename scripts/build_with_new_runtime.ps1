# Build script using the new runtime library

param(
    [Parameter(Mandatory=$true)]
    [string]$KujiFile,
    [string]$Output = "output.exe"
)

Write-Host "Building $KujiFile with new runtime..."

# Generate LLVM IR using the genir command
go run cmd/genir/main.go $KujiFile output.ll
if ($LASTEXITCODE -ne 0) {
    exit 1
}

# Compile LLVM IR to object
llc -filetype=obj output.ll -o output.o
if ($LASTEXITCODE -ne 0) {
    exit 1
}

# Link with new runtime library
gcc -static -O3 -s output.o runtime/libfuji_runtime.a -lm -o $Output
if ($LASTEXITCODE -ne 0) {
    exit 1
}

# Clean up
Remove-Item -Force output.ll, output.o -ErrorAction SilentlyContinue

Write-Host "Built: $Output"
