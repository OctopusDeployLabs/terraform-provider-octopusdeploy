package schemas

import (
	"github.com/OctopusDeploy/terraform-provider-octopusdeploy/octopusdeploy_framework/util"
	datasourceSchema "github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	resourceSchema "github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/mapplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type ProcessTemplatedChildStepSchema struct{}

var _ EntitySchema = ProcessTemplatedChildStepSchema{}

const ProcessTemplatedChildStepResourceName = "process_templated_child_step"

func (p ProcessTemplatedChildStepSchema) GetResourceSchema() resourceSchema.Schema {
	return resourceSchema.Schema{
		Description: "This resource manages a child step of a Runbook or Deployment process which based on existing custom step template",
		Attributes: map[string]resourceSchema.Attribute{
			"id":       GetIdResourceSchema(),
			"space_id": GetSpaceIdResourceSchema(ProcessTemplatedChildStepResourceName),
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
			"template_id": util.ResourceString().
				Description("Id of template this step will be based on.").
				Required().
				Build(),
			"template_version": util.ResourceInt32().
				Description("Version of the template this step will be based on.").
				Required().
				Build(),
			"name": GetNameResourceSchema(true),
			"type": util.ResourceString().
				Description("Execution type of the step copied from the template.").
				Computed().
				PlanModifiers(stringplanmodifier.UseStateForUnknown()).
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
			"container": resourceActionContainerAttribute(),
			"git_dependencies": resourceSchema.MapNestedAttribute{
				Description: "References of git dependencies copied from the template",
				Computed:    true,
				PlanModifiers: []planmodifier.Map{
					mapplanmodifier.UseStateForUnknown(),
				},
				NestedObject: resourceActionGitDependencyNestedAttribute(),
			},
			"packages": resourceSchema.MapNestedAttribute{
				Description: "Package references copied from the template",
				Computed:    true,
				PlanModifiers: []planmodifier.Map{
					mapplanmodifier.UseStateForUnknown(),
				},
				NestedObject: resourceActionPackageReferenceNestedAttribute(),
			},
			"parameters": util.ResourceMap(types.StringType).
				Description("Parameters required by template. Default value will be assigned when parameter has default value and parameter is not set.").
				Optional().
				Computed().
				DefaultEmpty().
				Build(),
			"unmanaged_parameters": util.ResourceMap(types.StringType).
				Description("Template parameters not configured by the practitioner (usually parameters with default value).").
				Computed().
				Build(),
			"template_properties": util.ResourceMap(types.StringType).
				Description("Template properties which will be copied to the step").
				Computed().
				Build(),
			"execution_properties": util.ResourceMap(types.StringType).
				Description("Action properties where the key is the property name and the value is its value.").
				Optional().
				Computed().
				DefaultEmpty().
				Validators(warnAboutReservedExecutionProperties()).
				Build(),
		},
	}
}

func (p ProcessTemplatedChildStepSchema) GetDatasourceSchema() datasourceSchema.Schema {
	return datasourceSchema.Schema{}
}

type ProcessTemplatedChildStepResourceModel struct {
	SpaceID         types.String `tfsdk:"space_id"`
	ProcessID       types.String `tfsdk:"process_id"`
	ParentID        types.String `tfsdk:"parent_id"`
	TemplateID      types.String `tfsdk:"template_id"`
	TemplateVersion types.Int32  `tfsdk:"template_version"`
	Name            types.String `tfsdk:"name"`

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
	Parameters           types.Map                        `tfsdk:"parameters"`
	UnmanagedParameters  types.Map                        `tfsdk:"unmanaged_parameters"`
	TemplateProperties   types.Map                        `tfsdk:"template_properties"`
	ExecutionProperties  types.Map                        `tfsdk:"execution_properties"`

	ResourceModel
}
