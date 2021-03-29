resource "octopusdeploy_helm_feed" "example" {
  feed_uri = "https://charts.helm.sh/stable"
  password = "test-password"
  name     = "Test Helm Feed (OK to Delete)"
  username = "test-username"
}
