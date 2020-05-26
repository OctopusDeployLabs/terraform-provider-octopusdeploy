package octopusdeploy

import (
	"fmt"

	"github.com/OctopusDeploy/go-octopusdeploy/octopusdeploy"
	"github.com/hashicorp/terraform/helper/schema"
)

func resourceChannel() *schema.Resource {
	return &schema.Resource{
		Create: resourceChannelCreate,
		Read:   resourceChannelRead,
		Update: resourceChannelUpdate,
		Delete: resourceChannelDelete,

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"project_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"lifecycle_id": {
				Type:     schema.TypeString,
				Optional: true,
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
						"version_range": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"tag": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"actions": {
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

func resourceChannelCreate(d *schema.ResourceData, m interface{}) error {
	client := m.(*octopusdeploy.Client)

	newChannel := buildChannelResource(d)
	channel, err := client.Channel.Add(newChannel)

	if err != nil {
		return fmt.Errorf("error creating channel %s: %s", newChannel.Name, err.Error())
	}

	d.SetId(channel.ID)

	return nil
}

func buildChannelResource(d *schema.ResourceData) *octopusdeploy.Channel {
	channel := &octopusdeploy.Channel{
		Name:        d.Get("name").(string),
		Description: d.Get("description").(string),
		ProjectID:   d.Get("project_id").(string),
		LifecycleID: d.Get("lifecycle_id").(string),
		IsDefault:   d.Get("is_default").(bool),
	}

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
		VersionRange: tfRule["version_range"].(string),
		Tag:          tfRule["tag"].(string),
		Actions:      getSliceFromTerraformTypeList(tfRule["actions"]),
	}

	return rule
}

func flattenRules(in []octopusdeploy.ChannelRule) []map[string]interface{} {
	var flattened = make([]map[string]interface{}, len(in))
	for i, v := range in {
		m := make(map[string]interface{})
		m["version_range"] = v.VersionRange
		m["tag"] = v.Tag
		m["actions"] = v.Actions

		flattened[i] = m
	}

	return flattened
}

func resourceChannelRead(d *schema.ResourceData, m interface{}) error {
	client := m.(*octopusdeploy.Client)

	channelID := d.Id()
	channel, err := client.Channel.Get(channelID)

	if err == octopusdeploy.ErrItemNotFound {
		d.SetId("")
		return nil
	}

	if err != nil {
		return fmt.Errorf("error reading channel %s: %s", channelID, err.Error())
	}

	d.Set("name", channel.Name)
	d.Set("project_id", channel.ProjectID)
	d.Set("description", channel.Description)
	d.Set("lifecycle_id", channel.LifecycleID)
	d.Set("is_default", channel.IsDefault)
	d.Set("rule", flattenRules(channel.Rules))

	return nil
}

func resourceChannelUpdate(d *schema.ResourceData, m interface{}) error {
	channel := buildChannelResource(d)
	channel.ID = d.Id() // set channel struct ID so octopus knows which channel to update

	client := m.(*octopusdeploy.Client)

	updatedChannel, err := client.Channel.Update(channel)

	if err != nil {
		return fmt.Errorf("error updating channel id %s: %s", d.Id(), err.Error())
	}

	d.SetId(updatedChannel.ID)

	return nil
}

func resourceChannelDelete(d *schema.ResourceData, m interface{}) error {
	client := m.(*octopusdeploy.Client)

	channelID := d.Id()

	err := client.Channel.Delete(channelID)

	if err != nil {
		return fmt.Errorf("error deleting channel id %s: %s", channelID, err.Error())
	}

	d.SetId("")

	return nil
}
