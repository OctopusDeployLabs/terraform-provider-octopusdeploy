package octopusdeploy

import (
	"github.com/OctopusDeploy/go-octopusdeploy/octopusdeploy"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func expandEndpointAuthentication(values interface{}) octopusdeploy.EndpointAuthentication {
	flattenedValues := values.([]interface{})
	flattenedEndpointAuthentication := flattenedValues[0].(map[string]interface{})

	return octopusdeploy.EndpointAuthentication{
		AccountID:          flattenedEndpointAuthentication["account_id"].(string),
		AuthenticationType: flattenedEndpointAuthentication["authentication_type"].(string),
		ClientCertificate:  flattenedEndpointAuthentication["client_certificate"].(string),
	}
}

func flattenEndpointAuthentication(endpointAuthentication octopusdeploy.EndpointAuthentication) []interface{} {
	return []interface{}{map[string]interface{}{
		"account_id":          endpointAuthentication.AccountID,
		"authentication_type": endpointAuthentication.AuthenticationType,
		"client_certificate":  endpointAuthentication.ClientCertificate,
	}}
}

func getEndpointAuthenticationSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"account_id": {
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
		"client_certificate": {
			Optional: true,
			Type:     schema.TypeString,
		},
	}
}
