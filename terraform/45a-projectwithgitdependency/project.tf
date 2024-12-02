data "octopusdeploy_lifecycles" "lifecycle_default_lifecycle" {
  ids          = null
  partial_name = "Default Lifecycle"
  skip         = 0
  take         = 1
}


resource "octopusdeploy_project" "deploy_frontend_project" {
  auto_create_release                  = false
  default_guided_failure_mode          = "EnvironmentDefault"
  default_to_skip_if_already_installed = false
  description                          = "Test project"
  discrete_channel_release             = false
  is_disabled                          = false
  is_discrete_channel_release          = false
  is_version_controlled                = false
  lifecycle_id                         = data.octopusdeploy_lifecycles.lifecycle_default_lifecycle.lifecycles[0].id
  name                                 = "Test"
  project_group_id                     = octopusdeploy_project_group.project_group_test.id
  tenanted_deployment_participation    = "Untenanted"
  space_id                             = var.octopus_space_id
  included_library_variable_sets       = []
  versioning_strategy {
    template = "#{Octopus.Version.LastMajor}.#{Octopus.Version.LastMinor}.#{Octopus.Version.LastPatch}.#{Octopus.Version.NextRevision}"
  }

  connectivity_policy {
    allow_deployments_to_no_targets = false
    exclude_unhealthy_targets       = false
    skip_machine_behavior           = "SkipUnavailableMachines"
  }
}

resource "octopusdeploy_channel" "test_channel_1" {
    name = "Test Channel 1"
    project_id = "${octopusdeploy_project.deploy_frontend_project.id}"
    space_id = var.octopus_space_id
    lifecycle_id = data.octopusdeploy_lifecycles.lifecycle_default_lifecycle.lifecycles[0].id
}

resource "octopusdeploy_deployment_process" "deployment_process_project_noopterraform" {
  project_id = "${octopusdeploy_project.deploy_frontend_project.id}"

  step {
    condition           = "Success"
    name                = "Git Dependency"
    package_requirement = "LetOctopusDecide"
    start_trigger       = "StartAfterPrevious"
    properties          = {}
    target_roles        = ["bread"]

    action {
      action_type                        = "Octopus.KubernetesDeployRawYaml"
      name                               = "Git Dependency"
      sort_order                         = 1
      condition                          = "Success"
      run_on_server                      = true
      is_disabled                        = false
      can_be_used_for_project_versioning = false
      is_required                        = true
      worker_pool_id                     = ""
      properties                         = {
        "Octopus.Action.GitRepository.Source": "External",
        "Octopus.Action.Kubernetes.DeploymentTimeout" = "180",
        "Octopus.Action.Kubernetes.ResourceStatusCheck" = "True",
        "Octopus.Action.Kubernetes.ServerSideApply.Enabled" = "True",
        "Octopus.Action.Kubernetes.ServerSideApply.ForceConflicts" = "True",
        "Octopus.Action.KubernetesContainers.CustomResourceYamlFileName" = "test",
        "Octopus.Action.Script.ScriptSource" = "GitRepository"
      }
      environments                       = []
      excluded_environments              = []
      channels                           = []
      tenant_tags                        = []
      features                           = []
      
      git_dependency {
          default_branch = "main"
          file_path_filters = [
            "test"
          ]
          git_credential_id = "GitCredentials-1"
          git_credential_type = "Library"
          repository_uri = "https://github.com/octopus/testing.git"
      }
    }
  }
}