package octopusdeploy

import (
	"fmt"

	"github.com/OctopusDeploy/go-octopusdeploy/client"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func dataTagSet() *schema.Resource {
	return &schema.Resource{
		Read: dataTagSetReadByName,

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"tag": getTagSchema(),
		},
	}
}

func dataTagSetReadByName(d *schema.ResourceData, m interface{}) error {
	apiClient := m.(*client.Client)

	TagSetName := d.Get("name")
	env, err := apiClient.TagSets.GetByName(TagSetName.(string))

	if err == client.ErrItemNotFound {
		return nil
	}

	if err != nil {
		return fmt.Errorf("error reading TagSet with name %s: %s", TagSetName, err.Error())
	}

	d.SetId(env.ID)

	d.Set("name", env.Name)

	return nil
}
