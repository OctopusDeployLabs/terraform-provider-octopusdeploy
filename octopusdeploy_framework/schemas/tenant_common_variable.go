package schemas

import (
	"github.com/OctopusDeploy/terraform-provider-octopusdeploy/octopusdeploy_framework/util"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
)

const (
	TenantCommonVariableResourceDescription = "Tenant Common Variable"
	TenantCommonVariableResourceName        = "tenant_common_variable"
)

func GetTenantCommonVariableResourceSchema() schema.Schema {
	return schema.Schema{
		Description: "Manages a tenant common variable in Octopus Deploy.",
		Attributes: map[string]schema.Attribute{
			"id":                      util.GetIdResourceSchema(),
			"space_id":                util.GetSpaceIdResourceSchema(TenantCommonVariableResourceDescription),
			"tenant_id":               util.GetRequiredStringResourceSchema("The ID of the tenant."),
			"library_variable_set_id": util.GetRequiredStringResourceSchema("The ID of the library variable set."),
			"template_id":             util.GetRequiredStringResourceSchema("The ID of the variable template."),
			"value": schema.StringAttribute{
				Required:    true,
				Description: "The value of the variable.",
				Sensitive:   true,
			},
		},
	}
}
