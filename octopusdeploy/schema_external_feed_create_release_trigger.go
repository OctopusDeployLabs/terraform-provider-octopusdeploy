package octopusdeploy

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func getExternalFeedCreateReleaseTriggerSchema() map[string]*schema.Schema {
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
		"package": {
			Description: "List of package references that will cause the trigger to fire. The triggering condition is if any of the packages are updated.",
			Required:    true,
			Type:        schema.TypeList,
			Elem:        &schema.Resource{Schema: getDeploymentActionSlugPackageSchema()},
		},
		"is_disabled": {
			Description: "Disables the trigger from being run when set.",
			Optional:    true,
			Default:     false,
			Type:        schema.TypeBool,
		},
	}
}
