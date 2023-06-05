data "octopusdeploy_feeds" "example" {
  feed_type    = "Helm"
  partial_name = "Helm"
  skip         = 0
  take         = 1
}

output "data_lookup" {
  value = data.octopusdeploy_feeds.example.feeds[0].id
}