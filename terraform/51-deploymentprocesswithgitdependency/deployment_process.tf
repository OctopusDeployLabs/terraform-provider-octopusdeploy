locals {
  git_uri             = "https://github.com/OctopusSamples/OctoPetShop.git"
  default_branch      = "main"
  git_credential_type = "Library"
}

resource "octopusdeploy_deployment_process" "test_deployment_process" {
  project_id = octopusdeploy_project.test_project.id

# Supported steps
  step {
    condition           = "Success"
    name                = "Generic Action"
    package_requirement = "LetOctopusDecide"
    start_trigger       = "StartAfterPrevious"
    action {
      action_type = "Octopus.Script"
      name        = "Generic Action"
      git_dependency {
        repository_uri      = local.git_uri
        default_branch      = local.default_branch
        git_credential_type = local.git_credential_type
        git_credential_id   = octopusdeploy_git_credential.test_git_credential.id
      }
      can_be_used_for_project_versioning = false
      condition                          = "Success"
      is_disabled                        = false
      is_required                        = true
      run_on_server                      = true
      properties                         = {
        "Octopus.Action.EnabledFeatures" : "Octopus.Features.SubstituteInFiles",
        "Octopus.Action.GitRepository.Source" : "External",
        "Octopus.Action.Script.ScriptFileName" : "Test.sh",
        "Octopus.Action.Script.ScriptSource" : "GitRepository",
        "Octopus.Action.SubstituteInFiles.Enabled" : "True"
      }
    }
  }

  step {
    condition           = "Success"
    name                = "Run Script"
    package_requirement = "LetOctopusDecide"
    start_trigger       = "StartAfterPrevious"
    run_script_action {
      git_dependency {
        repository_uri      = local.git_uri
        default_branch      = local.default_branch
        git_credential_type = local.git_credential_type
        git_credential_id   = octopusdeploy_git_credential.test_git_credential.id
      }
      script_source                      = "GitRepository"
      script_file_name                   = "Test.sh"
      can_be_used_for_project_versioning = false
      condition                          = "Success"
      is_disabled                        = false
      is_required                        = true
      name                               = "Run Script"
      run_on_server                      = true
    }
  }

  step {
    condition           = "Success"
    name                = "Run Script - Anon"
    package_requirement = "LetOctopusDecide"
    start_trigger       = "StartAfterPrevious"
    run_script_action {
      git_dependency {
        repository_uri      = local.git_uri
        default_branch      = local.default_branch
        git_credential_type = "Anonymous"
      }
      script_source                      = "GitRepository"
      script_file_name                   = "Test.sh"
      can_be_used_for_project_versioning = false
      condition                          = "Success"
      is_disabled                        = false
      is_required                        = true
      name                               = "Run Script - Anon"
      run_on_server                      = true
    }
  }

  step {
    condition           = "Success"
    name                = "Kubectl Action"
    package_requirement = "LetOctopusDecide"
    start_trigger       = "StartAfterPrevious"
    run_kubectl_script_action {
      name             = "Kubectl Action"
      run_on_server    = true
      script_source    = "GitRepository"
      script_file_name = "Test.sh"
      git_dependency {
        repository_uri      = local.git_uri
        default_branch      = local.default_branch
        git_credential_type = local.git_credential_type
        git_credential_id   = octopusdeploy_git_credential.test_git_credential.id
      }
    }
  }

  step {
    condition           = "Success"
    name                = "Raw Yaml Action"
    package_requirement = "LetOctopusDecide"
    start_trigger       = "StartAfterPrevious"
    properties          = { "Octopus.Action.TargetRoles" : "Qwerty" }
    action {
      action_type = "Octopus.KubernetesDeployRawYaml"
      name        = "Raw Yaml Action"
      git_dependency {
        repository_uri      = local.git_uri
        default_branch      = local.default_branch
        git_credential_type = local.git_credential_type
        git_credential_id   = octopusdeploy_git_credential.test_git_credential.id
      }
      can_be_used_for_project_versioning = false
      condition                          = "Success"
      is_disabled                        = false
      is_required                        = true
      run_on_server                      = true
      properties                         = {
        "Octopus.Action.Script.ScriptFileName" : "Test.sh",
        "Octopus.Action.Script.ScriptSource" : "GitRepository",
        "Octopus.Action.SubstituteInFiles.Enabled" : "True",
        "Octopus.Action.KubernetesContainers.CustomResourceYamlFileName" : "files/*",
        "Octopus.Action.KubernetesContainers.DeploymentWait" : "NoWait",
        "Octopus.Action.KubernetesContainers.Namespace" : "default"
      }
    }
  }

# Step that doesn't support git dependencies
  step {
    condition           = "Success"
    name                = "Deploy Secret"
    package_requirement = "LetOctopusDecide"
    start_trigger       = "StartAfterPrevious"
    properties          = { "Octopus.Action.TargetRoles" : "Qwerty" }
    deploy_kubernetes_secret_action {
      name          = "Deploy secret"
      secret_name   = "name"
      secret_values = {
        "val" = "123"
      }
      git_dependency {
        repository_uri      = local.git_uri
        default_branch      = local.default_branch
        git_credential_type = local.git_credential_type
        git_credential_id   = octopusdeploy_git_credential.test_git_credential.id
      }
    }
  }
}