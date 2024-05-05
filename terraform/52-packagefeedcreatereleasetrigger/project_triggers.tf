resource "octopusdeploy_external_feed_create_release_trigger" "external_feed_trigger_1" {
  name        = "My first trigger"
  space_id    = var.octopus_space_id
  project_id  = "${octopusdeploy_project.deploy_frontend_project.id}"
  package {
    deployment_action_slug = octopusdeploy_deployment_process.example.step[0].run_script_action[0].slug
    package_reference = "busybox"
  }
  package {
    deployment_action_slug = octopusdeploy_deployment_process.example.step[0].run_script_action[0].slug
    package_reference = "nginx"
  }
  channel_id = octopusdeploy_channel.test_channel_1.id

  depends_on = [
    octopusdeploy_deployment_process.example
  ]
}

resource "octopusdeploy_external_feed_create_release_trigger" "external_feed_trigger_2" {
  name        = "My second trigger"
  space_id    = var.octopus_space_id
  project_id  = "${octopusdeploy_project.deploy_frontend_project.id}"
  is_disabled = true
  package {
    deployment_action_slug = octopusdeploy_deployment_process.example.step[1].run_script_action[0].slug
    package_reference = "scratch"
  }
  channel_id = octopusdeploy_channel.test_channel_2.id

  depends_on = [
    octopusdeploy_deployment_process.example
  ]
}

resource "octopusdeploy_external_feed_create_release_trigger" "external_feed_trigger_3" {
  name        = "My third trigger"
  space_id    = var.octopus_space_id
  project_id  = "${octopusdeploy_project.deploy_frontend_project.id}"
  is_disabled = false
  package {
    deployment_action_slug = octopusdeploy_deployment_process.example.step[0].run_script_action[0].slug
    package_reference = "busybox"
  }
  primary_package {
    deployment_action_slug = octopusdeploy_deployment_process.example.step[2].action[0].slug
  }
  primary_package {
    deployment_action_slug = octopusdeploy_deployment_process.example.step[3].action[0].slug
  }
  channel_id = octopusdeploy_channel.test_channel_3.id

  depends_on = [
    octopusdeploy_deployment_process.example
  ]
}