package octopusdeploy

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func getGitTriggerSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"name": getNameSchema(true),
		"space_id": {
			Optional:         true,
			Description:      "The space ID associated with the project to attach the trigger.",
			Type:             schema.TypeString,
			ValidateDiagFunc: validation.ToDiagFunc(validation.StringIsNotWhiteSpace),
		},
		"project_id": {
			Description:      "The ID of the project to attach the trigger.",
			Required:         true,
			Type:             schema.TypeString,
			ValidateDiagFunc: validation.ToDiagFunc(validation.StringIsNotWhiteSpace),
		},
		"channel_id": {
			Description:      "The ID of the channel in which the release will be created if the action type is CreateRelease.",
			Required:         true,
			Type:             schema.TypeString,
			ValidateDiagFunc: validation.ToDiagFunc(validation.StringIsNotWhiteSpace),
		},
		"sources": {
			Description: "List of Git trigger sources. Contains details of the deployment action slug, the git dependency and what file paths to monitor.",
			Optional:    true,
			Type:        schema.TypeList,
			Elem:        &schema.Resource{Schema: getGitTriggerSourceSchema()},
		},
		"is_disabled": {
			Description: "Disables the trigger from being run when set.",
			Optional:    true,
			Default:     false,
			Type:        schema.TypeBool,
		},
	}
}
