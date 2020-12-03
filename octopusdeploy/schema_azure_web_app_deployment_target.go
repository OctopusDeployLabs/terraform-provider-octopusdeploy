package octopusdeploy

import (
	"context"

	"github.com/OctopusDeploy/go-octopusdeploy/octopusdeploy"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func expandAzureWebAppDeploymentTarget(d *schema.ResourceData) *octopusdeploy.DeploymentTarget {
	endpoint := octopusdeploy.NewAzureWebAppEndpoint()

	if v, ok := d.GetOk("account_id"); ok {
		endpoint.AccountID = v.(string)
	}

	if v, ok := d.GetOk("resource_group_name"); ok {
		endpoint.ResourceGroupName = v.(string)
	}

	if v, ok := d.GetOk("web_app_name"); ok {
		endpoint.WebAppName = v.(string)
	}

	if v, ok := d.GetOk("web_app_slot_name"); ok {
		endpoint.WebAppSlotName = v.(string)
	}

	deploymentTarget := expandDeploymentTarget(d)
	deploymentTarget.Endpoint = endpoint
	return deploymentTarget
}

func flattenAzureWebAppDeploymentTarget(deploymentTarget *octopusdeploy.DeploymentTarget) map[string]interface{} {
	if deploymentTarget == nil {
		return nil
	}

	flattenedDeploymentTarget := flattenDeploymentTarget(deploymentTarget)
	endpointResource, _ := octopusdeploy.ToEndpointResource(deploymentTarget.Endpoint)
	flattenedDeploymentTarget["account_id"] = endpointResource.AccountID
	flattenedDeploymentTarget["resource_group_name"] = endpointResource.ResourceGroupName
	flattenedDeploymentTarget["web_app_name"] = endpointResource.WebAppName
	flattenedDeploymentTarget["web_app_slot_name"] = endpointResource.WebAppSlotName
	return flattenedDeploymentTarget
}

func getAzureWebAppDeploymentTargetDataSchema() map[string]*schema.Schema {
	dataSchema := getAzureWebAppDeploymentTargetSchema()
	setDataSchema(&dataSchema)

	deploymentTargetDataSchema := getDeploymentTargetDataSchema()

	deploymentTargetDataSchema["azure_web_apps"] = &schema.Schema{
		Computed:    true,
		Description: "A list of Azure web app deployment targets that match the filter(s).",
		Elem:        &schema.Resource{Schema: dataSchema},
		Optional:    true,
		Type:        schema.TypeList,
	}

	delete(deploymentTargetDataSchema, "communication_styles")
	delete(deploymentTargetDataSchema, "deployment_targets")
	deploymentTargetDataSchema["id"] = getIDDataSchema()

	return deploymentTargetDataSchema
}

func getAzureWebAppDeploymentTargetSchema() map[string]*schema.Schema {
	azureWebAppDeploymentTargetSchema := getDeploymentTargetSchema()

	azureWebAppDeploymentTargetSchema["account_id"] = &schema.Schema{
		Required: true,
		Type:     schema.TypeString,
	}

	azureWebAppDeploymentTargetSchema["resource_group_name"] = &schema.Schema{
		Required: true,
		Type:     schema.TypeString,
	}

	azureWebAppDeploymentTargetSchema["web_app_name"] = &schema.Schema{
		Required: true,
		Type:     schema.TypeString,
	}

	azureWebAppDeploymentTargetSchema["web_app_slot_name"] = &schema.Schema{
		Optional: true,
		Type:     schema.TypeString,
	}

	return azureWebAppDeploymentTargetSchema
}

func setAzureWebAppDeploymentTarget(ctx context.Context, d *schema.ResourceData, deploymentTarget *octopusdeploy.DeploymentTarget) {
	if deploymentTarget == nil {
		return
	}

	endpointResource, err := octopusdeploy.ToEndpointResource(deploymentTarget.Endpoint)
	if err != nil {
		return
	}

	d.Set("account_id", endpointResource.AccountID)
	d.Set("resource_group_name", endpointResource.ResourceGroupName)
	d.Set("web_app_name", endpointResource.WebAppName)
	d.Set("web_app_slot_name", endpointResource.WebAppSlotName)

	setDeploymentTarget(ctx, d, deploymentTarget)
}
