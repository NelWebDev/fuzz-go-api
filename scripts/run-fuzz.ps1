[CmdletBinding()]
param(
    [string]$FuzzTime = "30s",
    [string]$ArtifactsRoot = "artifacts",
    [int]$KeepRuns = 10,
    [int]$MaxLogBytes = 8192,
    [string[]]$Targets = @(
        "FuzzGetEndpoint",
        "FuzzPostEndpoint",
        "FuzzPutEndpoint",
        "FuzzPatchEndpoint",
        "FuzzDeleteEndpoint"
    ),
    [switch]$StopOnFailure,
    [switch]$QuietRequests
)

$ErrorActionPreference = "Stop"
Set-StrictMode -Version Latest

$repoRoot = Resolve-Path (Join-Path $PSScriptRoot "..")
$timestamp = Get-Date -Format "yyyy-MM-dd_HH-mm-ss"
$artifactsRootPath = Join-Path $repoRoot $ArtifactsRoot
$runDir = Join-Path $artifactsRootPath $timestamp
$findingsPath = Join-Path $runDir "fuzz-findings.jsonl"
$summaryPath = Join-Path $runDir "summary.txt"

New-Item -ItemType Directory -Force -Path $runDir | Out-Null

$previousExternal = $env:FUZZ_API_EXTERNAL
$previousFindings = $env:FUZZ_API_FINDINGS
$previousMaxLogBytes = $env:FUZZ_API_MAX_LOG_BYTES
$previousLogRequests = $env:FUZZ_API_LOG_REQUESTS
$results = @()

try {
    Push-Location $repoRoot

    $env:FUZZ_API_EXTERNAL = "1"
    $env:FUZZ_API_FINDINGS = $findingsPath
    $env:FUZZ_API_MAX_LOG_BYTES = $MaxLogBytes.ToString()
    $env:FUZZ_API_LOG_REQUESTS = if ($QuietRequests) { "0" } else { "1" }

    @(
        "Fuzz run: $timestamp",
        "Repository: $repoRoot",
        "Fuzz time: $FuzzTime",
        "Findings: $findingsPath",
        "Max log bytes: $MaxLogBytes",
        "Request logs: $(if ($QuietRequests) { 'disabled' } else { 'enabled' })",
        "Targets: $($Targets -join ', ')",
        ""
    ) | Set-Content -Path $summaryPath

    foreach ($target in $Targets) {
        $logPath = Join-Path $runDir "$target.log"
        $commandText = "go test -run=^$ -fuzz=$target -fuzztime=$FuzzTime ./fuzz"

        Write-Host "Running $target for $FuzzTime"
        @(
            "Target: $target",
            "Command: $commandText",
            "Started: $(Get-Date -Format o)",
            ""
        ) | Set-Content -Path $logPath

        & go test "-run=^$" "-fuzz=$target" "-fuzztime=$FuzzTime" "./fuzz" *>&1 |
            Tee-Object -FilePath $logPath -Append

        $exitCode = $LASTEXITCODE
        $status = if ($exitCode -eq 0) { "PASS" } else { "FAIL" }
        $results += [pscustomobject]@{
            Target = $target
            Status = $status
            ExitCode = $exitCode
            Log = $logPath
        }

        Add-Content -Path $logPath -Value @(
            "",
            "Finished: $(Get-Date -Format o)",
            "Exit code: $exitCode"
        )

        if ($exitCode -ne 0 -and $StopOnFailure) {
            break
        }
    }

    Add-Content -Path $summaryPath -Value "Results:"
    foreach ($result in $results) {
        Add-Content -Path $summaryPath -Value ("- {0}: {1} (exit {2})" -f $result.Target, $result.Status, $result.ExitCode)
    }

    if (Test-Path $findingsPath) {
        $findingCount = (Get-Content -Path $findingsPath | Measure-Object -Line).Lines
    } else {
        $findingCount = 0
        New-Item -ItemType File -Path $findingsPath | Out-Null
    }
    Add-Content -Path $summaryPath -Value ""
    Add-Content -Path $summaryPath -Value "Findings: $findingCount"

    if ($KeepRuns -gt 0 -and (Test-Path $artifactsRootPath)) {
        Get-ChildItem -Path $artifactsRootPath -Directory |
            Where-Object { $_.Name -match '^\d{4}-\d{2}-\d{2}_\d{2}-\d{2}-\d{2}$' } |
            Sort-Object LastWriteTime -Descending |
            Select-Object -Skip $KeepRuns |
            Remove-Item -Recurse -Force
    }

    $failed = @($results | Where-Object { $_.ExitCode -ne 0 })
    Write-Host "Artifacts saved in $runDir"
    Write-Host "Findings written to $findingsPath"

    if ($failed.Count -gt 0) {
        exit 1
    }
}
finally {
    Pop-Location
    $env:FUZZ_API_EXTERNAL = $previousExternal
    $env:FUZZ_API_FINDINGS = $previousFindings
    $env:FUZZ_API_MAX_LOG_BYTES = $previousMaxLogBytes
    $env:FUZZ_API_LOG_REQUESTS = $previousLogRequests
}
