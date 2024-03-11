resource "octopusdeploy_helm_feed_create_release_trigger" "my_helm_trigger" {
  name        = "Helm trigger"
  space_id    = var.octopus_space_id
  project_id  = "${octopusdeploy_project.deploy_frontend_project.id}"
  package {
    deployment_action = "My Helm step"
    package_reference = "prometheus"
  }
  package {
    deployment_action = "My other Helm step"
    package_reference = "nginx"
  }
  channel_id = "Channels-1"
}

resource "octopusdeploy_container_feed_create_release_trigger" "my_container_trigger" {
  name        = "Container image trigger"
  space_id    = var.octopus_space_id
  project_id  = "${octopusdeploy_project.deploy_frontend_project.id}"
  is_disabled = true
  package {
    deployment_action = "My Docker step"
    package_reference = "busybox"
  }
  channel_id = "Channels-1"
}