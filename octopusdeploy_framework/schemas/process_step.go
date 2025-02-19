package schemas

import (
	"github.com/OctopusDeploy/terraform-provider-octopusdeploy/octopusdeploy_framework/util"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	datasourceSchema "github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	resourceSchema "github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/mapdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/objectdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type ProcessStepSchema struct{}

var _ EntitySchema = ProcessStepSchema{}

const ProcessStepResourceName = "process_step"

func (p ProcessStepSchema) GetResourceSchema() resourceSchema.Schema {
	return resourceSchema.Schema{
		Description: "This resource manages single step of execution process in Octopus Deploy.",
		Attributes: map[string]resourceSchema.Attribute{
			"id":         GetIdResourceSchema(),
			"space_id":   GetSpaceIdResourceSchema(ProcessStepResourceName),
			"process_id": util.ResourceString().Required().Description("Id of the process this step belongs to.").Build(),
			"name":       GetNameResourceSchema(true),
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
			"step_properties": util.ResourceMap(types.StringType).
				Description("A collection of process step properties where the key is the property name and the value is its value.").
				Optional().
				Computed().
				Default(mapdefault.StaticValue(types.MapValueMust(types.StringType, map[string]attr.Value{}))).
				Build(),

			"action_type": util.ResourceString().
				Description("Type of the step action.").
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
			"container": resourceSchema.SingleNestedAttribute{
				Description: "When set used to run step inside a container on the Octopus Server. Octopus Server must support container execution.",
				Attributes: map[string]resourceSchema.Attribute{
					"feed_id": util.ResourceString().
						Description("Feed where the container will be pulled from.").
						Optional().
						Build(),
					"image": util.ResourceString().
						Description("Image of the container with tag included.").
						Optional().
						Build(),
				},
				Optional: true,
				Computed: true,
				Default: objectdefault.StaticValue(
					types.ObjectValueMust(
						map[string]attr.Type{
							"feed_id": types.StringType,
							"image":   types.StringType,
						},
						map[string]attr.Value{
							"feed_id": types.StringValue(""),
							"image":   types.StringValue(""),
						},
					),
				),
			},
			"action_properties": util.ResourceMap(types.StringType).
				Description("A collection of step action properties where the key is the property name and the value is its value.").
				Optional().
				Computed().
				DefaultEmpty().
				Build(),
		},
	}
}

func (p ProcessStepSchema) GetDatasourceSchema() datasourceSchema.Schema {
	return datasourceSchema.Schema{}
}

type ProcessStepResourceModel struct {
	SpaceID            types.String `tfsdk:"space_id"`
	ProcessID          types.String `tfsdk:"process_id"`
	Name               types.String `tfsdk:"name"`
	StartTrigger       types.String `tfsdk:"start_trigger"`
	PackageRequirement types.String `tfsdk:"package_requirement"`
	Condition          types.String `tfsdk:"condition"`
	StepProperties     types.Map    `tfsdk:"step_properties"`

	ActionType           types.String                     `tfsdk:"action_type"`
	Slug                 types.String                     `tfsdk:"slug"`
	IsDisabled           types.Bool                       `tfsdk:"is_disabled"`
	IsRequired           types.Bool                       `tfsdk:"is_required"`
	Notes                types.String                     `tfsdk:"notes"`
	WorkerPoolId         types.String                     `tfsdk:"worker_pool_id"`
	WorkerPoolVariable   types.String                     `tfsdk:"worker_pool_variable"`
	TenantTags           types.Set                        `tfsdk:"tenant_tags"`
	Environments         types.Set                        `tfsdk:"environments"`
	ExcludedEnvironments types.Set                        `tfsdk:"excluded_environments"`
	Channels             types.Set                        `tfsdk:"channels"`
	Container            *ProcessStepActionContainerModel `tfsdk:"container"`
	ActionProperties     types.Map                        `tfsdk:"action_properties"`

	ResourceModel
}

type ProcessStepActionContainerModel struct {
	FeedId types.String `tfsdk:"feed_id"`
	Image  types.String `tfsdk:"image"`
}
