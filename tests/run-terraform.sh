echo "Finding plugin for ${1}"

DIR=$(find ./build -maxdepth 1 -name "*$1*" -type d)

echo "Plugin dir is set to ${DIR}"

cat <<EOF > examples/provider.tf
  provider "octopusdeploy" {
  address = "http://localhost:8080/"
  apikey  = "${OCTOPUS_APIKEY}"
}
EOF

terraform init -plugin-dir ${DIR} examples
terraform apply -auto-approve examples
