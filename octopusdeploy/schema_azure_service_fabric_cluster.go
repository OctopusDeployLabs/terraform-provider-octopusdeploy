package octopusdeploy

import (
	"github.com/OctopusDeploy/go-octopusdeploy/octopusdeploy"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func expandAzureServiceFabricCluster(flattenedMap map[string]interface{}) *octopusdeploy.AzureServiceFabricEndpoint {
	endpoint := octopusdeploy.NewAzureServiceFabricEndpoint()
	endpoint.AadClientCredentialSecret = flattenedMap["aad_client_credential_secret"].(string)
	endpoint.AadCredentialType = flattenedMap["aad_credential_type"].(string)
	endpoint.AadUserCredentialPassword = octopusdeploy.NewSensitiveValue(flattenedMap["aad_user_credential_password"].(string))
	endpoint.AadUserCredentialUsername = flattenedMap["aad_user_credential_username"].(string)
	endpoint.CertificateStoreLocation = flattenedMap["certificate_store_location"].(string)
	endpoint.CertificateStoreName = flattenedMap["certificate_store_name"].(string)
	endpoint.ClientCertificateVariable = flattenedMap["client_certificate_variable"].(string)
	endpoint.ConnectionEndpoint = flattenedMap["connection_endpoint"].(string)
	endpoint.ID = flattenedMap["id"].(string)
	endpoint.SecurityMode = flattenedMap["security_mode"].(string)
	endpoint.ServerCertificateThumbprint = flattenedMap["server_certificate_thumbprint"].(string)

	return endpoint
}

func flattenAzureServiceFabricCluster(endpoint *octopusdeploy.AzureServiceFabricEndpoint) []interface{} {
	if endpoint == nil {
		return nil
	}

	return []interface{}{map[string]interface{}{
		"aad_client_credential_secret":  endpoint.AadClientCredentialSecret,
		"aad_credential_type":           endpoint.AadCredentialType,
		"aad_user_credential_username":  endpoint.AadUserCredentialUsername,
		"certificate_store_location":    endpoint.CertificateStoreLocation,
		"certificate_store_name":        endpoint.CertificateStoreName,
		"client_certificate_variable":   endpoint.ClientCertificateVariable,
		"connection_endpoint":           endpoint.ConnectionEndpoint,
		"id":                            endpoint.GetID(),
		"security_mode":                 endpoint.SecurityMode,
		"server_certificate_thumbprint": endpoint.ServerCertificateThumbprint,
	}}
}

func getAzureServiceFabricClusterSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"aad_client_credential_secret": {
			Optional: true,
			Type:     schema.TypeString,
		},
		"aad_credential_type": {
			Optional: true,
			Type:     schema.TypeString,
		},
		"aad_user_credential_password": {
			Optional:  true,
			Sensitive: true,
			Type:      schema.TypeString,
		},
		"aad_user_credential_username": {
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
		"connection_endpoint": {
			Optional: true,
			Type:     schema.TypeString,
		},
		"id": getIDSchema(),
		"security_mode": {
			Optional: true,
			Type:     schema.TypeString,
		},
		"server_certificate_thumbprint": {
			Optional: true,
			Type:     schema.TypeString,
		},
	}
}
