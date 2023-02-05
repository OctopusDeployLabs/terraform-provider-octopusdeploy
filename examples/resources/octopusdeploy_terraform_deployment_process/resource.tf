# basic deployment process with 2 run a script steps
resource "octopusdeploy_deployment_process" "example" {
  project_id = "Projects-123"
  step {
    condition           = "Success"
    name                = "Apply a Terraform template"
    package_requirement = "LetOctopusDecide"
    start_trigger       = "StartAfterPrevious"

    action {
      action_type                        = "Octopus.TerraformApply"
      name                               = "Apply a Terraform template"
      condition                          = "Success"
      run_on_server                      = true
      is_disabled                        = false
      can_be_used_for_project_versioning = false
      is_required                        = false
      worker_pool_id                     = ""
      worker_pool_variable               = ""
      properties                         = {
        "Octopus.Action.Script.ScriptSource" = "Inline"
        "Octopus.Action.Terraform.AzureAccount" = "False"
        "Octopus.Action.Terraform.PlanJsonOutput" = "False"
        "Octopus.Action.Terraform.RunAutomaticFileSubstitution" = "True"
        "Octopus.Action.GoogleCloud.UseVMServiceAccount" = "True"
        "Octopus.Action.Terraform.GoogleCloudAccount" = "False"
        "Octopus.Action.Terraform.TemplateParameters" = jsonencode({})
        "Octopus.Action.Terraform.AllowPluginDownloads" = "True"
        "Octopus.Action.Terraform.ManagedAccount" = "None"
        "Octopus.Action.Terraform.Template" = "#test"
        "Octopus.Action.GoogleCloud.ImpersonateServiceAccount" = "False"
      }
      environments                       = []
      excluded_environments              = []
      channels                           = []
      tenant_tags                        = []
      features                           = []
    }

    properties   = {}
    target_roles = []
  }
}

