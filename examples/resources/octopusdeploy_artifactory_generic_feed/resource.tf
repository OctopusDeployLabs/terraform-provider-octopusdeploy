resource "octopusdeploy_artifactory_generic_feed" "example" {
  feed_uri                       = "https://example.jfrog.io"
  password                       = "test-password"
  name                           = "Test Artifactory Generic Feed (OK to Delete)"
  username                       = "test-username"
  repository                     = "repo"
  layout_regex                   = "this is regex"
}
