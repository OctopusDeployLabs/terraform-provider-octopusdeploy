package octopusdeploy

import (
	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/core"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func expandExternalSecurityGroups(externalSecurityGroups []interface{}) []core.NamedReferenceItem {
	expandedExternalSecurityGroups := []core.NamedReferenceItem{}
	for _, externalSecurityGroup := range externalSecurityGroups {
		if externalSecurityGroup != nil {
			rawExternalSecurityGroup := externalSecurityGroup.(map[string]interface{})

			displayIDAndName := false
			if rawExternalSecurityGroup["display_id_and_name"] != nil {
				displayIDAndName = rawExternalSecurityGroup["display_id_and_name"].(bool)
			}

			displayName := ""
			if rawExternalSecurityGroup["display_name"] != nil {
				displayName = rawExternalSecurityGroup["display_name"].(string)
			}

			id := ""
			if rawExternalSecurityGroup["id"] != nil {
				id = rawExternalSecurityGroup["id"].(string)
			}

			item := core.NamedReferenceItem{
				DisplayIDAndName: displayIDAndName,
				DisplayName:      displayName,
				ID:               id,
			}
			expandedExternalSecurityGroups = append(expandedExternalSecurityGroups, item)
		}
	}
	return expandedExternalSecurityGroups
}

func flattenExternalSecurityGroups(externalSecurityGroups []core.NamedReferenceItem) []interface{} {
	flattenedExternalSecurityGroups := []interface{}{}
	for _, v := range externalSecurityGroups {
		flattenedExternalSecurityGroups = append(flattenedExternalSecurityGroups, map[string]interface{}{
			"display_id_and_name": v.DisplayIDAndName,
			"display_name":        v.DisplayName,
			"id":                  v.ID,
		})
	}

	return flattenedExternalSecurityGroups
}

func getExternalSecurityGroupsSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"display_id_and_name": {
			Computed: true,
			Optional: true,
			Type:     schema.TypeBool,
		},
		"display_name": {
			Computed: true,
			Optional: true,
			Type:     schema.TypeString,
		},
		"id": getIDSchema(),
	}
}
