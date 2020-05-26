param (
    [string]$operatingSystem = "linux_amd64"
)

Write-Host (ls build)

$dir = ls build | ? {$_.Name -like "*$operatingSystem*"}

Write-Host "plugin dir is set to $($dir.Name)"

terraform init -plugin-dir "./build/$($dir.Name)" tests
terraform apply -auto-approve tests
