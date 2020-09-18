package octopusdeploy

import (
	"fmt"
	"log"

	"github.com/OctopusDeploy/go-octopusdeploy/client"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func dataLifecycle() *schema.Resource {
	return &schema.Resource{
		Read: dataLifecycleReadByName,

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},
		},
	}
}

func dataLifecycleReadByName(d *schema.ResourceData, m interface{}) error {
	apiClient := m.(*client.Client)

	lifecycleName := d.Get("name")

	lifecycle, err := apiClient.Lifecycles.GetByName(lifecycleName.(string))

	if err == client.ErrItemNotFound {
		return nil
	}

	if err != nil {
		return fmt.Errorf("error reading lifecycle name %s: %s", lifecycleName, err.Error())
	}

	d.SetId(lifecycle.ID)

	log.Printf("[DEBUG] lifecycle: %v", m)
	d.Set("name", lifecycle.Name)
	d.Set("description", lifecycle.Description)
	return nil
}
