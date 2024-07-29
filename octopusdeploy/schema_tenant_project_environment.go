package octopusdeploy

import "github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

func getTenantProjectSchema() map[string]*schema.Schema {
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
		"environment_ids": {
			Description: "The environment ID associated with this tenant.",
			Elem:        &schema.Schema{Type: schema.TypeString},
			Optional:    true,
			Type:        schema.TypeList,
			Required:    false,
		},
		"space_id": getSpaceIDSchema(),
	}
}
