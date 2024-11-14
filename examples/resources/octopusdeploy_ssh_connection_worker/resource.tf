resource "octopusdeploy_ssh_connection_worker" "minimum" {
  name              = "ssh_worker"
  machine_policy_id = "machine-policy-1"
  worker_pools      = ["worker-pools-1"]
  account_id        = "account-42"
  host              = "hostname"
  port              = 22
  fingerprint       = "SHA256: 12345abc"
  dotnet_platform   = "linux-x64"
}

resource "octopusdeploy_ssh_connection_worker" "optionals" {
  name              = "optional_ssh_worker"
  machine_policy_id = "machine-policy-1"
  worker_pools      = ["worker-pools-1", "worker-pools-2"]
  account_id        = "account-42"
  host              = "hostname"
  port              = 22
  fingerprint       = "SHA256: 12345abc"
  dotnet_platform   = "linux-x64"
  proxy_id          = "proxy-31"
  is_disabled       = true
}
