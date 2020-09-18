provider "octopusdeploy" {
  address = var.serverURL
  apikey  = var.apiKey
  space_id   = var.space
}

resource "octopusdeploy_channel" "newChannel" {
  name            = var.channelName
  project_id      = "Projects-1"
}
