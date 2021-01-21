package octopusdeploy

import (
	"github.com/OctopusDeploy/go-octopusdeploy/octopusdeploy"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func expandKubernetesAwsAuthentication(values interface{}) *octopusdeploy.KubernetesAwsAuthentication {
	flattenedValues := values.([]interface{})
	flattenedAuthentication := flattenedValues[0].(map[string]interface{})

	authentication := octopusdeploy.NewKubernetesAwsAuthentication()
	authentication.AccountID = flattenedAuthentication["account_id"].(string)
	authentication.AssumedRoleARN = flattenedAuthentication["assumed_role_arn"].(string)
	authentication.AssumedRoleSession = flattenedAuthentication["assumed_role_session"].(string)
	authentication.AssumeRole = flattenedAuthentication["assume_role"].(bool)
	authentication.AssumeRoleExternalID = flattenedAuthentication["assume_role_external_id"].(string)
	authentication.AssumeRoleSessionDuration = flattenedAuthentication["assume_role_session_duration"].(int)
	authentication.AuthenticationType = "KubernetesAws"
	authentication.ClusterName = flattenedAuthentication["cluster_name"].(string)
	authentication.UseInstanceRole = flattenedAuthentication["use_instance_role"].(bool)
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
			Required: true,
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
			Required: true,
			Type:     schema.TypeString,
		},
		"use_instance_role": {
			Optional: true,
			Type:     schema.TypeBool,
		},
	}
}
