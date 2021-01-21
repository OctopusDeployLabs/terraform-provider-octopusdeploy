package octopusdeploy

import (
	"github.com/OctopusDeploy/go-octopusdeploy/octopusdeploy"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func expandChannelRule(channelRule map[string]interface{}) octopusdeploy.ChannelRule {
	return octopusdeploy.ChannelRule{
		Actions:      getSliceFromTerraformTypeList(channelRule["actions"]),
		ID:           channelRule["id"].(string),
		Tag:          channelRule["tag"].(string),
		VersionRange: channelRule["version_range"].(string),
	}
}

func flattenChannelRules(channelRules []octopusdeploy.ChannelRule) []map[string]interface{} {
	var flattenedRules = make([]map[string]interface{}, len(channelRules))
	for key, channelRule := range channelRules {
		flattenedRules[key] = map[string]interface{}{
			"actions":       channelRule.Actions,
			"id":            channelRule.ID,
			"tag":           channelRule.Tag,
			"version_range": channelRule.VersionRange,
		}
	}

	return flattenedRules
}

func getChannelRuleSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"actions": {
			Elem:     &schema.Schema{Type: schema.TypeString},
			Optional: true,
			Type:     schema.TypeList,
		},
		"id": getIDSchema(),
		"tag": {
			Optional: true,
			Type:     schema.TypeString,
		},
		"version_range": {
			Optional: true,
			Type:     schema.TypeString,
		},
	}
}
