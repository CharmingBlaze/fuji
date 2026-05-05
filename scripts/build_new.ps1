# Build script for Fuji using the new codegen and runtime

param(
    [Parameter(Mandatory=$true)]
    [string]$KujiFile,
    [string]$Output = "output.exe"
)

Write-Host "Building $KujiFile..."

# Parse and generate LLVM IR
go run cmd/kuji/main.go check $KujiFile
if ($LASTEXITCODE -ne 0) {
    exit 1
}

# Generate LLVM IR
go run cmd/kuji/main.go disasm $KujiFile | Out-File -Encoding ASCII output.ll
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
