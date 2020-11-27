package octopusdeploy

import (
	"context"

	"github.com/OctopusDeploy/go-octopusdeploy/octopusdeploy"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func expandChannel(d *schema.ResourceData) *octopusdeploy.Channel {
	name := d.Get("name").(string)
	projectID := d.Get("project_id").(string)

	channel := octopusdeploy.NewChannel(name, projectID)
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

func getChannelDataSchema() map[string]*schema.Schema {
	channelSchema := getChannelSchema()
	for _, field := range channelSchema {
		field.Computed = true
		field.Default = nil
		field.MaxItems = 0
		field.MinItems = 0
		field.Optional = false
		field.Required = false
		field.ValidateDiagFunc = nil
	}

	return map[string]*schema.Schema{
		"channels": {
			Computed: true,
			Elem:     &schema.Resource{Schema: channelSchema},
			Type:     schema.TypeList,
		},
		"ids": {
			Description: "Query and/or search by a list of IDs",
			Elem:        &schema.Schema{Type: schema.TypeString},
			Optional:    true,
			Type:        schema.TypeList,
		},
		"partial_name": {
			Description: "Query and/or search by partial name",
			Optional:    true,
			Type:        schema.TypeString,
		},
		"skip": {
			Default:     0,
			Description: "Indicates the number of items to skip in the response",
			Type:        schema.TypeInt,
			Optional:    true,
		},
		"take": {
			Default:     1,
			Description: "Indicates the number of items to take (or return) in the response",
			Type:        schema.TypeInt,
			Optional:    true,
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
			Required:         true,
			Type:             schema.TypeString,
			ValidateDiagFunc: validateDiagFunc(validation.StringIsNotEmpty),
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

func setChannel(ctx context.Context, d *schema.ResourceData, channel *octopusdeploy.Channel) {
	d.Set("description", channel.Description)
	d.Set("id", channel.GetID())
	d.Set("is_default", channel.IsDefault)
	d.Set("lifecycle_id", channel.LifecycleID)
	d.Set("name", channel.Name)
	d.Set("project_id", channel.ProjectID)
	d.Set("rules", flattenChannelRules(channel.Rules))
	d.Set("space_id", channel.SpaceID)
	d.Set("tenant_tags", channel.TenantTags)
}
