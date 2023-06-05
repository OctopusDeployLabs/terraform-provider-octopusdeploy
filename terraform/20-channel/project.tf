data "octopusdeploy_lifecycles" "lifecycle_default_lifecycle" {
  ids          = null
  partial_name = "Default Lifecycle"
  skip         = 0
  take         = 1
}

data "octopusdeploy_worker_pools" "workerpool_default" {
  name = "Default Worker Pool"
  ids  = null
  skip = 0
  take = 1
}

data "octopusdeploy_feeds" "built_in_feed" {
  feed_type    = "BuiltIn"
  ids          = null
  partial_name = ""
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

resource "octopusdeploy_deployment_process" "deploy_backend" {
  project_id = octopusdeploy_project.deploy_frontend_project.id

  step {
    condition           = "Success"
    name                = "Test"
    package_requirement = "LetOctopusDecide"
    start_trigger       = "StartAfterPrevious"

    action {
      action_type                        = "Octopus.AwsRunCloudFormation"
      name                               = "Test"
      condition                          = "Success"
      run_on_server                      = true
      is_disabled                        = false
      can_be_used_for_project_versioning = false
      is_required                        = false
      worker_pool_id                     = data.octopusdeploy_worker_pools.workerpool_default.worker_pools[0].id
      properties                         = { "Octopus.Action.Aws.AssumeRole" = "False", "Octopus.Action.Aws.CloudFormation.Tags" = "[{\"key\":\"Environment\",\"value\":\"#{Octopus.Environment.Name | Replace \\\" .*\\\" \\\"\\\"}\"},{\"key\":\"DeploymentProject\",\"value\":\"API_Gateway\"}]", "Octopus.Action.Aws.CloudFormationStackName" = "OctopusBuilder-APIGateway-mcasperson-#{Octopus.Environment.Name | Replace \" .*\" \"\" | ToLower}", "Octopus.Action.Aws.CloudFormationTemplate" = "Resources:\n  RestApi:\n    Type: 'AWS::ApiGateway::RestApi'\n    Properties:\n      Description: My API Gateway\n      Name: Octopus Workflow Builder\n      BinaryMediaTypes:\n        - '*/*'\n      EndpointConfiguration:\n        Types:\n          - REGIONAL\n  Health:\n    Type: 'AWS::ApiGateway::Resource'\n    Properties:\n      RestApiId:\n        Ref: RestApi\n      ParentId:\n        'Fn::GetAtt':\n          - RestApi\n          - RootResourceId\n      PathPart: health\n  Api:\n    Type: 'AWS::ApiGateway::Resource'\n    Properties:\n      RestApiId:\n        Ref: RestApi\n      ParentId:\n        'Fn::GetAtt':\n          - RestApi\n          - RootResourceId\n      PathPart: api\n  Web:\n    Type: 'AWS::ApiGateway::Resource'\n    Properties:\n      RestApiId: !Ref RestApi\n      ParentId: !GetAtt\n        - RestApi\n        - RootResourceId\n      PathPart: '{proxy+}'\nOutputs:\n  RestApi:\n    Description: The REST API\n    Value: !Ref RestApi\n  RootResourceId:\n    Description: ID of the resource exposing the root resource id\n    Value:\n      'Fn::GetAtt':\n        - RestApi\n        - RootResourceId\n  Health:\n    Description: ID of the resource exposing the health endpoints\n    Value: !Ref Health\n  Api:\n    Description: ID of the resource exposing the api endpoint\n    Value: !Ref Api\n  Web:\n    Description: ID of the resource exposing the web app frontend\n    Value: !Ref Web\n", "Octopus.Action.Aws.CloudFormationTemplateParameters" = "[]", "Octopus.Action.Aws.CloudFormationTemplateParametersRaw" = "[]", "Octopus.Action.Aws.Region" = "ap-southeast-2", "Octopus.Action.Aws.TemplateSource" = "Inline", "Octopus.Action.Aws.WaitForCompletion" = "True", "Octopus.Action.AwsAccount.UseInstanceRole" = "False", "Octopus.Action.AwsAccount.Variable" = "AWS Account" }
      environments                       = []
      excluded_environments              = []
      channels                           = []
      tenant_tags                        = []
      features                           = []
      package {
        name                      = "test"
        package_id                = "test"
        feed_id                   = data.octopusdeploy_feeds.built_in_feed.feeds[0].id
        acquisition_location      = "Server"
        extract_during_deployment = true
      }
    }

    properties   = {}
    target_roles = []
  }
}