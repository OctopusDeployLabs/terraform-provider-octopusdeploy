data "octopusdeploy_feeds" "example" {
  feed_type    = "AwsElasticContainerRegistry"
  partial_name = "ECR"
  skip         = 0
  take         = 1
}

output "data_lookup" {
  value = data.octopusdeploy_feeds.example.feeds[0].id
}