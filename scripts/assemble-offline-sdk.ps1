# Assemble a self-contained Fuji SDK folder (and optional .zip) matching GitHub Release SDK layout
# for Windows amd64: fuji.exe, fujiwrap.exe, stdlib/, docs/, root *.md, wrappers/, examples/,
# third_party/raylib_static/stage (raylib 5.0 headers + libraylib.a + raylib.dll).
#
# Prerequisites:
#   • Release builds of fuji and fujiwrap (embedded Clang + llc + lld + runtime). Build with:
#       powershell -File scripts/build-release.ps1
#     then pass -FujiExe and -FujiwrapExe, or use -UseBuildRelease to run that script first.
#
# Usage (from repo root):
#   powershell -File scripts/assemble-offline-sdk.ps1 -UseBuildRelease
#   powershell -File scripts/assemble-offline-sdk.ps1 -FujiExe .\fuji-release.exe -FujiwrapExe .\fujiwrap.exe -Zip
#
# GitHub: push a tag v* — .github/workflows/release.yml builds all platforms and runs
# scripts/vendor-raylib-stage.sh + scripts/package-release-sdk.sh (Linux). This script is the
# Windows-local equivalent for testing or ad-hoc distribution.

param(
    [string]$FujiExe = "",
    [string]$FujiwrapExe = "",
    [switch]$UseBuildRelease,
    [string]$Version = "",
    [string]$OutputRoot = "dist",
    [switch]$SkipRaylib,
    [string]$RaylibZipPath = "",
    [switch]$Zip
)

$ErrorActionPreference = "Stop"
$RepoRoot = (Resolve-Path (Join-Path $PSScriptRoot "..")).Path
Set-Location $RepoRoot

function Get-FujiVersionFromSource {
    $mainGo = Join-Path $RepoRoot "cmd\fuji\main.go"
    $t = Get-Content -Raw $mainGo
    $m = [regex]::Match($t, 'var version = "([^"]+)"')
    if (-not $m.Success) { throw "Could not parse version from $mainGo" }
    return $m.Groups[1].Value
}

if ([string]::IsNullOrWhiteSpace($Version)) {
    $Version = Get-FujiVersionFromSource
}

if ($UseBuildRelease) {
    Write-Host "Running scripts/build-release.ps1 ..."
    & (Join-Path $PSScriptRoot "build-release.ps1")
    $FujiExe = Join-Path $RepoRoot "fuji-release.exe"
    $FujiwrapExe = Join-Path $RepoRoot "fujiwrap.exe"
}

if ([string]::IsNullOrWhiteSpace($FujiExe) -or [string]::IsNullOrWhiteSpace($FujiwrapExe)) {
    throw "Specify -FujiExe and -FujiwrapExe (release builds), or -UseBuildRelease."
}

$FujiExe = (Resolve-Path $FujiExe).Path
$FujiwrapExe = (Resolve-Path $FujiwrapExe).Path

$folderName = "fuji-$Version-sdk-windows-amd64"
$stageParent = Join-Path $RepoRoot $OutputRoot
$outDir = Join-Path $stageParent $folderName

if (Test-Path $outDir) {
    Remove-Item -Recurse -Force $outDir
}
New-Item -ItemType Directory -Force -Path $outDir | Out-Null

Copy-Item -Force $FujiExe (Join-Path $outDir "fuji.exe")
Copy-Item -Force $FujiwrapExe (Join-Path $outDir "fujiwrap.exe")

foreach ($d in @("stdlib", "docs", "wrappers", "examples")) {
    $src = Join-Path $RepoRoot $d
    if (Test-Path $src) {
        Copy-Item -Recurse -Force $src (Join-Path $outDir $d)
    }
}

Get-ChildItem -Path $RepoRoot -Filter "*.md" -File | ForEach-Object {
    Copy-Item -Force $_.FullName (Join-Path $outDir $_.Name)
}

