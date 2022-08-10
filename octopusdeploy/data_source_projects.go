package octopusdeploy

import (
	"context"
	"time"

	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/client"
	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/projects"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceProjects() *schema.Resource {
	return &schema.Resource{
		Description: "Provides information about existing projects.",
		ReadContext: dataSourceProjectsRead,
		Schema:      getProjectDataSchema(),
	}
}

func dataSourceProjectsRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	query := projects.ProjectsQuery{
		ClonedFromProjectID: d.Get("cloned_from_project_id").(string),
		IDs:                 expandArray(d.Get("ids").([]interface{})),
		IsClone:             d.Get("is_clone").(bool),
		Name:                d.Get("name").(string),
		PartialName:         d.Get("partial_name").(string),
		Skip:                d.Get("skip").(int),
		Take:                d.Get("take").(int),
	}

	client := m.(*client.Client)
	existingProjects, err := client.Projects.Get(query)
	if err != nil {
		return diag.FromErr(err)
	}

	flattenedProjects := []interface{}{}
	for _, project := range existingProjects.Items {
		flattenedProjects = append(flattenedProjects, flattenProject(project))
	}

	d.Set("projects", flattenedProjects)
	d.SetId("Projects " + time.Now().UTC().String())

	return nil
}
