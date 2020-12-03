package octopusdeploy

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// Provider is the plugin entry point for the Terraform Provider for Octopus
// Deploy.
func Provider() *schema.Provider {
	return &schema.Provider{
		DataSourcesMap: map[string]*schema.Resource{
			"octopusdeploy_accounts":                                        dataSourceAccounts(),
			"octopusdeploy_azure_cloud_service_deployment_targets":          dataSourceAzureCloudServiceDeploymentTargets(),
			"octopusdeploy_azure_service_fabric_cluster_deployment_targets": dataSourceAzureServiceFabricClusterDeploymentTargets(),
			"octopusdeploy_azure_web_app_deployment_targets":                dataSourceAzureWebAppDeploymentTargets(),
			"octopusdeploy_certificates":                                    dataSourceCertificates(),
			"octopusdeploy_cloud_region_deployment_targets":                 dataSourceCloudRegionDeploymentTargets(),
			"octopusdeploy_channels":                                        dataSourceChannels(),
			"octopusdeploy_environments":                                    dataSourceEnvironments(),
			"octopusdeploy_feeds":                                           dataSourceFeeds(),
			"octopusdeploy_kubernetes_cluster_deployment_targets":           dataSourceKubernetesClusterDeploymentTargets(),
			"octopusdeploy_library_variable_sets":                           dataSourceLibraryVariableSet(),
			"octopusdeploy_lifecycles":                                      dataSourceLifecycles(),
			"octopusdeploy_listening_tentacle_deployment_targets":           dataSourceListeningTentacleDeploymentTargets(),
			"octopusdeploy_machine_policies":                                dataSourceMachinePolicies(),
			"octopusdeploy_offline_package_drop_deployment_targets":         dataSourceOfflinePackageDropDeploymentTargets(),
			"octopusdeploy_polling_tentacle_deployment_targets":             dataSourcePollingTentacleDeploymentTargets(),
			"octopusdeploy_project_groups":                                  dataSourceProjectGroups(),
			"octopusdeploy_projects":                                        dataSourceProjects(),
			"octopusdeploy_spaces":                                          dataSourceSpaces(),
			"octopusdeploy_ssh_connection_deployment_targets":               dataSourceSSHConnectionDeploymentTargets(),
			"octopusdeploy_tag_sets":                                        dataSourceTagSets(),
			"octopusdeploy_users":                                           dataSourceUsers(),
			"octopusdeploy_user_roles":                                      dataSourceUserRoles(),
			"octopusdeploy_variables":                                       dataSourceVariable(),
		},
		ResourcesMap: map[string]*schema.Resource{
			"octopusdeploy_account":                                        resourceAccount(),
			"octopusdeploy_aws_account":                                    resourceAmazonWebServicesAccount(),
			"octopusdeploy_azure_cloud_service_deployment_target":          resourceAzureCloudServiceDeploymentTarget(),
			"octopusdeploy_azure_service_fabric_cluster_deployment_target": resourceAzureServiceFabricClusterDeploymentTarget(),
			"octopusdeploy_azure_service_principal":                        resourceAzureServicePrincipalAccount(),
			"octopusdeploy_azure_subscription_account":                     resourceAzureSubscriptionAccount(),
			"octopusdeploy_azure_web_app_deployment_target":                resourceAzureWebAppDeploymentTarget(),
			"octopusdeploy_certificate":                                    resourceCertificate(),
			"octopusdeploy_channel":                                        resourceChannel(),
			"octopusdeploy_cloud_region_deployment_target":                 resourceCloudRegionDeploymentTarget(),
			"octopusdeploy_deployment_process":                             resourceDeploymentProcess(),
			"octopusdeploy_environment":                                    resourceEnvironment(),
			"octopusdeploy_feed":                                           resourceFeed(),
			"octopusdeploy_kubernetes_cluster_deployment_target":           resourceKubernetesClusterDeploymentTarget(),
			"octopusdeploy_library_variable_set":                           resourceLibraryVariableSet(),
			"octopusdeploy_lifecycle":                                      resourceLifecycle(),
			"octopusdeploy_listening_tentacle_deployment_target":           resourceListeningTentacleDeploymentTarget(),
			"octopusdeploy_machine_policy":                                 resourceMachinePolicy(),
			"octopusdeploy_nuget_feed":                                     resourceNuGetFeed(),
			"octopusdeploy_offline_package_drop_deployment_target":         resourceOfflinePackageDropDeploymentTarget(),
			"octopusdeploy_polling_tentacle_deployment_target":             resourcePollingTentacleDeploymentTarget(),
			"octopusdeploy_project":                                        resourceProject(),
			"octopusdeploy_project_deployment_target_trigger":              resourceProjectDeploymentTargetTrigger(),
			"octopusdeploy_project_group":                                  resourceProjectGroup(),
			"octopusdeploy_space":                                          resourceSpace(),
			"octopusdeploy_ssh_connection_deployment_target":               resourceSSHConnectionDeploymentTarget(),
			"octopusdeploy_ssh_key_account":                                resourceSSHKeyAccount(),
			"octopusdeploy_tag_set":                                        resourceTagSet(),
			"octopusdeploy_token_account":                                  resourceTokenAccount(),
			"octopusdeploy_user":                                           resourceUser(),
			"octopusdeploy_user_role":                                      resourceUserRole(),
			"octopusdeploy_username_password_account":                      resourceUsernamePasswordAccount(),
			"octopusdeploy_variable":                                       resourceVariable(),
		},
		Schema: map[string]*schema.Schema{
			"address": {
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("OCTOPUS_URL", nil),
				Description: "The endpoint of the Octopus REST API",
			},
			"api_key": {
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("OCTOPUS_APIKEY", nil),
				Description: "The API key to use with the Octopus REST API",
			},
			"space_id": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("OCTOPUS_SPACE", ""),
				Description: "The space ID to target",
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
