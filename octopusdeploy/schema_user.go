package octopusdeploy

import (
	"context"

	"github.com/OctopusDeploy/go-octopusdeploy/octopusdeploy"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func expandUser(d *schema.ResourceData) *octopusdeploy.User {
	username := d.Get("username").(string)
	displayName := d.Get("display_name").(string)

	user := octopusdeploy.NewUser(username, displayName)
	user.ID = d.Id()

	if v, ok := d.GetOk("email_address"); ok {
		user.EmailAddress = v.(string)
	}

	if v, ok := d.GetOk("identity"); ok {
		user.Identities = expandIdentities(v.(*schema.Set).List())
	}

	if v, ok := d.GetOk("is_active"); ok {
		user.IsActive = v.(bool)
	}

	if v, ok := d.GetOk("is_requestor"); ok {
		user.IsRequestor = v.(bool)
	}

	if v, ok := d.GetOk("is_service"); ok {
		user.IsService = v.(bool)
	}

	if v, ok := d.GetOk("password"); ok {
		user.Password = v.(string)
	}

	return user
}

func flattenUser(ctx context.Context, d *schema.ResourceData, user *octopusdeploy.User) {
	d.Set("can_password_be_edited", user.CanPasswordBeEdited)
	d.Set("display_name", user.DisplayName)
	d.Set("email_address", user.EmailAddress)
	d.Set("identity", flattenIdentities(user.Identities))
	d.Set("is_active", user.IsActive)
	d.Set("is_requestor", user.IsRequestor)
	d.Set("is_service", user.IsService)
	d.Set("username", user.Username)

	d.SetId(user.GetID())
}

func getUserSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"can_password_be_edited": {
			Computed: true,
			Type:     schema.TypeBool,
		},
		"display_name": {
			Required: true,
			Type:     schema.TypeString,
		},
		"email_address": {
			Optional: true,
			Type:     schema.TypeString,
		},
		"identity": {
			Optional: true,
			Elem:     &schema.Resource{Schema: getIdentitySchema()},
			Type:     schema.TypeSet,
		},
		"is_active": {
			Optional: true,
			Type:     schema.TypeBool,
		},
		"is_requestor": {
			Computed: true,
			Type:     schema.TypeBool,
		},
		"is_service": {
			Optional: true,
			Type:     schema.TypeBool,
		},
		"username": {
			Required:  true,
			Sensitive: true,
			Type:      schema.TypeString,
		},
		"password": {
			Optional:  true,
			Sensitive: true,
			Type:      schema.TypeString,
		},
	}
}
