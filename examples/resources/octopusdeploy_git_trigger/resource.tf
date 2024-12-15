resource "octopusdeploy_git_trigger" "my_trigger" {
  name                = "My Git trigger"
  space_id            = "Spaces-1"
  description         = "My Git trigger description"
  project_id          = "Projects-1"
  channel_id          = "Channels-1"
  sources {
    deployment_action_slug = "deploy-action-slug"
    git_dependency_name    = ""
    include_file_paths     = [
      "include/me",
      "include/this/too"
    ]
    exclude_file_paths     = [
      "exclude/me",
      "exclude/this/too"
    ]
  }
}
