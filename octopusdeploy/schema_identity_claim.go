package octopusdeploy

import (
	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/users"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func expandIdentityClaims(claims []interface{}) map[string]users.IdentityClaim {
	expandedClaims := make(map[string]users.IdentityClaim, len(claims))
	for _, claim := range claims {
		claimMap := claim.(map[string]interface{})
		name := claimMap["name"].(string)
		identityClaim := users.IdentityClaim{
			IsIdentifyingClaim: claimMap["is_identifying_claim"].(bool),
			Value:              claimMap["value"].(string),
		}
		expandedClaims[name] = identityClaim
	}
	return expandedClaims
}

func flattenIdentityClaims(identityClaims map[string]users.IdentityClaim) []interface{} {
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

func getIdentityClaimSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"name": getNameSchema(true),
		"is_identifying_claim": {
			Required: true,
			Type:     schema.TypeBool,
		},
		"value": {
			Required: true,
			Type:     schema.TypeString,
		},
	}
}
