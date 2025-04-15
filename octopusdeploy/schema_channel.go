package octopusdeploy

import (
	"context"
	"fmt"

	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/channels"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func expandChannel(d *schema.ResourceData) *channels.Channel {
	name := d.Get("name").(string)
	projectID := d.Get("project_id").(string)

	channel := channels.NewChannel(name, projectID)
	channel.ID = d.Id()

	if v, ok := d.GetOk("description"); ok {
		channel.Description = v.(string)
	}

	if v, ok := d.GetOk("is_default"); ok {
		channel.IsDefault = v.(bool)
	}

	if v, ok := d.GetOk("lifecycle_id"); ok {
		channel.LifecycleID = v.(string)
	}

	if v, ok := d.GetOk("space_id"); ok {
		channel.SpaceID = v.(string)
	}

	if v, ok := d.GetOk("tenant_tags"); ok {
		channel.TenantTags = getSliceFromTerraformTypeList(v)
	}

	if v, ok := d.GetOk("rule"); ok {
		channelRules := v.([]interface{})
		for _, channelRule := range channelRules {
			rule := expandChannelRule(channelRule.(map[string]interface{}))
			channel.Rules = append(channel.Rules, rule)
		}
	}

	return channel
}

func flattenChannel(channel *channels.Channel) map[string]interface{} {
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
		"rule":         flattenChannelRules(channel.Rules),
		"space_id":     channel.SpaceID,
		"tenant_tags":  channel.TenantTags,
	}
}

func getChannelDataSchema() map[string]*schema.Schema {
	dataSchema := getChannelSchema()
	setDataSchema(&dataSchema)

	return map[string]*schema.Schema{
		"channels": {
			Computed:    true,
			Description: "A channel that matches the specified filter(s).",
			Elem:        &schema.Resource{Schema: dataSchema},
			Optional:    false,
			Type:        schema.TypeList,
		},
		"ids":          getQueryIDs(),
		"project_id":   getQueryProjectID(),
		"partial_name": getQueryPartialName(),
		"skip":         getQuerySkip(),
		"take":         getQueryTake(),
		"space_id":     getSpaceIDSchema(),
	}
}

func getChannelSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"description": getDescriptionSchema("channel"),
		"id":          getIDSchema(),
		"is_default": {
			Description: "Indicates if this is the default channel for the associated project.",
			Optional:    true,
			Type:        schema.TypeBool,
		},
		"lifecycle_id": {
			Description: "The lifecycle ID associated with this channel.",
			Optional:    true,
			Type:        schema.TypeString,
		},
		"name": getNameSchema(true),
		"project_id": {
			Description: "The project ID associated with this channel.",
			Required:    true,
			Type:        schema.TypeString,
		},
		"rule": {
			Description: "A list of rules associated with this channel.",
			Elem:        &schema.Resource{Schema: getChannelRuleSchema()},
			Optional:    true,
			Type:        schema.TypeList,
		},
		"space_id":    getSpaceIDSchema(),
		"tenant_tags": getTenantTagsSchema(),
	}
}

func setChannel(ctx context.Context, d *schema.ResourceData, channel *channels.Channel) error {
	d.Set("description", channel.Description)
	d.Set("is_default", channel.IsDefault)
	d.Set("lifecycle_id", channel.LifecycleID)
	d.Set("name", channel.Name)
	d.Set("project_id", channel.ProjectID)
	d.Set("space_id", channel.SpaceID)

	if err := d.Set("rule", flattenChannelRules(channel.Rules)); err != nil {
		return fmt.Errorf("error setting rule: %s", err)
	}

	if err := d.Set("tenant_tags", channel.TenantTags); err != nil {
		return fmt.Errorf("error setting tenant_tags: %s", err)
	}

	return nil
}
