package octopusdeploy

import (
	"github.com/OctopusDeploy/go-octopusdeploy/octopusdeploy"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func expandKubernetesGcpAuthentication(values interface{}) *octopusdeploy.KubernetesGcpAuthentication {
	flattenedValues := values.([]interface{})
	flattenedAuthentication := flattenedValues[0].(map[string]interface{})

	authentication := octopusdeploy.NewKubernetesGcpAuthentication()
	authentication.ClusterName = flattenedAuthentication["cluster_name"].(string)
	authentication.ImpersonateServiceAccount = flattenedAuthentication["impersonate_service_account"].(bool)
	authentication.Project = flattenedAuthentication["project"].(string)
	authentication.Region = flattenedAuthentication["region"].(string)
	authentication.ServiceAccountEmails = flattenedAuthentication["service_account_emails"].(string)
	authentication.UseVmServiceAccount = flattenedAuthentication["use_vm_service_account"].(bool)
	authentication.Zone = flattenedAuthentication["zone"].(string)
	return authentication
}

func flattenKubernetesGcpAuthentication(kubernetesGcpAuthentication *octopusdeploy.KubernetesGcpAuthentication) []interface{} {
	if kubernetesGcpAuthentication == nil {
		return nil
	}

	return []interface{}{map[string]interface{}{
		"cluster_name":                kubernetesGcpAuthentication.ClusterName,
		"impersonate_service_account": kubernetesGcpAuthentication.ImpersonateServiceAccount,
		"project":                     kubernetesGcpAuthentication.Project,
		"region":                      kubernetesGcpAuthentication.Region,
		"service_account_emails":      kubernetesGcpAuthentication.ServiceAccountEmails,
		"use_vm_service_account":      kubernetesGcpAuthentication.UseVmServiceAccount,
		"zone":                        kubernetesGcpAuthentication.Zone,
	}}
}

func getKubernetesGcpAuthenticationSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"cluster_name": {
			Required: true,
			Type:     schema.TypeString,
		},
		"impersonate_service_account": {
			Optional: true,
			Type:     schema.TypeBool,
		},
		"project": {
			Required: true,
			Type:     schema.TypeString,
		},
		"region": {
			Optional: true,
			Type:     schema.TypeString,
		},
		"service_account_emails": {
			Optional: true,
			Type:     schema.TypeString,
		},
		"use_vm_service_account": {
			Optional: true,
			Type:     schema.TypeBool,
		},
		"zone": {
			Optional: true,
			Type:     schema.TypeString,
		},
	}
}
