package octopusdeploy

import (
	"context"
	"log"
	"time"

	"github.com/OctopusDeploy/go-octopusdeploy/octopusdeploy"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceSpaces() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceSpacesRead,
		Schema: map[string]*schema.Schema{
			"ids": {
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Optional: true,
				Type:     schema.TypeList,
			},
			"name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"partial_name": {
				Type:     schema.TypeString,
				Optional: true,
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
			"spaces": {
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"description": {
							Optional: true,
							Type:     schema.TypeString,
						},
						"is_default": {
							Optional: true,
							Type:     schema.TypeBool,
						},
						"id": {
							Optional: true,
							Type:     schema.TypeString,
						},
						"name": {
							Optional: true,
							Type:     schema.TypeString,
						},
						"space_managers_team_members": {
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
							Optional: true,
							Type:     schema.TypeList,
						},
						"space_managers_teams": {
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
							Optional: true,
							Type:     schema.TypeList,
						},
						"task_queue_stopped": {
							Optional: true,
							Type:     schema.TypeBool,
						},
					},
				},
				Type: schema.TypeList,
			},
		},
	}
}

func dataSourceSpacesRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	ids := expandArray(d.Get("ids").([]interface{}))
	name := d.Get("name").(string)
	partialName := d.Get("partial_name").(string)
	skip := d.Get("skip").(int)
	take := d.Get("take").(int)

	query := octopusdeploy.SpacesQuery{
		IDs:         ids,
		Name:        name,
		PartialName: partialName,
		Skip:        skip,
		Take:        take,
	}

	client := m.(*octopusdeploy.Client)
	spaces, err := client.Spaces.Get(query)
	if err != nil {
		return diag.FromErr(err)
	}

	flattenedSpaces := []interface{}{}
	for _, space := range spaces.Items {
		flattenedSpace := map[string]interface{}{
			"description":                 space.Description,
			"id":                          space.GetID(),
			"is_default":                  space.IsDefault,
			"name":                        space.Name,
			"space_managers_team_members": space.SpaceManagersTeamMembers,
			"space_managers_teams":        space.SpaceManagersTeams,
			"task_queue_stopped":          space.TaskQueueStopped,
		}
		flattenedSpaces = append(flattenedSpaces, flattenedSpace)
	}

	d.Set("spaces", flattenedSpaces)
	d.SetId("Spaces " + time.Now().UTC().String())

	log.Println(flattenedSpaces)

	return nil
}
