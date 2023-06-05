resource "octopusdeploy_script_module" "library_variable_set_test2" {
  description = "Test script module"
  name        = "Test2"

  script {
    body   = "echo \"hi\""
    syntax = "PowerShell"
  }
}
