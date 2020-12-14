resource "octopusdeploy_project" "example" {
  description                 = "The development project."
  is_disabled                 = false
  is_discrete_channel_release = false
  is_version_controlled       = false
  lifecycle_id                = "Lifecycles-123"
  name                        = "Development Project (OK to Delete)"
  project_group_id            = "ProjectGroups-123"
}
