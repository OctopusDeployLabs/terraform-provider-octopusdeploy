package octopusdeploy

import (
	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/machines"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func expandKubernetesCertificateAuthentication(values interface{}) *machines.KubernetesCertificateAuthentication {
	flattenedValues := values.([]interface{})
	flattenedAuthentication := flattenedValues[0].(map[string]interface{})

	authentication := &machines.KubernetesCertificateAuthentication{
		ClientCertificate: flattenedAuthentication["client_certificate"].(string),
	}

	authentication.AuthenticationType = "KubernetesCertificate"

	return authentication
}

func flattenKubernetesCertificateAuthentication(kubernetesCertificateAuthentication *machines.KubernetesCertificateAuthentication) []interface{} {
	if kubernetesCertificateAuthentication == nil {
		return nil
	}

	return []interface{}{map[string]interface{}{
		"client_certificate": kubernetesCertificateAuthentication.ClientCertificate,
	}}
}

func getKubernetesCertificateAuthenticationSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"client_certificate": {
			Optional: true,
			Type:     schema.TypeString,
		},
	}
}
