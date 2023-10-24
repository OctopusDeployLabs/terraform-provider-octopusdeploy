resource "octopusdeploy_tenant" "example" {
  name = "Example Tenant"

  lifecycle {
    ignore_changes = [project_environment]
  }
}

resource "octopusdeploy_tenant_connection" "example" {
  tenant_id       = octopusdeploy_tenant.example.id
  project_id      = "Projects-123"
  environment_ids = ["Environments-1", "Environments-2", "Environments-3"]
}
