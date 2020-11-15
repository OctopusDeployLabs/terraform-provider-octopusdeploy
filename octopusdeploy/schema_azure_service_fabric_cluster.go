package octopusdeploy

import (
	"github.com/OctopusDeploy/go-octopusdeploy/octopusdeploy"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func expandAzureServiceFabricCluster(d *schema.ResourceData) *octopusdeploy.AzureServiceFabricEndpoint {
	endpoint := octopusdeploy.NewAzureServiceFabricEndpoint()
	endpoint.ID = d.Id()

	if v, ok := d.GetOk("aad_client_credential_secret"); ok {
		endpoint.AadClientCredentialSecret = v.(string)
	}

	if v, ok := d.GetOk("aad_credential_type"); ok {
		endpoint.AadCredentialType = v.(string)
	}

	if v, ok := d.GetOk("aad_user_credential_password"); ok {
		endpoint.AadUserCredentialPassword = octopusdeploy.NewSensitiveValue(v.(string))
	}

	if v, ok := d.GetOk("aad_user_credential_username"); ok {
		endpoint.AadUserCredentialUsername = v.(string)
	}

	if v, ok := d.GetOk("certificate_store_location"); ok {
		endpoint.CertificateStoreLocation = v.(string)
	}

	if v, ok := d.GetOk("certificate_store_name"); ok {
		endpoint.CertificateStoreName = v.(string)
	}

	if v, ok := d.GetOk("client_certificate_variable"); ok {
		endpoint.ClientCertificateVariable = v.(string)
	}

	if v, ok := d.GetOk("connection_endpoint"); ok {
		endpoint.ConnectionEndpoint = v.(string)
	}

	if v, ok := d.GetOk("security_mode"); ok {
		endpoint.SecurityMode = v.(string)
	}

	if v, ok := d.GetOk("server_certificate_thumbprint"); ok {
		endpoint.ServerCertificateThumbprint = v.(string)
	}

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
		"id": {
			Computed: true,
			Type:     schema.TypeString,
		},
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
