package octopusdeploy

import (
	"context"
	"time"

	"github.com/OctopusDeploy/go-octopusdeploy/octopusdeploy"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceChannels() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceChannelsRead,
		Schema: map[string]*schema.Schema{
			"ids": {
				Elem:     &schema.Schema{Type: schema.TypeString},
				Optional: true,
				Type:     schema.TypeList,
			},
			"partial_name": {
				Optional: true,
				Type:     schema.TypeString,
			},
			"skip": {
				Default:  0,
				Type:     schema.TypeInt,
				Optional: true,
			},
			"take": {
				Default:  1,
				Type:     schema.TypeInt,
				Optional: true,
			},
			"channels": {
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"description": {
							Optional: true,
							Type:     schema.TypeString,
						},
						"id": {
							Optional: true,
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
							Optional: true,
							Type:     schema.TypeString,
						},
						"space_id": {
							Optional: true,
							Type:     schema.TypeString,
						},
						"project_id": {
							Optional: true,
							Type:     schema.TypeString,
						},
						"rules": {
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"actions": {
										Elem:     &schema.Schema{Type: schema.TypeString},
										Optional: true,
										Type:     schema.TypeList,
									},
									"id": {
										Optional: true,
										Type:     schema.TypeString,
									},
									"tag": {
										Optional: true,
										Type:     schema.TypeString,
									},
									"version_range": {
										Optional: true,
										Type:     schema.TypeString,
									},
								},
							},
							Optional: true,
							Type:     schema.TypeList,
						},
						"tenant_tags": {
							Type:     schema.TypeList,
							Optional: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
					},
				},
				Type: schema.TypeList,
			},
		},
	}
}

func dataSourceChannelsRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	query := octopusdeploy.ChannelsQuery{
		IDs:         expandArray(d.Get("ids").([]interface{})),
		PartialName: d.Get("partial_name").(string),
		Skip:        d.Get("skip").(int),
		Take:        d.Get("take").(int),
	}

	client := m.(*octopusdeploy.Client)
	channels, err := client.Channels.Get(query)
	if err != nil {
		return diag.FromErr(err)
	}

	flattenedChannels := []interface{}{}
	for _, channel := range channels.Items {
		flattenedChannel := map[string]interface{}{
			"description":  channel.Description,
			"id":           channel.GetID(),
			"is_default":   channel.IsDefault,
			"lifecycle_id": channel.LifecycleID,
			"name":         channel.Name,
			"project_id":   channel.ProjectID,
			"rules":        flattenRules(channel.Rules),
			"space_id":     channel.SpaceID,
			"tenant_tags":  channel.TenantTags,
		}
		flattenedChannels = append(flattenedChannels, flattenedChannel)
	}

	d.Set("channels", flattenedChannels)
	d.SetId("Channels " + time.Now().UTC().String())

	return nil
}
