resource "octopusdeploy_feed" "example" {
  download_attempts                    = 10
  download_retry_backoff_seconds       = 60
  feed_uri                             = "https://repo.maven.apache.org/maven2/"
  name                                 = "Test Maven Feed (OK to Delete)"
  package_acquisition_location_options = ["Server", "ExecutionTarget"]
  password                             = "password123"
  space_id                             = "Spaces-123"
  username                             = "john.smith@example.com"
}