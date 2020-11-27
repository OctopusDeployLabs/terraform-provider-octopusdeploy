package octopusdeploy

import (
	"github.com/OctopusDeploy/go-octopusdeploy/octopusdeploy"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func expandKubernetesAwsAuthentication(values interface{}) *octopusdeploy.KubernetesAwsAuthentication {
	flattenedValues := values.([]interface{})
	flattenedAuthentication := flattenedValues[0].(map[string]interface{})

	authentication := &octopusdeploy.KubernetesAwsAuthentication{
		AssumedRoleARN:            flattenedAuthentication["assumed_role_arn"].(string),
		AssumedRoleSession:        flattenedAuthentication["assumed_role_session"].(string),
		AssumeRole:                flattenedAuthentication["assume_role"].(bool),
		AssumeRoleExternalID:      flattenedAuthentication["assume_role_external_id"].(string),
		AssumeRoleSessionDuration: flattenedAuthentication["assume_role_session_duration"].(int),
		ClusterName:               flattenedAuthentication["cluster_name"].(string),
		UseInstanceRole:           flattenedAuthentication["use_instance_role"].(bool),
	}

	authentication.AccountID = flattenedAuthentication["account_id"].(string)
	authentication.AuthenticationType = "KubernetesAws"

	return authentication
}

func flattenKubernetesAwsAuthentication(kubernetesAwsAuthentication *octopusdeploy.KubernetesAwsAuthentication) []interface{} {
	if kubernetesAwsAuthentication == nil {
		return nil
	}

	return []interface{}{map[string]interface{}{
		"account_id":                   kubernetesAwsAuthentication.AccountID,
		"assumed_role_arn":             kubernetesAwsAuthentication.AssumedRoleARN,
		"assumed_role_session":         kubernetesAwsAuthentication.AssumedRoleSession,
		"assume_role":                  kubernetesAwsAuthentication.AssumeRole,
		"assume_role_external_id":      kubernetesAwsAuthentication.AssumeRoleExternalID,
		"assume_role_session_duration": kubernetesAwsAuthentication.AssumeRoleSessionDuration,
		"cluster_name":                 kubernetesAwsAuthentication.ClusterName,
		"use_instance_role":            kubernetesAwsAuthentication.UseInstanceRole,
	}}
}

func getKubernetesAwsAuthenticationSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"account_id": {
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
		"assume_role_external_id": {
			Optional: true,
			Type:     schema.TypeString,
		},
		"assume_role_session_duration": {
			Optional: true,
			Type:     schema.TypeInt,
		},
		"cluster_name": {
			Optional: true,
			Type:     schema.TypeString,
		},
		"use_instance_role": {
			Optional: true,
			Type:     schema.TypeBool,
		},
	}
}
