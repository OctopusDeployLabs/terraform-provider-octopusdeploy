workflow "Publish binaries on release" {
  resolves = ["publish"]
  on = "push"
}

action "release-tag" {
  uses = "actions/bin/filter@master"
  args = "tag"
}

action "build" {
  uses = "sosedoff/actions/golang-build@master"
  env = {
    GO111MODULE = "on"
  }
  needs = ["release-tag"]
}

action "publish" {
  uses = "docker://moonswitch/github-upload-release:master"
  secrets = ["GITHUB_TOKEN"]
  needs = ["build"]
}
