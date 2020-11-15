package octopusdeploy

import (
	"github.com/OctopusDeploy/go-octopusdeploy/octopusdeploy"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func expandAzureWebApp(d *schema.ResourceData) *octopusdeploy.AzureWebAppEndpoint {
	endpoint := octopusdeploy.NewAzureWebAppEndpoint()
	endpoint.ID = d.Id()

	if v, ok := d.GetOk("resource_group_name"); ok {
		endpoint.ResourceGroupName = v.(string)
	}

	if v, ok := d.GetOk("web_app_name"); ok {
		endpoint.WebAppName = v.(string)
	}

	if v, ok := d.GetOk("web_app_slot_name"); ok {
		endpoint.WebAppSlotName = v.(int)
	}

	return endpoint
}

func flattenAzureWebApp(endpoint *octopusdeploy.AzureWebAppEndpoint) []interface{} {
	if endpoint == nil {
		return nil
	}

	return []interface{}{map[string]interface{}{
		"id":                  endpoint.GetID(),
		"resource_group_name": endpoint.ResourceGroupName,
		"web_app_name":        endpoint.WebAppName,
		"web_app_slot_name":   endpoint.WebAppSlotName,
	}}
}

func getAzureWebAppSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"id": {
			Computed: true,
			Type:     schema.TypeString,
		},
		"resource_group_name": {
			Optional: true,
			Type:     schema.TypeString,
		},
		"web_app_name": {
			Optional: true,
			Type:     schema.TypeString,
		},
		"web_app_slot_name": {
			Optional: true,
			Type:     schema.TypeInt,
		},
	}
}
