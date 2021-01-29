package octopusdeploy

import (
	"context"
	"time"

	"github.com/OctopusDeploy/go-octopusdeploy/octopusdeploy"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceTeams() *schema.Resource {
	return &schema.Resource{
		Description: "Provides information about existing users.",
		ReadContext: dataSourceTeamsRead,
		Schema:      getTeamDataSchema(),
	}
}

func dataSourceTeamsRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	query := octopusdeploy.TeamsQuery{
		IDs:           expandArray(d.Get("ids").([]interface{})),
		IncludeSystem: d.Get("include_system").(bool),
		PartialName:   d.Get("partial_name").(string),
		Skip:          d.Get("skip").(int),
		Take:          d.Get("take").(int),
	}

	client := meta.(*octopusdeploy.Client)
	users, err := client.Teams.Get(query)
	if err != nil {
		return diag.FromErr(err)
	}

	flattenedTeams := []interface{}{}
	for _, user := range users.Items {
		flattenedTeams = append(flattenedTeams, flattenTeam(user))
	}

	d.Set("teams", flattenedTeams)
	d.SetId("Teams " + time.Now().UTC().String())

	return nil
}
