package octopusdeploy

import (
	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/machines"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func expandKubernetesAzureAuthentication(values interface{}) *machines.KubernetesAzureAuthentication {
	flattenedValues := values.([]interface{})
	flattenedAuthentication := flattenedValues[0].(map[string]interface{})

	authentication := machines.NewKubernetesAzureAuthentication()
	authentication.AccountID = flattenedAuthentication["account_id"].(string)
	authentication.AuthenticationType = "KubernetesAzure"
	authentication.ClusterName = flattenedAuthentication["cluster_name"].(string)
	authentication.ClusterResourceGroup = flattenedAuthentication["cluster_resource_group"].(string)
	return authentication
}

func flattenKubernetesAzureAuthentication(kubernetesAzureAuthentication *machines.KubernetesAzureAuthentication) []interface{} {
	if kubernetesAzureAuthentication == nil {
		return nil
	}

	return []interface{}{map[string]interface{}{
		"account_id":             kubernetesAzureAuthentication.AccountID,
		"cluster_name":           kubernetesAzureAuthentication.ClusterName,
		"cluster_resource_group": kubernetesAzureAuthentication.ClusterResourceGroup,
	}}
}

func getKubernetesAzureAuthenticationSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"account_id": {
			Required: true,
			Type:     schema.TypeString,
		},
		"cluster_name": {
			Required: true,
			Type:     schema.TypeString,
		},
		"cluster_resource_group": {
			Required: true,
			Type:     schema.TypeString,
		},
	}
}
