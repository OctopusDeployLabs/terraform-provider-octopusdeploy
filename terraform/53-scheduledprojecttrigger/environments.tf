resource "octopusdeploy_environment" "env_1" {
  name = "Env1"
  space_id = var.octopus_space_id
}

resource "octopusdeploy_environment" "env_2" {
  name = "Env2"
  space_id = var.octopus_space_id
}

data "octopusdeploy_lifecycles" "lifecycle_default_lifecycle" {
  ids          = null
  partial_name = "Default Lifecycle"
  space_id = var.octopus_space_id
  skip         = 0
  take         = 1
}

resource "octopusdeploy_project_group" "project_group_test" {
  name        = "Test"
  space_id = var.octopus_space_id
  description = "Test Description"
}
