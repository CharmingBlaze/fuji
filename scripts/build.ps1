# Build kuji + kujiwrap (wrapper generator) into ./bin — one script for Windows authors.
$ErrorActionPreference = "Stop"
$root = Split-Path -Parent $PSScriptRoot
Set-Location $root

New-Item -ItemType Directory -Force -Path "bin" | Out-Null

Write-Host "Building kuji..." -ForegroundColor Cyan
go build -trimpath -ldflags "-s -w" -o "bin/kuji.exe" ./cmd/kuji

Write-Host "Building kujiwrap (from cmd/wrapgen)..." -ForegroundColor Cyan
go build -trimpath -ldflags "-s -w" -o "bin/kujiwrap.exe" ./cmd/wrapgen

Write-Host "Building C runtime (runtime/libfuji_runtime.a)..." -ForegroundColor Cyan
& "$PSScriptRoot\build-runtime.ps1"
if ($LASTEXITCODE -ne 0) {
    Write-Host "  runtime build failed (install MinGW gcc/ar or run scripts/build-runtime.ps1 — see CONTRIBUTING.md)" -ForegroundColor Yellow
}

Write-Host ""
Write-Host "Done. Add .\bin to PATH, or run:" -ForegroundColor Green
Write-Host "  .\bin\kuji.exe help" -ForegroundColor Gray
Write-Host "  .\bin\kujiwrap.exe -name mylib -headers .\mylib.h -out .\wrappers\mylib" -ForegroundColor Gray
