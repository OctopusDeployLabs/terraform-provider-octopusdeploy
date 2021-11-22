resource "octopusdeploy_script_module" "example" {
  description = "A script module to use."
  name        = "Hello Octopus Script Module"

  script {
    body   = "function Say-Hello()\r\n{\r\n    Write-Output \"Hello, Octopus!\"\r\n}\r\n"
    syntax = "PowerShell"
  }
}
