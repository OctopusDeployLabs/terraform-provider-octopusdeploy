resource "octopusdeploy_listening_tentacle_worker" "minimum" {
  name         = "listening_worker"
  machine_policy_id = "machine-policy-1"
  worker_pools = ["worker-pools-1", "worker-pools-2"]
  thumbprint   = "96203ED84246201C26A2F4360D7CBC36AC1D232D"
  uri          = "https://tentacle.listening/"
}

resource "octopusdeploy_listening_tentacle_worker" "optionals" {
  name              = "optional_worker"
  machine_policy_id = "machine-policy-1"
  worker_pools      = ["worker-pools-1"]
  thumbprint        = "96203ED84246201C26A2F4360D7CBC36AC1D232D"
  uri               = "https://tentacle.listening/"
  proxy_id          = "proxys-1"
  is_disabled       = true
}
