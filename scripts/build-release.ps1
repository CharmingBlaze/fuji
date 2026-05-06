# scripts/build-release.ps1
# Builds a local release binary for testing.
# Offline-only: this script never downloads dependencies.

param(
    [string]$Platform = "windows",
    [string]$LLVMBin = "C:\Program Files\LLVM\bin",
    [string]$MinGWBin = "C:\ProgramData\mingw64\mingw64\bin",
    [string]$OutputDir = "dist\offline-release"
)

$ErrorActionPreference = "Stop"
$RepoRoot = (Resolve-Path (Join-Path $PSScriptRoot "..")).Path

Write-Host "Building Fuji release binary for $Platform..."
Write-Host "Mode: offline (no downloads)"

if ($Platform -ne "windows") {
    throw "This script currently supports Platform=windows only."
}

$clangExe = Join-Path $LLVMBin "clang.exe"
$lldExe = Join-Path $LLVMBin "lld.exe"
$clangDriver = Join-Path $RepoRoot "scripts\clang-gnu.cmd"
$mingwArExe = Join-Path $MinGWBin "ar.exe"
$llcExe = Join-Path $LLVMBin "llc.exe"

if (!(Test-Path $clangExe)) { throw "Missing $clangExe" }
if (!(Test-Path $lldExe)) { throw "Missing $lldExe" }
if (!(Test-Path $mingwArExe)) { throw "Missing $mingwArExe" }
if (!(Test-Path $clangDriver)) { throw "Missing $clangDriver" }

$env:PATH = "$LLVMBin;$MinGWBin;$env:PATH"

$fujiMain = Join-Path $RepoRoot "cmd\fuji\main.go"
$fujiMainText = Get-Content -Raw $fujiMain
$versionMatch = [regex]::Match($fujiMainText, 'var version = "([^"]+)"')
if (!$versionMatch.Success) {
    throw "Could not extract Fuji version from $fujiMain"
}
$releaseVersion = $versionMatch.Groups[1].Value
Write-Host "Version: $releaseVersion"

Write-Host "Building C runtime..."
New-Item -ItemType Directory -Force runtime\obj_win | Out-Null
& cmd /c """$clangDriver"" -c runtime\src\value.c        -O2 -std=c11 -Iruntime\src -D_AMD64_=1 -o runtime\obj_win\value.o"
if ($LASTEXITCODE -ne 0) { throw "Failed compiling runtime\src\value.c" }
& cmd /c """$clangDriver"" -c runtime\src\object.c       -O2 -std=c11 -Iruntime\src -D_AMD64_=1 -o runtime\obj_win\object.o"
if ($LASTEXITCODE -ne 0) { throw "Failed compiling runtime\src\object.c" }
& cmd /c """$clangDriver"" -c runtime\src\gc.c           -O2 -std=c11 -Iruntime\src -D_AMD64_=1 -o runtime\obj_win\gc.o"
if ($LASTEXITCODE -ne 0) { throw "Failed compiling runtime\src\gc.c" }
& cmd /c """$clangDriver"" -c runtime\src\fuji_runtime.c -O2 -std=c11 -Iruntime\src -D_AMD64_=1 -o runtime\obj_win\fuji_runtime.o"
if ($LASTEXITCODE -ne 0) { throw "Failed compiling runtime\src\fuji_runtime.c" }
& $mingwArExe rcs runtime\libfuji_runtime.a runtime\obj_win\value.o runtime\obj_win\object.o runtime\obj_win\gc.o runtime\obj_win\fuji_runtime.o
if ($LASTEXITCODE -ne 0) { throw "Failed creating runtime\libfuji_runtime.a" }

Write-Host "Populating embed directory..."
New-Item -ItemType Directory -Force internal\embed\windows\amd64 | Out-Null
Copy-Item $clangExe internal\embed\windows\amd64\clang.exe
Copy-Item $lldExe   internal\embed\windows\amd64\lld.exe
if (Test-Path $llcExe) {
    Copy-Item $llcExe internal\embed\windows\amd64\llc.exe
}
Copy-Item runtime\libfuji_runtime.a  internal\embed\windows\amd64\

Write-Host "Building fuji.exe..."
go build -trimpath -tags release -ldflags="-s -w -X main.version=$releaseVersion" -o fuji-release.exe .\cmd\fuji
if ($LASTEXITCODE -ne 0) { throw "Failed building fuji-release.exe" }

Write-Host "Building kujiwrap/fujiwrap..."
go build -trimpath -ldflags="-s -w -X main.WrapgenVersion=$releaseVersion" -o fujiwrap.exe .\cmd\wrapgen
if ($LASTEXITCODE -ne 0) { throw "Failed building fujiwrap.exe" }
Copy-Item -Force .\fujiwrap.exe .\kujiwrap.exe

Write-Host "Assembling offline distribution..."
New-Item -ItemType Directory -Force $OutputDir | Out-Null
Copy-Item -Force .\fuji-release.exe (Join-Path $OutputDir "fuji.exe")
Copy-Item -Force .\fujiwrap.exe (Join-Path $OutputDir "fujiwrap.exe")
Copy-Item -Force .\kujiwrap.exe (Join-Path $OutputDir "kujiwrap.exe")
if (Test-Path .\stdlib)  { Copy-Item -Recurse -Force .\stdlib  (Join-Path $OutputDir "stdlib") }
if (Test-Path .\wrappers){ Copy-Item -Recurse -Force .\wrappers (Join-Path $OutputDir "wrappers") }
if (Test-Path .\runtime) { Copy-Item -Recurse -Force .\runtime  (Join-Path $OutputDir "runtime") }

@"
Fuji offline bundle
-------------------
This build is self-contained for fuji build/run usage.
No installer or network download is required by the compiler.
Included binaries:
- fuji.exe
- fujiwrap.exe
- kujiwrap.exe (legacy compatibility name)
"@ | Set-Content -Path (Join-Path $OutputDir "README_OFFLINE.txt")

Write-Host ""
Write-Host "Done. Test with:"
Write-Host "  .\fuji-release.exe build tests\hello.fuji -o hello.exe"
Write-Host "  .\hello.exe"
Write-Host ""
Write-Host "Offline bundle:"
Write-Host "  $OutputDir"
