resource "octopusdeploy_docker_container_registry" "docker_feed" {
  feed_uri      = "https://index.docker.io"
  name          = "Test Docker Container Registry"
}

resource "octopusdeploy_helm_feed" "feed_helm_charts" {
  name                                 = "Test Helm Charts"
  feed_uri                             = "https://charts.helm.sh/stable"
  package_acquisition_location_options = ["ExecutionTarget", "NotAcquired"]
}
