resource "octopusdeploy_tag_set" "example" {
  description = "Provides tenants with access to certain early access programs."
  name        = "Early Access Program (EAP)"

  tag {
    color = "#00FF00"
    name  = "Alpha"
  }

  tag {
    color = "#FF0000"
    name  = "Beta"
  }
}
