resource "octopusdeploy_tentacle_certificate" "example" {}

resource "octopusdeploy_tentacle_certificate" "example_with_dependencies" {
  dependencies = {
    "target" = octopusdeploy_kubernetes_agent_deployment_target.agent.id
  }
}

# Usage
resource "octopusdeploy_kubernetes_agent_deployment_target" "agent" {
  name         = "agent"
  environments = ["environments-1"]
  roles        = ["role-1", "role-2"]
  thumbprint   = octopusdeploy_tentacle_certificate.example_with_dependencies.thumbprint
  uri          = "poll://kcxzcv2fpsxkn6tk9u6d/"
}