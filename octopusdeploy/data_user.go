package octopusdeploy

import (
	"fmt"
	"log"

	"github.com/OctopusDeploy/go-octopusdeploy/octopusdeploy"
	"github.com/hashicorp/terraform/helper/schema"
)

func dataUser() *schema.Resource {
	return &schema.Resource{
		Read: dataUserReadByName,

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
		},
	}
}

func dataUserReadByName(d *schema.ResourceData, m interface{}) error {
	client := m.(*octopusdeploy.Client)

	userName := d.Get("UserName")

	user, err := client.User.GetByName(userName.(string))

	if err == octopusdeploy.ErrItemNotFound {
		return nil
	}

	if err != nil {
		return fmt.Errorf("error reading user name %s: %s", userName, err.Error())
	}

	d.SetId(user.ID)

	log.Printf("[DEBUG] user name: %v", m)
	d.Set("name", userName)
	return nil
}
