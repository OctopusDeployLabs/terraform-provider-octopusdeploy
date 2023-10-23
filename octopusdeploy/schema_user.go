package octopusdeploy

import (
	"context"
	"fmt"

	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/users"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func expandUser(d *schema.ResourceData) *users.User {
	username := d.Get("username").(string)
	displayName := d.Get("display_name").(string)

	user := users.NewUser(username, displayName)
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

func flattenUser(user *users.User) map[string]interface{} {
	if user == nil {
		return nil
	}

	return map[string]interface{}{
		"can_password_be_edited": user.CanPasswordBeEdited,
		"display_name":           user.DisplayName,
		"email_address":          user.EmailAddress,
		"id":                     user.GetID(),
		"identity":               flattenIdentities(user.Identities),
		"is_active":              user.IsActive,
		"is_service":             user.IsService,
		"username":               user.Username,
	}
}

func getUserDataSchema() map[string]*schema.Schema {
	dataSchema := getUserSchema()
	setDataSchema(&dataSchema)

	return map[string]*schema.Schema{
		"filter":   getQueryFilter(),
		"id":       getDataSchemaID(),
		"ids":      getQueryIDs(),
		"skip":     getQuerySkip(),
		"take":     getQueryTake(),
		"space_id": getQuerySpaceID(),
		"users": {
			Computed:    true,
			Description: "A list of users that match the filter(s).",
			Elem:        &schema.Resource{Schema: dataSchema},
			Optional:    true,
			Type:        schema.TypeList,
		},
	}
}

func getUserSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"can_password_be_edited": {
			Computed: true,
			Type:     schema.TypeBool,
		},
		"display_name":  getDisplayNameSchema(true),
		"email_address": getEmailAddressSchema(false),
		"id":            getIDSchema(),
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
		"username": getUsernameSchema(true),
		"password": getPasswordSchema(false),
	}
}

func setUser(ctx context.Context, d *schema.ResourceData, user *users.User) error {
	d.Set("can_password_be_edited", user.CanPasswordBeEdited)
	d.Set("display_name", user.DisplayName)
	d.Set("email_address", user.EmailAddress)
	d.Set("id", user.GetID())

	if err := d.Set("identity", flattenIdentities(user.Identities)); err != nil {
		return fmt.Errorf("error setting identity: %s", err)
	}

	d.Set("is_active", user.IsActive)
	d.Set("is_requestor", user.IsRequestor)
	d.Set("is_service", user.IsService)
	d.Set("username", user.Username)

	d.SetId(user.GetID())

	return nil
}
