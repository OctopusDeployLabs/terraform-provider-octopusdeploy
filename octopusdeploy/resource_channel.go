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
		},
	}
}

func resourceChannelCreate(d *schema.ResourceData, m interface{}) error {
	client := m.(*octopusdeploy.Client)

	newChannel := buildChannelResource(d)
	channel, err := client.Channel.Add(newChannel)

	if err != nil {
		return fmt.Errorf("error creating account %s: %s", newChannel.Name, err.Error())
	}

	d.SetId(channel.ID)

	return nil
}

func buildChannelResource(d *schema.ResourceData) *octopusdeploy.Channel {

	ch := &octopusdeploy.Channel{
		Name:        d.Get("name").(string),
		Description: d.Get("description").(string),
		ProjectID:   d.Get("project_id").(string),
		Rules: []octopusdeploy.ChannelRule{
			{
				Actions: []string{
					"deploy a package",
				},
				VersionRange: "1.0",
			},
		},
	}

	return ch
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
	/*d.Set("account_type", account.AccountType)
	d.Set("client_id", account.ClientId)
	d.Set("tenant_id", account.TenantId)
	d.Set("subscription_id", account.SubscriptionNumber)
	d.Set("client_secret", account.Password)
	d.Set("tenant_tags", account.TenantTags)
	d.Set("tenanted_deployment_participation", account.TenantedDeploymentParticipation)*/

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
