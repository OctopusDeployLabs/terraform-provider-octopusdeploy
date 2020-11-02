package octopusdeploy

import (
	"log"

	"github.com/OctopusDeploy/go-octopusdeploy/octopusdeploy"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
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
	channel := buildChannelResource(d)

	client := m.(*octopusdeploy.Client)
	resource, err := client.Channels.Add(channel)
	if err != nil {
		return createResourceOperationError(errorCreatingChannel, channel.Name, err)
	}

	if isEmpty(resource.GetID()) {
		log.Println("ID is nil")
	} else {
		d.SetId(resource.GetID())
	}

	return nil
}

func buildChannelResource(d *schema.ResourceData) *octopusdeploy.Channel {
	channel := &octopusdeploy.Channel{
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

func buildRulesResource(tfRule map[string]interface{}) octopusdeploy.ChannelRule {
	rule := octopusdeploy.ChannelRule{
		VersionRange: tfRule[constVersionRange].(string),
		Tag:          tfRule[constTag].(string),
		Actions:      getSliceFromTerraformTypeList(tfRule[constActions]),
	}

	return rule
}

func flattenRules(in []octopusdeploy.ChannelRule) []map[string]interface{} {
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
	id := d.Id()

	client := m.(*octopusdeploy.Client)
	resource, err := client.Channels.GetByID(id)
	if err != nil {
		return createResourceOperationError(errorReadingChannel, id, err)
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
	channel.ID = d.Id() // set ID so Octopus API knows which channel to update

	client := m.(*octopusdeploy.Client)
	resource, err := client.Channels.Update(*channel)
	if err != nil {
		return createResourceOperationError(errorUpdatingChannel, d.Id(), err)
	}

	d.SetId(resource.GetID())

	return nil
}

func resourceChannelDelete(d *schema.ResourceData, m interface{}) error {
	id := d.Id()

	client := m.(*octopusdeploy.Client)
	err := client.Channels.DeleteByID(id)
	if err != nil {
		return createResourceOperationError(errorDeletingChannel, id, err)
	}

	d.SetId(constEmptyString)
	return nil
}
