package octopusdeploy

import "github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

func getTenantProjectEnvironmentSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"tenant_id": {
			Description: "The tenant ID associated with this tenant.",
			Required:    true,
			Type:        schema.TypeString,
			ForceNew:    true,
		},
		"project_id": {
			Description: "The project ID associated with this tenant.",
			Required:    true,
			Type:        schema.TypeString,
			ForceNew:    true,
		},
		"environment_id": {
			Description: "The environment ID associated with this tenant.",
			Required:    true,
			Type:        schema.TypeString,
			ForceNew:    true,
		},
		"space_id": getSpaceIDSchema(),
	}
}
