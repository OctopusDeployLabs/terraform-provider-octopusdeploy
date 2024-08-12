package schemas

import (
	"github.com/OctopusDeploy/terraform-provider-octopusdeploy/octopusdeploy_framework/util"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
)

const (
	TenantProjectVariableResourceDescription = "Tenant Project Variable"
	TenantProjectVariableResourceName        = "tenant_project_variable"
)

func GetTenantProjectVariableResourceSchema() schema.Schema {
	return schema.Schema{
		Description: "Manages a tenant project variable in Octopus Deploy.",
		Attributes: map[string]schema.Attribute{
			"id": util.ResourceString().
				Computed().
				Description("The unique ID for this resource.").
				PlanModifiers(stringplanmodifier.UseStateForUnknown()).
				Build(),
			"space_id": util.ResourceString().
				Optional().
				Computed().
				Description("The space ID associated with this Tenant Project Variable.").
				PlanModifiers(stringplanmodifier.UseStateForUnknown()).
				Build(),
			"tenant_id": util.ResourceString().
				Required().
				Description("The ID of the tenant.").
				Build(),
			"project_id": util.ResourceString().
				Required().
				Description("The ID of the project.").
				Build(),
			"environment_id": util.ResourceString().
				Required().
				Description("The ID of the environment.").
				Build(),
			"template_id": util.ResourceString().
				Required().
				Description("The ID of the variable template.").
				Build(),
			"value": util.ResourceString().
				Optional().
				Sensitive().
				Description("The value of the variable.").
				Build(),
		},
	}
}
