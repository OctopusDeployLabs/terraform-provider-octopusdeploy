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
		DeleteContext: resourceChannelDelete,
		Importer:      getImporter(),
		ReadContext:   resourceChannelRead,
		Schema:        getChannelSchema(),
		UpdateContext: resourceChannelUpdate,
	}
}

func resourceChannelCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	channel := expandChannel(d)

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

func resourceChannelRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*octopusdeploy.Client)
	channel, err := client.Channels.GetByID(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}
	if channel == nil {
		d.SetId("")
		return nil
	}

	flattenChannel(ctx, d, channel)
	return nil
}

func resourceChannelUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	channel := expandChannel(d)

	client := m.(*octopusdeploy.Client)
	updatedChannel, err := client.Channels.Update(*channel)
	if err != nil {
		return diag.FromErr(err)
	}

	flattenChannel(ctx, d, updatedChannel)
	return nil
}

func resourceChannelDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*octopusdeploy.Client)
	err := client.Channels.DeleteByID(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId("")
	return nil
}
