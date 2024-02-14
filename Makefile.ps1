#!/usr/bin/env pwsh

$env:NAME = "tf"

## Read .env
if (Test-Path -Path ".env" -PathType Leaf) {
  Get-Content -Path ".env" | ForEach-Object {
    $key, $value = $_ -split '=', 2
    Set-Item -Path "env:$key" -Value $value
  }
}

# Parse command line arguments
$targets = @()
$currentTarget = ""

foreach ($arg in $args) {
  if ($arg -match '=') {
    ## Set variable
    $key, $value = $arg -split '=', 2
    Set-Item -Path "env:$key" -Value $value
  }
  else {
    ## Remember the target
    $targets += $arg
  }
}

if (-not $env:GO) { $env:GO = "go" }
if (-not $env:GORELEASER) { $env:GORELEASER = "goreleaser" }

if ($env:GOOS) {
  $GOOS = $env:GOOS
} else {
  $GOOS = & $env:GO env GOOS
}

if ($GOOS -eq "windows") {
  $env:EXE = ".exe"
}
else {
  $env:EXE = ""
}

if (-not $env:BIN) { $env:BIN = "$env:NAME$env:EXE" }

if ($env:OS -eq "Windows_NT") {
  if ($env:LOCALAPPDATA) {
    if (-not $env:BINDIR) { $env:BINDIR = "$env:LOCALAPPDATA\Microsoft\WindowsApps" }
  }
  else {
    if (-not $env:BINDIR) { $env:BINDIR = "C:\Windows\System32" }
  }
}
else {
  if (Test-Path "$env:HOME\.local\bin") {
    if (-not $env:BINDIR) { $env:BINDIR = "$env:HOME\.local\bin" }
  }
  elseif (Test-Path "$env:HOME\bin") {
    if (-not $env:BINDIR) { $env:BINDIR = "$env:HOME\bin" }
  }
  else {
    if (-not $env:BINDIR) { $env:BINDIR = "/usr/local/bin" }
  }
}

function Invoke-CommandWithEcho {
  param (
    [string]$Command,
    [string[]]$Arguments
  )
  Write-Host $Command $Arguments
  $processInfo = Start-Process -FilePath $Command -ArgumentList $Arguments -NoNewWindow -PassThru
  $processInfo.WaitForExit()
  if ($processInfo.ExitCode -and $processInfo.ExitCode -ne 0) {
    Write-Host "make: *** [$currentTarget] Error $($processInfo.ExitCode)"
    break
  }
}

function Invoke-ExpressionWithEcho {
  param (
    [string]$Command
  )
  Write-Host $Command
  Invoke-Expression -Command $Command
}

function Write-Target {
  param (
    [string]$Target
  )
  Write-Host "Executing target: $Target"
}

## TARGET build Build app binary for single target
function Invoke-Target-Build {
  Write-Target "build"
  # 	$(GORELEASER) build --clean --snapshot --single-target --output $(BIN)

  Invoke-CommandWithEcho $env:GORELEASER -Arguments "build", "--clean", "--snapshot", "--single-target", "--output", $env:BIN
}

## TARGET install Build and install app binary
function Invoke-Target-Install {
  if (-not (Test-Path -Path "$env:BIN" -PathType Leaf)) {
    Invoke-Target-Build
  }
  Write-Target "install"
  Invoke-ExpressionWithEcho "Copy-Item -Path '$env:BIN' -Destination $env:BINDIR -Force"
}

## TARGET uninstall Uninstall app binary
function Invoke-Target-Uninstall {
  Write-Target "uninstall"
  $path = Join-Path $env:BINDIR "$env:BIN"
  Invoke-ExpressionWithEcho -Command "Remove-Item $path -Force -ErrorAction SilentlyContinue"
}

## TARGET download Download Go modules
function Invoke-Target-Download {
  Write-Target "download"
  Invoke-CommandWithEcho $env:GO -Arguments "mod", "download"
}

## TARGET tidy Tidy Go modules
function Invoke-Target-Tidy {
  Write-Target "tidy"
  Invoke-CommandWithEcho $env:GO -Arguments "mod", "tidy"
}

## TARGET upgrade Upgrade Go modules
function Invoke-Target-Upgrade {
  Write-Target "upgrade"
  Invoke-CommandWithEcho $env:GO -Arguments "get", "-u"
}

## TARGET clean Clean working directory
function Invoke-Target-Clean {
  Write-Target "clean"
  Invoke-ExpressionWithEcho -Command "Remove-Item '$env:BIN' -Force -ErrorAction SilentlyContinue"
  Invoke-ExpressionWithEcho -Command "Remove-Item dist -Recurse -Force -ErrorAction SilentlyContinue"
}

function Invoke-Target-Help {
  Write-Host "Targets:"
  Get-Content $PSCommandPath |
  Select-String '^## TARGET ' |
  Sort-Object |
  ForEach-Object {
    $target, $description = $_ -split ' ', 4 | Select-Object -Skip 2
    Write-Output ("  {0,-20} {1}" -f $target, $description)
  }
}

## Run target
if ($targets.Count -eq 0) {
  & Invoke-Target-Help
}
else {
  foreach ($target in $targets) {
    $currentTarget = $target
    Invoke-Expression ("Invoke-Target-" + $target)
  }
}
