TEST?=$$(go list ./... |grep -v 'vendor')

default: build test

fmt:
	go fmt ./octopusdeploy/...

build: fmt
	go build

test: fmt
	go test -v -timeout 30s ./...

testacc:
	TF_ACC=1 go test $(TEST) -v -timeout 120m

tf_build: fmt build
	terraform init
	terraform plan
	terraform apply -auto-approve

tf_destroy:
	terraform destroy -auto-approve

tf_apply:
	terraform apply -auto-approve

tf_plan:
	terraform plan
