# Time-bounded fuji run invocations for Windows CI (no GNU coreutils `timeout` in default PATH).
$ErrorActionPreference = "Stop"
$fuji = $env:FUJI_CI_BIN
if (-not $fuji -or -not (Test-Path -LiteralPath $fuji)) {
    throw "FUJI_CI_BIN must point to an existing fuji executable (got: '$fuji')"
}
$repo = $env:GITHUB_WORKSPACE
if (-not $repo) { $repo = (Get-Location).Path }

function Invoke-FujiTimed([int] $Seconds, [string[]] $Arguments) {
    $p = Start-Process -FilePath $fuji -ArgumentList $Arguments -WorkingDirectory $repo `
        -PassThru -NoNewWindow
    if (-not $p.WaitForExit($Seconds * 1000)) {
        try { Stop-Process -Id $p.Id -Force -ErrorAction SilentlyContinue } catch { }
        throw "timeout after ${Seconds}s: fuji $($Arguments -join ' ')"
    }
    if ($p.ExitCode -ne 0) {
        throw "exit code $($p.ExitCode): fuji $($Arguments -join ' ')"
    }
}

Write-Host "==> GC soak (timed)"
Invoke-FujiTimed 90 @("run", "tests/gc_pressure_expr.fuji")
Invoke-FujiTimed 90 @("run", "tests/globals_perf.fuji")
Invoke-FujiTimed 90 @("run", "tests/gc_soak.fuji")
Invoke-FujiTimed 120 @("run", "--no-opt", "tests/nursery_test.fuji")
Invoke-FujiTimed 120 @("run", "--no-opt", "tests/incremental_gc_test.fuji")

Write-Host "==> Stress smoke (timed)"
Invoke-FujiTimed 90 @("run", "tests/stress/stress_mixed_alloc.fuji")
Invoke-FujiTimed 120 @("run", "--no-opt", "tests/stress/large_game_sim.fuji")

Write-Host "==> GC / stress timed runs OK"
