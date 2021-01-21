resource "octopusdeploy_environment" "example" {
  allow_dynamic_infrastructure = false
  description                  = "An environment for the development team."
  name                         = "Development Environment (OK to Delete)"
  use_guided_failure           = false
}
