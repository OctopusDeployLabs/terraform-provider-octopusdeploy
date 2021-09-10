package octopusdeploy

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func getProjectDeploymentTargetTriggerSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"name": getNameSchema(true),
		"project_id": {
			Description: "The ID of the project to attach the trigger.",
			Required:    true,
			Type:        schema.TypeString,
		},
		"should_redeploy": {
			Default:     false,
			Description: "Enable to re-deploy to the deployment targets even if they are already up-to-date with the current deployment.",
			Optional:    true,
			Type:        schema.TypeBool,
		},
		"event_groups": {
			Description: "Apply event group filters to restrict which deployment targets will actually cause the trigger to fire, and consequently, which deployment targets will be automatically deployed to.",
			Elem:        &schema.Schema{Type: schema.TypeString},
			Optional:    true,
			Type:        schema.TypeList,
		},
		"event_categories": {
			Description: "Apply event category filters to restrict which deployment targets will actually cause the trigger to fire, and consequently, which deployment targets will be automatically deployed to.",
			Elem:        &schema.Schema{Type: schema.TypeString},
			Optional:    true,
			Type:        schema.TypeList,
		},
		"roles": {
			Description: "Apply event role filters to restrict which deployment targets will actually cause the trigger to fire, and consequently, which deployment targets will be automatically deployed to.",
			Elem:        &schema.Schema{Type: schema.TypeString},
			Optional:    true,
			Type:        schema.TypeList,
		},
		"environment_ids": {
			Description: "Apply environment id filters to restrict which deployment targets will actually cause the trigger to fire, and consequently, which deployment targets will be automatically deployed to.",
			Elem:        &schema.Schema{Type: schema.TypeString},
			Optional:    true,
			Type:        schema.TypeList,
		},
	}
}
