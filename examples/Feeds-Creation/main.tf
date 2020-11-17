provider "octopusdeploy" {
  address  = var.serverURL
  apikey   = var.apiKey
  space_id = var.space
}

# Feed username and password are only needed if authentication is required

resource "octopusdeploy_feed" "newFeed" {
  name      = var.feedName
  feed_type = "GitHub"
  feed_uri  = var.feed_uri
  #username = github_username
  #password = github_password
}
