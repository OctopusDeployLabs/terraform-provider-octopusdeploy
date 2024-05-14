resource "octopusdeploy_kubernetes_agent_deployment_target" "minimal" {
  name         = "agent-minimal"
  environments = ["environments-1"]
  roles        = ["role-1", "role-2"]
  thumbprint   = "96203ED84246201C26A2F4360D7CBC36AC1D232D"
  uri          = "poll://kcxzcv2fpsxkn6tk9u6d/"
}

resource "octopusdeploy_kubernetes_agent_deployment_target" "optionals" {
  name              = "agent-optionals"
  environments      = ["environments-1"]
  roles             = ["role-1", "role-2"]
  thumbprint   = "96203ED84246201C26A2F4360D7CBC36AC1D232D"
  uri          = "poll://kcxzcv2fpsxkn6tk9u6d/"
  machine_policy_id = "machinepolicies-1"
  default_namespace = "kubernetes-namespace"
  upgrade_locked    = true
  is_disabled       = true
}

resource "octopusdeploy_kubernetes_agent_deployment_target" "tenanted_agent" {
  name                              = "agent-tenanted"
  environments                      = ["environments-1"]
  roles                             = ["role-1", "role-2"]
  thumbprint   = "96203ED84246201C26A2F4360D7CBC36AC1D232D"
  uri          = "poll://kcxzcv2fpsxkn6tk9u6d/"
  tenanted_deployment_participation = "Tenanted"
  tenants                           = ["tenants-1"]
  tenant_tags                       = ["TagSets-1/Tags-1"]
}