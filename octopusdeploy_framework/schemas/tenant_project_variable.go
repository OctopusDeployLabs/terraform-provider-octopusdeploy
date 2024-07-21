package schemas

import (
	"github.com/OctopusDeploy/terraform-provider-octopusdeploy/octopusdeploy_framework/util"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
)

const (
	TenantProjectVariableResourceDescription = "Tenant Project Variable"
	TenantProjectVariableResourceName        = "tenant_project_variable"
)

func GetTenantProjectVariableResourceSchema() schema.Schema {
	return schema.Schema{
		Description: "Manages a tenant project variable in Octopus Deploy.",
		Attributes: map[string]schema.Attribute{
			"id":             util.GetIdResourceSchema(),
			"space_id":       util.GetSpaceIdResourceSchema(TenantProjectVariableResourceDescription),
			"tenant_id":      util.GetRequiredStringResourceSchema("The ID of the tenant."),
			"project_id":     util.GetRequiredStringResourceSchema("The ID of the project."),
			"environment_id": util.GetRequiredStringResourceSchema("The ID of the environment."),
			"template_id":    util.GetRequiredStringResourceSchema("The ID of the variable template."),
			"value": schema.StringAttribute{
				Required:    true,
				Description: "The value of the variable.",
			},
		},
	}
}
