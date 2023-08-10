package octopusdeploy

import (
	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/machines"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func expandKubernetesAuthentication(values interface{}) machines.IKubernetesAuthentication {
	if values == nil {
		return nil
	}

	flattenedValues := values.(*schema.Set)
	if len(flattenedValues.List()) == 0 {
		return nil
	}

	flattenedMap := flattenedValues.List()[0].(map[string]interface{})

	authenticationType := flattenedMap["authentication_type"].(string)
	switch authenticationType {
	case "KubernetesAws":
		return expandKubernetesAwsAuthentication(flattenedMap)
	case "KubernetesAzure":
		return expandKubernetesAzureAuthentication(flattenedMap)
	case "KubernetesCertificate":
		return expandKubernetesCertificateAuthentication(flattenedMap)
	case "KubernetesGoogleCloud":
		return expandKubernetesGcpAuthentication(flattenedMap)
	case "KubernetesStandard":
		return expandKubernetesStandardAuthentication(flattenedMap)
	case "KubernetesPodService":
		return expandKubernetesPodAuthentication(flattenedMap)
	case "None":
		return expandKubernetesStandardAuthentication(flattenedMap)
	}

	return &machines.KubernetesAuthentication{
		AccountID:                 flattenedMap["account_id"].(string),
		AdminLogin:                flattenedMap["admin_login"].(string),
		AssumedRoleARN:            flattenedMap["assumed_role_arn"].(string),
		AssumedRoleSession:        flattenedMap["assumed_role_session"].(string),
		AssumeRole:                flattenedMap["assume_role"].(bool),
		AssumeRoleExternalID:      flattenedMap["assume_role_external_id"].(string),
		AssumeRoleSessionDuration: flattenedMap["assume_role_session_duration"].(int),
		AuthenticationType:        flattenedMap["authentication_type"].(string),
		ClientCertificate:         flattenedMap["client_certificate"].(string),
		ClusterName:               flattenedMap["cluster_name"].(string),
		ClusterResourceGroup:      flattenedMap["cluster_resource_group"].(string),
		ImpersonateServiceAccount: flattenedMap["impersonate_service_account"].(bool),
		Project:                   flattenedMap["project"].(string),
		Region:                    flattenedMap["region"].(string),
		ServiceAccountEmails:      flattenedMap["service_account_emails"].(string),
		UseVmServiceAccount:       flattenedMap["use_vm_service_account"].(bool),
		UseInstanceRole:           flattenedMap["use_instance_role"].(bool),
		Zone:                      flattenedMap["zone"].(string),
	}
}

func flattenKubernetesAuthentication(kubernetesAuthentication machines.IKubernetesAuthentication) []interface{} {
	if kubernetesAuthentication == nil {
		return nil
	}

	switch kubernetesAuthentication.GetAuthenticationType() {
	case "KubernetesAws":
		return flattenKubernetesAwsAuthentication(kubernetesAuthentication.(*machines.KubernetesAwsAuthentication))
	case "KubernetesAzure":
		return flattenKubernetesAzureAuthentication(kubernetesAuthentication.(*machines.KubernetesAzureAuthentication))
	case "KubernetesCertificate":
		return flattenKubernetesCertificateAuthentication(kubernetesAuthentication.(*machines.KubernetesCertificateAuthentication))
	case "KubernetesGoogleCloud":
		return flattenKubernetesGcpAuthentication(kubernetesAuthentication.(*machines.KubernetesGcpAuthentication))
	case "KubernetesStandard":
		return flattenKubernetesStandardAuthentication(kubernetesAuthentication.(*machines.KubernetesStandardAuthentication))
	case "KubernetesPodService":
		return flattenKubernetesPodAuthentication(kubernetesAuthentication.(*machines.KubernetesPodAuthentication))
	case "None":
		return flattenKubernetesStandardAuthentication(kubernetesAuthentication.(*machines.KubernetesStandardAuthentication))
	}

	return nil

	// authentication := kubernetesAuthentication.(*machines.KubernetesAuthentication)

	// return []interface{}{map[string]interface{}{
	// 	"account_id":                   authentication.AccountID,
	// 	"admin_login":                  authentication.AdminLogin,
	// 	"assumed_role_arn":             authentication.AssumedRoleARN,
	// 	"assumed_role_session":         authentication.AssumedRoleSession,
	// 	"assume_role":                  authentication.AssumeRole,
	// 	"assume_role_session_duration": authentication.AssumeRoleSessionDuration,
	// 	"assume_role_external_id":      authentication.AssumeRoleExternalID,
	// 	"authentication_type":          authentication.AuthenticationType,
	// 	"client_certificate":           authentication.ClientCertificate,
	// 	"cluster_name":                 authentication.ClusterName,
	// 	"cluster_resource_group":       authentication.ClusterResourceGroup,
	// 	"use_instance_role":            authentication.UseInstanceRole,
	// }}
}

func getKubernetesAuthenticationSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"account_id": {
			Optional: true,
			Type:     schema.TypeString,
		},
		"admin_login": {
			Optional: true,
			Type:     schema.TypeString,
		},
		"assumed_role_arn": {
			Optional: true,
			Type:     schema.TypeString,
		},
		"assumed_role_session": {
			Optional: true,
			Type:     schema.TypeString,
		},
		"assume_role": {
			Optional: true,
			Type:     schema.TypeBool,
		},
		"assume_role_session_duration": {
			Optional: true,
			Type:     schema.TypeInt,
		},
		"assume_role_external_id": {
			Optional: true,
			Type:     schema.TypeString,
		},
		"authentication_type": {
			Optional: true,
			Type:     schema.TypeString,
			ValidateDiagFunc: validation.ToDiagFunc(validation.StringInSlice([]string{
				"KubernetesAws",
				"KubernetesAzure",
				"KubernetesCertificate",
				"KubernetesGoogleCloud",
				"KubernetesStandard",
				"KubernetesPodService",
				"None",
			}, false)),
		},
		"client_certificate": {
			Optional: true,
			Type:     schema.TypeString,
		},
		"cluster_name": {
			Optional: true,
			Type:     schema.TypeString,
		},
		"cluster_resource_group": {
			Optional: true,
			Type:     schema.TypeString,
		},
		"impersonate_service_account": {
			Optional: true,
			Type:     schema.TypeBool,
		},
		"project": {
			Optional: true,
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
		"use_instance_role": {
			Optional: true,
			Type:     schema.TypeBool,
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
