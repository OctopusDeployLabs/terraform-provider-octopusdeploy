package octopusdeploy

import (
	"context"
	"log"

	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/client"
	"github.com/OctopusDeploy/terraform-provider-octopusdeploy/internal/errors"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceChannel() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceChannelCreate,
		DeleteContext: resourceChannelDelete,
		Description:   "This resource manages channels in Octopus Deploy.",
		Importer:      getImporter(),
		ReadContext:   resourceChannelRead,
		Schema:        getChannelSchema(),
		UpdateContext: resourceChannelUpdate,
	}
}

func resourceChannelCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	channel := expandChannel(d)

	log.Printf("[INFO] creating channel: %#v", channel)

	client := m.(*client.Client)
	createdChannel, err := client.Channels.Add(channel)
	if err != nil {
		return diag.FromErr(err)
	}

	if err := setChannel(ctx, d, createdChannel); err != nil {
		return diag.FromErr(err)
	}

	d.SetId(createdChannel.GetID())

	log.Printf("[INFO] channel created (%s)", d.Id())
	return nil
}

func resourceChannelDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[INFO] deleting channel (%s)", d.Id())

	client := m.(*client.Client)
	if err := client.Channels.DeleteByID(d.Id()); err != nil {
		return diag.FromErr(err)
	}

	d.SetId("")

	log.Printf("[INFO] channel deleted")
	return nil
}

func resourceChannelRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[INFO] reading channel (%s)", d.Id())

	client := m.(*client.Client)
	channel, err := client.Channels.GetByID(d.Id())
	if err != nil {
		return errors.ProcessApiError(ctx, d, err, "channel")
	}

	if err := setChannel(ctx, d, channel); err != nil {
		return diag.FromErr(err)
	}

	log.Printf("[INFO] channel read (%s)", d.Id())
	return nil
}

func resourceChannelUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[INFO] updating channel (%s)", d.Id())

	channel := expandChannel(d)
	client := m.(*client.Client)
	updatedChannel, err := client.Channels.Update(channel)
	if err != nil {
		return diag.FromErr(err)
	}

	if err := setChannel(ctx, d, updatedChannel); err != nil {
		return diag.FromErr(err)
	}

	log.Printf("[INFO] channel updated (%s)", d.Id())
	return nil
}
