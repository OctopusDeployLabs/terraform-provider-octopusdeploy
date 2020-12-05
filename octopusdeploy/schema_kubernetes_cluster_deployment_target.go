package octopusdeploy

import (
	"context"
	"net/url"

	"github.com/OctopusDeploy/go-octopusdeploy/octopusdeploy"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func expandKubernetesClusterDeploymentTarget(d *schema.ResourceData) *octopusdeploy.DeploymentTarget {
	clusterURL, _ := url.Parse(d.Get("cluster_url").(string))

	endpoint := octopusdeploy.NewKubernetesEndpoint(clusterURL)

	if v, ok := d.GetOk("aws_account_authentication"); ok {
		endpoint.Authentication = expandKubernetesAwsAuthentication(v)
	}

	if v, ok := d.GetOk("azure_service_principal_authentication"); ok {
		endpoint.Authentication = expandKubernetesAzureAuthentication(v)
	}

	if v, ok := d.GetOk("cluster_certificate"); ok {
		endpoint.ClusterCertificate = v.(string)
	}

	if v, ok := d.GetOk("container"); ok {
		endpoint.Container = expandDeploymentActionContainer(v)
	}

	if v, ok := d.GetOk("default_worker_pool_id"); ok {
		endpoint.DefaultWorkerPoolID = v.(string)
	}

	if v, ok := d.GetOk("namespace"); ok {
		endpoint.Namespace = v.(string)
	}

	if v, ok := d.GetOk("proxy_id"); ok {
		endpoint.ProxyID = v.(string)
	}

	if v, ok := d.GetOk("running_in_container"); ok {
		endpoint.RunningInContainer = v.(bool)
	}

	if v, ok := d.GetOk("skip_tls_verification"); ok {
		endpoint.SkipTLSVerification = v.(bool)
	}

	deploymentTarget := expandDeploymentTarget(d)
	deploymentTarget.Endpoint = endpoint
	return deploymentTarget
}

func flattenKubernetesClusterDeploymentTarget(deploymentTarget *octopusdeploy.DeploymentTarget) map[string]interface{} {
	if deploymentTarget == nil {
		return nil
	}

	flattenedDeploymentTarget := flattenDeploymentTarget(deploymentTarget)
	endpointResource, _ := octopusdeploy.ToEndpointResource(deploymentTarget.Endpoint)

	flattenedDeploymentTarget["cluster_certificate"] = endpointResource.ClusterCertificate
	flattenedDeploymentTarget["container"] = flattenDeploymentActionContainer(endpointResource.Container)
	flattenedDeploymentTarget["default_worker_pool_id"] = endpointResource.DefaultWorkerPoolID
	flattenedDeploymentTarget["namespace"] = endpointResource.Namespace
	flattenedDeploymentTarget["proxy_id"] = endpointResource.ProxyID
	flattenedDeploymentTarget["running_in_container"] = endpointResource.RunningInContainer
	flattenedDeploymentTarget["skip_tls_verification"] = endpointResource.SkipTLSVerification

	if endpointResource.ClusterURL != nil {
		flattenedDeploymentTarget["cluster_url"] = endpointResource.ClusterURL.String()
	}

	switch endpointResource.Authentication.GetAuthenticationType() {
	case "KubernetesAws":
		flattenedDeploymentTarget["aws_account_authentication"] = flattenKubernetesAwsAuthentication(endpointResource.Authentication.(*octopusdeploy.KubernetesAwsAuthentication))
	case "KubernetesAzure":
		flattenedDeploymentTarget["azure_service_principal_authentication"] = flattenKubernetesAzureAuthentication(endpointResource.Authentication.(*octopusdeploy.KubernetesAzureAuthentication))
	case "KubernetesCertificate":
		flattenedDeploymentTarget["certificate_authentication"] = flattenKubernetesCertificateAuthentication(endpointResource.Authentication.(*octopusdeploy.KubernetesCertificateAuthentication))
	case "KubernetesStandard":
		flattenedDeploymentTarget["authentication"] = flattenKubernetesStandardAuthentication(endpointResource.Authentication.(*octopusdeploy.KubernetesStandardAuthentication))
	case "None":
		flattenedDeploymentTarget["authentication"] = flattenKubernetesStandardAuthentication(endpointResource.Authentication.(*octopusdeploy.KubernetesStandardAuthentication))
	}

	return flattenedDeploymentTarget
}

func getKubernetesClusterDeploymentTargetDataSchema() map[string]*schema.Schema {
	dataSchema := getKubernetesClusterDeploymentTargetSchema()
	setDataSchema(&dataSchema)

	deploymentTargetDataSchema := getDeploymentTargetDataSchema()

	deploymentTargetDataSchema["kubernetes_cluster_deployment_target"] = &schema.Schema{
		Computed:    true,
		Description: "A list of Kubernetes cluster deployment targets that match the filter(s).",
		Elem:        &schema.Resource{Schema: dataSchema},
		Optional:    true,
		Type:        schema.TypeList,
	}

	delete(deploymentTargetDataSchema, "communication_styles")
	delete(deploymentTargetDataSchema, "deployment_targets")
	deploymentTargetDataSchema["id"] = getDataSchemaID()

	return deploymentTargetDataSchema
}

func getKubernetesClusterDeploymentTargetSchema() map[string]*schema.Schema {
	kubernetesClusterDeploymentTargetSchema := getDeploymentTargetSchema()

	kubernetesClusterDeploymentTargetSchema["aws_account_authentication"] = &schema.Schema{
		Computed:     true,
		Elem:         &schema.Resource{Schema: getKubernetesAwsAuthenticationSchema()},
		ExactlyOneOf: []string{"aws_account_authentication", "azure_service_principal_authentication", "certificate_authentication", "token_authentication", "username_password_authentication"},
		MaxItems:     1,
		MinItems:     0,
		Optional:     true,
		Type:         schema.TypeList,
	}

	kubernetesClusterDeploymentTargetSchema["azure_service_principal_authentication"] = &schema.Schema{
		Computed:     true,
		Elem:         &schema.Resource{Schema: getKubernetesAzureAuthenticationSchema()},
		ExactlyOneOf: []string{"aws_account_authentication", "azure_service_principal_authentication", "certificate_authentication", "token_authentication", "username_password_authentication"},
		MaxItems:     1,
		MinItems:     0,
		Optional:     true,
		Type:         schema.TypeList,
	}

	kubernetesClusterDeploymentTargetSchema["certificate_authentication"] = &schema.Schema{
		Computed:     true,
		Elem:         &schema.Resource{Schema: getKubernetesCertificateAuthenticationSchema()},
		ExactlyOneOf: []string{"aws_account_authentication", "azure_service_principal_authentication", "certificate_authentication", "token_authentication", "username_password_authentication"},
		MaxItems:     1,
		MinItems:     0,
		Optional:     true,
		Type:         schema.TypeSet,
	}

	kubernetesClusterDeploymentTargetSchema["cluster_certificate"] = &schema.Schema{
		Optional: true,
		Type:     schema.TypeString,
	}

	kubernetesClusterDeploymentTargetSchema["cluster_url"] = &schema.Schema{
		Required: true,
		Type:     schema.TypeString,
	}

	kubernetesClusterDeploymentTargetSchema["container"] = &schema.Schema{
		Computed: true,
		Elem:     &schema.Resource{Schema: getDeploymentActionContainerSchema()},
		Optional: true,
		Type:     schema.TypeList,
	}

	kubernetesClusterDeploymentTargetSchema["default_worker_pool_id"] = &schema.Schema{
		Optional: true,
		Type:     schema.TypeString,
	}

	kubernetesClusterDeploymentTargetSchema["namespace"] = &schema.Schema{
		Optional: true,
		Type:     schema.TypeString,
	}

	kubernetesClusterDeploymentTargetSchema["proxy_id"] = &schema.Schema{
		Optional: true,
		Type:     schema.TypeString,
	}

	kubernetesClusterDeploymentTargetSchema["running_in_container"] = &schema.Schema{
		Optional: true,
		Type:     schema.TypeBool,
	}

	kubernetesClusterDeploymentTargetSchema["skip_tls_verification"] = &schema.Schema{
		Optional: true,
		Type:     schema.TypeBool,
	}

	kubernetesClusterDeploymentTargetSchema["token_authentication"] = &schema.Schema{
		Computed:     true,
		Elem:         &schema.Resource{Schema: getKubernetesStandardAuthenticationSchema()},
		ExactlyOneOf: []string{"aws_account_authentication", "azure_service_principal_authentication", "certificate_authentication", "token_authentication", "username_password_authentication"},
		MaxItems:     1,
		MinItems:     0,
		Optional:     true,
		Type:         schema.TypeSet,
	}

	kubernetesClusterDeploymentTargetSchema["username_password_authentication"] = &schema.Schema{
		Computed:     true,
		Elem:         &schema.Resource{Schema: getKubernetesStandardAuthenticationSchema()},
		ExactlyOneOf: []string{"aws_account_authentication", "azure_service_principal_authentication", "certificate_authentication", "token_authentication", "username_password_authentication"},
		MaxItems:     1,
		MinItems:     0,
		Optional:     true,
		Type:         schema.TypeSet,
	}

	return kubernetesClusterDeploymentTargetSchema
}

func setKubernetesClusterDeploymentTarget(ctx context.Context, d *schema.ResourceData, deploymentTarget *octopusdeploy.DeploymentTarget) {
	if deploymentTarget == nil {
		return
	}

	endpointResource, err := octopusdeploy.ToEndpointResource(deploymentTarget.Endpoint)
	if err != nil {
		return
	}

	d.Set("cluster_certificate", endpointResource.ClusterCertificate)
	d.Set("container", flattenDeploymentActionContainer(endpointResource.Container))
	d.Set("default_worker_pool_id", endpointResource.DefaultWorkerPoolID)
	d.Set("namespace", endpointResource.Namespace)
	d.Set("proxy_id", endpointResource.ProxyID)
	d.Set("running_in_container", endpointResource.RunningInContainer)
	d.Set("skip_tls_verification", endpointResource.SkipTLSVerification)

	if endpointResource.ClusterURL != nil {
		d.Set("cluster_url", endpointResource.ClusterURL.String())
	}

	switch endpointResource.Authentication.GetAuthenticationType() {
	case "KubernetesAws":
		d.Set("aws_account_authentication", flattenKubernetesAwsAuthentication(endpointResource.Authentication.(*octopusdeploy.KubernetesAwsAuthentication)))
	case "KubernetesAzure":
		d.Set("azure_service_principal_authentication", flattenKubernetesAzureAuthentication(endpointResource.Authentication.(*octopusdeploy.KubernetesAzureAuthentication)))
	case "KubernetesCertificate":
		d.Set("certificate_authentication", flattenKubernetesCertificateAuthentication(endpointResource.Authentication.(*octopusdeploy.KubernetesCertificateAuthentication)))
	case "KubernetesStandard":
		d.Set("authentication", flattenKubernetesStandardAuthentication(endpointResource.Authentication.(*octopusdeploy.KubernetesStandardAuthentication)))
	case "None":
		d.Set("authentication", flattenKubernetesStandardAuthentication(endpointResource.Authentication.(*octopusdeploy.KubernetesStandardAuthentication)))
	}

	setDeploymentTarget(ctx, d, deploymentTarget)
}
