resource "octopusdeploy_environment" "development_environment" {
  allow_dynamic_infrastructure = true
  description                  = "A test environment"
  name                         = "Development"
  use_guided_failure           = false
}

resource "octopusdeploy_environment" "test_environment" {
  allow_dynamic_infrastructure = true
  description                  = "A test environment"
  name                         = "Test"
  use_guided_failure           = false
}

resource "octopusdeploy_environment" "production_environment" {
  allow_dynamic_infrastructure = true
  description                  = "A test environment"
  name                         = "Production"
  use_guided_failure           = false
}