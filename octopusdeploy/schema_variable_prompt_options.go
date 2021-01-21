package octopusdeploy

import "github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

func getVariablePromptOptionsSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"description": getDescriptionSchema(),
		"is_required": {
			Type:     schema.TypeBool,
			Optional: true,
		},
		"label": {
			Type:     schema.TypeString,
			Optional: true,
		},
	}
}
