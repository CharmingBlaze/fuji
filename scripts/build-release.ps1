# scripts/build-release.ps1
# Builds a local release binary for testing. Requires LLVM on PATH.

param(
    [string]$Platform = "windows"
)

Write-Host "Building Fuji release binary for $Platform..."

Write-Host "Building C runtime..."
New-Item -ItemType Directory -Force runtime\obj_win | Out-Null
clang -c runtime\src\value.c        -O2 -std=c11 -Iruntime\src -D_AMD64_=1 -o runtime\obj_win\value.o
clang -c runtime\src\object.c       -O2 -std=c11 -Iruntime\src -D_AMD64_=1 -o runtime\obj_win\object.o
clang -c runtime\src\gc.c           -O2 -std=c11 -Iruntime\src -D_AMD64_=1 -o runtime\obj_win\gc.o
clang -c runtime\src\fuji_runtime.c -O2 -std=c11 -Iruntime\src -D_AMD64_=1 -o runtime\obj_win\fuji_runtime.o
llvm-ar rcs runtime\libfuji_runtime.a runtime\obj_win\*.o

Write-Host "Populating embed directory..."
New-Item -ItemType Directory -Force internal\embed\windows\amd64 | Out-Null
Copy-Item (Get-Command clang).Source internal\embed\windows\amd64\clang.exe
Copy-Item (Get-Command lld).Source   internal\embed\windows\amd64\lld.exe
Copy-Item runtime\libfuji_runtime.a  internal\embed\windows\amd64\

Write-Host "Building fuji.exe..."
go build -trimpath -tags release -ldflags="-s -w" -o fuji-release.exe .\cmd\fuji

Write-Host ""
Write-Host "Done. Test with:"
Write-Host "  .\fuji-release.exe build tests\hello.fuji -o hello.exe"
Write-Host "  .\hello.exe"
