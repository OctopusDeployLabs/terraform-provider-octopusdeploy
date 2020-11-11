package octopusdeploy

import (
	"context"

	"github.com/OctopusDeploy/go-octopusdeploy/octopusdeploy"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func expandChannel(d *schema.ResourceData) *octopusdeploy.Channel {
	description := d.Get("description").(string)
	name := d.Get("name").(string)
	projectID := d.Get("project_id").(string)

	channel := octopusdeploy.NewChannel(name, description, projectID)
	channel.ID = d.Id()

	channel.IsDefault = d.Get("is_default").(bool)
	channel.LifecycleID = d.Get("lifecycle_id").(string)

	if v, ok := d.GetOk("rule"); ok {
		channelRules := v.([]interface{})
		for _, channelRule := range channelRules {
			rule := expandRules(channelRule.(map[string]interface{}))
			channel.Rules = append(channel.Rules, rule)
		}
	}

	return channel
}

func expandRules(channelRule map[string]interface{}) octopusdeploy.ChannelRule {
	return octopusdeploy.ChannelRule{
		Actions:      getSliceFromTerraformTypeList(channelRule["actions"]),
		ID:           channelRule["id"].(string),
		Tag:          channelRule["tag"].(string),
		VersionRange: channelRule["version_range"].(string),
	}
}

func flattenChannel(ctx context.Context, d *schema.ResourceData, channel *octopusdeploy.Channel) {
	d.Set("description", channel.Description)
	d.Set("is_default", channel.IsDefault)
	d.Set("lifecycle_id", channel.LifecycleID)
	d.Set("name", channel.Name)
	d.Set("project_id", channel.ProjectID)
	d.Set("rule", flattenRules(channel.Rules))

	d.SetId(channel.GetID())
}

func flattenRules(channelRules []octopusdeploy.ChannelRule) []map[string]interface{} {
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

func getChannelSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"description": {
			Optional: true,
			Type:     schema.TypeString,
		},
		"is_default": {
			Optional: true,
			Type:     schema.TypeBool,
		},
		"lifecycle_id": {
			Optional: true,
			Type:     schema.TypeString,
		},
		"name": {
			Required:     true,
			Type:         schema.TypeString,
			ValidateFunc: validation.StringIsNotEmpty,
		},
		"project_id": {
			Required: true,
			Type:     schema.TypeString,
		},
		"rule": {
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"actions": {
						Elem:     &schema.Schema{Type: schema.TypeString},
						Optional: true,
						Type:     schema.TypeList,
					},
					"id": {
						Optional: true,
						Type:     schema.TypeString,
					},
					"tag": {
						Optional: true,
						Type:     schema.TypeString,
					},
					"version_range": {
						Optional: true,
						Type:     schema.TypeString,
					},
				},
			},
			Optional: true,
			Type:     schema.TypeList,
		},
		"tenant_tags": {
			Type:     schema.TypeList,
			Optional: true,
			Elem:     &schema.Schema{Type: schema.TypeString},
		},
	}
}
