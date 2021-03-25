package octopusdeploy

import (
	"github.com/OctopusDeploy/go-octopusdeploy/octopusdeploy"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/jinzhu/copier"
)

func expandEndpoint(values interface{}) octopusdeploy.IEndpoint {
	if values == nil {
		return nil
	}

	flattenedValues := values.([]interface{})
	if len(flattenedValues) == 0 {
		return nil
	}

	flattenedEndpoint := flattenedValues[0].(map[string]interface{})

	communicationStyle := flattenedEndpoint["communication_style"].(string)
	switch communicationStyle {
	case "AzureCloudService":
		return expandAzureCloudService(flattenedEndpoint)
	case "AzureServiceFabricCluster":
		return expandAzureServiceFabricCluster(flattenedEndpoint)
	case "AzureWebApp":
		return expandAzureWebApp(flattenedEndpoint)
	case "Kubernetes":
		return expandKubernetesCluster(flattenedEndpoint)
	case "None":
		return expandCloudRegion(flattenedEndpoint)
	case "OfflineDrop":
		return expandOfflinePackageDrop(flattenedEndpoint)
	case "Ssh":
		return expandSSHConnection(flattenedEndpoint)
	case "TentacleActive":
		return expandPollingTentacle(flattenedEndpoint)
	case "TentaclePassive":
		return expandListeningTentacle(flattenedEndpoint)
	}

	return nil
}

func flattenEndpoint(endpoint *octopusdeploy.EndpointResource) []interface{} {
	if endpoint == nil {
		return nil
	}

	switch endpoint.CommunicationStyle {
	case "AzureCloudService":
		azureCloudServiceEndpoint := octopusdeploy.NewAzureCloudServiceEndpoint()
		copier.Copy(azureCloudServiceEndpoint, endpoint)
		return flattenAzureCloudService(azureCloudServiceEndpoint)
	case "AzureServiceFabricCluster":
		azureServiceFabricEndpoint := octopusdeploy.NewAzureServiceFabricEndpoint()
		copier.Copy(azureServiceFabricEndpoint, endpoint)
		return flattenAzureServiceFabricCluster(azureServiceFabricEndpoint)
	case "AzureWebApp":
		azureWebAppEndpoint := octopusdeploy.NewAzureWebAppEndpoint()
		copier.Copy(azureWebAppEndpoint, endpoint)
		return flattenAzureWebApp(azureWebAppEndpoint)
	case "Kubernetes":
		kubernetesEndpoint := octopusdeploy.NewKubernetesEndpoint(endpoint.ClusterURL)
		copier.Copy(kubernetesEndpoint, endpoint)
		return flattenKubernetesCluster(kubernetesEndpoint)
	case "None":
		cloudRegionEndpoint := octopusdeploy.NewCloudRegionEndpoint()
		copier.Copy(cloudRegionEndpoint, endpoint)
		return flattenCloudRegion(cloudRegionEndpoint)
	case "OfflineDrop":
		offlinePackageDropEndpoint := octopusdeploy.NewOfflinePackageDropEndpoint()
		copier.Copy(offlinePackageDropEndpoint, endpoint)
		return flattenOfflinePackageDrop(offlinePackageDropEndpoint)
	case "Ssh":
		sshEndpoint := octopusdeploy.NewSSHEndpoint(endpoint.Host, endpoint.Port, endpoint.Fingerprint)
		copier.Copy(sshEndpoint, endpoint)
		return flattenSSHConnection(sshEndpoint)
	case "TentacleActive":
		pollingTentacleEndpoint := octopusdeploy.NewPollingTentacleEndpoint(endpoint.URI, endpoint.Thumbprint)
		copier.Copy(pollingTentacleEndpoint, endpoint)
		return flattenPollingTentacle(pollingTentacleEndpoint)
	case "TentaclePassive":
		listeningTentacleEndpoint := octopusdeploy.NewListeningTentacleEndpoint(endpoint.URI, endpoint.Thumbprint)
		copier.Copy(listeningTentacleEndpoint, endpoint)
		return flattenListeningTentacle(listeningTentacleEndpoint)
	}

	rawEndpoint := map[string]interface{}{
		"aad_client_credential_secret":    endpoint.AadClientCredentialSecret,
		"aad_credential_type":             endpoint.AadCredentialType,
		"aad_user_credential_username":    endpoint.AadUserCredentialUsername,
		"account_id":                      endpoint.AccountID,
		"applications_directory":          endpoint.ApplicationsDirectory,
		"authentication":                  flattenKubernetesAuthentication(endpoint.Authentication),
		"certificate_signature_algorithm": endpoint.CertificateSignatureAlgorithm,
		"certificate_store_location":      endpoint.CertificateStoreLocation,
		"certificate_store_name":          endpoint.CertificateStoreName,
		"client_certificate_variable":     endpoint.ClientCertificateVariable,
		"cloud_service_name":              endpoint.CloudServiceName,
		"cluster_certificate":             endpoint.ClusterCertificate,
		"communication_style":             endpoint.CommunicationStyle,
		"connection_endpoint":             endpoint.ConnectionEndpoint,
		"container":                       flattenDeploymentActionContainer(endpoint.Container),
		"default_worker_pool_id":          endpoint.DefaultWorkerPoolID,
		"destination":                     flattenOfflinePackageDropDestination(endpoint.Destination),
		"dot_net_core_platform":           endpoint.DotNetCorePlatform,
		"fingerprint":                     endpoint.Fingerprint,
		"host":                            endpoint.Host,
		"id":                              endpoint.GetID(),
		"namespace":                       endpoint.Namespace,
		"proxy_id":                        endpoint.ProxyID,
		"port":                            endpoint.Port,
		"resource_group_name":             endpoint.ResourceGroupName,
		"running_in_container":            endpoint.RunningInContainer,
		"security_mode":                   endpoint.SecurityMode,
		"server_certificate_thumbprint":   endpoint.ServerCertificateThumbprint,
		"skip_tls_verification":           endpoint.SkipTLSVerification,
		"slot":                            endpoint.Slot,
		"storage_account_name":            endpoint.StorageAccountName,
		"swap_if_possible":                endpoint.SwapIfPossible,
		"tentacle_version_details":        flattenTentacleVersionDetails(endpoint.TentacleVersionDetails),
		"thumbprint":                      endpoint.Thumbprint,
		"working_directory":               endpoint.WorkingDirectory,
		"use_current_instance_count":      endpoint.UseCurrentInstanceCount,
		"web_app_name":                    endpoint.WebAppName,
		"web_app_slot_name":               endpoint.WebAppSlotName,
	}

	if endpoint.ClusterURL != nil {
		rawEndpoint["cluster_url"] = endpoint.ClusterURL.String()
	}

	if endpoint.URI != nil {
		rawEndpoint["uri"] = endpoint.URI.String()
	}

	return []interface{}{rawEndpoint}
}

