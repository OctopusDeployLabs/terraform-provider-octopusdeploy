param (
    [string]$operatingSystem = "linux_amd64"
)

Write-Host "Finding plugin for $operatingSystem"

$dir = Get-Childitem build | ? {$_.Name -like "*$operatingSystem*"} | Select -First 1
Write-Host "Plugin dir is set to $($dir.Name)"

terraform init -plugin-dir "./build/$($dir.Name)" tests
terraform apply -auto-approve tests
