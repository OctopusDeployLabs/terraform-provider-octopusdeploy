resource "octopusdeploy_docker_container_registry" "docker_feed" {
  feed_uri      = "https://index.docker.io"
  name          = "Test Docker Container Registry"
}