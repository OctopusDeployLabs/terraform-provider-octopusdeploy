package schemas

import (
	"github.com/OctopusDeploy/terraform-provider-octopusdeploy/octopusdeploy_framework/util"
	datasourceSchema "github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	resourceSchema "github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type ProjectVersioningStrategySchema struct{}

var _ EntitySchema = ProjectVersioningStrategySchema{}

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
				Required().
				Build(),
			"donor_package_step_id": util.ResourceString().
				Description("The associated donor package step ID.").
				Optional().
				Build(),
			"template": util.ResourceString().
				Optional().
				Computed().
				Build(),
			"donor_package": resourceSchema.SingleNestedAttribute{
				Required:    true,
				Description: "Donor Packages.",
				Attributes: map[string]resourceSchema.Attribute{
					"deployment_action": util.ResourceString().
						Description("Deployment action.").
						Optional().
						Build(),
					"package_reference": util.ResourceString().
						Description("Package reference.").
						Optional().
						Build(),
				},
			},
		},
	}
}

func (p ProjectVersioningStrategySchema) GetDatasourceSchema() datasourceSchema.Schema {
	// no datasource required, returned as part of project datasource
	return datasourceSchema.Schema{}
}

type ProjectVersioningStrategyModel struct {
	ProjectID          types.String      `tfsdk:"project_id"`
	SpaceID            types.String      `tfsdk:"space_id"`
	DonorPackageStepID types.String      `tfsdk:"donor_package_step_id"`
	Template           types.String      `tfsdk:"template"`
	DonorPackage       DonorPackageModel `tfsdk:"donor_package"`
}

type DonorPackageModel struct {
	DeploymentAction types.String `tfsdk:"deployment_action"`
	PackageReference types.String `tfsdk:"package_reference"`
}