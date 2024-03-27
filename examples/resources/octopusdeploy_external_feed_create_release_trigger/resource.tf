resource "octopusdeploy_external_feed_create_release_trigger" "my_trigger" {
  name        = "My feed trigger"
  space_id    = "Spaces-1"
  project_id  = "Projects-2"
  package {
    deployment_action = "My Helm step"
    package_reference = "nginx"
  }
  package {
    deployment_action = "My container step"
    package_reference = "busybox"
  }
  channel_id = "Channels-21"
}