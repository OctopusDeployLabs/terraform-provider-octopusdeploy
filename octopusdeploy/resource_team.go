package octopusdeploy

import (
	"context"
	"log"

	"github.com/OctopusDeploy/go-octopusdeploy/octopusdeploy"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceTeam() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceTeamCreate,
		DeleteContext: resourceTeamDelete,
		Description:   "This resource manages teams in Octopus Deploy.",
		Importer:      getImporter(),
		ReadContext:   resourceTeamRead,
		Schema:        getTeamSchema(),
		UpdateContext: resourceTeamUpdate,
	}
}

func resourceTeamCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	team := expandTeam(d)

	log.Printf("[INFO] creating team: %#v", team)

	client := m.(*octopusdeploy.Client)
	createdTeam, err := client.Teams.Add(team)
	if err != nil {
		return diag.FromErr(err)
	}

	if err := setTeam(ctx, d, createdTeam); err != nil {
		return diag.FromErr(err)
	}

	d.SetId(createdTeam.GetID())

	log.Printf("[INFO] team created (%s)", d.Id())
	return nil
}

func resourceTeamDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[INFO] deleting team (%s)", d.Id())

	client := m.(*octopusdeploy.Client)
	if err := client.Teams.DeleteByID(d.Id()); err != nil {
		return diag.FromErr(err)
	}

	d.SetId("")

	log.Printf("[INFO] team deleted")
	return nil
}

func resourceTeamRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[INFO] reading team (%s)", d.Id())

	client := m.(*octopusdeploy.Client)
	team, err := client.Teams.GetByID(d.Id())
	if err != nil {
		if apiError, ok := err.(*octopusdeploy.APIError); ok {
			if apiError.StatusCode == 404 {
				log.Printf("[INFO] team (%s) not found; deleting from state", d.Id())
				d.SetId("")
				return nil
			}
		}
		return diag.FromErr(err)
	}

	if err := setTeam(ctx, d, team); err != nil {
		return diag.FromErr(err)
	}

	log.Printf("[INFO] team read (%s)", d.Id())
	return nil
}

func resourceTeamUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[INFO] updating team (%s)", d.Id())

	team := expandTeam(d)
	client := m.(*octopusdeploy.Client)
	updatedTeam, err := client.Teams.Update(team)
	if err != nil {
		return diag.FromErr(err)
	}

	if err := setTeam(ctx, d, updatedTeam); err != nil {
		return diag.FromErr(err)
	}

	log.Printf("[INFO] team updated (%s)", d.Id())
	return nil
}
