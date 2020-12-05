package octopusdeploy

import (
	"context"

	"github.com/OctopusDeploy/go-octopusdeploy/octopusdeploy"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func expandAzureServiceFabricClusterDeploymentTarget(d *schema.ResourceData) *octopusdeploy.DeploymentTarget {
	endpoint := octopusdeploy.NewAzureServiceFabricEndpoint()

	if v, ok := d.GetOk("aad_client_credential_secret"); ok {
		endpoint.AadClientCredentialSecret = v.(string)
	}

	if v, ok := d.GetOk("aad_credential_type"); ok {
		endpoint.AadCredentialType = v.(string)
	}

	if v, ok := d.GetOk("aad_user_credential_password"); ok {
		endpoint.AadUserCredentialPassword = octopusdeploy.NewSensitiveValue(v.(string))
	}

	if v, ok := d.GetOk("aad_user_credential_username"); ok {
		endpoint.AadUserCredentialUsername = v.(string)
	}

	if v, ok := d.GetOk("certificate_store_location"); ok {
		endpoint.CertificateStoreLocation = v.(string)
	}

	if v, ok := d.GetOk("certificate_store_name"); ok {
		endpoint.CertificateStoreName = v.(string)
	}

	if v, ok := d.GetOk("client_certificate_variable"); ok {
		endpoint.ClientCertificateVariable = v.(string)
	}

	if v, ok := d.GetOk("connection_endpoint"); ok {
		endpoint.ConnectionEndpoint = v.(string)
	}

	if v, ok := d.GetOk("security_mode"); ok {
		endpoint.SecurityMode = v.(string)
	}

	if v, ok := d.GetOk("server_certificate_thumbprint"); ok {
		endpoint.ServerCertificateThumbprint = v.(string)
	}

	deploymentTarget := expandDeploymentTarget(d)
	deploymentTarget.Endpoint = endpoint
	return deploymentTarget
}

func flattenAzureServiceFabricClusterDeploymentTarget(deploymentTarget *octopusdeploy.DeploymentTarget) map[string]interface{} {
	if deploymentTarget == nil {
		return nil
	}

	flattenedDeploymentTarget := flattenDeploymentTarget(deploymentTarget)
	endpointResource, _ := octopusdeploy.ToEndpointResource(deploymentTarget.Endpoint)
	flattenedDeploymentTarget["aad_client_credential_secret"] = endpointResource.AadClientCredentialSecret
	flattenedDeploymentTarget["aad_credential_type"] = endpointResource.AadCredentialType
	flattenedDeploymentTarget["aad_user_credential_username"] = endpointResource.AadUserCredentialUsername
	flattenedDeploymentTarget["certificate_store_location"] = endpointResource.CertificateStoreLocation
	flattenedDeploymentTarget["certificate_store_name"] = endpointResource.CertificateStoreName
	flattenedDeploymentTarget["client_certificate_variable"] = endpointResource.ClientCertificateVariable
	flattenedDeploymentTarget["connection_endpoint"] = endpointResource.ConnectionEndpoint
	flattenedDeploymentTarget["security_mode"] = endpointResource.SecurityMode
	flattenedDeploymentTarget["server_certificate_thumbprint"] = endpointResource.ServerCertificateThumbprint
	return flattenedDeploymentTarget
}

func getAzureServiceFabricClusterDeploymentTargetDataSchema() map[string]*schema.Schema {
	dataSchema := getAzureServiceFabricClusterDeploymentTargetSchema()
	setDataSchema(&dataSchema)

	deploymentTargetDataSchema := getDeploymentTargetDataSchema()

	deploymentTargetDataSchema["azure_service_fabric_cluster_deployment_target"] = &schema.Schema{
		Computed:    true,
		Description: "A list of Azure service fabric cluster deployment targets that match the filter(s).",
		Elem:        &schema.Resource{Schema: dataSchema},
		Optional:    true,
		Type:        schema.TypeList,
	}

	delete(deploymentTargetDataSchema, "communication_styles")
	delete(deploymentTargetDataSchema, "deployment_targets")
	deploymentTargetDataSchema["id"] = getDataSchemaID()

	return deploymentTargetDataSchema
}

func getAzureServiceFabricClusterDeploymentTargetSchema() map[string]*schema.Schema {
	azureServiceFabricClusterDeploymentTargetSchema := getDeploymentTargetSchema()

	azureServiceFabricClusterDeploymentTargetSchema["aad_client_credential_secret"] = &schema.Schema{
		Computed: true,
		Optional: true,
		Type:     schema.TypeString,
	}

	azureServiceFabricClusterDeploymentTargetSchema["aad_credential_type"] = &schema.Schema{
		Computed: true,
		Optional: true,
		Type:     schema.TypeString,
	}

	azureServiceFabricClusterDeploymentTargetSchema["aad_user_credential_password"] = &schema.Schema{
		Optional:  true,
		Sensitive: true,
		Type:      schema.TypeString,
	}

	azureServiceFabricClusterDeploymentTargetSchema["aad_user_credential_username"] = &schema.Schema{
		Computed: true,
		Optional: true,
		Type:     schema.TypeString,
	}

	azureServiceFabricClusterDeploymentTargetSchema["certificate_store_location"] = &schema.Schema{
		Computed: true,
		Optional: true,
		Type:     schema.TypeString,
	}

	azureServiceFabricClusterDeploymentTargetSchema["certificate_store_name"] = &schema.Schema{
		Computed: true,
		Optional: true,
		Type:     schema.TypeString,
	}

	azureServiceFabricClusterDeploymentTargetSchema["client_certificate_variable"] = &schema.Schema{
		Computed: true,
		Optional: true,
		Type:     schema.TypeString,
	}

	azureServiceFabricClusterDeploymentTargetSchema["connection_endpoint"] = &schema.Schema{
		Required: true,
		Type:     schema.TypeString,
	}

	azureServiceFabricClusterDeploymentTargetSchema["security_mode"] = &schema.Schema{
		Computed: true,
		Optional: true,
		Type:     schema.TypeString,
	}

	azureServiceFabricClusterDeploymentTargetSchema["server_certificate_thumbprint"] = &schema.Schema{
		Computed: true,
		Optional: true,
		Type:     schema.TypeString,
	}

	return azureServiceFabricClusterDeploymentTargetSchema
}

func setAzureServiceFabricClusterDeploymentTarget(ctx context.Context, d *schema.ResourceData, deploymentTarget *octopusdeploy.DeploymentTarget) {
	if deploymentTarget == nil {
		return
	}

	endpointResource, err := octopusdeploy.ToEndpointResource(deploymentTarget.Endpoint)
	if err != nil {
		return
	}

	d.Set("aad_client_credential_secret", endpointResource.AadClientCredentialSecret)
	d.Set("aad_credential_type", endpointResource.AadCredentialType)
	d.Set("aad_user_credential_username", endpointResource.AadUserCredentialUsername)
	d.Set("certificate_store_location", endpointResource.CertificateStoreLocation)
	d.Set("certificate_store_name", endpointResource.CertificateStoreName)
	d.Set("client_certificate_variable", endpointResource.ClientCertificateVariable)
	d.Set("connection_endpoint", endpointResource.ConnectionEndpoint)
	d.Set("security_mode", endpointResource.SecurityMode)
	d.Set("server_certificate_thumbprint", endpointResource.ServerCertificateThumbprint)

	setDeploymentTarget(ctx, d, deploymentTarget)
}
