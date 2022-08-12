package octopusdeploy

import (
	"context"
	"time"

	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/client"
	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/teams"
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
	query := teams.TeamsQuery{
		IDs:           expandArray(d.Get("ids").([]interface{})),
		IncludeSystem: d.Get("include_system").(bool),
		PartialName:   d.Get("partial_name").(string),
		Spaces:        expandArray(d.Get("spaces").([]interface{})),
		Skip:          d.Get("skip").(int),
		Take:          d.Get("take").(int),
	}

	client := meta.(*client.Client)
	existingTeams, err := client.Teams.Get(query)
	if err != nil {
		return diag.FromErr(err)
	}

	flattenedTeams := []interface{}{}
	for _, team := range existingTeams.Items {
		flattenedTeams = append(flattenedTeams, flattenTeam(team))
	}

	d.Set("teams", flattenedTeams)
	d.SetId("Teams " + time.Now().UTC().String())

	return nil
}
