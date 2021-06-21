package octopusdeploy

import (
	"bytes"
	"context"
	"fmt"
	"log"
	"sort"

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

	if err := resourceTeamUpdateUserRoles(ctx, d, m, createdTeam); err != nil {
		return diag.FromErr(err)
	}

	if err := setTeam(ctx, d, createdTeam); err != nil {
		return diag.FromErr(err)
	}

	d.SetId(createdTeam.GetID())

	log.Printf("[INFO] team created (%s)", d.Id())
	return resourceTeamRead(ctx, d, m)
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
		apiError := err.(*octopusdeploy.APIError)
		if apiError.StatusCode == 404 {
			log.Printf("[INFO] team (%s) not found; deleting from state", d.Id())
			d.SetId("")
			return nil
		}
		return diag.FromErr(err)
	}

	userRoles, err := client.Teams.GetScopedUserRolesByID(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}
	remoteUserRoles := flattenScopedUserRoles(userRoles.Items)
	d.Set("user_role", remoteUserRoles)

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

	if err := resourceTeamUpdateUserRoles(ctx, d, m, team); err != nil {
		return diag.FromErr(err)
	}

	if err := setTeam(ctx, d, updatedTeam); err != nil {
		return diag.FromErr(err)
	}

	log.Printf("[INFO] team updated (%s)", d.Id())
	return resourceTeamRead(ctx, d, m)
}

func expandUserRoles(team *octopusdeploy.Team, userRoles []interface{}) []*octopusdeploy.ScopedUserRole {
	values := make([]*octopusdeploy.ScopedUserRole, 0, len(userRoles))
	for _, rawUserRole := range userRoles {
		userRole := rawUserRole.(map[string]interface{})
		scopedUserRole := octopusdeploy.NewScopedUserRole(userRole["user_role_id"].(string))
		scopedUserRole.TeamID = team.ID
		scopedUserRole.SpaceID = userRole["space_id"].(string)

		if v, ok := userRole["id"]; ok {
			scopedUserRole.ID = v.(string)
		} else {
			scopedUserRole.ID = ""
		}

		if v, ok := userRole["environment_ids"]; ok {
			scopedUserRole.EnvironmentIDs = getSliceFromTerraformTypeList(v)
		}

		if v, ok := userRole["project_group_ids"]; ok {
			scopedUserRole.ProjectGroupIDs = getSliceFromTerraformTypeList(v)
		}

		if v, ok := userRole["project_ids"]; ok {
			scopedUserRole.ProjectIDs = getSliceFromTerraformTypeList(v)
		}

		if v, ok := userRole["tenant_ids"]; ok {
			scopedUserRole.TenantIDs = getSliceFromTerraformTypeList(v)
		}
		values = append(values, scopedUserRole)
	}
	return values
}
func resourceTeamUpdateUserRoles(ctx context.Context, d *schema.ResourceData, m interface{}, team *octopusdeploy.Team) error {
	log.Printf("[INFO] updating team user roles (%s)", d.Id())
	if d.HasChange("user_role") {
		log.Printf("[INFO] user role has changes (%s)", d.Id())
		o, n := d.GetChange("user_role")
		if o == nil {
			o = new(schema.Set)
		}
		if n == nil {
			n = new(schema.Set)
		}

		os := o.(*schema.Set)
		ns := n.(*schema.Set)
		remove := expandUserRoles(team, os.Difference(ns).List())
		add := expandUserRoles(team, ns.Difference(os).List())

		if len(remove) > 0 || len(add) > 0 {
			log.Printf("[INFO] user role found diff (%s)", d.Id())
			client := m.(*octopusdeploy.Client)
			if len(remove) > 0 {
				log.Printf("[INFO] removing user roles from team (%s)", d.Id())
				for _, userRole := range remove {
					if userRole.ID != "" {
						err := client.ScopedUserRoles.DeleteByID(userRole.ID)
						if err != nil {
							apiError := err.(*octopusdeploy.APIError)
							if apiError.StatusCode != 404 {
								// It's already been deleted, maybe mixing with the independent resource?
								return fmt.Errorf("Error removing user role %s from team %s: %s", userRole.ID, team.ID, err)
							}
						}
					}
				}
			}
			if len(add) > 0 {
				log.Printf("[INFO] adding new user roles to team (%s)", d.Id())
				for _, userRole := range add {
					_, err := client.ScopedUserRoles.Add(userRole)
					if err != nil {
						return fmt.Errorf("Error creating user role for team %s: %s", team.ID, err)
					}
				}
			}
		}
	}
	return nil
}

func resourceTeamUserRoleListSetHash(buf *bytes.Buffer, v interface{}) {
	vs := v.(*schema.Set).List()
	s := make([]string, len(vs))
	for i, raw := range vs {
		s[i] = raw.(string)
	}
	sort.Strings(s)
	for _, v := range s {
		buf.WriteString(fmt.Sprintf("%s-", v))
	}
}

func resourceTeamUserRoleSetHash(v interface{}) int {
	var buf bytes.Buffer

	m := v.(map[string]interface{})
	buf.WriteString(fmt.Sprintf("%s-", m["user_role_id"].(string)))
	buf.WriteString(fmt.Sprintf("%s-", m["space_id"].(string)))

	if v, ok := m["environment_ids"]; ok {
		resourceTeamUserRoleListSetHash(&buf, v)
	}
	if v, ok := m["project_group_ids"]; ok {
		resourceTeamUserRoleListSetHash(&buf, v)
	}
	if v, ok := m["project_ids"]; ok {
		resourceTeamUserRoleListSetHash(&buf, v)
	}
	if v, ok := m["tenant_ids"]; ok {
		resourceTeamUserRoleListSetHash(&buf, v)
	}

	return stringHashCode(buf.String())
}
