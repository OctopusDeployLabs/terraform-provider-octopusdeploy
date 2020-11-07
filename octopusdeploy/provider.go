package octopusdeploy

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// Provider is the plugin entry point
func Provider() *schema.Provider {
	return &schema.Provider{
		DataSourcesMap: map[string]*schema.Resource{
			constOctopusDeployAccount:            dataAccount(),
			constOctopusDeployAWSAccount:         dataAwsAccount(),
			constOctopusDeployEnvironment:        dataSourceEnvironment(),
			constOctopusDeployFeed:               dataFeed(),
			constOctopusDeployLibraryVariableSet: dataSourceLibraryVariableSet(),
			constOctopusDeployLifecycle:          dataSourceLifecycle(),
			constOctopusDeployMachine:            dataMachine(),
			constOctopusDeployMachinePolicy:      dataMachinePolicy(),
			constOctopusDeployProject:            dataProject(),
			constOctopusDeploySpace:              dataSpace(),
			constOctopusDeployTokenAccount:       dataSourceTokenAccount(),
			constOctopusDeployUser:               dataUser(),
			constOctopusDeployVariable:           dataVariable(),
		},
		ResourcesMap: map[string]*schema.Resource{
			constOctopusDeployAccount:                        resourceAccount(),
			constOctopusDeployAWSAccount:                     resourceAWSAccount(),
			constOctopusDeployAzureServicePrincipal:          resourceAzureServicePrincipal(),
			constOctopusDeployCertificate:                    resourceCertificate(),
			constOctopusDeployChannel:                        resourceChannel(),
			constOctopusDeployDeploymentProcess:              resourceDeploymentProcess(),
			constOctopusDeployEnvironment:                    resourceEnvironment(),
			constOctopusDeployFeed:                           resourceFeed(),
			constOctopusDeployLibraryVariableSet:             resourceLibraryVariableSet(),
			constOctopusDeployLifecycle:                      resourceLifecycle(),
			constOctopusDeployMachine:                        resourceMachine(),
			constOctopusDeployNuGetFeed:                      resourceNugetFeed(),
			constOctopusDeployProject:                        resourceProject(),
			constOctopusDeployProjectDeploymentTargetTrigger: resourceProjectDeploymentTargetTrigger(),
			constOctopusDeployProjectGroup:                   resourceProjectGroup(),
			constOctopusDeploySpace:                          resourceSpace(),
			constOctopusDeploySSHKeyAccount:                  resourceSSHKey(),
			constOctopusDeployTagSet:                         resourceTagSet(),
			constOctopusDeployTokenAccount:                   resourceTokenAccount(),
			constOctopusDeployUser:                           resourceUser(),
			constOctopusDeployUsernamePasswordAccount:        resourceUsernamePassword(),
			constOctopusDeployVariable:                       resourceVariable(),
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

		ConfigureContextFunc: providerConfigure,
	}
}

func providerConfigure(ctx context.Context, d *schema.ResourceData) (interface{}, diag.Diagnostics) {
	config := Config{
		Address: d.Get(constAddress).(string),
		APIKey:  d.Get(constAPIKey).(string),
		Space:   d.Get(constSpaceID).(string),
	}

	return config.Client()
}
