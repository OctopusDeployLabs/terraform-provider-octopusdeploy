package schemas

import (
	"github.com/OctopusDeploy/terraform-provider-octopusdeploy/octopusdeploy_framework/util"
	"github.com/hashicorp/terraform-plugin-framework-validators/resourcevalidator"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	datasourceSchema "github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	resourceSchema "github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/objectplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type ProjectVersioningStrategySchema struct{}

var _ EntitySchemaWithResourceValidators = ProjectVersioningStrategySchema{}

const ProjectVersioningStrategyResourceName = "project_versioning_strategy"

func (p ProjectVersioningStrategySchema) GetResourceSchema() resourceSchema.Schema {
	return resourceSchema.Schema{
		Attributes: map[string]resourceSchema.Attribute{
			"project_id": util.ResourceString().
				Description("The associated project ID.").
				PlanModifiers(stringplanmodifier.RequiresReplace()).
				Required().
				Build(),
			"space_id": util.ResourceString().
				Description("Space ID of the associated project.").
				Optional().
				Build(),
			"donor_package_step_id": util.ResourceString().
				Description("The associated donor package step ID.").
				Optional().
				Computed().
				PlanModifiers(stringplanmodifier.RequiresReplace()).
				Build(),
			"template": util.ResourceString().
				Optional().
				Computed().
				PlanModifiers(stringplanmodifier.UseStateForUnknown()).
				Build(),
			"donor_package": resourceSchema.SingleNestedAttribute{
				Optional:    true,
				Computed:    true,
				Description: "Donor Packages.",
				Attributes: map[string]resourceSchema.Attribute{
					"deployment_action": util.ResourceString().
						Description("Deployment action.").
						Required().
						Build(),
					"package_reference": util.ResourceString().
						Description("Package reference.").
						Required().
						Build(),
				},
				PlanModifiers: []planmodifier.Object{
					objectplanmodifier.UseStateForUnknown(),
				},
			},
		},
	}
}

func (p ProjectVersioningStrategySchema) GetResourceConfigValidators() []resource.ConfigValidator {
	return []resource.ConfigValidator{
		resourcevalidator.Conflicting(
			path.MatchRelative().AtName("donor_package").AtName("deployment_action"),
			path.MatchRoot("donor_package_step_id"),
		),
		resourcevalidator.Conflicting(
			path.MatchRoot("template"),
			path.MatchRoot("donor_package"),
		),
		resourcevalidator.Conflicting(
			path.MatchRoot("template"),
			path.MatchRoot("donor_package_step_id"),
		),
	}
}

func (p ProjectVersioningStrategySchema) GetDatasourceSchema() datasourceSchema.Schema {
	// no datasource required, returned as part of project datasource
	return datasourceSchema.Schema{}
}

type ProjectVersioningStrategyModel struct {
	ProjectID          types.String `tfsdk:"project_id"`
	SpaceID            types.String `tfsdk:"space_id"`
	DonorPackageStepID types.String `tfsdk:"donor_package_step_id"`
	Template           types.String `tfsdk:"template"`
	DonorPackage       types.Object `tfsdk:"donor_package"`
}

type DonorPackageModel struct {
	DeploymentAction types.String `tfsdk:"deployment_action"`
	PackageReference types.String `tfsdk:"package_reference"`
}

func ProjectVersioningStrategyDonorPackageAttributeTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"deployment_action": types.StringType,
		"package_reference": types.StringType,
	}
}
