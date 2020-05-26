param (
    [string]$operatingSystem = "linux_amd64"
)

$dir = ls build | ? {$_.Name -like "*$operatingSystem*"}

terraform init -plugin-dir "./build/$($dir.Name)" tests
terraform apply -auto-approve tests
