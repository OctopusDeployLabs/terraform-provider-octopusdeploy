$dir = ls build | ? {$_.Name -like "*linux_amd64*"}

terraform init -plugin-dir "./build/$($dir.Name)" tests
terraform apply -auto-approve tests
