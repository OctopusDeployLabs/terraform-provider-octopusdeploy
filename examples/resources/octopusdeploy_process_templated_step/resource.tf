resource "octopusdeploy_environment" "development" {
  name = "Development"
}

resource "octopusdeploy_environment" "production" {
  name = "Production"
}

# Template
resource "octopusdeploy_step_template" "my_script" {
  action_type     = "Octopus.Script"
  name            = "Run My Script"
  step_package_id = "Octopus.Script"
  packages        = [
    {
      package_id = "example.package"
      acquisition_location = "Server"
      feed_id = "Feeds-12"
      name = "util.one"
      properties = {
        extract = "True"
        purpose = ""
        selection_mode = "immediate"
      }
    }
  ]

  parameters = [
    {
      name      = "Parameter.One"
      id = "10001000-0000-0000-0000-100010001001"
      label     = "First Parameter"
      default_value = "default-value-one"
      display_settings = {
        "Octopus.ControlType" : "SingleLineText"
      }
    },
    {
      name      = "Parameter.Two"
      id = "10001000-0000-0000-0000-100010001002"
      label     = "Second Parameter"
      display_settings = {
        "Octopus.ControlType" : "SingleLineText"
      }
    },
  ]
    
  properties = {
    "Octopus.Action.Script.ScriptBody" : "echo '1.#{Parameter.One} ... 2.#{Parameter.Two} ...'"
    "Octopus.Action.Script.ScriptSource" : "Inline"
    "Octopus.Action.Script.Syntax" : "Bash"
  }
}

resource "octopusdeploy_project" "example" {
  project_group_id = "ProjectGroups-1"
  lifecycle_id = "Lifecycles-1"
  name = "Example"
}

resource "octopusdeploy_process" "example" {
  project_id  = octopusdeploy_project.example.id
}

# Templated script step
resource "octopusdeploy_templated_process_step" "script" {
  process_id  = octopusdeploy_process.example.id
  name = "Templated Step"
  template_id = octopusdeploy_step_template.my_script.id
  template_version = octopusdeploy_step_template.my_script.version

  # Parameter's default value is used when not provided in configuration
  parameters = {
    "Parameter.Two" = "my-example-value"
  }

  execution_properties = {
    "Octopus.Action.RunOnServer" = "True",
    "My.Custom.Property" = "Something",
  }
}

