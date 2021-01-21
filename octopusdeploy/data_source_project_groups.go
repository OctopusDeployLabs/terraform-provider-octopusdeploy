package octopusdeploy

import (
	"context"
	"time"

	"github.com/OctopusDeploy/go-octopusdeploy/octopusdeploy"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceProjectGroups() *schema.Resource {
	return &schema.Resource{
		Description: "Provides information about existing project groups.",
		ReadContext: dataSourceProjectGroupsRead,
		Schema:      getProjectGroupDataSchema(),
	}
}

func dataSourceProjectGroupsRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	query := octopusdeploy.ProjectGroupsQuery{
		IDs:         expandArray(d.Get("ids").([]interface{})),
		PartialName: d.Get("partial_name").(string),
		Skip:        d.Get("skip").(int),
		Take:        d.Get("take").(int),
	}

	client := m.(*octopusdeploy.Client)
	projectGroups, err := client.ProjectGroups.Get(query)
	if err != nil {
		return diag.FromErr(err)
	}

	flattenedProjectGroups := []interface{}{}
	for _, projectGroup := range projectGroups.Items {
		flattenedProjectGroups = append(flattenedProjectGroups, flattenProjectGroup(projectGroup))
	}

	d.Set("project_groups", flattenedProjectGroups)
	d.SetId("ProjectGroups " + time.Now().UTC().String())

	return nil
}
