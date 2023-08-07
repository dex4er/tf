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

if (-not $env:DOCKER) { $env:DOCKER = "docker" }
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

function Get-Version {
  try {
    $exactMatch = git describe --tags --exact-match 2>$null
    if (-not [string]::IsNullOrEmpty($exactMatch)) { 
      $version = $exactMatch
    }
    else {
      $tags = git describe --tags 2>$null; 
      if ([string]::IsNullOrEmpty($tags)) { 
        $commitHash = (git rev-parse --short=8 HEAD).Trim();
        $version = "0.0.0-0-g$commitHash" 
      }
      else { 
        $version = $tags -replace '-[0-9][0-9]*-g', '-SNAPSHOT-' 
      }
    }
    $version = $version -replace '^v', ''
    return $version 
  }
  catch { 
    return "0.0.0" 
  }
}

function Get-Revision {
  $revision = git rev-parse HEAD
  return $revision
}

function Get-Builddate {
  $datetime = Get-Date
  $utc = $datetime.ToUniversalTime()
  return $utc.tostring("yyyy-MM-ddTHH:mm:ssZ")
}

if (-not $env:VERSION) { $env:VERSION = (& get-version) }
if (-not $env:REVISION) { $env:REVISION = (& get-revision) }
if (-not $env:BUILDDATE) { $env:BUILDDATE = (& get-builddate) }

if (-not $env:CGO_ENABLED) { $env:CGO_ENABLED = "0" }

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
  Invoke-CommandWithEcho $env:GO -Arguments "build", "-trimpath", "-ldflags=`"-s -w -X main.version=$env:VERSION`""
}

## TARGET goreleaser Build app binary for all targets
function Invoke-Target-Goreleaser {
  Write-Target "goreleaser"
  Invoke-CommandWithEcho $env:GORELEASER -Arguments "release", "--auto-snapshot", "--clean", "--skip-publish"
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

## TARGET version Show version
function Invoke-Target-Version {
  Write-Host $env:VERSION
}

## TARGET revision Show revision
function Invoke-Target-Revision {
  Write-Host $env:REVISION
}

## TARGET builddate Show build date
function Invoke-Target-Builddate {
  Write-Host $env:BUILDDATE
}

if (-not $env:DOCKERFILE) { $env:DOCKERFILE = "Dockerfile" }
if (-not $env:IMAGE_NAME) { $env:IMAGE_NAME = $env:BIN }
if (-not $env:LOCAL_REPO) { $env:LOCAL_REPO = "localhost:5000/$env:IMAGE_NAME" }
if (-not $env:DOCKER_REPO) { $env:DOCKER_REPO = "localhost:5000/$env:IMAGE_NAME" }

if (-not $env:PLATFORM) {
  if ($env:PROCESSOR_ARCHITECTURE -eq "ARM64") {
    $env:PLATFORM = "linux/arm64"
  }
  elseif ((uname -m) -eq "arm64") {
    $env:PLATFORM = "linux/arm64"
  }
  elseif ((uname -m) -eq "aarch64") {
    $env:PLATFORM = "linux/arm64"
  }
  elseif ((uname -s) -match "ARM64") {
    $env:PLATFORM = "linux/arm64"
  }
  else {
    $env:PLATFORM = "linux/amd64"
  }
}

## TARGET image Build a local image without publishing artifacts.
function Invoke-Target-Image {
  $env:GOOS = "linux"
  Invoke-Target-Build
  Write-Target "image"
  Invoke-CommandWithEcho $env:DOCKER -Arguments "buildx", "build", "--file=$env:DOCKERFILE",
  "--platform=$env:PLATFORM",
  "--build-arg", "VERSION=$env:VERSION",
  "--build-arg", "REVISION=$env:REVISION",
  "--build-arg", "BUILDDATE=$env:BUILDDATE",
  "--tag", $env:LOCAL_REPO,
  "--load",
  "."
}

## TARGET push Publish to container registry.
function Invoke-Target-Push {
  Write-Target "push"
  Invoke-CommandWithEcho $env:DOCKER -Arguments "tag", $env:LOCAL_REPO, "$($env:DOCKER_REPO):v$($env:VERSION)-$($env:PLATFORM -replace '/','-')"
  Invoke-CommandWithEcho $env:DOCKER -Arguments "push", "$($env:DOCKER_REPO):v$($env:VERSION)-$($env:PLATFORM -replace '/','-')"
}

## TARGET test-image Test local image
function Invoke-Target-Test-Image {
  Write-Target "test-image"
  Invoke-CommandWithEcho $env:DOCKER -Arguments "run", "--platform=$env:PLATFORM", "--rm", "-t", $env:LOCAL_REPO, "-v"
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
