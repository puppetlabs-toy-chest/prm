#!/usr/bin/env pwsh

[CmdletBinding()]
param (
  [Parameter()]
  [ValidateSet('build', 'quick', 'package')]
  [string]
  $Target = 'build'
)
$Env:WORKINGDIR = $PSScriptRoot

$arch = go env GOHOSTARCH
$platform = go env GOHOSTOS
$binPath = Join-Path $PSScriptRoot 'dist' "prm_${platform}_${arch}"
$binPath2 = Join-Path $PSScriptRoot 'dist' "notel_prm_${platform}_${arch}"

switch ($Target) {
  'build' {
    # Set goreleaser to build for current platform only
    # Add environment variables for honeycomb if not already loaded
    if (!(Test-Path ENV:\HONEYCOMB_API_KEY)) {
      $ENV:HONEYCOMB_API_KEY = 'not_set'
    }
    if (!(Test-Path ENV:\HONEYCOMB_DATASET)) {
      $ENV:HONEYCOMB_DATASET = 'not_set'
    }
    goreleaser build --snapshot --rm-dist --single-target
    git clone -b main --depth 1 --single-branch https://github.com/puppetlabs/fantastic-adventure (Join-Path $binPath 'prm-tools')
    Copy-Item (Join-Path $binPath 'prm-tools') -Destination (Join-Path $binPath2 'prm-tools') -Recurse
  }
  'quick' {
    If ($Env:OS -match '^Windows') {
      go build -o "$binPath/prm.exe" -tags telemetry
      go build -o "$binPath2/prm.exe"
    }
    else {
      go build -o "$binPath/prm" -tags telemetry
      go build -o "$binPath2/prm"
    }
  }
  'package' {
    git clone -b main --depth 1 --single-branch https://github.com/puppetlabs/fantastic-adventure  'prm-tools'
    goreleaser --skip-publish --snapshot --rm-dist
  }
}
