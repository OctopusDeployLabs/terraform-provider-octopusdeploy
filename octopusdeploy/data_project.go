package octopusdeploy

import (
	"github.com/OctopusDeploy/go-octopusdeploy/client"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataProject() *schema.Resource {
	return &schema.Resource{
		Read: dataProjectReadByName,

		Schema: map[string]*schema.Schema{
			constName: {
				Type:     schema.TypeString,
				Required: true,
			},
			constDescription: {
				Type:     schema.TypeString,
				Optional: true,
			},
			constLifecycleID: {
				Type:     schema.TypeString,
				Optional: true,
			},
			constProjectGroupID: {
				Type:     schema.TypeString,
				Optional: true,
			},
			constDefaultFailureMode: {
				Type:     schema.TypeString,
				Optional: true,
			},
			constSkipMachineBehavior: {
				Type:     schema.TypeString,
				Optional: true,
			},
		},
	}
}

func dataProjectReadByName(d *schema.ResourceData, m interface{}) error {
	name := d.Get(constName).(string)

	apiClient := m.(*client.Client)
	resource, err := apiClient.Projects.GetByName(name)
	if err != nil {
		return createResourceOperationError(errorReadingProject, name, err)
	}
	if resource == nil {
		return nil
	}

	logResource(constProject, m)

	d.SetId(resource.ID)
	d.Set(constName, resource.Name)
	d.Set(constDescription, resource.Description)
	d.Set(constLifecycleID, resource.LifecycleID)
	d.Set(constProjectGroupID, resource.ProjectGroupID)
	d.Set(constDefaultFailureMode, resource.DefaultGuidedFailureMode)
	d.Set(constSkipMachineBehavior, resource.ProjectConnectivityPolicy.SkipMachineBehavior)

	return nil
}
