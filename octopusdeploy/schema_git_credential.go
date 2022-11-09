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

func flattenGitCredential(credential *credentials.Resource) map[string]interface{} {
	if credential == nil {
		return nil
	}

	return map[string]interface{}{
		"id":          credential.GetID(),
		"name":        credential.Name,
		"description": credential.Description,
		"type":        credential.Details.Type(),
	}
}

func getGitCredentialDataSchema() map[string]*schema.Schema {
	dataSchema := getGitCredentialSchema()
	setDataSchema(&dataSchema)

	return map[string]*schema.Schema{
		"git_credentials": {
			Computed:    true,
			Description: "A list of Git Credentials that match the filter(s).",
			Elem:        &schema.Resource{Schema: dataSchema},
			Optional:    true,
			Type:        schema.TypeList,
		},
		"first_result": getQueryFirstResult(),
		"name":         getQueryName(),
		"skip":         getQuerySkip(),
		"take":         getQueryTake(),
	}
}

func getGitCredentialSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"space_id": getSpaceIDSchema(),
		"name": {
			Description:      "The name of the Git credential. This name must be unique.",
			Required:         true,
			Type:             schema.TypeString,
			ValidateDiagFunc: validation.ToDiagFunc(validation.StringIsNotWhiteSpace),
		},
		"description": {
			Description: "The description of this Git credential.",
			Optional:    true,
			Type:        schema.TypeString,
		},
		"username": {
			Description:      "The username for the Git credential.",
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
	}
}

func setGitCredential(ctx context.Context, d *schema.ResourceData, resource *credentials.Resource) error {
	d.Set("space_id", resource.SpaceID)
	d.Set("name", resource.GetName())
	d.Set("description", resource.Description)

	usernamePassword := resource.Details.(*credentials.UsernamePassword)
	d.Set("username", usernamePassword.Username)

	return nil
}
