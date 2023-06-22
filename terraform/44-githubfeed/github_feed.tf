resource "octopusdeploy_github_repository_feed" "example" {
  download_attempts              = 1
  download_retry_backoff_seconds = 30
  feed_uri                       = "https://api.github.com"
  password                       = "test-password"
  name                           = "Github"
  username                       = "test-username"
}