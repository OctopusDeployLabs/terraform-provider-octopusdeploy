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
	channel.SpaceID = d.Get("space_id").(string)
	channel.TenantTags = getSliceFromTerraformTypeList(d.Get("tenant_tags"))

	if v, ok := d.GetOk("rule"); ok {
		channelRules := v.([]interface{})
		for _, channelRule := range channelRules {
			rule := expandChannelRule(channelRule.(map[string]interface{}))
			channel.Rules = append(channel.Rules, rule)
		}
	}

	return channel
}

func expandChannelRule(channelRule map[string]interface{}) octopusdeploy.ChannelRule {
	return octopusdeploy.ChannelRule{
		Actions:      getSliceFromTerraformTypeList(channelRule["actions"]),
		ID:           channelRule["id"].(string),
		Tag:          channelRule["tag"].(string),
		VersionRange: channelRule["version_range"].(string),
	}
}

func flattenChannel(channel *octopusdeploy.Channel) map[string]interface{} {
	if channel == nil {
		return nil
	}

	return map[string]interface{}{
		"description":  channel.Description,
		"id":           channel.GetID(),
		"is_default":   channel.IsDefault,
		"lifecycle_id": channel.LifecycleID,
		"name":         channel.Name,
		"project_id":   channel.ProjectID,
		"rules":        flattenChannelRules(channel.Rules),
		"space_id":     channel.SpaceID,
		"tenant_tags":  channel.TenantTags,
	}
}

func setChannel(ctx context.Context, d *schema.ResourceData, channel *octopusdeploy.Channel) {
	d.Set("description", channel.Description)
	d.Set("is_default", channel.IsDefault)
	d.Set("lifecycle_id", channel.LifecycleID)
	d.Set("name", channel.Name)
	d.Set("project_id", channel.ProjectID)
	d.Set("rules", flattenChannelRules(channel.Rules))
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

func getChannelDataSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"ids": {
			Elem:     &schema.Schema{Type: schema.TypeString},
			Optional: true,
			Type:     schema.TypeList,
		},
		"partial_name": {
			Optional: true,
			Type:     schema.TypeString,
		},
		"skip": {
			Default:  0,
			Type:     schema.TypeInt,
			Optional: true,
		},
		"take": {
			Default:  1,
			Type:     schema.TypeInt,
			Optional: true,
		},
		"channels": {
			Computed: true,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"description": {
						Computed: true,
						Type:     schema.TypeString,
					},
					"id": {
						Computed: true,
						Type:     schema.TypeString,
					},
					"is_default": {
						Computed: true,
						Type:     schema.TypeBool,
					},
					"lifecycle_id": {
						Computed: true,
						Type:     schema.TypeString,
					},
					"name": {
						Computed: true,
						Type:     schema.TypeString,
					},
					"project_id": {
						Computed: true,
						Type:     schema.TypeString,
					},
					"rules": {
						Elem:     &schema.Resource{Schema: getChannelRuleSchema()},
						Computed: true,
						Type:     schema.TypeList,
					},
					"space_id": {
						Computed: true,
						Type:     schema.TypeString,
					},
					"tenant_tags": {
						Type:     schema.TypeList,
						Computed: true,
						Elem:     &schema.Schema{Type: schema.TypeString},
					},
				},
			},
			Type: schema.TypeList,
		},
	}
}

func getChannelSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"description": {
			Optional: true,
			Type:     schema.TypeString,
		},
		"id": {
			Computed: true,
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
		"rules": {
			Elem:     &schema.Resource{Schema: getChannelRuleSchema()},
			Optional: true,
			Type:     schema.TypeList,
		},
		"space_id": {
			Optional: true,
			Type:     schema.TypeString,
		},
		"tenant_tags": {
			Type:     schema.TypeList,
			Optional: true,
			Elem:     &schema.Schema{Type: schema.TypeString},
		},
	}
}

func getChannelRuleSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
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
	}
}
