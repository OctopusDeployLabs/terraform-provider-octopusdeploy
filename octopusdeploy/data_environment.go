package octopusdeploy

import (
	"github.com/OctopusDeploy/go-octopusdeploy/client"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataEnvironment() *schema.Resource {
	return &schema.Resource{
		Read: dataEnvironmentReadByName,

		Schema: map[string]*schema.Schema{
			constName: {
				Type:     schema.TypeString,
				Required: true,
			},
			constDescription: {
				Type:     schema.TypeString,
				Computed: true,
			},
			constUseGuidedFailure: {
				Type:     schema.TypeBool,
				Computed: true,
			},
			constAllowDynamicInfrastructure: {
				Type:     schema.TypeBool,
				Computed: true,
			},
		},
	}
}

func dataEnvironmentReadByName(d *schema.ResourceData, m interface{}) error {
	name := d.Get(constName).(string)

	apiClient := m.(*client.Client)
	resource, err := apiClient.Environments.GetByName(name)
	if err != nil {
		return createResourceOperationError(errorReadingEnvironment, name, err)
	}
	if resource == nil {
		return nil
	}

	logResource(constEnvironment, m)

	d.SetId(resource.ID)
	d.Set(constName, resource.Name)
	d.Set(constDescription, resource.Description)
	d.Set(constUseGuidedFailure, resource.UseGuidedFailure)
	d.Set(constAllowDynamicInfrastructure, resource.AllowDynamicInfrastructure)

	return nil
}
