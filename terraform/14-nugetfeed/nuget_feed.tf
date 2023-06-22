resource "octopusdeploy_nuget_feed" "feed_nuget" {
  name                                 = "Nuget"
  feed_uri                             = "https://index.docker.io"
  username                             = "username"
  password                             = "password"
  is_enhanced_mode                     = true
  package_acquisition_location_options = ["Server", "ExecutionTarget"]
  download_attempts                    = 5
  download_retry_backoff_seconds       = 10
}
