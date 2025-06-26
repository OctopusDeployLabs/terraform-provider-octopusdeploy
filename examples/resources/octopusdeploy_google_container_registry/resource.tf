resource "octopusdeploy_google_container_registry" "example" {
  name          = "Test Google Container Registry (OK to Delete)"
  feed_uri      = "https://google.docker.test"
  registry_path = "testing/test-image"
  password      = "google authentication key file contents (json)"
}

resource "octopusdeploy_google_container_registry" "example_with_oidc" {
  name          = "Test Google Container Registry (OK to Delete)"
  feed_uri      = "https://google.docker.test"
  registry_path = "testing/test-image"
  oidc_authentication = {
    audience      = "audience"
    subject_keys = ["feed", "space"]
  }
}
