package octopusdeploy

import (
	"context"
	"fmt"

	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/client"
	"github.com/OctopusDeploy/terraform-provider-octopusdeploy/internal/errors"
	"github.com/hashicorp/terraform-plugin-log/tflog"
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
	mutex.Lock()
	defer mutex.Unlock()

	channel := expandChannel(d)

	tflog.Info(ctx, fmt.Sprintf("creating channel: %#v", channel))

	client := m.(*client.Client)
	createdChannel, err := client.Channels.Add(channel)
	if err != nil {
		return diag.FromErr(err)
	}

	if err := setChannel(ctx, d, createdChannel); err != nil {
		return diag.FromErr(err)
	}

	d.SetId(createdChannel.GetID())

	tflog.Info(ctx, fmt.Sprintf("channel created (%s)", d.Id()))
	return nil
}

func resourceChannelDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	mutex.Lock()
	defer mutex.Unlock()

	tflog.Info(ctx, fmt.Sprintf("deleting channel (%s)", d.Id()))

	client := m.(*client.Client)
	if err := client.Channels.DeleteByID(d.Id()); err != nil {
		return diag.FromErr(err)
	}

	d.SetId("")

	tflog.Info(ctx, "channel deleted")
	return nil
}

func resourceChannelRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	tflog.Info(ctx, fmt.Sprintf("reading channel (%s)", d.Id()))

	client := m.(*client.Client)
	channel, err := client.Channels.GetByID(d.Id())
	if err != nil {
		return errors.ProcessApiError(ctx, d, err, "channel")
	}

	if err := setChannel(ctx, d, channel); err != nil {
		return diag.FromErr(err)
	}

	tflog.Info(ctx, fmt.Sprintf("channel read (%s)", d.Id()))
	return nil
}

func resourceChannelUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	mutex.Lock()
	defer mutex.Unlock()

	tflog.Info(ctx, fmt.Sprintf("updating channel (%s)", d.Id()))

	channel := expandChannel(d)
	client := m.(*client.Client)
	updatedChannel, err := client.Channels.Update(channel)
	if err != nil {
		return diag.FromErr(err)
	}

	if err := setChannel(ctx, d, updatedChannel); err != nil {
		return diag.FromErr(err)
	}

	tflog.Info(ctx, fmt.Sprintf("channel updated (%s)", d.Id()))
	return nil
}
