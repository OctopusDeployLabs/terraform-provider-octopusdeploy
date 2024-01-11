resource "octopusdeploy_environment" "development_environment" {
  allow_dynamic_infrastructure = true
  description                  = "A test environment"
  name                         = "Development"
  use_guided_failure           = false
}