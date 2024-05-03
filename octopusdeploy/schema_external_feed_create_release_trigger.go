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
			Description: "List of referenced package that will cause the trigger to fire. New versions of any of the packages you select will trigger release creation.
 ",
			Optional:    true,
			Type:        schema.TypeList,
			Elem:        &schema.Resource{Schema: getDeploymentActionSlugPackageSchema()},
		},
		"primary_package": {
			Description: "List of deployment actions for which the primary packages will cause the trigger to fire. New versions of any of the packages you select will trigger release creation.",
			Optional:    true,
			Type:        schema.TypeList,
			Elem:        &schema.Resource{Schema: getDeploymentActionSlugPrimaryPackageSchema()},
		},
		"is_disabled": {
			Description: "Disables the trigger from being run when set.",
			Optional:    true,
			Default:     false,
			Type:        schema.TypeBool,
		},
	}
}
