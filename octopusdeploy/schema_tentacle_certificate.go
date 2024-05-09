package octopusdeploy

import "github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

func getTentacleCertificateSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"base64": {
			Computed:    true,
			Sensitive:   true,
			Description: "The base64 encoded pfx certificate.",
			Type:        schema.TypeString,
		},
		"thumbprint": {
			Computed:    true,
			Description: "The SHA1 sum of the certificate represented in hexadecimal.",
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
