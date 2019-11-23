package octopusdeploy

import (
	"log"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/terraform"
)

//Provider is the plugin entry point
func Provider() terraform.ResourceProvider {
	return &schema.Provider{
		DataSourcesMap: map[string]*schema.Resource{
			"octopusdeploy_project":              dataProject(),
			"octopusdeploy_environment":          dataEnvironment(),
			"octopusdeploy_variable":             dataVariable(),
			"octopusdeploy_machinepolicy":        dataMachinePolicy(),
			"octopusdeploy_machine":              dataMachine(),
			"octopusdeploy_library_variable_set": dataLibraryVariableSet(),
			"octopusdeploy_lifecycle":            dataLifecycle(),
			"octopusdeploy_feed":                 dataFeed(),
			"octopusdeploy_account":              dataAccount(),
		},
		ResourcesMap: map[string]*schema.Resource{
			"octopusdeploy_project":                           resourceProject(),
			"octopusdeploy_project_group":                     resourceProjectGroup(),
			"octopusdeploy_project_deployment_target_trigger": resourceProjectDeploymentTargetTrigger(),
			"octopusdeploy_environment":                       resourceEnvironment(),
			"octopusdeploy_account":                           resourceAccount(),
			"octopusdeploy_feed":                              resourceFeed(),
			"octopusdeploy_variable":                          resourceVariable(),
			"octopusdeploy_machine":                           resourceMachine(),
			"octopusdeploy_library_variable_set":              resourceLibraryVariableSet(),
			"octopusdeploy_lifecycle":                         resourceLifecycle(),
			"octopusdeploy_deployment_process":                resourceDeploymentProcess(),
			"octopusdeploy_tag_set":                           resourceTagSet(),
			"octopusdeploy_certificate":                       resourceCertificate(),
			"octopusdeploy_channel":                           resourceChannel(),
			"octopusdeploy_nuget_feed":                        resourceNugetFeed(),
			"octopusdeploy_azure_service_principal":           resourceAzureServicePrincipal(),
		},
		Schema: map[string]*schema.Schema{
			"address": {
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("OCTOPUS_URL", nil),
				Description: "The URL of the Octopus Deploy server",
			},
			"apikey": {
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("OCTOPUS_APIKEY", nil),
				Description: "The API to use with the Octopus Deploy server.",
			},
			"space": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("OCTOPUS_SPACE", ""),
				Description: "The name of the Space in Octopus Deploy server",
			},
		},

		ConfigureFunc: providerConfigure,
	}
}

func providerConfigure(d *schema.ResourceData) (interface{}, error) {
	config := Config{
		Address: d.Get("address").(string),
		APIKey:  d.Get("apikey").(string),
		Space:   d.Get("space").(string),
	}

	log.Println("[INFO] Initializing Octopus Deploy client")
	client, err := config.Client()

	return client, err
}
