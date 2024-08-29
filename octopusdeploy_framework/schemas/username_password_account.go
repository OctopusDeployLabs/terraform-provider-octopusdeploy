package schemas

import (
	"github.com/OctopusDeploy/terraform-provider-octopusdeploy/octopusdeploy_framework/util"
	datasourceSchema "github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type UsernamePasswordAccountSchema struct{}

var _ EntitySchema = UsernamePasswordAccountSchema{}

func (u UsernamePasswordAccountSchema) GetDatasourceSchema() datasourceSchema.Schema {
	return datasourceSchema.Schema{}
}

func (u UsernamePasswordAccountSchema) GetResourceSchema() schema.Schema {
	return schema.Schema{
		Description: "This resource manages username-password accounts in Octopus Deploy.",
		Attributes: map[string]schema.Attribute{
			"id":                                util.ResourceString().Optional().Computed().PlanModifiers(stringplanmodifier.UseStateForUnknown()).Description("The unique ID for this resource.").Build(),
			"space_id":                          util.ResourceString().Optional().Computed().PlanModifiers(stringplanmodifier.UseStateForUnknown()).Description("The space ID associated with this resource.").Build(),
			"name":                              util.ResourceString().Required().Description("The name of the username-password account.").Build(),
			"description":                       util.ResourceString().Optional().Computed().PlanModifiers(stringplanmodifier.UseStateForUnknown()).Default("").Description("The description of this username/password account.").Build(),
			"environments":                      util.ResourceList(types.StringType).Optional().Computed().Description("A list of environment IDs associated with this resource.").Build(),
			"password":                          util.ResourceString().Optional().Sensitive().Description("The password associated with this resource.").Build(),
			"tenanted_deployment_participation": util.ResourceString().Optional().Optional().Computed().PlanModifiers(stringplanmodifier.UseStateForUnknown()).Description("The tenanted deployment mode of the resource. Valid account types are `Untenanted`, `TenantedOrUntenanted`, or `Tenanted`.").Build(),
			"tenants":                           util.ResourceList(types.StringType).Optional().Computed().Description("A list of tenant IDs associated with this resource.").Build(),
			"tenant_tags":                       util.ResourceList(types.StringType).Optional().Computed().Description("A list of tenant tags associated with this resource.").Build(),
			"username":                          util.ResourceString().Required().Sensitive().Description("The username associated with this resource.").Build(),
		},
	}
}

type UsernamePasswordAccountResourceModel struct {
	SpaceID                         types.String `tfsdk:"space_id"`
	Name                            types.String `tfsdk:"name"`
	Description                     types.String `tfsdk:"description"`
	Environments                    types.List   `tfsdk:"environments"`
	Password                        types.String `tfsdk:"password"`
	TenantedDeploymentParticipation types.String `tfsdk:"tenanted_deployment_participation"`
	Tenants                         types.List   `tfsdk:"tenants"`
	TenantTags                      types.List   `tfsdk:"tenant_tags"`
	Username                        types.String `tfsdk:"username"`

	ResourceModel
}
