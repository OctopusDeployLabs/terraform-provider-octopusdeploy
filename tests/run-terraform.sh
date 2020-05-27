DEFAULTARCH=linux_amd64
DEFAULTHOST=http://localhost:8080

ARCH=${1:-$DEFAULTARCH}
HOST=${2:-$DEFAULTHOST}

echo "Finding plugin for ${ARCH}"

DIR=$(find ./build -maxdepth 1 -name "*${ARCH}*" -type d)

echo "Plugin dir is set to ${DIR}"

cat <<EOF > examples/provider.tf
  provider "octopusdeploy" {
  address = "${HOST}"
  apikey  = "${OCTOPUS_APIKEY}"
}
EOF

terraform init -plugin-dir ${DIR} examples
terraform apply -auto-approve examples