$raylibStage = Join-Path $outDir "third_party\raylib_static\stage"
if (-not $SkipRaylib) {
    New-Item -ItemType Directory -Force -Path $raylibStage | Out-Null
    $zipLocal = $false
    $zipFile = $null
    $downloadTmp = $null
    if (-not [string]::IsNullOrWhiteSpace($RaylibZipPath)) {
        $zipFile = (Resolve-Path $RaylibZipPath).Path
        $zipLocal = $true
    }
    if (-not $zipLocal) {
        $raylibVer = "5.0"
        $url = "https://github.com/raysan5/raylib/releases/download/$raylibVer/raylib-${raylibVer}_win64_mingw-w64.zip"
        $downloadTmp = Join-Path ([System.IO.Path]::GetTempPath()) ("fuji-raylib-" + [Guid]::NewGuid().ToString("n"))
        New-Item -ItemType Directory -Force -Path $downloadTmp | Out-Null
        $zipFile = Join-Path $downloadTmp "raylib.zip"
        Write-Host "Downloading raylib $raylibVer Windows prebuild..."
        try {
            [Net.ServicePointManager]::SecurityProtocol = [Net.SecurityProtocolType]::Tls12
            Invoke-WebRequest -Uri $url -OutFile $zipFile -UseBasicParsing
        } catch {
            Remove-Item -Recurse -Force $downloadTmp -ErrorAction SilentlyContinue
            throw "Raylib download failed: $($_.Exception.Message). Use -RaylibZipPath or -SkipRaylib."
        }
    }

    $extract = Join-Path ([System.IO.Path]::GetTempPath()) ("fuji-raylib-extract-" + [Guid]::NewGuid().ToString("n"))
    New-Item -ItemType Directory -Force -Path $extract | Out-Null
    try {
        Expand-Archive -LiteralPath $zipFile -DestinationPath $extract -Force
        $inner = $null
        Get-ChildItem -Path $extract -Directory | ForEach-Object {
            $inc = Join-Path $_.FullName "include"
            $lib = Join-Path $_.FullName "lib"
            if ((Test-Path $inc) -and (Test-Path $lib)) { $inner = $_.FullName }
        }
        if (-not $inner) {
            throw "Raylib zip layout unexpected (expected one top folder with include/ and lib/)."
        }
        Copy-Item -Recurse -Force (Join-Path $inner "include") (Join-Path $raylibStage "include")
        Copy-Item -Recurse -Force (Join-Path $inner "lib") (Join-Path $raylibStage "lib")
        foreach ($x in @("LICENSE", "CHANGELOG", "README.md")) {
            $p = Join-Path $inner $x
            if (Test-Path $p) { Copy-Item -Force $p (Join-Path $raylibStage $x) }
        }
    } finally {
        Remove-Item -Recurse -Force $extract -ErrorAction SilentlyContinue
        if ($downloadTmp) {
            Remove-Item -Recurse -Force $downloadTmp -ErrorAction SilentlyContinue
        }
    }

    $need = @(
        (Join-Path $raylibStage "include\raylib.h"),
        (Join-Path $raylibStage "lib\libraylib.a"),
        (Join-Path $raylibStage "lib\raylib.dll")
    )
    foreach ($n in $need) {
        if (-not (Test-Path $n)) { throw "Raylib stage incomplete after extract: missing $n" }
    }
    Write-Host "Raylib stage OK under third_party\raylib_static\stage"
} else {
    Write-Host "Skipped raylib vendoring (-SkipRaylib). Set FUJI_RAYLIB_STAGE or add stage/ manually for raylib builds."
}

$sdkReadme = @"
Fuji SDK $Version (windows-amd64)
==============================

Offline layout: compiler, fujiwrap (C header -> .fuji + wrapper.c), stdlib, docs, examples,
wrappers (raylib + glue), and (unless -SkipRaylib) third_party/raylib_static/stage with raylib 5.0.

End users: fuji does not download LLVM, Raylib, or anything else at compile time. Embedded Clang/llc/runtime unpack to a local temp directory on first use only.

Use
---
  Keep this folder together so stdlib sits next to fuji.exe:

    .\fuji.exe version
    .\fuji.exe run examples\games\brick_breaker.fuji

  Raylib static linking: when third_party/raylib_static/stage exists next to fuji.exe, the
  compiler adds include + libraylib.a automatically (see docs/guides/raylib.md).

  For distribution of YOUR game: use fuji bundle, and on Windows copy raylib.dll next to the
  .exe if your link uses the DLL (official prebuild ships both .a and .dll).

Wrapper tool
------------
  .\fujiwrap.exe -help
  .\fuji.exe wrap -help

GitHub releases
---------------
  Tag v* on the default branch to run CI: per-OS SDK zips are attached automatically.
  This folder matches that zip contents for Windows amd64.

"@

Set-Content -Path (Join-Path $outDir "SDK_README.txt") -Value $sdkReadme -Encoding UTF8

Write-Host ""
Write-Host "Assembled SDK folder:"
Write-Host "  $outDir"

if ($Zip) {
    $zipOut = Join-Path $stageParent "$folderName.zip"
    if (Test-Path $zipOut) { Remove-Item -Force $zipOut }
    Compress-Archive -Path $outDir -DestinationPath $zipOut -Force
    Write-Host "Wrote zip: $zipOut"
    Write-Host "(Upload this to a GitHub Release asset alongside binaries, or share the folder as-is.)"
}
