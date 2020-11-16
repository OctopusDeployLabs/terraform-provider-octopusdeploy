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
			"octopusdeploy_accounts":              dataSourceAccounts(),
			"octopusdeploy_certificates":          dataSourceCertificates(),
			"octopusdeploy_channels":              dataSourceChannels(),
			"octopusdeploy_deployment_targets":    dataSourceDeploymentTargets(),
			"octopusdeploy_environments":          dataSourceEnvironments(),
			"octopusdeploy_feeds":                 dataSourceFeeds(),
			"octopusdeploy_library_variable_sets": dataSourceLibraryVariableSet(),
			"octopusdeploy_lifecycles":            dataSourceLifecycles(),
			"octopusdeploy_machine_policies":      dataSourceMachinePolicies(),
			"octopusdeploy_project_groups":        dataSourceProjectGroups(),
			"octopusdeploy_projects":              dataSourceProjects(),
			"octopusdeploy_spaces":                dataSourceSpaces(),
			"octopusdeploy_tag_sets":              dataSourceTagSets(),
			"octopusdeploy_users":                 dataSourceUsers(),
			"octopusdeploy_variables":             dataSourceVariable(),
		},
		ResourcesMap: map[string]*schema.Resource{
			"octopusdeploy_account":                           resourceAccount(),
			"octopusdeploy_aws_account":                       resourceAmazonWebServicesAccount(),
			"octopusdeploy_azure_service_principal":           resourceAzureServicePrincipalAccount(),
			"octopusdeploy_azure_subscription_account":        resourceAzureSubscriptionAccount(),
			"octopusdeploy_certificate":                       resourceCertificate(),
			"octopusdeploy_channel":                           resourceChannel(),
			"octopusdeploy_deployment_target":                 resourceDeploymentTarget(),
			"octopusdeploy_deployment_process":                resourceDeploymentProcess(),
			"octopusdeploy_environment":                       resourceEnvironment(),
			"octopusdeploy_feed":                              resourceFeed(),
			"octopusdeploy_library_variable_set":              resourceLibraryVariableSet(),
			"octopusdeploy_lifecycle":                         resourceLifecycle(),
			"octopusdeploy_machine_policy":                    resourceMachinePolicy(),
			"octopusdeploy_nuget_feed":                        resourceNuGetFeed(),
			"octopusdeploy_project":                           resourceProject(),
			"octopusdeploy_project_deployment_target_trigger": resourceProjectDeploymentTargetTrigger(),
			"octopusdeploy_project_group":                     resourceProjectGroup(),
			"octopusdeploy_space":                             resourceSpace(),
			"octopusdeploy_ssh_key_account":                   resourceSSHKeyAccount(),
			"octopusdeploy_tag_set":                           resourceTagSet(),
			"octopusdeploy_token_account":                     resourceTokenAccount(),
			"octopusdeploy_user":                              resourceUser(),
			"octopusdeploy_username_password_account":         resourceUsernamePasswordAccount(),
			"octopusdeploy_variable":                          resourceVariable(),
		},
		Schema: map[string]*schema.Schema{
			"address": {
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("OCTOPUS_URL", nil),
				Description: "The URL of the Octopus Deploy server",
			},
			"api_key": {
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("OCTOPUS_APIKEY", nil),
				Description: "The API to use with the Octopus Deploy server.",
			},
			"space_id": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("OCTOPUS_SPACE", ""),
				Description: "The name of the Space in Octopus Deploy server",
			},
		},

		ConfigureContextFunc: providerConfigure,
	}
}

func providerConfigure(ctx context.Context, d *schema.ResourceData) (interface{}, diag.Diagnostics) {
	config := Config{
		Address: d.Get("address").(string),
		APIKey:  d.Get("api_key").(string),
		Space:   d.Get("space_id").(string),
	}

	return config.Client()
}
