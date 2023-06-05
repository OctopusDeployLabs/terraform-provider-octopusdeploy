data "octopusdeploy_feeds" "example" {
  feed_type    = "GitHub"
  partial_name = "Github"
  skip         = 0
  take         = 1
}

output "data_lookup" {
  value = data.octopusdeploy_feeds.example.feeds[0].id
}