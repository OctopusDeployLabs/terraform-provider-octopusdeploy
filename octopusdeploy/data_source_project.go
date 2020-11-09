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
		return nil
	}

	d.SetId(project.GetID())
	d.Set("name", project.Name)
	d.Set("description", project.Description)
	d.Set("lifecycle_id", project.LifecycleID)
	d.Set(constProjectGroupID, project.ProjectGroupID)
	d.Set(constDefaultFailureMode, project.DefaultGuidedFailureMode)
	d.Set(constSkipMachineBehavior, project.ProjectConnectivityPolicy.SkipMachineBehavior)

	return nil
}
