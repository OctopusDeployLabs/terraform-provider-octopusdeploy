# Go Builder

This action is adapted from https://github.com/sosedoff/actions/tree/master/golang-build

Github Action to cross-compile Go project binaries for multiple platforms in a single run.

Uses `golang:1.11` Docker image with `CGO_ENABLED=0` flag.

## Usage

Basic usage:

```
action "build" {
  uses = "./actions/build"
}
```

Basic workflow configuration will compile binaries for the following platforms:

- linux: 386/amd64
- darwin: 386/amd64
- windows: 386/amd64 

Alternatively you can provide a list of target architectures in `arg`:

```
action "build" {
  uses = "./actions/build"
  args = "linux/amd64 darwin/amd64"
}
```

Example output:

```
----> Setting up Go repository
----> Building project for: darwin/amd64
  adding: test-go-action_darwin_amd64 (deflated 50%)
----> Building project for: darwin/386
  adding: test-go-action_darwin_386 (deflated 45%)
----> Building project for: linux/amd64
  adding: test-go-action_linux_amd64 (deflated 50%)
----> Building project for: linux/386
  adding: test-go-action_linux_386 (deflated 45%)
----> Building project for: windows/amd64
  adding: test-go-action_windows_amd64 (deflated 50%)
----> Building project for: windows/386
  adding: test-go-action_windows_386 (deflated 46%)
----> Build is complete. List of files at /github/workspace/.release:
total 16436
drwxr-xr-x 2 root root    4096 Feb  5 00:03 .
drwxr-xr-x 5 root root    4096 Feb  5 00:02 ..
-rw-r--r-- 1 root root  978566 Feb  5 00:02 test-go-action_darwin_386.zip
-rw-r--r-- 1 root root 1008819 Feb  5 00:02 test-go-action_darwin_amd64.zip
-rw-r--r-- 1 root root  918555 Feb  5 00:02 test-go-action_linux_386.zip
-rw-r--r-- 1 root root  952985 Feb  5 00:02 test-go-action_linux_amd64.zip
-rw-r--r-- 1 root root  930942 Feb  5 00:03 test-go-action_windows_386.zip
-rw-r--r-- 1 root root  972286 Feb  5 00:02 test-go-action_windows_amd64.zip
```
