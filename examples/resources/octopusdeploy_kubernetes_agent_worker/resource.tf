resource "octopusdeploy_kubernetes_agent_worker" "minimal" {
  name         = "agent-minimal"
  worker_pools = ["worker-pools-1"]
  thumbprint   = "96203ED84246201C26A2F4360D7CBC36AC1D232D"
  uri          = "poll://kcxzcv2fpsxkn6tk9u6d/"
}

resource "octopusdeploy_kubernetes_agent_worker" "optionals" {
  name         = "agent-optionals"
  worker_pools = ["worker-pools-1", "worker-pools-3"]
  thumbprint   = "96203ED84246201C26A2F4360D7CBC36AC1D232D"
  uri          = "poll://kcxzcv2fpsxkn6tk9u6d/"
  machine_policy_id = "machinepolicies-1"
  upgrade_locked    = true
  is_disabled       = true
}