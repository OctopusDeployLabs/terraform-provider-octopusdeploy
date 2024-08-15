package schemas

import (
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	resourceSchema "github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
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

var ActionsAttributeToActionTypeMap = map[string]string{
	DeploymentProcessAction:                       "",
	DeploymentProcessRunScriptAction:              "Octopus.Script",
	DeploymentProcessRunKubectlScriptAction:       "Octopus.KubernetesRunScript",
	DeploymentProcessPackageAction:                "Octopus.TentaclePackage",
	DeploymentProcessWindowsServiceAction:         "Octopus.WindowsService",
	DeploymentProcessManualInterventionAction:     "Octopus.Manual",
	DeploymentProcessApplyKubernetesSecretAction:  "Octopus.KubernetesDeploySecret",
	DeploymentProcessApplyTerraformTemplateAction: "Octopus.TerraformApply",
}

type DeploymentProcessResourceModel struct {
	ID             types.String `tfsdk:"id"`
	SpaceID        types.String `tfsdk:"space_id"`
	ProjectID      types.String `tfsdk:"project_id"`
	Version        types.String `tfsdk:"version"`
	LastSnapshotID types.String `tfsdk:"last_snapshot_id"`
	Branch         types.String `tfsdk:"branch"`
	Steps          types.List   `tfsdk:"step"`
}

func GetDeploymentProcessResourceSchema() resourceSchema.Schema {
	return resourceSchema.Schema{
		Attributes: map[string]resourceSchema.Attribute{
			"id":                    GetIdResourceSchema(),
			"space_id":              GetSpaceIdResourceSchema(DeploymentProcessDescription),
			DeploymentProcessBranch: GetBranchResourceSchema(DeploymentProcessDescription),
			DeploymentProcessLastSnapshotId: resourceSchema.StringAttribute{
				Optional: true,
				Computed: true,
				Default:  stringdefault.StaticString(""),
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			DeploymentProcessProjectId: GetProjectIdResourceSchema(DeploymentProcessDescription),
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
	return resourceSchema.ListNestedBlock{
		NestedObject: resourceSchema.NestedBlockObject{
			Attributes: map[string]resourceSchema.Attribute{
				"id":   GetIdResourceSchema(),
				"name": GetNameResourceSchema(true),
				DeploymentProcessCondition: resourceSchema.StringAttribute{
					Default:     stringdefault.StaticString("Success"),
					Description: "When to run the step, one of 'Success', 'Failure', 'Always' or 'Variable'",
					Optional:    true,
					Computed:    true,
					Validators: []validator.String{
						stringvalidator.OneOf("Always",
							"Failure",
							"Success",
							"Variable"),
					},
				},
				DeploymentProcessConditionExpression: getConditionExpressionResourceSchema(),
				DeploymentProcessPackageRequirement:  getPackageRequirementResourceSchema(),
				DeploymentProcessProperties:          getPropertiesResourceSchema(),
				DeploymentProcessStartTrigger:        getStartTriggerResourceSchema(),
				DeploymentProcessTargetRoles:         getTargetRolesResourceSchema(),
				DeploymentProcessWindowSize:          getWindowSizeResourceSchema(),
			},
			Blocks: map[string]resourceSchema.Block{
				DeploymentProcessAction:                       NewActionResourceSchemaBuilder().WithActionType().WithExecutionLocation().WithWorkerPool().WithWorkerPoolVariable().WithGitDependency().Build(),
				DeploymentProcessRunScriptAction:              NewActionResourceSchemaBuilder().WithExecutionLocation().WithScriptFromPackage().WithWorkerPool().WithWorkerPoolVariable().WithScript().WithVariableSubstitutionInFiles().Build(),
				DeploymentProcessRunKubectlScriptAction:       NewActionResourceSchemaBuilder().WithExecutionLocation().WithScriptFromPackage().WithPackages().WithWorkerPool().WithWorkerPoolVariable().WithScript().WithVariableSubstitutionInFiles().WithNamespace().Build(),
				DeploymentProcessApplyTerraformTemplateAction: NewActionResourceSchemaBuilder().WithExecutionLocation().WithPrimaryPackage().WithWorkerPool().WithWorkerPoolVariable().WithTerraform().Build(),
				DeploymentProcessApplyKubernetesSecretAction:  NewActionResourceSchemaBuilder().WithExecutionLocation().WithWorkerPool().WithWorkerPoolVariable().WithKubernetesSecret().Build(),
				DeploymentProcessPackageAction:                NewActionResourceSchemaBuilder().WithPrimaryPackage().WithWindowsServiceFeature().Build(),
				DeploymentProcessWindowsServiceAction:         NewActionResourceSchemaBuilder().WithPrimaryPackage().WithWindowsService().Build(),
				DeploymentProcessManualInterventionAction:     NewActionResourceSchemaBuilder().WithManualIntervention().Build(),
			},
		},
	}
}
