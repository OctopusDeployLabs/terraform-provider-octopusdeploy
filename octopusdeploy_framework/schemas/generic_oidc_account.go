package schemas

import (
	"github.com/OctopusDeploy/terraform-provider-octopusdeploy/octopusdeploy_framework/util"
	datasourceSchema "github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	resourceSchema "github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type GenericOidcAccountSchema struct{}

var _ EntitySchema = GenericOidcAccountSchema{}

func (a GenericOidcAccountSchema) GetDatasourceSchema() datasourceSchema.Schema {
	return datasourceSchema.Schema{}
}

func (a GenericOidcAccountSchema) GetResourceSchema() resourceSchema.Schema {
	return resourceSchema.Schema{
		Description: "This resource manages a Generic OIDC Account in Octopus Deploy.",
		Attributes: map[string]resourceSchema.Attribute{
			"description": util.ResourceString().
				Optional().
				Computed().
				PlanModifiers(stringplanmodifier.UseStateForUnknown()).
				Default("").
				Description("The description of this generic oidc account.").
				Build(),
			"environments": util.ResourceList(types.StringType).
				Optional().
				Computed().
				Description("A list of environment IDs associated with this resource.").
				Build(),
			"id": GetIdResourceSchema(),
			"name": util.ResourceString().
				Required().
				Description("The name of the generic oidc account.").
				Build(),
			"space_id": util.ResourceString().
				Optional().
				Computed().
				PlanModifiers(stringplanmodifier.UseStateForUnknown()).
				Description("The space ID associated with this resource.").
				Build(),
			"tenanted_deployment_participation": util.ResourceString().
				Optional().
				Computed().
				PlanModifiers(stringplanmodifier.UseStateForUnknown()).
				Description("The tenanted deployment mode of the resource. Valid account types are `Untenanted`, `TenantedOrUntenanted`, or `Tenanted`.").
				Build(),
			"tenants": util.ResourceList(types.StringType).
				Optional().
				Computed().
				Description("A list of tenant IDs associated with this resource.").
				Build(),
			"tenant_tags": util.ResourceList(types.StringType).
				Optional().
				Computed().
				Description("A list of tenant tags associated with this resource.").
				Build(),
			"execution_subject_keys": util.ResourceList(types.StringType).
				Optional().
				Description("Keys to include in a deployment or runbook. Valid options are `space`, `environment`, `project`, `tenant`, `runbook`, `account`, `type`.").
				Build(),
			"audience": util.ResourceString().
				Optional().
				Description("The audience associated with this resource.").
				Build(),
		},
	}
}

type GenericOidcAccountResourceModel struct {
	Description                     types.String `tfsdk:"description"`
	Environments                    types.List   `tfsdk:"environments"`
	Name                            types.String `tfsdk:"name"`
	SpaceID                         types.String `tfsdk:"space_id"`
	TenantedDeploymentParticipation types.String `tfsdk:"tenanted_deployment_participation"`
	Tenants                         types.List   `tfsdk:"tenants"`
	TenantTags                      types.List   `tfsdk:"tenant_tags"`
	ExecutionSubjectKeys            types.List   `tfsdk:"execution_subject_keys"`
	Audience                        types.String `tfsdk:"audience"`

	ResourceModel
}
