package octopusdeploy

import (
	"fmt"

	"github.com/OctopusDeploy/go-octopusdeploy/client"
	"github.com/OctopusDeploy/go-octopusdeploy/model"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceChannel() *schema.Resource {
	return &schema.Resource{
		Create: resourceChannelCreate,
		Read:   resourceChannelRead,
		Update: resourceChannelUpdate,
		Delete: resourceChannelDelete,

		Schema: map[string]*schema.Schema{
			constName: {
				Type:     schema.TypeString,
				Required: true,
			},
			constDescription: {
				Type:     schema.TypeString,
				Optional: true,
			},
			constProjectID: {
				Type:     schema.TypeString,
				Required: true,
			},
			constLifecycleID: {
				Type:     schema.TypeString,
				Optional: true,
			},
			constIsDefault: {
				Type:     schema.TypeBool,
				Optional: true,
			},
			constRule: {
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

func resourceChannelCreate(d *schema.ResourceData, m interface{}) error {
	apiClient := m.(*client.Client)

	newChannel := buildChannelResource(d)
	channel, err := apiClient.Channels.Add(newChannel)

	if err != nil {
		return fmt.Errorf("error creating channel %s: %s", newChannel.Name, err.Error())
	}

	d.SetId(channel.ID)

	return nil
}

func buildChannelResource(d *schema.ResourceData) *model.Channel {
	channel := &model.Channel{
		Name:        d.Get(constName).(string),
		Description: d.Get(constDescription).(string),
		ProjectID:   d.Get(constProjectID).(string),
		LifecycleID: d.Get(constLifecycleID).(string),
		IsDefault:   d.Get(constIsDefault).(bool),
	}

	if attr, ok := d.GetOk(constRule); ok {
		tfRules := attr.([]interface{})

		for _, tfrule := range tfRules {
			rule := buildRulesResource(tfrule.(map[string]interface{}))
			channel.Rules = append(channel.Rules, rule)
		}
	}

	return channel
}

func buildRulesResource(tfRule map[string]interface{}) model.ChannelRule {
	rule := model.ChannelRule{
		VersionRange: tfRule[constVersionRange].(string),
		Tag:          tfRule[constTag].(string),
		Actions:      getSliceFromTerraformTypeList(tfRule[constActions]),
	}

	return rule
}

func flattenRules(in []model.ChannelRule) []map[string]interface{} {
	var flattened = make([]map[string]interface{}, len(in))
	for i, v := range in {
		m := make(map[string]interface{})
		m[constVersionRange] = v.VersionRange
		m[constTag] = v.Tag
		m[constActions] = v.Actions

		flattened[i] = m
	}

	return flattened
}

func resourceChannelRead(d *schema.ResourceData, m interface{}) error {
	apiClient := m.(*client.Client)
	channelID := d.Id()
	resource, err := apiClient.Channels.GetByID(channelID)

	if err != nil {
		return createResourceOperationError(errorReadingChannel, channelID, err)
	}
	if resource == nil {
		d.SetId(constEmptyString)
		return nil
	}

	logResource(constChannel, m)

	d.Set(constName, resource.Name)
	d.Set(constProjectID, resource.ProjectID)
	d.Set(constDescription, resource.Description)
	d.Set(constLifecycleID, resource.LifecycleID)
	d.Set(constIsDefault, resource.IsDefault)
	d.Set(constRule, flattenRules(resource.Rules))

	return nil
}

func resourceChannelUpdate(d *schema.ResourceData, m interface{}) error {
	channel := buildChannelResource(d)
	channel.ID = d.Id() // set channel struct ID so octopus knows which channel to update

	apiClient := m.(*client.Client)

	updatedChannel, err := apiClient.Channels.Update(*channel)

	if err != nil {
		return createResourceOperationError(errorUpdatingChannel, d.Id(), err)
	}

	d.SetId(updatedChannel.ID)

	return nil
}

func resourceChannelDelete(d *schema.ResourceData, m interface{}) error {
	apiClient := m.(*client.Client)

	channelID := d.Id()

	err := apiClient.Channels.DeleteByID(channelID)

	if err != nil {
		return createResourceOperationError(errorDeletingChannel, channelID, err)
	}

	d.SetId(constEmptyString)

	return nil
}
