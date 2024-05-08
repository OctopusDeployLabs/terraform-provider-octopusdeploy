resource "octopusdeploy_polling_subscription_id" "example" {}

resource "octopusdeploy_polling_subscription_id" "example_with_dependencies" {
  dependencies = {
    "target" = octopusdeploy_kubernetes_agent_deployment_target.example.id
  }
}

# Usage
resource "octopusdeploy_kubernetes_agent_deployment_target" "agent" {
  name         = "agent"
  environments = ["environments-1"]
  roles        = ["role-1", "role-2"]
  thumbprint   = "96203ED84246201C26A2F4360D7CBC36AC1D232D"
  uri          = octopusdeploy_polling_subscription_id.example_with_dependencies.polling_uri
}