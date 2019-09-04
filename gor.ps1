
<#
.SYNOPSIS
  Name: gor.ps1
  The purpose of this script is to run Go actions
  
.DESCRIPTION
  This script gets round the limitiations in PowerShell, when compared to BASH
  that means you have to specific every .go file you want to pass to the 
  Go command.

.PARAMETER Action
  The action you want Go to execute, e.g. build, test or run.

.NOTES
    Updated: 2019-04-13      Initial release.
    Release: 2019-04-13
   
    Author : Simon Buckner

.EXAMPLE
  Run the Go program in the current directory.
  gor -Action run

.EXAMPLE 
  Run the Go program in the current directory.
  gor run

  .EXAMPLE 
  Build the Go program in the current directory.
  gor build
# Comment-based Help tags were introduced in PS 2.0
#requires -version 2
#>

[CmdletBinding()]

PARAM ( 
  [ValidateSet("run", "test", "build")]
  [string]$Action = "run"
)


$go = "go"

$modules = Get-ChildItem -Filter "*.go"

$cmd = "$go $Action"
$params = @($Action)
foreach ($m in $modules) {
  $params += $m.Name
  $cmd += " $($m.Name)"
}

Write-Host $cmd
& $go $params