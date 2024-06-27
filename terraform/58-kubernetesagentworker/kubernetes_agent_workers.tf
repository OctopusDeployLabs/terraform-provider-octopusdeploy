resource "octopusdeploy_kubernetes_agent_worker" "agent_with_minimum" {
  name         = "minimum-agent"
  worker_pool_ids = [octopusdeploy_static_worker_pool.workerpool_docker.id]
  uri          = "poll://kcxzcv2fpsxkn6tk9u6d/"
  thumbprint   = "96203ED84246201C26A2F4360D7CBC36AC1D232D"
}

resource "octopusdeploy_kubernetes_agent_worker" "agent_with_optional" {
  name               = "optional-agent"
  machine_policy_id  = octopusdeploy_machine_policy.machinepolicy_testing.id
  worker_pool_ids    = [octopusdeploy_static_worker_pool.workerpool_docker.id]
  communication_mode = "Polling"
  uri                = "poll://kcxzcv2fpsxkn6tk9u6d/"
  thumbprint         = "96203ED84246201C26A2F4360D7CBC36AC1D232D"
  is_disabled        = true
  upgrade_locked     = false
}