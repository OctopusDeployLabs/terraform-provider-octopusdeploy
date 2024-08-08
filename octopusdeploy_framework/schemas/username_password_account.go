package schemas

import (
	"github.com/OctopusDeploy/terraform-provider-octopusdeploy/octopusdeploy_framework/util"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func GetUsernamePasswordAccountResourceSchema() schema.Schema {
	return schema.Schema{
		Description: "This resource manages username-password accounts in Octopus Deploy.",
		Attributes: map[string]schema.Attribute{
			"id":                                util.ResourceString().Optional().Computed().PlanModifiers(stringplanmodifier.UseStateForUnknown()).Description("The unique ID for this resource.").Build(),
			"space_id":                          util.ResourceString().Optional().Computed().PlanModifiers(stringplanmodifier.UseStateForUnknown()).Description("The space ID associated with this resource.").Build(),
			"name":                              util.ResourceString().Required().Description("The name of the username-password account.").Build(),
			"description":                       util.ResourceString().Optional().Description("The description of this username/password resource.").Build(),
			"environments":                      util.ResourceList(types.StringType).Optional().Description("A list of environment IDs associated with this resource.").Build(),
			"password":                          util.ResourceString().Optional().Sensitive().Description("The password associated with this resource.").Build(),
			"tenanted_deployment_participation": util.ResourceString().Optional().Description("The tenanted deployment mode of the resource. Valid account types are `Untenanted`, `TenantedOrUntenanted`, or `Tenanted`.").Build(),
			"tenants":                           util.ResourceList(types.StringType).Optional().Description("A list of tenant IDs associated with this resource.").Build(),
			"tenant_tags":                       util.ResourceList(types.StringType).Optional().Description("A list of tenant tags associated with this resource.").Build(),
			"username":                          util.ResourceString().Required().Sensitive().Description("The username associated with this resource.").Build(),
		},
	}
}
