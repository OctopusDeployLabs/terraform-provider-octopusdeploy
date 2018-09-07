package octopusdeploy

import (
	"fmt"

	"github.com/MattHodge/go-octopusdeploy/octopusdeploy"
	"github.com/hashicorp/terraform/helper/schema"
)

func dataEnvironment() *schema.Resource {
	return &schema.Resource{
		Read: dataEnvironmentReadByName,

		Schema: map[string]*schema.Schema{
			"name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
		},
	}
}

func dataEnvironmentReadByName(d *schema.ResourceData, m interface{}) error {
	client := m.(*octopusdeploy.Client)

	environmentName := d.Get("name")
	env, err := client.Environment.GetByName(environmentName.(string))

	if err == octopusdeploy.ErrItemNotFound {
		return nil
	}

	if err != nil {
		return fmt.Errorf("error reading environment with name %s: %s", environmentName, err.Error())
	}

	d.SetId(env.ID)
	d.Set("name", env.Name)
	d.Set("description", env.Description)
	d.Set("useguidedfailure", env.UseGuidedFailure)

	return nil
}
