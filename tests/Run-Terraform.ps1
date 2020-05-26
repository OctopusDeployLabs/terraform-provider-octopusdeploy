param (
    [string]$operatingSystem = "linux_amd64"
)

Write-Host "finding plugin for $operatingSystem"

Write-Host "The contents of the build directory:"
Write-Host (ls build)

$dir = ls build | ? {$_.Name -like "*$operatingSystem*"}
Write-Host "Plugin dir is set to $($dir.Name)"

terraform init -plugin-dir "./build/$($dir.Name)" tests
terraform apply -auto-approve tests
