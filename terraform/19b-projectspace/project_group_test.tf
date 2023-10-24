resource "octopusdeploy_project_group" "project_group_test" {
  name        = "Test"
  description = "Test Description"
  space_id    = octopusdeploy_space.octopus_project_space_test.id
}
