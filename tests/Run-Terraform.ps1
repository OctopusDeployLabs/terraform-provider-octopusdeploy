$dir = ls build | ? {$_.Name -like "*windows_amd64*"}

terraform init -plugin-dir "./build/$($dir.Name)" tests
terraform apply -auto-approve tests
