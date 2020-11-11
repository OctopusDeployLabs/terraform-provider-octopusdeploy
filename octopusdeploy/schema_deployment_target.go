package octopusdeploy

import (
	"context"
	"net/url"

	"github.com/OctopusDeploy/go-octopusdeploy/octopusdeploy"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func expandDeploymentTarget(d *schema.ResourceData) *octopusdeploy.DeploymentTarget {
	var name string
	if v, ok := d.GetOk("name"); ok {
		name = v.(string)
	}

	environments := getSliceFromTerraformTypeList(d.Get("environments"))
	roles := getSliceFromTerraformTypeList(d.Get("roles"))
	deploymentMode := octopusdeploy.TenantedDeploymentMode(d.Get("tenanted_deployment_participation").(string))
	tenantIDs := getSliceFromTerraformTypeList(d.Get("tenants"))
	tenantTags := getSliceFromTerraformTypeList(d.Get("tenant_tags"))

	tfSchemaSetInterface, ok := d.GetOk("endpoint")
	if !ok {
		return nil
	}
	tfSchemaSet := tfSchemaSetInterface.(*schema.Set)
	if len(tfSchemaSet.List()) == 0 {
		return nil
	}
	// Get the first element in the list, which is a map of the interfaces
	tfSchemaList := tfSchemaSet.List()[0].(map[string]interface{})

	var proxyID string
	if tfSchemaList["proxy_id"] != nil {
		proxyString := tfSchemaList["proxy_id"].(string)
		proxyID = proxyString
	}

	communicationStyle := octopusdeploy.CommunicationStyle(tfSchemaList["communication_style"].(string))

	var endpoint octopusdeploy.IEndpoint
	switch communicationStyle {
	case "AzureCloudService":
		azureCloudServiceEndpoint := octopusdeploy.NewAzureCloudServiceEndpoint()
		azureCloudServiceEndpoint.DefaultWorkerPoolID = tfSchemaList["default_worker_pool_id"].(string)
		endpoint = azureCloudServiceEndpoint
	case "AzureServiceFabricCluster":
		endpoint = octopusdeploy.NewServiceFabricEndpoint()
	case "AzureWebApp":
		endpoint = octopusdeploy.NewAzureWebAppEndpoint()
	case "Kubernetes":
		clusterURL := d.Get("cluster_url").(url.URL)
		kubernetesEndpoint := octopusdeploy.NewKubernetesEndpoint(clusterURL)
		kubernetesEndpoint.ClusterCertificate = tfSchemaList["cluster_certificate"].(string)
		kubernetesEndpoint.ClusterURL, _ = url.Parse(tfSchemaList["cluster_url"].(string))
		kubernetesEndpoint.Namespace = tfSchemaList["namespace"].(string)
		kubernetesEndpoint.ProxyID = proxyID
		kubernetesEndpoint.SkipTLSVerification = tfSchemaList["skip_tls_verification"].(bool)
		endpoint = kubernetesEndpoint
	case "None":
		endpoint = octopusdeploy.NewCloudRegionEndpoint()
	case "OfflineDrop":
		endpoint = octopusdeploy.NewOfflineDropEndpoint()
	case "Ssh":
		host := d.Get("host").(string)
		port := d.Get("port").(int)
		fingerprint := d.Get("fingerprint").(string)
		sshEndpoint := octopusdeploy.NewSSHEndpoint(host, port, fingerprint)
		sshEndpoint.ProxyID = proxyID
		endpoint = sshEndpoint
	case "TentacleActive":
		uri, _ := url.Parse(tfSchemaList["uri"].(string))
		thumbprint := tfSchemaList["thumbprint"].(string)
		endpoint = octopusdeploy.NewPollingTentacleEndpoint(uri, thumbprint)
	case "TentaclePassive":
		uri, _ := url.Parse(tfSchemaList["uri"].(string))
		thumbprint := tfSchemaList["thumbprint"].(string)
		endpoint = octopusdeploy.NewListeningTentacleEndpoint(uri, thumbprint)
	}

	deploymentTarget := octopusdeploy.NewDeploymentTarget(name, endpoint, environments, roles)
	deploymentTarget.ID = d.Id()

	if v, ok := d.GetOk("machine_policy_id"); ok {
		deploymentTarget.MachinePolicyID = v.(string)
	}

	if v, ok := d.GetOk("is_disabled"); ok {
		deploymentTarget.IsDisabled = v.(bool)
	}

	deploymentTarget.TenantedDeploymentMode = deploymentMode
	deploymentTarget.TenantIDs = tenantIDs
	deploymentTarget.TenantTags = tenantTags
	deploymentTarget.Thumbprint = tfSchemaList["thumbprint"].(string)
	deploymentTarget.URI = tfSchemaList["uri"].(string)

	tfAuthenticationSchemaSetInterface, ok := tfSchemaList["authentication"]
	if ok {
		tfAuthenticationSchemaSet := tfAuthenticationSchemaSetInterface.(*schema.Set)
		if len(tfAuthenticationSchemaSet.List()) == 1 {
			// Get the first element in the list, which is a map of the interfaces
			// tfAuthenticationSchemaList := tfAuthenticationSchemaSet.List()[0].(map[string]interface{})

			// deploymentTarget.Endpoint.Authentication = &octopusdeploy.DeploymentTargetEndpointAuthentication{
			// 	AccountID:          tfAuthenticationSchemaList["account_id"].(string),
			// 	ClientCertificate:  tfAuthenticationSchemaList["client_certificate"].(string),
			// 	AuthenticationType: tfAuthenticationSchemaList["authentication_type"].(string),
			// }
		}
	}

	return deploymentTarget
}

func flattenEndpoint(endpoint octopusdeploy.IEndpoint) map[string]interface{} {
	if endpoint == nil {
		return nil
	}

	flattenedEndpoint := map[string]interface{}{
		"communication_style": endpoint.GetCommunicationStyle(),
	}

	switch endpoint.GetCommunicationStyle() {
	case "AzureCloudService":
		azureCloudServiceEndpoint := endpoint.(*octopusdeploy.AzureCloudServiceEndpoint)
		flattenedEndpoint["account_id"] = azureCloudServiceEndpoint.AccountID
	case "AzureServiceFabricCluster":
	case "AzureWebApp":
	case "Kubernetes":
	case "None":
	case "OfflineDrop":
	case "Ssh":
	case "TentacleActive":
	case "TentaclePassive":
	}

	return flattenedEndpoint
}

func flattenDeploymentTarget(ctx context.Context, d *schema.ResourceData, deploymentTarget *octopusdeploy.DeploymentTarget) {
	d.Set("endpoint", flattenEndpoint(deploymentTarget.Endpoint))
	d.Set("environments", deploymentTarget.EnvironmentIDs)
	d.Set("has_latest_calamari", deploymentTarget.HasLatestCalamari)
	d.Set("health_status", deploymentTarget.HealthStatus)
	d.Set("is_disabled", deploymentTarget.IsDisabled)
	d.Set("is_in_process", deploymentTarget.IsInProcess)
	d.Set("machine_policy_id", deploymentTarget.MachinePolicyID)
	d.Set("name", deploymentTarget.Name)
	d.Set("operating_system", deploymentTarget.OperatingSystem)
	d.Set("roles", deploymentTarget.Roles)
	d.Set("shell_name", deploymentTarget.ShellName)
	d.Set("shell_version", deploymentTarget.ShellVersion)
	d.Set("space_id", deploymentTarget.SpaceID)
	d.Set("status", deploymentTarget.Status)
	d.Set("status_summary", deploymentTarget.StatusSummary)
	d.Set("tenanted_deployment_participation", deploymentTarget.TenantedDeploymentMode)
	d.Set("tenants", deploymentTarget.TenantIDs)
	d.Set("tenant_tags", deploymentTarget.TenantTags)
	d.Set("thumbprint", deploymentTarget.Thumbprint)
	d.Set("uri", deploymentTarget.URI)

	d.SetId(deploymentTarget.GetID())
}

func getDeploymentTargetSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"name": &schema.Schema{
			Required:     true,
			Type:         schema.TypeString,
			ValidateFunc: validation.StringIsNotEmpty,
		},
		"endpoint": {
			Type:     schema.TypeSet,
			MaxItems: 1,
			MinItems: 1,
			Required: true,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"communication_style": {
						Type:     schema.TypeString,
						Required: true,
						ValidateDiagFunc: validateDiagFunc(validation.StringInSlice([]string{
							"None",
							"TentaclePassive",
							"TentacleActive",
							"Ssh",
							"OfflineDrop",
							"AzureWebApp",
							"Ftp",
							"AzureCloudService",
							"Kubernetes",
						}, false)),
					},
					"proxy_id": {
						Optional: true,
						Type:     schema.TypeString,
					},
					"thumbprint": {
						Required: true,
						Type:     schema.TypeString,
					},
					"uri": {
						Required: true,
						Type:     schema.TypeString,
					},
					"cluster_certificate": {
						Optional: true,
						Type:     schema.TypeString,
					},
					"cluster_url": {
						Optional: true,
						Type:     schema.TypeString,
					},
					"namespace": {
						Optional: true,
						Type:     schema.TypeString,
					},
					"skip_tls_verification": {
						Optional: true,
						Type:     schema.TypeBool,
					},
					"default_worker_pool_id": {
						Optional: true,
						Type:     schema.TypeString,
					},
					"authentication": {
						Elem: &schema.Resource{
							Schema: map[string]*schema.Schema{
								"account_id": {
									Optional: true,
									Type:     schema.TypeString,
								},
								"client_certificate": {
									Optional: true,
									Type:     schema.TypeString,
								},
								"authentication_type": {
									Optional: true,
									Type:     schema.TypeString,
									ValidateDiagFunc: validateDiagFunc(validation.StringInSlice([]string{
										"KubernetesCertificate",
										"KubernetesStandard",
									}, false)),
								},
							},
						},
						MaxItems: 1,
						MinItems: 0,
						Optional: true,
						Type:     schema.TypeSet,
					},
				},
			},
		},
		"environments": {
			Elem: &schema.Schema{
				Type: schema.TypeString,
			},
			Required: true,
			Type:     schema.TypeList,
		},
		"has_latest_calamari": {
			Computed: true,
			Type:     schema.TypeBool,
		},
		"health_status": {
			Type:     schema.TypeString,
			Required: true,
			ValidateDiagFunc: validateDiagFunc(validation.StringInSlice([]string{
				"HasWarnings",
				"Healthy",
				"Unavailable",
				"Unhealthy",
				"Unknown",
			}, false)),
		},
		"is_disabled": {
			Required: true,
			Type:     schema.TypeBool,
		},
		"is_in_process": {
			Computed: true,
			Type:     schema.TypeBool,
		},
		"machine_policy_id": {
			Required: true,
			Type:     schema.TypeString,
		},
		"operating_system": {
			Optional: true,
			Type:     schema.TypeString,
		},
		"roles": {
			Elem: &schema.Schema{
				Type:     schema.TypeString,
				MinItems: 1,
			},
			Required: true,
			Type:     schema.TypeList,
		},
		"shell_name": {
			Optional: true,
			Type:     schema.TypeString,
		},
		"shell_version": {
			Optional: true,
			Type:     schema.TypeString,
		},
		"space_id": {
			Computed: true,
			Type:     schema.TypeString,
		},
		"status": {
			Computed: true,
			Type:     schema.TypeString,
		},
		"status_summary": {
			Computed: true,
			Type:     schema.TypeString,
		},
		"tenanted_deployment_participation": {
			Required: true,
			Type:     schema.TypeString,
			ValidateDiagFunc: validateDiagFunc(validation.StringInSlice([]string{
				"Untenanted",
				"TenantedOrUntenanted",
				"Tenanted",
			}, false)),
		},
		"tenants": {
			Elem: &schema.Schema{
				Type: schema.TypeString,
			},
			Optional: true,
			Type:     schema.TypeList,
		},
		"tenant_tags": {
			Elem: &schema.Schema{
				Type: schema.TypeString,
			},
			Optional: true,
			Type:     schema.TypeList,
		},
		"thumbprint": {
			Optional: true,
			Type:     schema.TypeString,
		},
		"uri": {
			Optional: true,
			Type:     schema.TypeString,
		},
	}
}
