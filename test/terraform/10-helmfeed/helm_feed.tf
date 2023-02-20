resource "octopusdeploy_helm_feed" "feed_helm" {
  name                                 = "Helm"
  password                             = "password"
  feed_uri                             = "https://charts.helm.sh/stable/"
  username                             = "username"
  package_acquisition_location_options = ["ExecutionTarget", "NotAcquired"]
}
