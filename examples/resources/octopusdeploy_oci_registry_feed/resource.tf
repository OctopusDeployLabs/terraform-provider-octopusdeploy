resource "octopusdeploy_oci_registry_feed" "example" {
  feed_uri                       = "oci://test-registry.docker.io"
  password                       = "test-password"
  name                           = "Test oci Registry Feed (OK to Delete)"
  username                       = "test-username"
}
