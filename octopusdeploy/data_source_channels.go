package octopusdeploy

import (
	"context"
	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/resources"
	"time"

	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/channels"
	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/client"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceChannels() *schema.Resource {
	return &schema.Resource{
		Description: "Provides information about existing channels.",
		ReadContext: dataSourceChannelsRead,
		Schema:      getChannelDataSchema(),
	}
}

func dataSourceChannelsRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	query := channels.Query{
		IDs:         expandArray(d.Get("ids").([]interface{})),
		PartialName: d.Get("partial_name").(string),
		Skip:        d.Get("skip").(int),
		Take:        d.Get("take").(int),
	}
	projectID := d.Get("project_id").(string)
	spaceID := d.Get("space_id").(string)

	client := m.(*client.Client)

	var existingChannels *resources.Resources[*channels.Channel] = nil
	var err error
	if projectID != "" {
		queryByProjectID := channels.QueryByProjectID{
			ProjectID:   projectID,
			PartialName: query.PartialName,
			Skip:        query.Skip,
			Take:        query.Take,
		}
		existingChannels, err = channels.GetByProjectID(client, spaceID, queryByProjectID)
		if len(query.IDs) > 0 {
			filteredChannels := make([]*channels.Channel, 0)
			for _, channel := range existingChannels.Items {
				for _, id := range query.IDs {
					if channel.ID == id {

						filteredChannels = append(filteredChannels, channel)
						break
					}
				}
			}
			existingChannels.Items = filteredChannels
		}
	} else {
		existingChannels, err = channels.Get(client, spaceID, query)
	}

	if err != nil {
		return diag.FromErr(err)
	}

	flattenedChannels := []interface{}{}
	for _, channel := range existingChannels.Items {
		flattenedChannels = append(flattenedChannels, flattenChannel(channel))
	}

	d.Set("channels", flattenedChannels)
	d.SetId("Channels " + time.Now().UTC().String())

	return nil
}
