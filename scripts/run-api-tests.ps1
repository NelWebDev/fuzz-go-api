[CmdletBinding()]
param(
    [string]$FuzzTime = "30s",
    [int]$LoadUsers = 2,
    [int]$LoadRequests = 10,
    [string]$LoadDuration = "",
    [string]$ArtifactsRoot = "artifacts/api-tests",
    [int]$KeepRuns = 10,
    [switch]$QuietFuzzRequests
)

$ErrorActionPreference = "Stop"
Set-StrictMode -Version Latest

$repoRoot = Resolve-Path (Join-Path $PSScriptRoot "..")
$timestamp = Get-Date -Format "yyyy-MM-dd_HH-mm-ss"
$combinedRoot = Join-Path $repoRoot $ArtifactsRoot
$runDir = Join-Path $combinedRoot $timestamp
$summaryPath = Join-Path $runDir "summary.txt"
$relativeFuzzRoot = Join-Path $ArtifactsRoot "$timestamp/fuzz"
$relativeLoadRoot = Join-Path $ArtifactsRoot "$timestamp/load"
$results = @()

New-Item -ItemType Directory -Force -Path $runDir | Out-Null

try {
    Push-Location $repoRoot

    @(
        "Combined API test run: $timestamp",
        "Repository: $repoRoot",
        "Fuzz time: $FuzzTime",
        "Load users: $LoadUsers",
        "Load requests: $LoadRequests",
        "Load duration: $LoadDuration",
        ""
    ) | Set-Content -Path $summaryPath

    $powerShellExe = (Get-Process -Id $PID).Path

    $fuzzArgs = @("-NoProfile", "-ExecutionPolicy", "Bypass", "-File", "$PSScriptRoot\run-fuzz.ps1", "-FuzzTime", $FuzzTime, "-ArtifactsRoot", $relativeFuzzRoot, "-KeepRuns", "0")
    if ($QuietFuzzRequests) {
        $fuzzArgs += "-QuietRequests"
    }
    & $powerShellExe @fuzzArgs
    $fuzzExit = $LASTEXITCODE
    $results += [pscustomobject]@{ Name = "fuzz"; ExitCode = $fuzzExit }

    $loadArgs = @("-NoProfile", "-ExecutionPolicy", "Bypass", "-File", "$PSScriptRoot\run-load.ps1", "-ArtifactsRoot", $relativeLoadRoot, "-Users", $LoadUsers.ToString(), "-Requests", $LoadRequests.ToString(), "-KeepRuns", "0")
    if ($LoadDuration -ne "") {
        $loadArgs += @("-Duration", $LoadDuration)
    }
    & $powerShellExe @loadArgs
    $loadExit = $LASTEXITCODE
    $results += [pscustomobject]@{ Name = "load"; ExitCode = $loadExit }

    Add-Content -Path $summaryPath -Value "Results:"
    foreach ($result in $results) {
        $status = if ($result.ExitCode -eq 0) { "PASS" } else { "FAIL" }
        Add-Content -Path $summaryPath -Value ("- {0}: {1} (exit {2})" -f $result.Name, $status, $result.ExitCode)
    }

    if ($KeepRuns -gt 0 -and (Test-Path $combinedRoot)) {
        Get-ChildItem -Path $combinedRoot -Directory |
            Where-Object { $_.Name -match '^\d{4}-\d{2}-\d{2}_\d{2}-\d{2}-\d{2}$' } |
            Sort-Object LastWriteTime -Descending |
            Select-Object -Skip $KeepRuns |
            Remove-Item -Recurse -Force
    }

    Write-Host "Combined API test artifacts saved in $runDir"
    if (($results | Where-Object { $_.ExitCode -ne 0 }).Count -gt 0) {
        exit 1
    }
}
finally {
    Pop-Location
}
