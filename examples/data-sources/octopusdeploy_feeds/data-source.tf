data "octopusdeploy_feeds" "example" {
  feed_type    = "NuGet"
  ids          = ["Feeds-123", "Feeds-321"]
  partial_name = "Develop"
  skip         = 5
  take         = 100
}