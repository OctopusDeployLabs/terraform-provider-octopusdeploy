package schemas

import (
	"github.com/OctopusDeploy/terraform-provider-octopusdeploy/octopusdeploy_framework/util"
	datasourceSchema "github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	resourceSchema "github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type ProcessChildStepSchema struct{}

var _ EntitySchema = ProcessChildStepSchema{}

const ProcessChildStepResourceName = "process_child_step"

func (p ProcessChildStepSchema) GetResourceSchema() resourceSchema.Schema {
	return resourceSchema.Schema{
		Description: "This resource manages a child step in execution process in Octopus Deploy.",
		Attributes: map[string]resourceSchema.Attribute{
			"id":       GetIdResourceSchema(),
			"space_id": GetSpaceIdResourceSchema(ProcessChildStepResourceName),
			"process_id": util.ResourceString().
				Description("Id of the process this step belongs to.").
				Required().
				PlanModifiers(stringplanmodifier.RequiresReplace()).
				Build(),
			"parent_id": util.ResourceString().
				Description("Id of the process step this step belongs to.").
				Required().
				PlanModifiers(stringplanmodifier.RequiresReplace()).
				Build(),
			"name": GetNameResourceSchema(true),
			"type": util.ResourceString().
				Description("Execution type of the step.").
				Required().
				Build(),
			"slug": util.ResourceString().
				Description("The human-readable unique identifier for the step.").
				Optional().
				Computed().
				PlanModifiers(stringplanmodifier.UseStateForUnknown()).
				Build(),
			"is_disabled": util.ResourceBool().
				Description("Indicates the disabled status of this step.").
				Optional().
				Computed().
				Default(false).
				Build(),
			"is_required": util.ResourceBool().
				Description("Indicates the required status of this step.").
				Optional().
				Computed().
				Default(false).
				Build(),
			"condition": util.ResourceString().
				Description("When to run the step, can be 'Success' - run when previous child step succeed or variable expression - run when the expression evaluates to true").
				Optional().
				Computed().
				Default("Success").
				Build(),
			"notes": util.ResourceString().
				Description("The notes associated with this step.").
				Optional().
				Computed().
				Default("").
				Build(),
			"worker_pool_id": util.ResourceString().
				Description("The worker pool associated with this step.").
				Optional().
				Computed().
				Default("").
				Build(),
			"worker_pool_variable": util.ResourceString().
				Description("The worker pool variable associated with this step.").
				Optional().
				Computed().
				Default("").
				Build(),
			"tenant_tags": util.ResourceSet(types.StringType).
				Description("A set of tenant tags associated with this step.").
				Optional().
				Computed().
				DefaultEmpty().
				Build(),
			"environments": util.ResourceSet(types.StringType).
				Description("A set of environments within which this step will run.").
				Optional().
				Computed().
				DefaultEmpty().
				Build(),
			"excluded_environments": util.ResourceSet(types.StringType).
				Description("A set of environments that this step will be skipped in.").
				Optional().
				Computed().
				DefaultEmpty().
				Build(),
			"channels": util.ResourceSet(types.StringType).
				Description("A set of channels associated with this step.").
				Optional().
				Computed().
				DefaultEmpty().
				Build(),
			"container":        resourceActionContainerAttribute(),
			"git_dependencies": resourceActionGitDependenciesAttribute(),
			"packages":         resourceActionPackageReferencesAttribute(),
			"execution_properties": util.ResourceMap(types.StringType).
				Description("A collection of step execution properties where the key is the property name and the value is its value.").
				Optional().
				Computed().
				DefaultEmpty().
				Build(),
		},
	}
}

func (p ProcessChildStepSchema) GetDatasourceSchema() datasourceSchema.Schema {
	return datasourceSchema.Schema{}
}

type ProcessChildStepResourceModel struct {
	SpaceID   types.String `tfsdk:"space_id"`
	ProcessID types.String `tfsdk:"process_id"`
	ParentID  types.String `tfsdk:"parent_id"`
	Name      types.String `tfsdk:"name"`

	Type                 types.String                     `tfsdk:"type"`
	Slug                 types.String                     `tfsdk:"slug"`
	IsDisabled           types.Bool                       `tfsdk:"is_disabled"`
	IsRequired           types.Bool                       `tfsdk:"is_required"`
	Condition            types.String                     `tfsdk:"condition"`
	Notes                types.String                     `tfsdk:"notes"`
	WorkerPoolID         types.String                     `tfsdk:"worker_pool_id"`
	WorkerPoolVariable   types.String                     `tfsdk:"worker_pool_variable"`
	TenantTags           types.Set                        `tfsdk:"tenant_tags"`
	Environments         types.Set                        `tfsdk:"environments"`
	ExcludedEnvironments types.Set                        `tfsdk:"excluded_environments"`
	Channels             types.Set                        `tfsdk:"channels"`
	Container            *ProcessStepActionContainerModel `tfsdk:"container"`
	GitDependencies      types.Map                        `tfsdk:"git_dependencies"`
	Packages             types.Map                        `tfsdk:"packages"`
	ExecutionProperties  types.Map                        `tfsdk:"execution_properties"`

	ResourceModel
}
