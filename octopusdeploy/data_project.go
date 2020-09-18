package octopusdeploy

import (
	"fmt"
	"log"

	"github.com/OctopusDeploy/go-octopusdeploy/client"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func dataProject() *schema.Resource {
	return &schema.Resource{
		Read: dataProjectReadByName,

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"lifecycle_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"project_group_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"default_failure_mode": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"skip_machine_behavior": {
				Type:     schema.TypeString,
				Optional: true,
			},
		},
	}
}

func dataProjectReadByName(d *schema.ResourceData, m interface{}) error {
	apiClient := m.(*client.Client)

	projectName := d.Get("name")

	project, err := apiClient.Projects.GetByName(projectName.(string))

	if err == client.ErrItemNotFound {
		return nil
	}

	if err != nil {
		return fmt.Errorf("error reading project name %s: %s", projectName, err.Error())
	}

	d.SetId(project.ID)

	log.Printf("[DEBUG] project: %v", m)
	d.Set("name", project.Name)
	d.Set("description", project.Description)
	d.Set("lifecycle_id", project.LifecycleID)
	d.Set("project_group_id", project.ProjectGroupID)
	d.Set("default_failure_mode", project.DefaultGuidedFailureMode)
	d.Set("skip_machine_behavior", project.ProjectConnectivityPolicy.SkipMachineBehavior)

	return nil
}
