package octopusdeploy

import (
	"context"

	"github.com/OctopusDeploy/go-octopusdeploy/octopusdeploy"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func expandAzureCloudServiceDeploymentTarget(d *schema.ResourceData) *octopusdeploy.DeploymentTarget {
	endpoint := octopusdeploy.NewAzureCloudServiceEndpoint()

	if v, ok := d.GetOk("account_id"); ok {
		endpoint.AccountID = v.(string)
	}

	if v, ok := d.GetOk("cloud_service_name"); ok {
		endpoint.CloudServiceName = v.(string)
	}

	if v, ok := d.GetOk("default_worker_pool_id"); ok {
		endpoint.DefaultWorkerPoolID = v.(string)
	}

	if v, ok := d.GetOk("slot"); ok {
		endpoint.Slot = v.(string)
	}

	if v, ok := d.GetOk("storage_account_name"); ok {
		endpoint.StorageAccountName = v.(string)
	}

	if v, ok := d.GetOk("swap_if_possible"); ok {
		endpoint.SwapIfPossible = v.(bool)
	}

	if v, ok := d.GetOk("use_current_instance_count"); ok {
		endpoint.UseCurrentInstanceCount = v.(bool)
	}

	deploymentTarget := expandDeploymentTarget(d)
	deploymentTarget.Endpoint = endpoint
	return deploymentTarget
}

func flattenAzureCloudServiceDeploymentTarget(deploymentTarget *octopusdeploy.DeploymentTarget) map[string]interface{} {
	if deploymentTarget == nil {
		return nil
	}

	flattenedDeploymentTarget := flattenDeploymentTarget(deploymentTarget)
	endpointResource, _ := octopusdeploy.ToEndpointResource(deploymentTarget.Endpoint)
	flattenedDeploymentTarget["account_id"] = endpointResource.AccountID
	flattenedDeploymentTarget["cloud_service_name"] = endpointResource.CloudServiceName
	flattenedDeploymentTarget["default_worker_pool_id"] = endpointResource.DefaultWorkerPoolID
	flattenedDeploymentTarget["slot"] = endpointResource.Slot
	flattenedDeploymentTarget["storage_account_name"] = endpointResource.StorageAccountName
	flattenedDeploymentTarget["swap_if_possible"] = endpointResource.SwapIfPossible
	flattenedDeploymentTarget["use_current_instance_count"] = endpointResource.UseCurrentInstanceCount
	return flattenedDeploymentTarget
}

func getAzureCloudServiceDeploymentTargetDataSchema() map[string]*schema.Schema {
	dataSchema := getAzureCloudServiceDeploymentTargetSchema()
	setDataSchema(&dataSchema)

	deploymentTargetDataSchema := getDeploymentTargetDataSchema()

	deploymentTargetDataSchema["azure_cloud_services"] = &schema.Schema{
		Computed:    true,
		Description: "A list of Azure cloud service deployment targets that match the filter(s).",
		Elem:        &schema.Resource{Schema: dataSchema},
		Optional:    true,
		Type:        schema.TypeList,
	}

	delete(deploymentTargetDataSchema, "communication_styles")
	delete(deploymentTargetDataSchema, "deployment_targets")
	deploymentTargetDataSchema["id"] = getIDDataSchema()

	return deploymentTargetDataSchema
}

func getAzureCloudServiceDeploymentTargetSchema() map[string]*schema.Schema {
	azureCloudServiceDeploymentTargetSchema := getDeploymentTargetSchema()

	azureCloudServiceDeploymentTargetSchema["account_id"] = &schema.Schema{
		Required: true,
		Type:     schema.TypeString,
	}

	azureCloudServiceDeploymentTargetSchema["cloud_service_name"] = &schema.Schema{
		Required: true,
		Type:     schema.TypeString,
	}

	azureCloudServiceDeploymentTargetSchema["default_worker_pool_id"] = &schema.Schema{
		Optional: true,
		Type:     schema.TypeString,
	}

	azureCloudServiceDeploymentTargetSchema["slot"] = &schema.Schema{
		Optional: true,
		Type:     schema.TypeString,
	}

	azureCloudServiceDeploymentTargetSchema["storage_account_name"] = &schema.Schema{
		Required: true,
		Type:     schema.TypeString,
	}

	azureCloudServiceDeploymentTargetSchema["swap_if_possible"] = &schema.Schema{
		Optional: true,
		Type:     schema.TypeBool,
	}

	azureCloudServiceDeploymentTargetSchema["use_current_instance_count"] = &schema.Schema{
		Optional: true,
		Type:     schema.TypeBool,
	}

	return azureCloudServiceDeploymentTargetSchema
}

func setAzureCloudServiceDeploymentTarget(ctx context.Context, d *schema.ResourceData, deploymentTarget *octopusdeploy.DeploymentTarget) {
	if deploymentTarget == nil {
		return
	}

	endpointResource, err := octopusdeploy.ToEndpointResource(deploymentTarget.Endpoint)
	if err != nil {
		return
	}

	d.Set("account_id", endpointResource.AccountID)
	d.Set("cloud_service_name", endpointResource.CloudServiceName)
	d.Set("default_worker_pool_id", endpointResource.DefaultWorkerPoolID)
	d.Set("slot", endpointResource.Slot)
	d.Set("storage_account_name", endpointResource.StorageAccountName)
	d.Set("swap_if_possible", endpointResource.SwapIfPossible)
	d.Set("use_current_instance_count", endpointResource.UseCurrentInstanceCount)

	setDeploymentTarget(ctx, d, deploymentTarget)
}
