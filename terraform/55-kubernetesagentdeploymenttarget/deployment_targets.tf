resource "octopusdeploy_kubernetes_agent_deployment_target" "agent_with_minimum" {
  name         = "minimum-agent"
  environments = [octopusdeploy_environment.development_environment.id]
  roles        = ["role-1", "role-2"]
  uri          = "poll://kcxzcv2fpsxkn6tk9u6d/"
  thumbprint   = "96203ED84246201C26A2F4360D7CBC36AC1D232D"
}

resource "octopusdeploy_kubernetes_agent_deployment_target" "agent_with_optional" {
  name               = "optional-agent"
  environments       = [octopusdeploy_environment.development_environment.id]
  roles              = ["role-1", "role-2"]
  machine_policy_id  = octopusdeploy_machine_policy.machinepolicy_testing.id
  communication_mode = "Polling"
  uri                = "poll://kcxzcv2fpsxkn6tk9u6d/"
  thumbprint         = "96203ED84246201C26A2F4360D7CBC36AC1D232D"
  default_namespace  = "kubernetes-namespace"
  is_disabled        = true
  upgrade_locked     = true
}

resource "octopusdeploy_kubernetes_agent_deployment_target" "tenanted_agent" {
  name                              = "tenanted-agent"
  environments                      = [octopusdeploy_environment.development_environment.id]
  roles                             = ["role-1", "role-2"]
  uri                               = "poll://kcxzcv2fpsxkn6tk9u6d/"
  thumbprint                        = "96203ED84246201C26A2F4360D7CBC36AC1D232D"
  tenanted_deployment_participation = "Tenanted"
  tenants                           = [octopusdeploy_tenant.agent_tenant.id]
  tenant_tags                       = [
    octopusdeploy_tag.tag_a.canonical_tag_name, octopusdeploy_tag.tag_b.canonical_tag_name
  ]
}