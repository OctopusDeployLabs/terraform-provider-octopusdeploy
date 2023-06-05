resource "octopusdeploy_maven_feed" "feed_maven" {
  name                                 = "Maven"
  feed_uri                             = "https://repo.maven.apache.org/maven2/"
  username                             = "username"
  password                             = "password"
  package_acquisition_location_options = ["Server", "ExecutionTarget"]
  download_attempts                    = 5
  download_retry_backoff_seconds       = 10
}
