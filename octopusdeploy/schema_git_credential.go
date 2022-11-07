package octopusdeploy

import (
	"context"

	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/core"
	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/credentials"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func expandGitCredential(d *schema.ResourceData) *credentials.Resource {
	password := core.NewSensitiveValue(d.Get("password").(string))
	name := d.Get("name").(string)
	username := d.Get("username").(string)

	usernamePassword := credentials.NewUsernamePassword(username, password)

	resource := credentials.NewResource(name, usernamePassword)
	resource.ID = d.Id()

	if v, ok := d.GetOk("description"); ok {
		resource.Description = v.(string)
	}

	if v, ok := d.GetOk("space_id"); ok {
		resource.SpaceID = v.(string)
	}

	return resource
}

func getGitCredentialSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"description": {
			Description: "The description of this Git credential.",
			Optional:    true,
			Type:        schema.TypeString,
		},
		"name": {
			Description:      "The name of the Git credential. This name must be unique.",
			Required:         true,
			Type:             schema.TypeString,
			ValidateDiagFunc: validation.ToDiagFunc(validation.StringIsNotWhiteSpace),
		},
		"password": {
			Description:      "The password for the Git credential.",
			Required:         true,
			Sensitive:        true,
			Type:             schema.TypeString,
			ValidateDiagFunc: validation.ToDiagFunc(validation.StringIsNotEmpty),
		},
		"space_id": getSpaceIDSchema(),
		"username": {
			Description:      "The username for the Git credential.",
			Required:         true,
			Type:             schema.TypeString,
			ValidateDiagFunc: validation.ToDiagFunc(validation.StringIsNotWhiteSpace),
		},
	}
}

func setGitCredential(ctx context.Context, d *schema.ResourceData, resource *credentials.Resource) error {
	d.Set("description", resource.Description)
	d.Set("name", resource.GetName())
	d.Set("space_id", resource.SpaceID)

	usernamePassword := resource.Details.(*credentials.UsernamePassword)

	d.Set("username", usernamePassword.Username)

	return nil
}
