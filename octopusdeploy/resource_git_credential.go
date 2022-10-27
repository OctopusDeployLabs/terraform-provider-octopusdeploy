package octopusdeploy

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func getGitCredentialSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"description": {
			Description: "The description of this Git credential.",
			Optional:    true,
			Type:        schema.TypeString,
		},
		"details": {
			Description: "Contains the credentials for the Git credential.",
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"password": {
						Description:      "The password for the Git credential.",
						Required:         true,
						Sensitive:        true,
						Type:             schema.TypeString,
						ValidateDiagFunc: validation.ToDiagFunc(validation.StringIsNotEmpty),
					},
					"username": {
						Description:      "The username for the Git credential.",
						Required:         true,
						Type:             schema.TypeString,
						ValidateDiagFunc: validation.ToDiagFunc(validation.StringIsNotWhiteSpace),
					},
				},
			},
		},
		"id": getIDSchema(),
		"name": {
			Description:      "The name of the Git credential. This name must be unique.",
			Required:         true,
			Type:             schema.TypeString,
			ValidateDiagFunc: validation.ToDiagFunc(validation.StringIsNotWhiteSpace),
		},
	}
}
