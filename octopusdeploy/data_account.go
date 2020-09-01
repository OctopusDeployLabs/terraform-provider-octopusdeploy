package octopusdeploy

import (
	"fmt"
	"log"

	"github.com/OctopusDeploy/go-octopusdeploy/client"
	"github.com/hashicorp/terraform/helper/schema"
)

func dataAccount() *schema.Resource {
	return &schema.Resource{
		Read: dataAccountReadByName,

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
		},
	}
}

func dataAccountReadByName(d *schema.ResourceData, m interface{}) error {
	apiapiClient := m.(*client.Client)

	accountName := d.Get("name")

	account, err := apiapiClient.Accounts.GetByName(accountName.(string))

	if err == client.ErrItemNotFound {
		return nil
	}

	if err != nil {
		return fmt.Errorf("error reading account name %s: %s", accountName, err.Error())
	}

	d.SetId(account.ID)

	log.Printf("[DEBUG] account: %v", m)
	d.Set("name", account.Name)
	return nil
}
