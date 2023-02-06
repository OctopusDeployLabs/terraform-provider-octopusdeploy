resource "octopusdeploy_library_variable_set" "library_variable_set_octopus_variables" {
  name        = "Octopus Variables"

  template {
    name             = "template"
    label            = "a"
    help_text        = "a"
    default_value    = "a"
    display_settings = { "Octopus.ControlType" = "SingleLineText" }
  }
}
