package octopusdeploy

import (
	"github.com/OctopusDeploy/go-octopusdeploy/client"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataTagSet() *schema.Resource {
	return &schema.Resource{
		Read: dataTagSetReadByName,

		Schema: map[string]*schema.Schema{
			constName: {
				Type:     schema.TypeString,
				Required: true,
			},
			constTag: getTagSchema(),
		},
	}
}

func dataTagSetReadByName(d *schema.ResourceData, m interface{}) error {
	name := d.Get(constName).(string)

	apiClient := m.(*client.Client)
	resource, err := apiClient.TagSets.GetByName(name)
	if err != nil {
		return createResourceOperationError(errorReadingTagSet, name, err)
	}
	if resource == nil {
		return nil
	}

	logResource(constTagSet, m)

	d.SetId(resource.ID)
	d.Set(constName, resource.Name)

	return nil
}
