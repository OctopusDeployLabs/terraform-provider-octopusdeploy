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

	flattenedIdentities := make([]interface{}, len(identities))
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
