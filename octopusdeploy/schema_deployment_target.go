package octopusdeploy

import (
	"context"

	"github.com/OctopusDeploy/go-octopusdeploy/octopusdeploy"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func expandDeploymentTarget(d *schema.ResourceData) *octopusdeploy.DeploymentTarget {
	deploymentMode := octopusdeploy.TenantedDeploymentMode(d.Get("tenanted_deployment_participation").(string))
	environments := getSliceFromTerraformTypeList(d.Get("environments"))
	name := d.Get("name").(string)
	roles := getSliceFromTerraformTypeList(d.Get("roles"))
	tenantIDs := getSliceFromTerraformTypeList(d.Get("tenants"))
	tenantTags := getSliceFromTerraformTypeList(d.Get("tenant_tags"))

	endpoint := expandEndpoint(d.Get("endpoint"))

	// tfSchemaSetInterface, ok := d.GetOk("endpoint")
	// if !ok {
	// 	return nil
	// }
	// tfSchemaSet := tfSchemaSetInterface.([]interface{})
	// if len(tfSchemaSet) == 0 {
	// 	return nil
	// }
	// // Get the first element in the list, which is a map of the interfaces
	// tfSchemaList := tfSchemaSet[0].(map[string]interface{})

	// communicationStyle := octopusdeploy.CommunicationStyle(tfSchemaList["communication_style"].(string))

	// var endpoint octopusdeploy.IEndpoint
	// switch communicationStyle {
	// case "AzureCloudService":
	// 	endpoint = expandAzureCloudService(d)
	// case "AzureServiceFabricCluster":
	// 	endpoint = expandAzureServiceFabricCluster(d)
	// case "AzureWebApp":
	// 	endpoint = expandAzureWebApp(d)
	// case "Kubernetes":
	// 	endpoint = expandKubernetesCluster(d)
	// case "None":
	// 	endpoint = expandCloudRegion(d)
	// case "OfflineDrop":
	// 	endpoint = expandOfflineDrop(d)
	// case "Ssh":
	// 	endpoint = expandSSHConnection(d)
	// case "TentacleActive":
	// 	endpoint = expandPollingTentacle(d)
	// case "TentaclePassive":
	// 	endpoint = expandListeningTentacle(d)
	// }

	deploymentTarget := octopusdeploy.NewDeploymentTarget(name, endpoint, environments, roles)
	deploymentTarget.ID = d.Id()
	deploymentTarget.TenantedDeploymentMode = deploymentMode
	deploymentTarget.TenantIDs = tenantIDs
	deploymentTarget.TenantTags = tenantTags

	if v, ok := d.GetOk("machine_policy_id"); ok {
		deploymentTarget.MachinePolicyID = v.(string)
	}

	if v, ok := d.GetOk("is_disabled"); ok {
		deploymentTarget.IsDisabled = v.(bool)
	}

	if v, ok := d.GetOk("thumbprint"); ok {
		deploymentTarget.Thumbprint = v.(string)
	}

	if v, ok := d.GetOk("uri"); ok {
		deploymentTarget.URI = v.(string)
	}

	return deploymentTarget
}

func flattenDeploymentTarget(deploymentTarget *octopusdeploy.DeploymentTarget) map[string]interface{} {
	if deploymentTarget == nil {
		return nil
	}

	// endpointResource, _ := octopusdeploy.ToEndpointResource(deploymentTarget.Endpoint)

	flattenedDeploymentTarget := map[string]interface{}{
		// "endpoint":                          flattenEndpoint(endpointResource),
		"environments":                      deploymentTarget.EnvironmentIDs,
		"has_latest_calamari":               deploymentTarget.HasLatestCalamari,
		"health_status":                     deploymentTarget.HealthStatus,
		"id":                                deploymentTarget.GetID(),
		"is_disabled":                       deploymentTarget.IsDisabled,
		"is_in_process":                     deploymentTarget.IsInProcess,
		"machine_policy_id":                 deploymentTarget.MachinePolicyID,
		"name":                              deploymentTarget.Name,
		"operating_system":                  deploymentTarget.OperatingSystem,
		"roles":                             deploymentTarget.Roles,
		"shell_name":                        deploymentTarget.ShellName,
		"shell_version":                     deploymentTarget.ShellVersion,
		"space_id":                          deploymentTarget.SpaceID,
		"status":                            deploymentTarget.Status,
		"status_summary":                    deploymentTarget.StatusSummary,
		"tenanted_deployment_participation": deploymentTarget.TenantedDeploymentMode,
		"tenants":                           deploymentTarget.TenantIDs,
		"tenant_tags":                       deploymentTarget.TenantTags,
		"thumbprint":                        deploymentTarget.Thumbprint,
		"uri":                               deploymentTarget.URI,
	}

	switch deploymentTarget.Endpoint.GetCommunicationStyle() {
	case "AzureCloudService":
		flattenedDeploymentTarget["azure_cloud_service"] = flattenAzureCloudService(deploymentTarget.Endpoint.(*octopusdeploy.AzureCloudServiceEndpoint))
	case "AzureServiceFabricCluster":
		flattenedDeploymentTarget["azure_service_fabric_cluster"] = flattenAzureServiceFabricCluster(deploymentTarget.Endpoint.(*octopusdeploy.AzureServiceFabricEndpoint))
	case "AzureWebApp":
		flattenedDeploymentTarget["azure_web_app"] = flattenAzureWebApp(deploymentTarget.Endpoint.(*octopusdeploy.AzureWebAppEndpoint))
	case "Kubernetes":
		flattenedDeploymentTarget["kubernetes_cluster"] = flattenKubernetesCluster(deploymentTarget.Endpoint.(*octopusdeploy.KubernetesEndpoint))
	case "None":
		flattenedDeploymentTarget["cloud_region"] = flattenCloudRegion(deploymentTarget.Endpoint.(*octopusdeploy.CloudRegionEndpoint))
	case "OfflineDrop":
		flattenedDeploymentTarget["offline_drop"] = flattenOfflineDrop(deploymentTarget.Endpoint.(*octopusdeploy.OfflineDropEndpoint))
	case "Ssh":
		flattenedDeploymentTarget["ssh_connection"] = flattenSSHConnection(deploymentTarget.Endpoint.(*octopusdeploy.SSHEndpoint))
	case "TentacleActive":
		flattenedDeploymentTarget["polling_tentacle"] = flattenPollingTentacle(deploymentTarget.Endpoint.(*octopusdeploy.PollingTentacleEndpoint))
	case "TentaclePassive":
		flattenedDeploymentTarget["listening_tentacle"] = flattenListeningTentacle(deploymentTarget.Endpoint.(*octopusdeploy.ListeningTentacleEndpoint))
	}

	return flattenedDeploymentTarget
}

func getDeploymentTargetDataSchema() map[string]*schema.Schema {
	deploymentTargetsSchema := getDeploymentTargetSchema()
	for _, field := range deploymentTargetsSchema {
		field.Computed = true
		field.Default = nil
		field.MaxItems = 0
		field.MinItems = 0
		field.Optional = false
		field.Required = false
		field.ValidateDiagFunc = nil
	}

	return map[string]*schema.Schema{
		"communication_styles": {
			Description: "A list of deployment target communication styles to be matched",
			Elem:        &schema.Schema{Type: schema.TypeString},
			Optional:    true,
			Type:        schema.TypeList,
		},
		"deployment_id": {
			Description: "A deployment ID to be matched",
			Optional:    true,
			Type:        schema.TypeString,
		},
		"deployment_targets": {
			Computed:    true,
			Description: "A computed list of deployment targets that are matched based on the criteria set for this data source",
			Elem:        &schema.Resource{Schema: deploymentTargetsSchema},
			Type:        schema.TypeList,
		},
		"environments": {
			Description: "A list of environments to be matched",
			Elem:        &schema.Schema{Type: schema.TypeString},
			Optional:    true,
			Type:        schema.TypeList,
		},
		"health_statuses": {
			Description: "A list of deployment target health statuses to be matched",
			Elem:        &schema.Schema{Type: schema.TypeString},
			Optional:    true,
			Type:        schema.TypeList,
		},
		"ids": {
			Description: "A list of deployment target IDs to be matched",
			Elem:        &schema.Schema{Type: schema.TypeString},
			Optional:    true,
			Type:        schema.TypeList,
		},
		"is_disabled": {
			Description: "The state of deployment targets to be matched",
			Optional:    true,
			Type:        schema.TypeBool,
		},
		"name": {
			Description: "The name of the deployment target to be matched",
			Optional:    true,
			Type:        schema.TypeString,
		},
		"partial_name": {
			Description: "The partial name of a deployment target to be matched",
			Optional:    true,
			Type:        schema.TypeString,
		},
		"roles": {
			Description: "A list of role IDs to be matched",
			Elem:        &schema.Schema{Type: schema.TypeString},
			Optional:    true,
			Type:        schema.TypeList,
		},
		"shell_names": {
			Description: "A list of shell names to be matched",
			Elem:        &schema.Schema{Type: schema.TypeString},
			Optional:    true,
			Type:        schema.TypeList,
		},
		"skip": {
			Default:  0,
			Optional: true,
			Type:     schema.TypeInt,
		},
		"take": {
			Default:  1,
			Optional: true,
			Type:     schema.TypeInt,
		},
		"tenants": {
			Elem:     &schema.Schema{Type: schema.TypeString},
			Optional: true,
			Type:     schema.TypeList,
		},
		"tenant_tags": {
			Elem:     &schema.Schema{Type: schema.TypeString},
			Optional: true,
			Type:     schema.TypeList,
		},
		"thumbprint": {
			Optional: true,
			Type:     schema.TypeString,
		},
	}
}

func getDeploymentTargetSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"azure_cloud_service": {
			Computed: true,
			Elem:     &schema.Resource{Schema: getAzureCloudServiceSchema()},
			MaxItems: 1,
			Optional: true,
			Type:     schema.TypeList,
		},
		"azure_service_fabric_cluster": {
			Computed: true,
			Elem:     &schema.Resource{Schema: getAzureServiceFabricClusterSchema()},
			MaxItems: 1,
			Optional: true,
			Type:     schema.TypeList,
		},
		"azure_web_app": {
			Computed: true,
			Elem:     &schema.Resource{Schema: getAzureWebAppSchema()},
			MaxItems: 1,
			Optional: true,
			Type:     schema.TypeList,
		},
		"cloud_region": {
			Computed: true,
			Elem:     &schema.Resource{Schema: getCloudRegionSchema()},
			MaxItems: 1,
			Optional: true,
			Type:     schema.TypeList,
		},
		"endpoint": {
			Deprecated: "use endpoint-specific attribute instead (i.e. azure_cloud_service, azure_service_fabric_cluster, azure_web_app, cloud_region, kubernetes_cluster, offline_drop, ssh_connection, polling_tentacle, listening_tentacle)",
			Elem:       &schema.Resource{Schema: getEndpointSchema()},
			MaxItems:   1,
			Optional:   true,
			Type:       schema.TypeList,
		},
		"environments": {
			Elem:     &schema.Schema{Type: schema.TypeString},
			Required: true,
			Type:     schema.TypeList,
		},
		"has_latest_calamari": {
			Computed: true,
			Type:     schema.TypeBool,
		},
		"health_status": {
			Computed: true,
			Optional: true,
			Type:     schema.TypeString,
			ValidateDiagFunc: validateDiagFunc(validation.StringInSlice([]string{
				"HasWarnings",
				"Healthy",
				"Unavailable",
				"Unhealthy",
				"Unknown",
			}, false)),
		},
		"id": {
			Computed: true,
			Type:     schema.TypeString,
		},
		"is_disabled": {
			Computed: true,
			Optional: true,
			Type:     schema.TypeBool,
		},
		"is_in_process": {
			Computed: true,
			Type:     schema.TypeBool,
		},
		"kubernetes_cluster": {
			Computed: true,
			Elem:     &schema.Resource{Schema: getKubernetesClusterSchema()},
			MaxItems: 1,
			Optional: true,
			Type:     schema.TypeList,
		},
		"listening_tentacle": {
			Computed: true,
			Elem:     &schema.Resource{Schema: getListeningTentacleSchema()},
			MaxItems: 1,
			Optional: true,
			Type:     schema.TypeList,
		},
		"machine_policy_id": {
			Computed: true,
			Optional: true,
			Type:     schema.TypeString,
		},
		"name": {
			Required:         true,
			Type:             schema.TypeString,
			ValidateDiagFunc: validateDiagFunc(validation.StringIsNotEmpty),
		},
		"offline_drop": {
			Computed: true,
			Elem:     &schema.Resource{Schema: getOfflineDropSchema()},
			MaxItems: 1,
			Optional: true,
			Type:     schema.TypeList,
		},
		"operating_system": {
			Computed: true,
			Optional: true,
			Type:     schema.TypeString,
		},
		"polling_tentacle": {
			Computed: true,
			Elem:     &schema.Resource{Schema: getPollingTentacleSchema()},
			MaxItems: 1,
			Optional: true,
			Type:     schema.TypeList,
		},
		"roles": {
			Elem:     &schema.Schema{Type: schema.TypeString},
			MinItems: 1,
			Required: true,
			Type:     schema.TypeList,
		},
		"shell_name": {
			Computed: true,
			Optional: true,
			Type:     schema.TypeString,
		},
		"shell_version": {
			Computed: true,
			Optional: true,
			Type:     schema.TypeString,
		},
		"space_id": {
			Computed: true,
			Type:     schema.TypeString,
		},
		"ssh_connection": {
			Computed: true,
			Elem:     &schema.Resource{Schema: getSSHConnectionSchema()},
			MaxItems: 1,
			Optional: true,
			Type:     schema.TypeList,
		},
		"status": {
			Computed: true,
			Optional: true,
			Type:     schema.TypeString,
		},
		"status_summary": {
			Computed: true,
			Optional: true,
			Type:     schema.TypeString,
		},
		"tenanted_deployment_participation": getTenantedDeploymentSchema(),
		"tenants": {
			Elem:     &schema.Schema{Type: schema.TypeString},
			Optional: true,
			Type:     schema.TypeList,
		},
		"tenant_tags": {
			Elem:     &schema.Schema{Type: schema.TypeString},
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

func setDeploymentTarget(ctx context.Context, d *schema.ResourceData, deploymentTarget *octopusdeploy.DeploymentTarget) {
	endpointResource, err := octopusdeploy.ToEndpointResource(deploymentTarget.Endpoint)
	if err != nil {
		return
	}

	switch deploymentTarget.Endpoint.GetCommunicationStyle() {
	case "AzureCloudService":
		d.Set("azure_cloud_service", flattenAzureCloudService(deploymentTarget.Endpoint.(*octopusdeploy.AzureCloudServiceEndpoint)))
	case "AzureServiceFabricCluster":
		d.Set("azure_service_fabric_cluster", flattenAzureServiceFabricCluster(deploymentTarget.Endpoint.(*octopusdeploy.AzureServiceFabricEndpoint)))
	case "AzureWebApp":
		d.Set("azure_web_app", flattenAzureWebApp(deploymentTarget.Endpoint.(*octopusdeploy.AzureWebAppEndpoint)))
	case "Kubernetes":
		d.Set("kubernetes_cluster", flattenKubernetesCluster(deploymentTarget.Endpoint.(*octopusdeploy.KubernetesEndpoint)))
	case "None":
		d.Set("cloud_region", flattenCloudRegion(deploymentTarget.Endpoint.(*octopusdeploy.CloudRegionEndpoint)))
	case "OfflineDrop":
		d.Set("offline_drop", flattenOfflineDrop(deploymentTarget.Endpoint.(*octopusdeploy.OfflineDropEndpoint)))
	case "Ssh":
		d.Set("ssh_connection", flattenSSHConnection(deploymentTarget.Endpoint.(*octopusdeploy.SSHEndpoint)))
	case "TentacleActive":
		d.Set("polling_tentacle", flattenPollingTentacle(deploymentTarget.Endpoint.(*octopusdeploy.PollingTentacleEndpoint)))
	case "TentaclePassive":
		d.Set("listening_tentacle", flattenListeningTentacle(deploymentTarget.Endpoint.(*octopusdeploy.ListeningTentacleEndpoint)))
	}

	d.Set("endpoint", flattenEndpoint(endpointResource))
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
