package octopusdeploy

import (
	"github.com/OctopusDeploy/go-octopusdeploy/octopusdeploy"
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
