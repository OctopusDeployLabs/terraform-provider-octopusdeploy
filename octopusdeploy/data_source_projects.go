package octopusdeploy

import (
	"context"

	"github.com/OctopusDeploy/go-octopusdeploy/octopusdeploy"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceProject() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceProjectReadByName,
		Schema:      getProjectDataSchema(),
	}
}

func dataSourceProjectReadByName(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	name := d.Get("name").(string)

	client := m.(*octopusdeploy.Client)
	project, err := client.Projects.GetByName(name)
	if err != nil {
		return diag.FromErr(err)
	}
	if project == nil {
		return diag.Errorf("unable to retrieve project (name: %s)", name)
	}

	flattenProject(ctx, d, project)
	return nil
}