func flattenEndpointResource(endpoint *octopusdeploy.EndpointResource) []interface{} {
	if endpoint == nil {
		return nil
	}

	rawEndpoint := map[string]interface{}{
		"aad_client_credential_secret":    endpoint.AadClientCredentialSecret,
		"aad_credential_type":             endpoint.AadCredentialType,
		"aad_user_credential_username":    endpoint.AadUserCredentialUsername,
		"account_id":                      endpoint.AccountID,
		"applications_directory":          endpoint.ApplicationsDirectory,
		"authentication":                  flattenKubernetesAuthentication(endpoint.Authentication),
		"certificate_signature_algorithm": endpoint.CertificateSignatureAlgorithm,
		"certificate_store_location":      endpoint.CertificateStoreLocation,
		"certificate_store_name":          endpoint.CertificateStoreName,
		"client_certificate_variable":     endpoint.ClientCertificateVariable,
		"cloud_service_name":              endpoint.CloudServiceName,
		"cluster_certificate":             endpoint.ClusterCertificate,
		"communication_style":             endpoint.CommunicationStyle,
		"connection_endpoint":             endpoint.ConnectionEndpoint,
		"container":                       flattenDeploymentActionContainer(endpoint.Container),
		"default_worker_pool_id":          endpoint.DefaultWorkerPoolID,
		"destination":                     flattenOfflinePackageDropDestination(endpoint.Destination),
		"dot_net_core_platform":           endpoint.DotNetCorePlatform,
		"fingerprint":                     endpoint.Fingerprint,
		"host":                            endpoint.Host,
		"id":                              endpoint.GetID(),
		"namespace":                       endpoint.Namespace,
		"proxy_id":                        endpoint.ProxyID,
		"port":                            endpoint.Port,
		"resource_group_name":             endpoint.ResourceGroupName,
		"running_in_container":            endpoint.RunningInContainer,
		"security_mode":                   endpoint.SecurityMode,
		"server_certificate_thumbprint":   endpoint.ServerCertificateThumbprint,
		"skip_tls_verification":           endpoint.SkipTLSVerification,
		"slot":                            endpoint.Slot,
		"storage_account_name":            endpoint.StorageAccountName,
		"swap_if_possible":                endpoint.SwapIfPossible,
		"tentacle_version_details":        flattenTentacleVersionDetails(endpoint.TentacleVersionDetails),
		"thumbprint":                      endpoint.Thumbprint,
		"working_directory":               endpoint.WorkingDirectory,
		"use_current_instance_count":      endpoint.UseCurrentInstanceCount,
		"web_app_name":                    endpoint.WebAppName,
		"web_app_slot_name":               endpoint.WebAppSlotName,
	}

	if endpoint.ClusterURL != nil {
		rawEndpoint["cluster_url"] = endpoint.ClusterURL.String()
	}

	if endpoint.URI != nil {
		rawEndpoint["uri"] = endpoint.URI.String()
	}

	return []interface{}{rawEndpoint}
}

func getEndpointSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"aad_client_credential_secret": {
			Optional: true,
			Type:     schema.TypeString,
		},
		"aad_credential_type": {
			Optional: true,
			Type:     schema.TypeString,
		},
		"aad_user_credential_username": {
			Optional: true,
			Type:     schema.TypeString,
		},
		"account_id": {
			Optional: true,
			Type:     schema.TypeString,
		},
		"applications_directory": {
			Optional: true,
			Type:     schema.TypeString,
		},
		"authentication": {
			Computed: true,
			Elem:     &schema.Resource{Schema: getKubernetesAuthenticationSchema()},
			MaxItems: 1,
			MinItems: 0,
			Optional: true,
			Type:     schema.TypeSet,
		},
		"certificate_signature_algorithm": {
			Optional: true,
			Type:     schema.TypeString,
		},
		"certificate_store_location": {
			Optional: true,
			Type:     schema.TypeString,
		},
		"certificate_store_name": {
			Optional: true,
			Type:     schema.TypeString,
		},
		"client_certificate_variable": {
			Optional: true,
			Type:     schema.TypeString,
		},
		"cloud_service_name": {
			Optional: true,
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
		"communication_style": {
			Type:     schema.TypeString,
			Required: true,
			ValidateDiagFunc: validation.ToDiagFunc(validation.StringInSlice([]string{
				"AzureCloudService",
				"AzureWebApp",
				"Ftp",
				"Kubernetes",
				"None",
				"OfflineDrop",
				"Ssh",
				"TentacleActive",
				"TentaclePassive",
			}, false)),
		},
		"connection_endpoint": {
			Optional: true,
			Type:     schema.TypeString,
		},
		"container": {
			Computed: true,
			Elem:     &schema.Resource{Schema: getDeploymentActionContainerSchema()},
			Optional: true,
			Type:     schema.TypeList,
		},
		"default_worker_pool_id": {
			Optional: true,
			Type:     schema.TypeString,
		},
		"destination": {
			Computed: true,
			Elem:     &schema.Resource{Schema: getOfflinePackageDropDestinationSchema()},
			Optional: true,
			Type:     schema.TypeList,
		},
		"dot_net_core_platform": {
			Optional: true,
			Type:     schema.TypeString,
		},
		"fingerprint": {
			Optional: true,
			Type:     schema.TypeString,
		},
		"host": {
			Optional: true,
			Type:     schema.TypeString,
		},
		"id": getIDSchema(),
		"namespace": {
			Optional: true,
			Type:     schema.TypeString,
		},
		"proxy_id": {
			Optional: true,
			Type:     schema.TypeString,
		},
		"port": {
			Optional: true,
			Type:     schema.TypeInt,
		},
		"resource_group_name": {
			Optional: true,
			Type:     schema.TypeString,
		},
		"running_in_container": {
			Optional: true,
			Type:     schema.TypeBool,
		},
		"security_mode": {
			Optional: true,
			Type:     schema.TypeString,
		},
		"server_certificate_thumbprint": {
			Optional: true,
			Type:     schema.TypeString,
		},
		"skip_tls_verification": {
			Optional: true,
			Type:     schema.TypeBool,
		},
		"slot": {
			Optional: true,
			Type:     schema.TypeString,
		},
		"storage_account_name": {
			Optional: true,
			Type:     schema.TypeString,
		},
		"swap_if_possible": {
			Optional: true,
			Type:     schema.TypeBool,
		},
		"tentacle_version_details": {
			Computed: true,
			Elem:     &schema.Resource{Schema: getTentacleVersionDetailsSchema()},
			Optional: true,
			Type:     schema.TypeList,
		},
		"thumbprint": {
			Optional: true,
			Type:     schema.TypeString,
		},
		"working_directory": {
			Optional: true,
			Type:     schema.TypeString,
		},
		"use_current_instance_count": {
			Optional: true,
			Type:     schema.TypeBool,
		},
		"uri": {
			Optional: true,
			Type:     schema.TypeString,
		},
		"web_app_name": {
			Optional: true,
			Type:     schema.TypeString,
		},
		"web_app_slot_name": {
			Optional: true,
			Type:     schema.TypeString,
		},
	}
}
