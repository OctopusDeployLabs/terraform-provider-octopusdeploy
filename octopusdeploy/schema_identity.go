package octopusdeploy

import (
	"github.com/OctopusDeploy/go-octopusdeploy/octopusdeploy"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

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
				Claims:               expandIdentityClaims(rawIdentity["claim"].(*schema.Set).List()),
			}
			expandedIdentities = append(expandedIdentities, i)
		}
	}
	return expandedIdentities
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

func getIdentitySchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"provider": {
			Optional: true,
			Type:     schema.TypeString,
		},
		"claim": {
			Optional: true,
			Elem:     &schema.Resource{Schema: getIdentityClaimSchema()},
			Type:     schema.TypeSet,
		},
	}
}
