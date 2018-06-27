TEST?=$$(go list ./... |grep -v 'vendor')

default: build test

fmt:
	go fmt github.com/MattHodge/terraform-provider-octopusdeploy/...

build: fmt
	go build

test: fmt
	go test -v -timeout 30s ./...

testacc:
	TF_ACC=1 go test $(TEST) -v -timeout 120m
