package octopusdeploy

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// Provider is the plugin entry point for the Terraform provider for Octopus Deploy.
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
			"octopusdeploy_deployment_targets":                              dataSourceDeploymentTargets(),
			"octopusdeploy_kubernetes_agent_deployment_targets":             dataSourceKubernetesAgentDeploymentTargets(),
			"octopusdeploy_kubernetes_cluster_deployment_targets":           dataSourceKubernetesClusterDeploymentTargets(),
			"octopusdeploy_listening_tentacle_deployment_targets":           dataSourceListeningTentacleDeploymentTargets(),
			"octopusdeploy_machine":                                         dataSourceMachine(),
			"octopusdeploy_machine_policies":                                dataSourceMachinePolicies(),
			"octopusdeploy_offline_package_drop_deployment_targets":         dataSourceOfflinePackageDropDeploymentTargets(),
			"octopusdeploy_polling_tentacle_deployment_targets":             dataSourcePollingTentacleDeploymentTargets(),
			"octopusdeploy_script_modules":                                  dataSourceScriptModules(),
			"octopusdeploy_ssh_connection_deployment_targets":               dataSourceSSHConnectionDeploymentTargets(),
			"octopusdeploy_tag_sets":                                        dataSourceTagSets(),
			"octopusdeploy_teams":                                           dataSourceTeams(),
			"octopusdeploy_tenants":                                         dataSourceTenants(),
			"octopusdeploy_users":                                           dataSourceUsers(),
			"octopusdeploy_user_roles":                                      dataSourceUserRoles(),
			"octopusdeploy_worker_pools":                                    dataSourceWorkerPools(),
		},
		ResourcesMap: map[string]*schema.Resource{
			"octopusdeploy_aws_account":                                    resourceAmazonWebServicesAccount(),
			"octopusdeploy_aws_openid_connect_account":                     resourceAmazonWebServicesOpenIDConnectAccount(),
			"octopusdeploy_azure_cloud_service_deployment_target":          resourceAzureCloudServiceDeploymentTarget(),
			"octopusdeploy_azure_service_fabric_cluster_deployment_target": resourceAzureServiceFabricClusterDeploymentTarget(),
			"octopusdeploy_azure_service_principal":                        resourceAzureServicePrincipalAccount(),
			"octopusdeploy_azure_openid_connect":                           resourceAzureOpenIDConnectAccount(),
			"octopusdeploy_azure_subscription_account":                     resourceAzureSubscriptionAccount(),
			"octopusdeploy_azure_web_app_deployment_target":                resourceAzureWebAppDeploymentTarget(),
			"octopusdeploy_certificate":                                    resourceCertificate(),
			"octopusdeploy_channel":                                        resourceChannel(),
			"octopusdeploy_cloud_region_deployment_target":                 resourceCloudRegionDeploymentTarget(),
			"octopusdeploy_docker_container_registry":                      resourceDockerContainerRegistry(),
			"octopusdeploy_dynamic_worker_pool":                            resourceDynamicWorkerPool(),
			"octopusdeploy_gcp_account":                                    resourceGoogleCloudPlatformAccount(),
			"octopusdeploy_kubernetes_agent_deployment_target":             resourceKubernetesAgentDeploymentTarget(),
			"octopusdeploy_kubernetes_cluster_deployment_target":           resourceKubernetesClusterDeploymentTarget(),
			"octopusdeploy_listening_tentacle_deployment_target":           resourceListeningTentacleDeploymentTarget(),
			"octopusdeploy_machine_policy":                                 resourceMachinePolicy(),
			"octopusdeploy_offline_package_drop_deployment_target":         resourceOfflinePackageDropDeploymentTarget(),
			"octopusdeploy_polling_tentacle_deployment_target":             resourcePollingTentacleDeploymentTarget(),
			"octopusdeploy_polling_subscription_id":                        resourcePollingSubscriptionId(),
			"octopusdeploy_project_deployment_target_trigger":              resourceProjectDeploymentTargetTrigger(),
			"octopusdeploy_external_feed_create_release_trigger":           resourceExternalFeedCreateReleaseTrigger(),
			"octopusdeploy_project_scheduled_trigger":                      resourceProjectScheduledTrigger(),
			"octopusdeploy_runbook":                                        resourceRunbook(),
			"octopusdeploy_runbook_process":                                resourceRunbookProcess(),
			"octopusdeploy_scoped_user_role":                               resourceScopedUserRole(),
			"octopusdeploy_script_module":                                  resourceScriptModule(),
			"octopusdeploy_ssh_connection_deployment_target":               resourceSSHConnectionDeploymentTarget(),
			"octopusdeploy_ssh_key_account":                                resourceSSHKeyAccount(),
			"octopusdeploy_static_worker_pool":                             resourceStaticWorkerPool(),
			"octopusdeploy_tag":                                            resourceTag(),
			"octopusdeploy_tag_set":                                        resourceTagSet(),
			"octopusdeploy_team":                                           resourceTeam(),
			"octopusdeploy_tenant":                                         resourceTenant(),
			"octopusdeploy_tentacle_certificate":                           resourceTentacleCertificate(),
			"octopusdeploy_token_account":                                  resourceTokenAccount(),
			"octopusdeploy_user":                                           resourceUser(),
			"octopusdeploy_user_role":                                      resourceUserRole(),
			"octopusdeploy_username_password_account":                      resourceUsernamePasswordAccount(),
		},
		Schema: map[string]*schema.Schema{
			"address": {
				DefaultFunc: schema.EnvDefaultFunc("OCTOPUS_URL", nil),
				Description: "The endpoint of the Octopus REST API",
				Optional:    true,
				Type:        schema.TypeString,
			},
			"api_key": {
				DefaultFunc: schema.EnvDefaultFunc("OCTOPUS_APIKEY", nil),
				Description: "The API key to use with the Octopus REST API",
				Optional:    true,
				Type:        schema.TypeString,
			},
			"space_id": {
				Description: "The space ID to target",
				Optional:    true,
				Type:        schema.TypeString,
			},
		},

		ConfigureContextFunc: providerConfigure,
	}
}

func providerConfigure(ctx context.Context, d *schema.ResourceData) (interface{}, diag.Diagnostics) {
	config := Config{
		Address: d.Get("address").(string),
		APIKey:  d.Get("api_key").(string),
	}
	if spaceID, ok := d.GetOk("space_id"); ok {
		config.SpaceID = spaceID.(string)
	}

	return config.Client()
}
