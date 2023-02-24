resource "octopusdeploy_tag_set" "example" {
  description = "Provides tenants with access to certain early access programs."
  name        = "Early Access Program (EAP)"
}

# tags are distinct resources and associated with tag sets through tag_set_id

resource "octopusdeploy_tag" "alpha" {
  color      = "#00FF00"
  name       = "Alpha"
  tag_set_id = octopusdeploy_tag_set.example.id
}

resource "octopusdeploy_tag" "beta" {
  color      = "#FF0000"
  name       = "Beta"
  tag_set_id = octopusdeploy_tag_set.example.id
}
