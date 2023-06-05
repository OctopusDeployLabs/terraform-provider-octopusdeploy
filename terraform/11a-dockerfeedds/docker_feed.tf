data "octopusdeploy_feeds" "example" {
  feed_type    = "Docker"
  partial_name = "Docker"
  skip         = 0
  take         = 1
}

output "data_lookup" {
  value = data.octopusdeploy_feeds.example.feeds[0].id
}