package octopusdeploy

import (
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// Provider is the plugin entry point
func Provider() *schema.Provider {
	log.Println("[INFO] Initializing Resource Provider")
	return &schema.Provider{
		DataSourcesMap: map[string]*schema.Resource{
			constOctopusDeployProject:            dataProject(),
			constOctopusDeployEnvironment:        dataEnvironment(),
			constOctopusDeployVariable:           dataVariable(),
			constOctopusDeployMachinePolicy:      dataMachinePolicy(),
			constOctopusDeployMachine:            dataMachine(),
			constOctopusDeployLibraryVariableSet: dataLibraryVariableSet(),
			constOctopusDeployLifecycle:          dataLifecycle(),
			constOctopusDeployFeed:               dataFeed(),
			constOctopusDeployAccount:            dataAccount(),
		},
		ResourcesMap: map[string]*schema.Resource{
			constOctopusDeployProject:                        resourceProject(),
			constOctopusDeployProjectGroup:                   resourceProjectGroup(),
			constOctopusDeployProjectDeploymentTargetTrigger: resourceProjectDeploymentTargetTrigger(),
			constOctopusDeployEnvironment:                    resourceEnvironment(),
			// constOctopusDeployAccount:                           resourceAccount(),
			constOctopusDeployFeed:                    resourceFeed(),
			constOctopusDeployVariable:                resourceVariable(),
			constOctopusDeployMachine:                 resourceMachine(),
			constOctopusDeployLibraryVariableSet:      resourceLibraryVariableSet(),
			constOctopusDeployLifecycle:               resourceLifecycle(),
			constOctopusDeployDeploymentProcess:       resourceDeploymentProcess(),
			constOctopusDeployTagSet:                  resourceTagSet(),
			constOctopusDeployCertificate:             resourceCertificate(),
			constOctopusDeployChannel:                 resourceChannel(),
			constOctopusDeployNuGetFeed:               resourceNugetFeed(),
			constOctopusDeployAzureServicePrincipal:   resourceAzureServicePrincipal(),
			constOctopusDeployUsernamePasswordAccount: resourceUsernamePassword(),
			constOctopusDeploySSHKeyAccount:           resourceSSHKey(),
			constOctopusDeployAWSAccount:              resourceAmazonWebServicesAccount(),
		},
		Schema: map[string]*schema.Schema{
			constAddress: {
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("OCTOPUS_URL", nil),
				Description: "The URL of the Octopus Deploy server",
			},
			constAPIKey: {
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("OCTOPUS_APIKEY", nil),
				Description: "The API to use with the Octopus Deploy server.",
			},
			constSpaceID: {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("OCTOPUS_SPACE", constEmptyString),
				Description: "The name of the Space in Octopus Deploy server",
			},
		},

		ConfigureFunc: providerConfigure,
	}
}

func providerConfigure(d *schema.ResourceData) (interface{}, error) {
	log.Println("[INFO] Parsing Client Configuration")
	config := Config{
		Address: d.Get(constAddress).(string),
		APIKey:  d.Get(constAPIKey).(string),
		Space:   d.Get(constSpaceID).(string),
	}

	log.Println("[INFO] Initializing Octopus Deploy client")
	client, err := config.Client()

	return client, err
}
