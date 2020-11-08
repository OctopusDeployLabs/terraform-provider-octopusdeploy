package octopusdeploy

import (
	"context"
	"time"

	"github.com/OctopusDeploy/go-octopusdeploy/octopusdeploy"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceUser() *schema.Resource {
	resourceUserImporter := &schema.ResourceImporter{
		StateContext: schema.ImportStatePassthroughContext,
	}
	resourceUserSchema := map[string]*schema.Schema{
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
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"provider": {
						Optional: true,
						Type:     schema.TypeString,
					},
					"claim": {
						Optional: true,
						Elem: &schema.Resource{
							Schema: map[string]*schema.Schema{
								"name": {
									Required: true,
									Type:     schema.TypeString,
								},
								"is_identifying_claim": {
									Required: true,
									Type:     schema.TypeBool,
								},
								"value": {
									Required: true,
									Type:     schema.TypeString,
								},
							},
						},
						Type: schema.TypeSet,
					},
				},
			},
			Type: schema.TypeSet,
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
		"modified_by": {
			Computed: true,
			Type:     schema.TypeString,
		},
		"modified_on": {
			Computed: true,
			Type:     schema.TypeString,
		},
		"username": {
			Required: true,
			Type:     schema.TypeString,
		},
		"password": {
			Optional:  true,
			Sensitive: true,
			Type:      schema.TypeString,
		},
	}

	return &schema.Resource{
		CreateContext: resourceUserCreate,
		DeleteContext: resourceUserDelete,
		Importer:      resourceUserImporter,
		ReadContext:   resourceUserRead,
		Schema:        resourceUserSchema,
		UpdateContext: resourceUserUpdate,
	}
}

func resourceUserRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*octopusdeploy.Client)
	user, err := client.Users.GetByID(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	flattenUser(ctx, d, user)
	return nil
}

func resourceUserCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	user := expandUser(d)

	client := m.(*octopusdeploy.Client)
	createdUser, err := client.Users.Add(user)
	if err != nil {
		return diag.FromErr(err)
	}

	flattenUser(ctx, d, createdUser)
	return nil
}

func resourceUserUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	user := expandUser(d)

	client := m.(*octopusdeploy.Client)
	updatedUser, err := client.Users.Update(user)
	if err != nil {
		return diag.FromErr(err)
	}

	flattenUser(ctx, d, updatedUser)
	return nil
}

func resourceUserDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*octopusdeploy.Client)
	err := client.Users.DeleteByID(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId("")
	return nil
}

func expandUser(d *schema.ResourceData) *octopusdeploy.User {
	var username string
	if v, ok := d.GetOk("username"); ok {
		username = v.(string)
	}

	var displayName string
	if v, ok := d.GetOk("display_name"); ok {
		displayName = v.(string)
	}

	user := octopusdeploy.NewUser(username, displayName)
	user.ID = d.Id()

	if v, ok := d.GetOk("can_password_be_edited"); ok {
		user.CanPasswordBeEdited = v.(bool)
	}

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

	if v, ok := d.GetOk("modified_by"); ok {
		user.ModifiedBy = v.(string)
	}

	if v, ok := d.GetOk("modified_on"); ok {
		modifiedOnTime, _ := time.Parse(time.RFC3339, v.(string))
		user.ModifiedOn = &modifiedOnTime
	}

	if v, ok := d.GetOk("password"); ok {
		user.Password = v.(string)
	}

	return user
}

func expandClaims(claims []interface{}) map[string]octopusdeploy.IdentityClaim {
	expandedClaims := make(map[string]octopusdeploy.IdentityClaim, len(claims))
	for _, claim := range claims {
		claimMap := claim.(map[string]interface{})
		name := claimMap["name"].(string)
		identityClaim := octopusdeploy.IdentityClaim{
			IsIdentifyingClaim: claimMap["is_identifying_claim"].(bool),
			Value:              claimMap["value"].(string),
		}
		expandedClaims[name] = identityClaim
	}
	return expandedClaims
}

func expandIdentities(identities []interface{}) []octopusdeploy.Identity {
	expandedIdentities := make([]octopusdeploy.Identity, 0, len(identities))
	for _, identity := range identities {
		if identity != nil {
			rawIdentity := identity.(map[string]interface{})

			identityProviderName := ""
			if rawIdentity["provider"] != nil {
				identityProviderName = rawIdentity["provider"].(string)
			}

			i := octopusdeploy.Identity{
				IdentityProviderName: identityProviderName,
				Claims:               expandClaims(rawIdentity["claim"].(*schema.Set).List()),
			}
			expandedIdentities = append(expandedIdentities, i)
		}
	}
	return expandedIdentities
}

func flattenIdentityClaims(identityClaims map[string]octopusdeploy.IdentityClaim) []interface{} {
	if identityClaims == nil {
		return nil
	}

	flattenedIdentityClaims := []interface{}{}
	for key, identityClaim := range identityClaims {
		rawIdentityClaim := map[string]interface{}{
			"is_identifying_claim": identityClaim.IsIdentifyingClaim,
			"name":                 key,
			"value":                identityClaim.Value,
		}

		flattenedIdentityClaims = append(flattenedIdentityClaims, rawIdentityClaim)
	}

	return flattenedIdentityClaims
}

func flattenIdentities(identities []octopusdeploy.Identity) []interface{} {
	if identities == nil {
		return nil
	}

	var flattenedIdentities = make([]interface{}, len(identities))
	for i, identity := range identities {
		rawIdentity := map[string]interface{}{
			"provider": identity.IdentityProviderName,
		}
		if identity.Claims != nil {
			rawIdentity["claim"] = flattenIdentityClaims(identity.Claims)
		}

		flattenedIdentities[i] = rawIdentity
	}

	return flattenedIdentities
}

func flattenUser(ctx context.Context, d *schema.ResourceData, user *octopusdeploy.User) {
	d.Set("can_password_be_edited", user.CanPasswordBeEdited)
	d.Set("display_name", user.DisplayName)
	d.Set("email_address", user.EmailAddress)
	d.Set("identity", flattenIdentities(user.Identities))
	d.Set("is_active", user.IsActive)
	d.Set("is_requestor", user.IsRequestor)
	d.Set("is_service", user.IsService)
	d.Set("modified_by", user.ModifiedBy)

	if user.ModifiedOn != nil {
		d.Set("modified_on", user.ModifiedOn.Format(time.RFC3339))
	}

	d.Set("password", user.Password)
	d.Set("username", user.Username)

	d.SetId(user.GetID())
}
