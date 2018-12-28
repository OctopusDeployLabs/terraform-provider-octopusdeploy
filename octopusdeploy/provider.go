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
		},
		ResourcesMap: map[string]*schema.Resource{
			"octopusdeploy_project":                           resourceProject(),
			"octopusdeploy_project_group":                     resourceProjectGroup(),
			"octopusdeploy_project_deployment_target_trigger": resourceProjectDeploymentTargetTrigger(),
			"octopusdeploy_environment":                       resourceEnvironment(),
			"octopusdeploy_variable":                          resourceVariable(),
			"octopusdeploy_machine":                           resourceMachine(),
			"octopusdeploy_library_variable_set":              resourceLibraryVariableSet(),
			"octopusdeploy_lifecycle":                         resourceLifecycle(),
			"octopusdeploy_deployment_process":                resourceDeploymentProcess(),
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
		},

		ConfigureFunc: providerConfigure,
	}
}

func providerConfigure(d *schema.ResourceData) (interface{}, error) {
	config := Config{
		Address: d.Get("address").(string),
		APIKey:  d.Get("apikey").(string),
	}

	log.Println("[INFO] Initializing Octopus Deploy client")
	client := config.Client()

	return client, nil
}
