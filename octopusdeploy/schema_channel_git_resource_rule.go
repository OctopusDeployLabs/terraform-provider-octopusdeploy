package octopusdeploy

import (
	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/channels"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func expandChannelGitResourceRules(ChannelGitResourceRule map[string]interface{}) channels.ChannelGitResourceRule {
	if len(ChannelGitResourceRule) == 0 {
		return channels.ChannelGitResourceRule{}
	}

	return channels.ChannelGitResourceRule{
		Id:                   ChannelGitResourceRule["id"].(string),
		GitDependencyActions: expandDeploymentActionGitDependencies(ChannelGitResourceRule["git_dependency_actions"]),
		Rules:                ChannelGitResourceRule["rules"].([]string),
	}
}

func flattenChannelGitResourceRules(ChannelGitResourceRules []channels.ChannelGitResourceRule) []map[string]interface{} {
	if len(ChannelGitResourceRules) == 0 {
		return []map[string]interface{}{}
	}

	var flattenedRules = make([]map[string]interface{}, len(ChannelGitResourceRules))
	for key, ChannelGitResourceRule := range ChannelGitResourceRules {
		flattenedRules[key] = map[string]interface{}{
			"id":                     ChannelGitResourceRule.Id,
			"git_dependency_actions": flattenDeploymentActionGitDependencies(ChannelGitResourceRule.GitDependencyActions),
			"rules":                  ChannelGitResourceRule.Rules,
		}
	}

	return flattenedRules
}

func getChannelGitResourceRuleSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"id": getIDSchema(),
		"git_dependency_actions": {
			Elem:     &schema.Resource{Schema: getDeploymentActionGitDependencySchema()},
			Required: true,
			Type:     schema.TypeList,
		},
		"rules": {
			Elem:     &schema.Schema{Type: schema.TypeString},
			Optional: true,
			Type:     schema.TypeList,
		},
	}
}
