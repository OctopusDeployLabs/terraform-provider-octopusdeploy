resource "octopusdeploy_tenant_project_variable" "example" {
  environment_id = "environment-123"
  project_id = "project-123"
  template_id = "template-123"
  tenant_id = "tenant-123"
  value = "my-tenant-project-variable-value"
}