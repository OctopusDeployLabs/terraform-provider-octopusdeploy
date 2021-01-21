package octopusdeploy

import (
	"github.com/OctopusDeploy/go-octopusdeploy/octopusdeploy"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func expandIdentityClaims(claims []interface{}) map[string]octopusdeploy.IdentityClaim {
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
