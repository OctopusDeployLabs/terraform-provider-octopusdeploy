resource "octopusdeploy_user_role" "example" {
  can_be_deleted             = true
  description                = "Responsible for all development-related operations."
  granted_space_permissions  = ["DeploymentCreate", "DeploymentDelete", "DeploymentView"]
  granted_system_permissions = ["SpaceCreate"]
  name                       = "Developer Managers"
  space_permission_descriptions = [
    "Delete deployments (restrictable to Environments, Projects, Tenants)",
    "Deploy releases to target environments (restrictable to Environments, Projects, Tenants)",
    "View deployments (restrictable to Environments, Projects, Tenants)"
  ]
}
