# Build fuji + fujiwrap (C header → .fuji + wrapper.c) into ./bin — one script for Windows authors.
$ErrorActionPreference = "Stop"
$root = Split-Path -Parent $PSScriptRoot
Set-Location $root

New-Item -ItemType Directory -Force -Path "bin" | Out-Null

Write-Host "Building fuji..." -ForegroundColor Cyan
go build -trimpath -ldflags "-s -w" -o "bin/fuji.exe" ./cmd/fuji

Write-Host "Building fujiwrap (cmd/wrapgen)..." -ForegroundColor Cyan
go build -trimpath -ldflags "-s -w" -o "bin/fujiwrap.exe" ./cmd/wrapgen

Write-Host "Building wrapgen.exe (same tool, legacy name)..." -ForegroundColor DarkGray
go build -trimpath -ldflags "-s -w" -o "bin/wrapgen.exe" ./cmd/wrapgen

Write-Host "Building C runtime (runtime/libfuji_runtime.a)..." -ForegroundColor Cyan
& "$PSScriptRoot\build-runtime.ps1"
if ($LASTEXITCODE -ne 0) {
    Write-Host "  runtime build failed (install MinGW gcc/ar or run scripts/build-runtime.ps1 — see CONTRIBUTING.md)" -ForegroundColor Yellow
}

Write-Host ""
Write-Host "Done. Add .\bin to PATH, or run:" -ForegroundColor Green
Write-Host "  .\bin\fuji.exe help" -ForegroundColor Gray
Write-Host "  .\bin\fujiwrap.exe -name mylib -headers .\mylib.h -out .\wrappers\mylib" -ForegroundColor Gray
