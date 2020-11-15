package octopusdeploy

import "github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

func getExtensionSettingsSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"extension_id": {
			Computed: true,
			Type:     schema.TypeString,
		},
		"values": {
			Computed: true,
			Elem:     &schema.Schema{Type: schema.TypeString},
			Type:     schema.TypeList,
		},
	}
}
