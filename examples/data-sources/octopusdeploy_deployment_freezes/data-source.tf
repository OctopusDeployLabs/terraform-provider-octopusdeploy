data "octopusdeploy_deployment_freezes" "test_freeze" {
  ids          = null
  partial_name = "Freeze Name"
  skip         = 5
  take         = 100
}


data "octopusdeploy_deployment_freezes" "project_freezes" {
  project_ids = ["projects-1"]
  skip        = 0
  take        = 5
}

data "octopusdeploy_deployment_freezes" "tenant_freezes" {
  tenant_ids = ["tenants-1"]
  skip       = 0
  take       = 10
}
