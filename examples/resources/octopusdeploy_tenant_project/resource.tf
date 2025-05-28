resource "octopusdeploy_tenant_project" "example" {
  project_id = "project-123"
  tenant_id = "tenant-123"
  environment_ids = [
    "environment-123"
  ]
}