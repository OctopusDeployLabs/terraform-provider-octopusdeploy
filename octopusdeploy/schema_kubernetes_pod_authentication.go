package octopusdeploy

import (
	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/machines"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func expandKubernetesPodAuthentication(values interface{}) *machines.KubernetesPodAuthentication {
	flattenedValues := values.([]interface{})
	flattenedAuthentication := flattenedValues[0].(map[string]interface{})

	return &machines.KubernetesPodAuthentication{
		AuthenticationType: "KubernetesPodService",
		TokenPath:          flattenedAuthentication["token_path"].(string),
	}
}

func flattenKubernetesPodAuthentication(KubernetesPodAuthentication *machines.KubernetesPodAuthentication) []interface{} {
	if KubernetesPodAuthentication == nil {
		return nil
	}

	return []interface{}{map[string]interface{}{
		"token_path": KubernetesPodAuthentication.TokenPath,
	}}
}

func getKubernetesPodAuthenticationSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"token_path": {
			Required: true,
			Type:     schema.TypeString,
		},
	}
}
