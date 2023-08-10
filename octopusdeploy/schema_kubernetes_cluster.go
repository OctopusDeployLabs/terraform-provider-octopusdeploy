package octopusdeploy

import (
	"net/url"

	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/machines"
)

func expandKubernetesCluster(flattenedMap map[string]interface{}) *machines.KubernetesEndpoint {
	clusterURL, _ := url.Parse(flattenedMap["cluster_url"].(string))

	endpoint := machines.NewKubernetesEndpoint(clusterURL)
	endpoint.Authentication = expandKubernetesAuthentication(flattenedMap["authentication"])
	endpoint.ClusterCertificate = flattenedMap["cluster_certificate"].(string)
	endpoint.ClusterCertificatePath = flattenedMap["cluster_certificate_path"].(string)
	endpoint.Container = expandContainer(flattenedMap["container"])
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

	// authenticationType := machines.CommunicationStyle(tfSchemaList["authentication_type"].(string))

	// var kubernetesAuthentication machines.IKubernetesAuthentication
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
