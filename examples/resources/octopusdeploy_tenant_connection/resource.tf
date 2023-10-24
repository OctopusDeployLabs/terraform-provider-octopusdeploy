resource "octopusdeploy_tenant_connection" "example" {
  tenant_id       = "Tenants-123"
  project_id      = "Projects-123"
  environment_ids = ["Environments-1", "Environments-2", "Environments-3"]
}
