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
		Description: "This resource manages a single step of a Runbook or Deployment Process in Octopus Deploy.",
		Attributes: map[string]resourceSchema.Attribute{
			"id":       GetIdResourceSchema(),
			"space_id": GetSpaceIdResourceSchema(ProcessStepResourceName),
			"process_id": util.ResourceString().
				Description("Id of the process this step belongs to.").
				Required().
				PlanModifiers(stringplanmodifier.RequiresReplace()).
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
	ExecutionProperties  types.Map                        `tfsdk:"execution_properties"`

	ResourceModel
}

type ProcessStepActionContainerModel struct {
	FeedID types.String `tfsdk:"feed_id"`
	Image  types.String `tfsdk:"image"`
}

type ProcessStepPackageReferenceResourceModel struct {
	PackageID           types.String `tfsdk:"package_id"`
	FeedID              types.String `tfsdk:"feed_id"`
	AcquisitionLocation types.String `tfsdk:"acquisition_location"`
	Properties          types.Map    `tfsdk:"properties"`

	ResourceModel
}

type ProcessStepGitDependencyResourceModel struct {
	RepositoryUri     types.String `tfsdk:"repository_uri"`
	DefaultBranch     types.String `tfsdk:"default_branch"`
	GitCredentialType types.String `tfsdk:"git_credential_type"`
	FilePathFilters   types.Set    `tfsdk:"file_path_filters"`
	GitCredentialID   types.String `tfsdk:"git_credential_id"`
}

func ProcessStepPackageReferenceObjectType() types.ObjectType {
	return types.ObjectType{
		AttrTypes: ProcessStepPackageReferenceAttributeTypes(),
	}
}

func ProcessStepPackageReferenceAttributeTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"id":                   types.StringType,
		"package_id":           types.StringType,
		"feed_id":              types.StringType,
		"acquisition_location": types.StringType,
		"properties":           types.MapType{ElemType: types.StringType},
	}
}

func ProcessStepGitDependencyObjectType() types.ObjectType {
	return types.ObjectType{
		AttrTypes: ProcessStepGitDependencyAttributeTypes(),
	}
}

func ProcessStepGitDependencyAttributeTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"repository_uri":      types.StringType,
		"default_branch":      types.StringType,
		"git_credential_type": types.StringType,
		"file_path_filters":   types.SetType{ElemType: types.StringType},
		"git_credential_id":   types.StringType,
	}
}

func resourceActionContainerAttribute() resourceSchema.SingleNestedAttribute {
	return resourceSchema.SingleNestedAttribute{
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
	}
}

func resourceActionGitDependenciesAttribute() resourceSchema.MapNestedAttribute {
	return resourceSchema.MapNestedAttribute{
		Description:  "References of git dependencies for this step where key is a name of the reference (can be empty). Is the Git equivalent of packages",
		Optional:     true,
		Computed:     true,
		Default:      mapdefault.StaticValue(types.MapValueMust(ProcessStepGitDependencyObjectType(), map[string]attr.Value{})),
		NestedObject: resourceActionGitDependencyNestedAttribute(),
	}
}

func resourceActionGitDependencyNestedAttribute() resourceSchema.NestedAttributeObject {
	return resourceSchema.NestedAttributeObject{
		Attributes: map[string]resourceSchema.Attribute{
			"repository_uri": util.ResourceString().
				Description("The Git URI for the repository where this resource is sourced from").
				Required().
				Build(),
			"default_branch": util.ResourceString().
				Description("Name of the default branch of the repository.").
				Required().
				Build(),
			"git_credential_type": util.ResourceString().
				Description("The Git credential authentication type.").
				Required().
				Build(),
			"file_path_filters": util.ResourceSet(types.StringType).
				Description("List of file path filters used to narrow down the directory where files are to be sourced from. Supports glob patten syntax.").
				Optional().
				Computed().
				DefaultEmpty().
				Build(),
			"git_credential_id": util.ResourceString().
				Description("ID of an existing Git credential.").
				Optional().
				Computed().
				Default("").
				Build(),
		},
	}
}

func resourceActionPackageReferencesAttribute() resourceSchema.MapNestedAttribute {
	return resourceSchema.MapNestedAttribute{
		Description:  "Package references associated with this step where key is a name of the package reference (use empty name for primary package)",
		Optional:     true,
		Computed:     true,
		Default:      mapdefault.StaticValue(types.MapValueMust(ProcessStepPackageReferenceObjectType(), map[string]attr.Value{})),
		NestedObject: resourceActionPackageReferenceNestedAttribute(),
	}
}

func resourceActionPackageReferenceNestedAttribute() resourceSchema.NestedAttributeObject {
	return resourceSchema.NestedAttributeObject{
		Attributes: map[string]resourceSchema.Attribute{
			"id": GetIdResourceSchema(),
			"package_id": util.ResourceString().
				Description("Package ID or a variable-expression").
				Required().
				Build(),
			"feed_id": util.ResourceString().
				Description("The feed ID associated with this package reference").
				Optional().
				Computed().
				Default("feeds-builtin").
				Build(),
			"acquisition_location": util.ResourceString().
				Description("Whether to acquire this package on the server ('Server'), target ('ExecutionTarget') or not at all ('NotAcquired'). Can be an expression").
				Optional().
				Computed().
				Default("Server").
				Build(),
			"properties": util.ResourceMap(types.StringType).
				Description("A collection of properties associated with this package").
				Optional().
				Computed().
				DefaultEmpty().
				Build(),
		},
	}
}
