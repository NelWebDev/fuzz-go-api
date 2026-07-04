[CmdletBinding()]
param(
    [string]$Config = "config/load.json",
    [string]$ArtifactsRoot = "artifacts/load",
    [int]$Users = 2,
    [string]$Duration = "",
    [int]$Requests = 10,
    [string]$Scenario = "",
    [int]$KeepRuns = 10
)

$ErrorActionPreference = "Stop"
Set-StrictMode -Version Latest

$repoRoot = Resolve-Path (Join-Path $PSScriptRoot "..")
$timestamp = Get-Date -Format "yyyy-MM-dd_HH-mm-ss"
$artifactsRootPath = Join-Path $repoRoot $ArtifactsRoot
$runDir = Join-Path $artifactsRootPath $timestamp
$consoleLogPath = Join-Path $runDir "load-console.log"

New-Item -ItemType Directory -Force -Path $runDir | Out-Null

$previousGoCache = $env:GOCACHE

try {
    Push-Location $repoRoot

    $env:GOCACHE = Join-Path $repoRoot ".gocache"

    $argsList = @(
        "run",
        "./cmd/loadtest",
        "-config", $Config,
        "-artifacts", $runDir,
        "-users", $Users.ToString(),
        "-requests", $Requests.ToString()
    )

    if ($Duration -ne "") {
        $argsList += @("-duration", $Duration)
    }
    if ($Scenario -ne "") {
        $argsList += @("-scenario", $Scenario)
    }

    @(
        "Load run: $timestamp",
        "Repository: $repoRoot",
        "Config: $Config",
        "Users: $Users",
        "Requests: $Requests",
        "Duration: $Duration",
        "Scenario: $Scenario",
        ""
    ) | Set-Content -Path $consoleLogPath

    & go @argsList *>&1 | Tee-Object -FilePath $consoleLogPath -Append
    $exitCode = $LASTEXITCODE

    if ($KeepRuns -gt 0 -and (Test-Path $artifactsRootPath)) {
        Get-ChildItem -Path $artifactsRootPath -Directory |
            Where-Object { $_.Name -match '^\d{4}-\d{2}-\d{2}_\d{2}-\d{2}-\d{2}$' } |
            Sort-Object LastWriteTime -Descending |
            Select-Object -Skip $KeepRuns |
            Remove-Item -Recurse -Force
    }

    Write-Host "Load artifacts saved in $runDir"
    exit $exitCode
}
finally {
    Pop-Location
    $env:GOCACHE = $previousGoCache
}
