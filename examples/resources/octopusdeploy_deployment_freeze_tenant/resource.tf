# Deployment freeze
resource "octopusdeploy_deployment_freeze" "example" {
  name = "Summer break"
  start = "2024-06-25T00:00:00+10:00"
  end = "2024-06-27T00:00:00+08:00"
}

# Freeze with ids sourced from resources and data sources. 
# Tenants can be sourced from different spaces, a single scope can only reference resources from the same space.

resource "octopusdeploy_deployment_freeze_tenant" "production_freeze" {
  deploymentfreeze_id = octopusdeploy_deployment_freeze.example.id
  tenant_id           = resource.octopusdeploy_tenant.example.id
  project_id          = resource.octopusdeploy_project.example.id
  environment_id      = data.environments.production.id
}