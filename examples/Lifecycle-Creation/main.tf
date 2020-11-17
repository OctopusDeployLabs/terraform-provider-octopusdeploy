provider "octopusdeploy" {
  address  = var.serverURL
  apikey   = var.apiKey
  space_id = var.space
}

resource "octopusdeploy_lifecycle" "newLifecycle" {
  name = var.lifecycleName

  release_retention_policy {
    quantity_to_keep = 5
    unit             = "Days"
  }
}
