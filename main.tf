provider "octopusdeploy" {
  address = "http://localhost:8081/"
  apikey  = "API-RBAR7M04RXMUC7JDGWHJQWB1EDC"
}

resource "octopusdeploy_project_group" "team1" {
  description = "The Best Team"
  name        = "Team #1"
}

resource "octopusdeploy_project" "test_project" {
  description      = "A really groundbreaking app"
  lifecycle_id     = "Lifecycles-1"
  name             = "Epic App"
  project_group_id = "${octopusdeploy_project_group.team1.id}"

  deployment_step_windows_service {
    executable_path = "bin\\wicked_service1.exe"
    service_name    = "WickedService"
    step_name       = "1 - Deploy Wicked Windows Service"
    step_order      = 1

    target_roles = [
      "MyRole1",
      "MyRole2",
    ]
  }

  deployment_step_windows_service {
    executable_path          = "bin\\lame_service.exe"
    service_name             = "Lame Service"
    step_name                = "2 - Deploy Lame Windows Service"
    service_start_mode       = "demand"
    configuration_transforms = false
    step_condition           = "failure"
    step_order               = 30

    target_roles = [
      "MyRole1",
    ]
  }

  deployment_step_windows_service {
    executable_path          = "bin\\lame_service.exe"
    service_name             = "Lame Service 3"
    step_name                = "3 - Deploy Lame Windows Service"
    service_start_mode       = "demand"
    configuration_transforms = false
    step_condition           = "failure"
    step_order               = 3

    target_roles = [
      "MyRole1",
    ]
  }
}

#   deployment_steps = [
#     "${octopusdeploy_deployment_step_windows_service.windows_service.id}",
#   ]


# resource "octopusdeploy_deployment_step_windows_service" "windows_service" {
#   configuration_transforms = true
#   configuration_variables  = true
#   executable_path          = "bin\\wicked_service.exe"
#   feed_id                  = "feeds-builtin"
#   service_account          = "LocalSystem"
#   service_name             = "Wicked Service"
#   service_start_mode       = "auto"
#   step_condition           = "success"
#   step_name                = "asd"
#   step_start_trigger       = "startafterprevious"


#   target_roles = [
#     "Testing",
#     "Acceptance",
#   ]
# }


# resource "octopusdeploy_deployment_step_custom" "custom_service" {
#   name        = "My Step Name"
#   action_type = "Octopus.WindowsService"


#   target_roles = [
#     "octopus-server",
#     "some-other-role",
#   ]


#   condition     = "success"            // Success, Failure, Always, Variable
#   start_trigger = "StartAfterPrevious"
# }

