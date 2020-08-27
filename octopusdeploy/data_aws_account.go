package octopusdeploy

import (
	"fmt"

	"github.com/OctopusDeploy/go-octopusdeploy/octopusdeploy"
	"github.com/hashicorp/terraform/helper/schema"
)

func dataAwsAccount() *schema.Resource {
	return &schema.Resource{
		Read: dataAwsAccountReadByName,

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"account_type": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "AWS",
			},
			"environments": {
				Type: schema.TypeList,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Optional: true,
			},
			"tenant_tags": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"tenanted_deployment_participation": getTenantedDeploymentSchema(),
			"secret_key": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"access_key": {
				Type:     schema.TypeString,
				Optional: true,
			},
		},
	}
}

func dataAwsAccountReadByName(d *schema.ResourceData, m interface{}) error {
	client := m.(*octopusdeploy.Client)

	AwsAccountName := d.Get("name")
	env, err := client.Account.GetByName(AwsAccountName.(string))

	if err == octopusdeploy.ErrItemNotFound {
		return nil
	}

	if err != nil {
		return fmt.Errorf("error reading Aws Account with name %s: %s", AwsAccountName, err.Error())
	}

	d.SetId(env.ID)

	d.Set("name", env.Name)
	d.Set("description", env.Description)

	return nil
}
