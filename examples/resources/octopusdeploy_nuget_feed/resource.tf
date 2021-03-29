resource "octopusdeploy_nuget_feed" "example" {
  download_attempts              = 1
  download_retry_backoff_seconds = 30
  feed_uri                       = "https://api.nuget.org/v3/index.json"
  is_enhanced_mode               = true
  password                       = "test-password"
  name                           = "Test NuGet Feed (OK to Delete)"
  username                       = "test-username"
}
