
$go = "go"

$modules = Get-ChildItem -Filter "*.go"

$cmd = "$go"
$params = @("build")
foreach ($m in $modules) {
    $params += $m.Name
    $cmd += " $($m.Name)"
}

Write-Host "$cmd"
& $go $params