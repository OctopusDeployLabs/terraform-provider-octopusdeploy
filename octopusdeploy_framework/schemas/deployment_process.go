package schemas

import (
	"github.com/OctopusDeploy/terraform-provider-octopusdeploy/octopusdeploy_framework/util"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	resourceSchema "github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
)

const (
	DeploymentProcessDescription                  = "deployment process"
	DeploymentProcessBranch                       = "branch"
	DeploymentProcessLastSnapshotId               = "last_snapshot_id"
	DeploymentProcessProjectId                    = "project_id"
	DeploymentProcessVersion                      = "version"
	DeploymentProcessStep                         = "step"
	DeploymentProcessAction                       = "action"
	DeploymentProcessApplyTerraformTemplateAction = "apply_terraform_template_action"
	DeploymentProcessApplyKubernetesSecretAction  = "deploy_kubernetes_secret_action"
	DeploymentProcessPackageAction                = "deploy_package_action"
	DeploymentProcessManualInterventionAction     = "manual_intervention_action"
	DeploymentProcessWindowsServiceAction         = "deploy_windows_service_action"
	DeploymentProcessRunKubectlScriptAction       = "run_kubectl_script_action"
	DeploymentProcessRunScriptAction              = "run_script_action"
	DeploymentProcessCondition                    = "condition"
	DeploymentProcessConditionExpression          = "condition_expression"
	DeploymentProcessPackageRequirement           = "package_requirement"
	DeploymentProcessProperties                   = "properties"
	DeploymentProcessStartTrigger                 = "start_trigger"
	DeploymentProcessTargetRoles                  = "target_roles"
	DeploymentProcessWindowSize                   = "window_size"
)

type DeploymentProcessResourceModel struct {
}

func GetDeploymentProcessResourceSchema() resourceSchema.Schema {
	return resourceSchema.Schema{
		Attributes: map[string]resourceSchema.Attribute{
			"id":                            GetIdResourceSchema(),
			"space_id":                      GetSpaceIdResourceSchema(DeploymentProcessDescription),
			DeploymentProcessBranch:         GetBranchResourceSchema(DeploymentProcessDescription),
			DeploymentProcessLastSnapshotId: resourceSchema.StringAttribute{Optional: true},
			DeploymentProcessProjectId:      GetProjectIdResourceSchema(DeploymentProcessDescription),
			DeploymentProcessVersion: resourceSchema.StringAttribute{
				Computed:    true,
				Optional:    true,
				Description: "The version number of this deployment process.",
			},
		},
		Blocks: map[string]resourceSchema.Block{
			DeploymentProcessStep: getStepResourceBlockSchema(DeploymentProcessDescription),
		},
	}
}

func getStepResourceBlockSchema(resourceDescription string) resourceSchema.ListNestedBlock {
	actionBuilder := util.NewActionResourceSchemaBuilder()
	return resourceSchema.ListNestedBlock{
		NestedObject: resourceSchema.NestedBlockObject{
			Attributes: map[string]resourceSchema.Attribute{
				"id":   GetIdResourceSchema(),
				"name": GetNameResourceSchema(true),

				//DeploymentProcessApplyTerraformTemplateAction: getApplyTerraformTemplateActionSchema(),
				DeploymentProcessCondition: resourceSchema.StringAttribute{
					Default:     stringdefault.StaticString("Success"),
					Description: "When to run the step, one of 'Success', 'Failure', 'Always' or 'Variable'",
					Optional:    true,
					Validators: []validator.String{
						stringvalidator.OneOf("Always",
							"Failure",
							"Success",
							"Variable"),
					},
				},
				DeploymentProcessConditionExpression: getConditionExpressionResourceSchema(),
				//DeploymentProcessApplyKubernetesSecretAction: getDeployKubernetesSecretActionResourceSchema(),
				//DeploymentProcessPackageAction: getPackageActionResourceSchema(),
				//DeploymentProcessWindowsServiceAction: getWindowsServiceActionResourceSchema(),
				//DeploymentProcessManualInterventionAction: getManualInterventionActionResourceSchema(),
				DeploymentProcessPackageRequirement: getPackageRequirementResourceSchema(),
				DeploymentProcessProperties:         getPropertiesResourceSchema(),
				//DeploymentProcessRunKubectlScriptAction: getKubectlActionResourceSchema(),
				DeploymentProcessStartTrigger: getStartTriggerResourceSchema(),
				DeploymentProcessTargetRoles:  getTargetRolesResourceSchema(),
				DeploymentProcessWindowSize:   getWindowSizeResourceSchema(),
			},
			Blocks: map[string]resourceSchema.Block{
				DeploymentProcessAction: actionBuilder.WithActionType().Build(),
			},
		},
	}
}
