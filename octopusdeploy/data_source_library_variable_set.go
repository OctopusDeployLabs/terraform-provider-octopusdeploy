package octopusdeploy

import (
	"context"

	"github.com/OctopusDeploy/go-octopusdeploy/octopusdeploy"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceLibraryVariableSet() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceLibraryVariableSetReadByName,
		Schema: map[string]*schema.Schema{
			"name": {
				Required: true,
				Type:     schema.TypeString,
			},
		},
	}
}

func dataSourceLibraryVariableSetReadByName(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	name := d.Get("name").(string)

	client := m.(*octopusdeploy.Client)
	resourceList, err := client.LibraryVariableSets.GetByPartialName(name)
	if err != nil {
		return diag.FromErr(err)
	}
	if len(resourceList) == 0 {
		return nil
	}

	logResource(constLibraryVariableSet, m)

	// NOTE: two or more library variables can have the same name in Octopus.
	// Therefore, a better search criteria needs to be implemented below.

	for _, resource := range resourceList {
		if resource.Name == name {
			logResource(constLibraryVariableSet, m)

			d.SetId(resource.GetID())
			d.Set("name", resource.Name)
			d.Set("description", resource.Description)
			d.Set(constVariableSetID, resource.VariableSetID)

			return nil
		}
	}

	return nil
}
