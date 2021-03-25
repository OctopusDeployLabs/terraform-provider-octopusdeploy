package octopusdeploy

import (
	"github.com/OctopusDeploy/go-octopusdeploy/octopusdeploy"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func expandKubernetesAuthentication(values interface{}) octopusdeploy.IKubernetesAuthentication {
	flattenedValues := values.(*schema.Set)
	flattenedMap := flattenedValues.List()[0].(map[string]interface{})

	authenticationType := flattenedMap["authentication_type"].(string)
	switch authenticationType {
	case "KubernetesAws":
		return expandKubernetesAwsAuthentication(flattenedMap)
	case "KubernetesAzure":
		return expandKubernetesAzureAuthentication(flattenedMap)
	case "KubernetesCertificate":
		return expandKubernetesCertificateAuthentication(flattenedMap)
	case "KubernetesStandard":
		return expandKubernetesStandardAuthentication(flattenedMap)
	case "None":
		return expandKubernetesStandardAuthentication(flattenedMap)
	}

	return &octopusdeploy.KubernetesAuthentication{
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
		UseInstanceRole:           flattenedMap["use_instance_role"].(bool),
	}
}

func flattenKubernetesAuthentication(kubernetesAuthentication octopusdeploy.IKubernetesAuthentication) []interface{} {
	if kubernetesAuthentication == nil {
		return nil
	}

	switch kubernetesAuthentication.GetAuthenticationType() {
	case "KubernetesAws":
		return flattenKubernetesAwsAuthentication(kubernetesAuthentication.(*octopusdeploy.KubernetesAwsAuthentication))
	case "KubernetesAzure":
		return flattenKubernetesAzureAuthentication(kubernetesAuthentication.(*octopusdeploy.KubernetesAzureAuthentication))
	case "KubernetesCertificate":
		return flattenKubernetesCertificateAuthentication(kubernetesAuthentication.(*octopusdeploy.KubernetesCertificateAuthentication))
	case "KubernetesStandard":
		return flattenKubernetesStandardAuthentication(kubernetesAuthentication.(*octopusdeploy.KubernetesStandardAuthentication))
	case "None":
		return flattenKubernetesStandardAuthentication(kubernetesAuthentication.(*octopusdeploy.KubernetesStandardAuthentication))
	}

	return nil

	// authentication := kubernetesAuthentication.(*octopusdeploy.KubernetesAuthentication)

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
				"KubernetesStandard",
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
		"use_instance_role": {
			Optional: true,
			Type:     schema.TypeBool,
		},
	}
}
