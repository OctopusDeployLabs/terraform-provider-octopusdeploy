package octopusdeploy

import (
	"context"
	"time"

	"github.com/OctopusDeploy/go-octopusdeploy/octopusdeploy"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceProjects() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceProjectsRead,
		Schema:      getProjectDataSchema(),
	}
}

func dataSourceProjectsRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	query := octopusdeploy.ProjectsQuery{
		ClonedFromProjectID: d.Get("cloned_from_project_id").(string),
		IDs:                 expandArray(d.Get("ids").([]interface{})),
		IsClone:             d.Get("is_clone").(bool),
		Name:                d.Get("name").(string),
		PartialName:         d.Get("partial_name").(string),
		Skip:                d.Get("skip").(int),
		Take:                d.Get("take").(int),
	}

	client := m.(*octopusdeploy.Client)
	projects, err := client.Projects.Get(query)
	if err != nil {
		return diag.FromErr(err)
	}

	flattenedProjects := []interface{}{}
	for _, project := range projects.Items {
		flattenedProjects = append(flattenedProjects, flattenProject(project))
	}

	d.Set("projects", flattenedProjects)
	d.SetId("Projects " + time.Now().UTC().String())

	return nil
}
