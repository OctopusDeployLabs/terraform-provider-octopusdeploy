package octopusdeploy

import "github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

func getPollingSubscriptionIDSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"id": {
			Computed:    true,
			Description: "The generated polling subscription ID.",
			Type:        schema.TypeString,
		},
		"polling_uri": {
			Computed:    true,
			Description: "The URI of the polling subscription ID.",
			Type:        schema.TypeString,
		},
		"dependencies": {
			Optional:    true,
			Type:        schema.TypeMap,
			Description: "Optional map of dependencies that when modified will trigger a re-creation of this resource.",
			ForceNew:    true,
		},
	}
}
