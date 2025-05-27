package schemas

import (
	"github.com/OctopusDeploy/terraform-provider-octopusdeploy/octopusdeploy_framework/util"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	datasourceSchema "github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	resourceSchema "github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/mapplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type ProcessTemplatedStepSchema struct{}

var _ EntitySchema = ProcessTemplatedStepSchema{}

const ProcessTemplatedStepResourceName = "process_templated_step"

func (p ProcessTemplatedStepSchema) GetResourceSchema() resourceSchema.Schema {
	return resourceSchema.Schema{
		Description: "This resource manages a single step of a Runbook or Deployment Process which based on existing custom step template",
		Attributes: map[string]resourceSchema.Attribute{
			"id":       GetIdResourceSchema(),
			"space_id": GetSpaceIdResourceSchema(ProcessTemplatedStepResourceName),
			"process_id": util.ResourceString().
				Description("Id of the process this step belongs to.").
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
			"start_trigger": util.ResourceString().
				Description("Whether to run this step after the previous step ('StartAfterPrevious') or at the same time as the previous step ('StartWithPrevious').").
				Optional().
				Computed().
				Default("StartAfterPrevious").
				Validators(stringvalidator.OneOf("StartAfterPrevious", "StartWithPrevious")).
				Build(),
			"package_requirement": util.ResourceString().
				Description("Whether to run this step before or after package acquisition (if possible).").
				Optional().
				Computed().
				Default("LetOctopusDecide").
				Validators(stringvalidator.OneOf("LetOctopusDecide", "AfterPackageAcquisition", "BeforePackageAcquisition")).
				Build(),
			"condition": util.ResourceString().
				Description("When to run the step, one of 'Success', 'Failure', 'Always' or 'Variable'").
				Optional().
				Computed().
				Default("Success").
				Validators(stringvalidator.OneOf("Success", "Failure", "Always", "Variable")).
				Build(),
			"properties": util.ResourceMap(types.StringType).
				Description("A collection of process step properties where the key is the property name and the value is its value.").
				Optional().
				Computed().
				DefaultEmpty().
				Build(),

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
				Description("Properties copied from the template").
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

func (p ProcessTemplatedStepSchema) GetDatasourceSchema() datasourceSchema.Schema {
	return datasourceSchema.Schema{}
}

type ProcessTemplatedStepResourceModel struct {
	SpaceID            types.String `tfsdk:"space_id"`
	ProcessID          types.String `tfsdk:"process_id"`
	TemplateID         types.String `tfsdk:"template_id"`
	TemplateVersion    types.Int32  `tfsdk:"template_version"`
	Name               types.String `tfsdk:"name"`
	StartTrigger       types.String `tfsdk:"start_trigger"`
	PackageRequirement types.String `tfsdk:"package_requirement"`
	Condition          types.String `tfsdk:"condition"`
	Properties         types.Map    `tfsdk:"properties"`

	Type                 types.String                     `tfsdk:"type"`
	Slug                 types.String                     `tfsdk:"slug"`
	IsDisabled           types.Bool                       `tfsdk:"is_disabled"`
	IsRequired           types.Bool                       `tfsdk:"is_required"`
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

type ProcessTemplatedStepGroupedPropertyValues struct {
	TemplateID          types.String
	TemplateVersion     types.Int32
	Parameters          types.Map
	UnmanagedParameters types.Map
	TemplateProperties  types.Map
	ExecutionProperties types.Map
}
