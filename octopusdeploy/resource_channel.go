package octopusdeploy

import (
	"context"
	"log"

	"github.com/OctopusDeploy/go-octopusdeploy/octopusdeploy"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceChannel() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceChannelCreate,
		ReadContext:   resourceChannelRead,
		UpdateContext: resourceChannelUpdate,
		DeleteContext: resourceChannelDelete,

		Schema: map[string]*schema.Schema{
			"name": {
				Required: true,
				Type:     schema.TypeString,
			},
			"description": {
				Optional: true,
				Type:     schema.TypeString,
			},
			"project_id": {
				Required: true,
				Type:     schema.TypeString,
			},
			"lifecycle_id": {
				Optional: true,
				Type:     schema.TypeString,
			},
			"is_default": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"rule": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						constVersionRange: {
							Type:     schema.TypeString,
							Optional: true,
						},
						constTag: {
							Type:     schema.TypeString,
							Optional: true,
						},
						constActions: {
							Type:     schema.TypeList,
							Optional: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
					},
				},
			},
		},
	}
}

func resourceChannelCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	channel := expandChannelResource(d)

	client := m.(*octopusdeploy.Client)
	resource, err := client.Channels.Add(channel)
	if err != nil {
		return diag.FromErr(err)
	}

	if isEmpty(resource.GetID()) {
		log.Println("ID is nil")
	} else {
		d.SetId(resource.GetID())
	}

	return nil
}

func expandChannelResource(d *schema.ResourceData) *octopusdeploy.Channel {
	channel := &octopusdeploy.Channel{
		Name:        d.Get(constName).(string),
		Description: d.Get("description").(string),
		ProjectID:   d.Get("project_id").(string),
		LifecycleID: d.Get("lifecycle_id").(string),
		IsDefault:   d.Get("is_default").(bool),
	}
	channel.ID = d.Id()

	if attr, ok := d.GetOk("rule"); ok {
		tfRules := attr.([]interface{})

		for _, tfrule := range tfRules {
			rule := buildRulesResource(tfrule.(map[string]interface{}))
			channel.Rules = append(channel.Rules, rule)
		}
	}

	return channel
}

func buildRulesResource(tfRule map[string]interface{}) octopusdeploy.ChannelRule {
	rule := octopusdeploy.ChannelRule{
		VersionRange: tfRule[constVersionRange].(string),
		Tag:          tfRule[constTag].(string),
		Actions:      getSliceFromTerraformTypeList(tfRule[constActions]),
	}

	return rule
}

func flattenRules(channelRules []octopusdeploy.ChannelRule) []map[string]interface{} {
	var flattenedRules = make([]map[string]interface{}, len(channelRules))
	for key, channelRule := range channelRules {
		m := make(map[string]interface{})
		m[constVersionRange] = channelRule.VersionRange
		m[constTag] = channelRule.Tag
		m[constActions] = channelRule.Actions

		flattenedRules[key] = m
	}

	return flattenedRules
}

func resourceChannelRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*octopusdeploy.Client)
	resource, err := client.Channels.GetByID(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}
	if resource == nil {
		d.SetId("")
		return nil
	}

	logResource(constChannel, m)

	d.Set("name", resource.Name)
	d.Set("project_id", resource.ProjectID)
	d.Set("description", resource.Description)
	d.Set("lifecycle_id", resource.LifecycleID)
	d.Set("is_default", resource.IsDefault)
	d.Set("rule", flattenRules(resource.Rules))

	return nil
}

func resourceChannelUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	channel := expandChannelResource(d)

	client := m.(*octopusdeploy.Client)
	resource, err := client.Channels.Update(*channel)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(resource.GetID())

	return nil
}

func resourceChannelDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	id := d.Id()

	client := m.(*octopusdeploy.Client)
	err := client.Channels.DeleteByID(id)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId("")
	return nil
}
