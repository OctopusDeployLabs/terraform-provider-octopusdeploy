package octopusdeploy

import (
	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/machines"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func expandKubernetesStandardAuthentication(values interface{}) *machines.KubernetesStandardAuthentication {
	flattenedValues := values.([]interface{})
	flattenedAuthentication := flattenedValues[0].(map[string]interface{})

	authentication := &machines.KubernetesStandardAuthentication{
		AccountID: flattenedAuthentication["account_id"].(string),
	}

	authentication.AuthenticationType = "KubernetesStandard"

	return authentication
}

func flattenKubernetesStandardAuthentication(kubernetesStandardAuthentication *machines.KubernetesStandardAuthentication) []interface{} {
	if kubernetesStandardAuthentication == nil {
		return nil
	}

	return []interface{}{map[string]interface{}{
		"account_id": kubernetesStandardAuthentication.AccountID,
	}}
}

func getKubernetesStandardAuthenticationSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"account_id": {
			Optional: true,
			Type:     schema.TypeString,
		},
	}
}
