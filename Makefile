TEST?=$$(go list ./... | grep -v 'vendor')
HOSTNAME=octopus.com
NAMESPACE=com
NAME=octopusdeploy
BINARY=terraform-provider-${NAME}
VERSION=0.7.102

ifeq ($(OS), Windows_NT)
	OS_ARCH?=windows_386
	PROFILE=${APPDATA}/terraform.d
	EXT=.exe
else
 	PROFILE=~/.terraform.d
	UNAME_S := $(shell uname -s)
	ifeq ($(UNAME_S), Linux)
		OS_ARCH?=linux_amd64
	else
		UNAME_P := $(shell uname -p)
		ifeq ($(UNAME_P), arm)
			OS_ARCH?=darwin_arm64
		else
			OS_ARCH?=darwin_amd64
		endif
	endif
endif

.PHONY: default
default: install

.PHONY: build
build:
	go build -o ${BINARY}${EXT}

.PHONY: docs
docs:
	go generate main.go

.PHONY: release
release:
	GOOS=darwin GOARCH=amd64 go build -o ./bin/${BINARY}_${VERSION}_darwin_amd64
	GOOS=darwin GOARCH=arm64 go build -o ./bin/${BINARY}_${VERSION}_darwin_arm64
	GOOS=freebsd GOARCH=386 go build -o ./bin/${BINARY}_${VERSION}_freebsd_386
	GOOS=freebsd GOARCH=amd64 go build -o ./bin/${BINARY}_${VERSION}_freebsd_amd64
	GOOS=freebsd GOARCH=arm go build -o ./bin/${BINARY}_${VERSION}_freebsd_arm
	GOOS=linux GOARCH=386 go build -o ./bin/${BINARY}_${VERSION}_linux_386
	GOOS=linux GOARCH=amd64 go build -o ./bin/${BINARY}_${VERSION}_linux_amd64
	GOOS=linux GOARCH=arm go build -o ./bin/${BINARY}_${VERSION}_linux_arm
	GOOS=openbsd GOARCH=386 go build -o ./bin/${BINARY}_${VERSION}_openbsd_386
	GOOS=openbsd GOARCH=amd64 go build -o ./bin/${BINARY}_${VERSION}_openbsd_amd64
	GOOS=solaris GOARCH=amd64 go build -o ./bin/${BINARY}_${VERSION}_solaris_amd64
	GOOS=windows GOARCH=386 go build -o ./bin/${BINARY}_${VERSION}_windows_386
	GOOS=windows GOARCH=amd64 go build -o ./bin/${BINARY}_${VERSION}_windows_amd64

.PHONY: install
install: build
	mkdir -p $(PROFILE)/plugins/${HOSTNAME}/${NAMESPACE}/${NAME}/${VERSION}/${OS_ARCH}
	mv ${BINARY}${EXT} $(PROFILE)/plugins/${HOSTNAME}/${NAMESPACE}/${NAME}/${VERSION}/${OS_ARCH}

.PHONY: test
test:
	go test -i $(TEST) || exit 1
	echo $(TEST) | xargs -t -n4 go test $(TESTARGS) -timeout=30s -parallel=4

.PHONY: testacc
testacc:
	TF_ACC=1 go test $(TEST) -v $(TESTARGS) -timeout 120m