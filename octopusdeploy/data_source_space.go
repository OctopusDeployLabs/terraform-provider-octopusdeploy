package octopusdeploy

import (
	"context"

	"github.com/OctopusDeploy/go-octopusdeploy/octopusdeploy"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceSpace() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceSpaceReadByName,
		Schema: map[string]*schema.Schema{
			"name": {
				Required: true,
				Type:     schema.TypeString,
			},
			"description": {
				Computed: true,
				Type:     schema.TypeString,
			},
			"is_default": {
				Computed: true,
				Type:     schema.TypeBool,
			},
			"space_managers_team_members": {
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Computed: true,
				Type:     schema.TypeList,
			},
			"space_managers_teams": {
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Computed: true,
				Type:     schema.TypeList,
			},
			"task_queue_stopped": {
				Computed: true,
				Type:     schema.TypeBool,
			},
		},
	}
}

func dataSourceSpaceReadByName(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	name := d.Get("name").(string)

	client := m.(*octopusdeploy.Client)
	space, err := client.Spaces.GetByName(name)
	if err != nil {
		return diag.FromErr(err)
	}
	if space == nil {
		return diag.Errorf("unable to retrieve space (name: %s)", name)
	}

	flattenSpace(ctx, d, space)
	return nil
}
