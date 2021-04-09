resource "octopusdeploy_maven_feed" "example" {
  download_attempts = 10
  download_retry_backoff_seconds = 20
  feed_uri = "https://repo.maven.apache.org/maven2/"
  password = "test-password"
  name     = "Test Maven Feed (OK to Delete)"
  username = "test-username"
}
