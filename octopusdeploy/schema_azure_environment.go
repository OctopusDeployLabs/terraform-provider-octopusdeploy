package octopusdeploy

import "github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

func getAzureEnvironmentDataSchema() *schema.Schema {
	return &schema.Schema{
		Computed: true,
		Type:     schema.TypeString,
	}
}

func getAzureEnvironmentSchema() *schema.Schema {
	return &schema.Schema{
		Default:  "AzureCloud",
		Optional: true,
		Type:     schema.TypeString,
		ValidateDiagFunc: validateValueFunc([]string{
			"AzureCloud",
			"AzureChinaCloud",
			"AzureGermanCloud",
			"AzureUSGovernment",
		}),
	}
}
