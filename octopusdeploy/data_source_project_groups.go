package octopusdeploy

import (
	"context"

	"github.com/OctopusDeploy/go-octopusdeploy/octopusdeploy"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceProjectGroup() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceProjectGroupReadByName,
		Schema:      getProjectGroupDataSchema(),
	}
}

func dataSourceProjectGroupReadByName(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*octopusdeploy.Client)
	name := d.Get("name").(string)
	projectGroups, err := client.ProjectGroups.GetByPartialName(name)
	if err != nil {
		return diag.FromErr(err)
	}
	if len(projectGroups) == 0 {
		d.SetId("")
		return diag.Errorf("unable to retrieve project group (partial name: %s)", name)
	}

	projectGroup := projectGroups[0]

	flattenProjectGroup(ctx, d, projectGroup)
	return nil
}
