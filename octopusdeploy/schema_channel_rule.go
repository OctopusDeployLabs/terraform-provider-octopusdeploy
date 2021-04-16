package octopusdeploy

import (
	"github.com/OctopusDeploy/go-octopusdeploy/octopusdeploy"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func expandChannelRule(channelRule map[string]interface{}) octopusdeploy.ChannelRule {
	return octopusdeploy.ChannelRule{
		ActionPackages: expandDeploymentActionPackages(channelRule["action_package"]),
		ID:             channelRule["id"].(string),
		Tag:            channelRule["tag"].(string),
		VersionRange:   channelRule["version_range"].(string),
	}
}

func flattenChannelRules(channelRules []octopusdeploy.ChannelRule) []map[string]interface{} {
	var flattenedRules = make([]map[string]interface{}, len(channelRules))
	for key, channelRule := range channelRules {
		flattenedRules[key] = map[string]interface{}{
			"action_package": flattenDeploymentActionPackages(channelRule.ActionPackages),
			"id":             channelRule.ID,
			"tag":            channelRule.Tag,
			"version_range":  channelRule.VersionRange,
		}
	}

	return flattenedRules
}

func getChannelRuleSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"action_package": {
			Elem:     &schema.Resource{Schema: getDeploymentActionPackageSchema()},
			Required: true,
			Type:     schema.TypeList,
		},
		"id": getIDSchema(),
		"tag": {
			Optional:      true,
			Type:          schema.TypeString,
		},
		"version_range": {
			Optional: true,
			Type:     schema.TypeString,
		},
	}
}
