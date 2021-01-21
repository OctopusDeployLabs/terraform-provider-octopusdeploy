package octopusdeploy

import (
	"net/url"

	"github.com/OctopusDeploy/go-octopusdeploy/octopusdeploy"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func expandKubernetesCluster(flattenedMap map[string]interface{}) *octopusdeploy.KubernetesEndpoint {
	clusterURL, _ := url.Parse(flattenedMap["cluster_url"].(string))

	endpoint := octopusdeploy.NewKubernetesEndpoint(clusterURL)
	endpoint.Authentication = expandKubernetesAuthentication(flattenedMap["authentication"])
	endpoint.ClusterCertificate = flattenedMap["cluster_certificate"].(string)
	endpoint.Container = expandDeploymentActionContainer(flattenedMap["container"])
	endpoint.DefaultWorkerPoolID = flattenedMap["default_worker_pool_id"].(string)
	endpoint.ID = flattenedMap["id"].(string)
	endpoint.Namespace = flattenedMap["namespace"].(string)
	endpoint.ProxyID = flattenedMap["proxy_id"].(string)
	endpoint.RunningInContainer = flattenedMap["running_in_container"].(bool)
	endpoint.SkipTLSVerification = flattenedMap["skip_tls_verification"].(bool)

	// tfSchemaSetInterface, ok := d.GetOk("endpoint")
	// if !ok {
	// 	return nil
	// }
	// tfSchemaSet := tfSchemaSetInterface.([]interface{})
	// if len(tfSchemaSet) == 0 {
	// 	return nil
	// }
	// // Get the first element in the list, which is a map of the interfaces
	// tfSchemaList := tfSchemaSet[0].(map[string]interface{})

	// authenticationType := octopusdeploy.CommunicationStyle(tfSchemaList["authentication_type"].(string))

	// var kubernetesAuthentication octopusdeploy.IKubernetesAuthentication
	// switch authenticationType {
	// case "KubernetesAws":
	// 	kubernetesAuthentication = expandKubernetesAwsAuthentication(d)
	// case "KubernetesAzure":
	// 	kubernetesAuthentication = expandKubernetesAzureAuthentication(d)
	// case "KubernetesCertificate":
	// 	kubernetesAuthentication = expandKubernetesCertificateAuthentication(d)
	// case "KubernetesStandard":
	// 	kubernetesAuthentication = expandKubernetesStandardAuthentication(d)
	// case "None":
	// 	kubernetesAuthentication = expandKubernetesStandardAuthentication(d)
	// }

	// endpoint.Authentication = kubernetesAuthentication

	// if v, ok := d.GetOk("aws_account_authentication"); ok {
	// 	endpoint.Authentication = expandKubernetesAwsAuthentication(v)
	// }

	// if v, ok := d.GetOk("azure_service_principal_authentication"); ok {
	// 	endpoint.Authentication = expandKubernetesAzureAuthentication(v)
	// }

	// if v, ok := d.GetOk("certificate_authentication"); ok {
	// 	endpoint.Authentication = expandKubernetesCertificateAuthentication(v)
	// }

	return endpoint
}

func flattenKubernetesCluster(endpoint *octopusdeploy.KubernetesEndpoint) []interface{} {
	if endpoint == nil {
		return nil
	}

	flattenedEndpoint := map[string]interface{}{
		"cluster_certificate":    endpoint.ClusterCertificate,
		"container":              flattenDeploymentActionContainer(endpoint.Container),
		"default_worker_pool_id": endpoint.DefaultWorkerPoolID,
		"id":                     endpoint.GetID(),
		"namespace":              endpoint.Namespace,
		"proxy_id":               endpoint.ProxyID,
		"running_in_container":   endpoint.RunningInContainer,
		"skip_tls_verification":  endpoint.SkipTLSVerification,
	}

	if endpoint.ClusterURL != nil {
		flattenedEndpoint["cluster_url"] = endpoint.ClusterURL.String()
	}

	switch endpoint.Authentication.GetAuthenticationType() {
	case "KubernetesAws":
		flattenedEndpoint["aws_account_authentication"] = flattenKubernetesAwsAuthentication(endpoint.Authentication.(*octopusdeploy.KubernetesAwsAuthentication))
	case "KubernetesAzure":
		flattenedEndpoint["azure_service_principal_authentication"] = flattenKubernetesAzureAuthentication(endpoint.Authentication.(*octopusdeploy.KubernetesAzureAuthentication))
	case "KubernetesCertificate":
		flattenedEndpoint["certificate_authentication"] = flattenKubernetesCertificateAuthentication(endpoint.Authentication.(*octopusdeploy.KubernetesCertificateAuthentication))
	case "KubernetesStandard":
		flattenedEndpoint["authentication"] = flattenKubernetesStandardAuthentication(endpoint.Authentication.(*octopusdeploy.KubernetesStandardAuthentication))
	case "None":
		flattenedEndpoint["authentication"] = flattenKubernetesStandardAuthentication(endpoint.Authentication.(*octopusdeploy.KubernetesStandardAuthentication))
	}

	return []interface{}{flattenedEndpoint}
}

func getKubernetesClusterSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"authentication": {
			Computed:     true,
			Elem:         &schema.Resource{Schema: getKubernetesAuthenticationSchema()},
			ExactlyOneOf: []string{"authentication", "aws_account_authentication", "azure_service_principal_authentication", "certificate_authentication"},
			MaxItems:     1,
			MinItems:     0,
			Optional:     true,
			Type:         schema.TypeList,
		},
		"aws_account_authentication": {
			Computed:     true,
			Elem:         &schema.Resource{Schema: getKubernetesAwsAuthenticationSchema()},
			ExactlyOneOf: []string{"authentication", "aws_account_authentication", "azure_service_principal_authentication", "certificate_authentication"},
			MaxItems:     1,
			MinItems:     0,
			Optional:     true,
			Type:         schema.TypeList,
		},
		"azure_service_principal_authentication": {
			Computed:     true,
			Elem:         &schema.Resource{Schema: getKubernetesAzureAuthenticationSchema()},
			ExactlyOneOf: []string{"authentication", "aws_account_authentication", "azure_service_principal_authentication", "certificate_authentication"},
			MaxItems:     1,
			MinItems:     0,
			Optional:     true,
			Type:         schema.TypeList,
		},
		"certificate_authentication": {
			Computed:     true,
			Elem:         &schema.Resource{Schema: getKubernetesCertificateAuthenticationSchema()},
			ExactlyOneOf: []string{"authentication", "aws_account_authentication", "azure_service_principal_authentication", "certificate_authentication"},
			MaxItems:     1,
			MinItems:     0,
			Optional:     true,
			Type:         schema.TypeList,
		},
		"cluster_certificate": {
			Optional: true,
			Type:     schema.TypeString,
		},
		"cluster_url": {
			Required: true,
			Type:     schema.TypeString,
		},
		"container": {
			Computed: true,
			Elem:     &schema.Resource{Schema: getDeploymentActionContainerSchema()},
			Optional: true,
			Type:     schema.TypeList,
		},
		"default_worker_pool_id": {
			Optional: true,
			Type:     schema.TypeString,
		},
		"id": getIDSchema(),
		"namespace": {
			Optional: true,
			Type:     schema.TypeString,
		},
		"proxy_id": {
			Optional: true,
			Type:     schema.TypeString,
		},
		"running_in_container": {
			Optional: true,
			Type:     schema.TypeBool,
		},
		"skip_tls_verification": {
			Optional: true,
			Type:     schema.TypeBool,
		},
	}
}
