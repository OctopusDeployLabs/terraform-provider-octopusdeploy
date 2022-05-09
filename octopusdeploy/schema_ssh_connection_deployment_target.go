package octopusdeploy

import (
	"context"

	"github.com/OctopusDeploy/go-octopusdeploy/octopusdeploy"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func expandSSHConnectionDeploymentTarget(d *schema.ResourceData) *octopusdeploy.DeploymentTarget {
	host := d.Get("host").(string)
	port := d.Get("port").(int)
	fingerprint := d.Get("fingerprint").(string)

	endpoint := octopusdeploy.NewSSHEndpoint(host, port, fingerprint)

	if v, ok := d.GetOk("account_id"); ok {
		endpoint.AccountID = v.(string)
	}

	if v, ok := d.GetOk("dot_net_core_platform"); ok {
		endpoint.DotNetCorePlatform = v.(string)
	}

	if v, ok := d.GetOk("proxy_id"); ok {
		endpoint.ProxyID = v.(string)
	}

	deploymentTarget := expandDeploymentTarget(d)
	deploymentTarget.Endpoint = endpoint
	return deploymentTarget
}

func flattenSSHConnectionDeploymentTarget(deploymentTarget *octopusdeploy.DeploymentTarget) map[string]interface{} {
	if deploymentTarget == nil {
		return nil
	}

	flattenedDeploymentTarget := flattenDeploymentTarget(deploymentTarget)
	endpointResource, _ := octopusdeploy.ToEndpointResource(deploymentTarget.Endpoint)
	flattenedDeploymentTarget["account_id"] = endpointResource.AccountID
	flattenedDeploymentTarget["dot_net_core_platform"] = endpointResource.DotNetCorePlatform
	flattenedDeploymentTarget["fingerprint"] = endpointResource.Fingerprint
	flattenedDeploymentTarget["host"] = endpointResource.Host
	flattenedDeploymentTarget["proxy_id"] = endpointResource.ProxyID
	flattenedDeploymentTarget["port"] = endpointResource.Port
	return flattenedDeploymentTarget
}

func getSSHConnectionDeploymentTargetDataSchema() map[string]*schema.Schema {
	dataSchema := getSSHConnectionDeploymentTargetSchema()
	setDataSchema(&dataSchema)

	deploymentTargetDataSchema := getDeploymentTargetDataSchema()

	deploymentTargetDataSchema["ssh_connection_deployment_targets"] = &schema.Schema{
		Computed:    true,
		Description: "A list of SSH connection deployment targets that match the filter(s).",
		Elem:        &schema.Resource{Schema: dataSchema},
		Optional:    true,
		Type:        schema.TypeList,
	}

	delete(deploymentTargetDataSchema, "communication_styles")
	delete(deploymentTargetDataSchema, "deployment_targets")
	deploymentTargetDataSchema["id"] = getDataSchemaID()

	return deploymentTargetDataSchema
}

func getSSHConnectionDeploymentTargetSchema() map[string]*schema.Schema {
	sshConnectionDeploymentTargetSchema := getDeploymentTargetSchema()

	sshConnectionDeploymentTargetSchema["account_id"] = &schema.Schema{
		Required: true,
		Type:     schema.TypeString,
	}

	sshConnectionDeploymentTargetSchema["dot_net_core_platform"] = &schema.Schema{
		Optional: true,
		Type:     schema.TypeString,
	}

	sshConnectionDeploymentTargetSchema["fingerprint"] = &schema.Schema{
		Required: true,
		Type:     schema.TypeString,
	}

	sshConnectionDeploymentTargetSchema["host"] = &schema.Schema{
		Required: true,
		Type:     schema.TypeString,
	}

	sshConnectionDeploymentTargetSchema["port"] = &schema.Schema{
		Default:  22,
		Optional: true,
		Type:     schema.TypeInt,
	}

	sshConnectionDeploymentTargetSchema["proxy_id"] = &schema.Schema{
		Optional: true,
		Type:     schema.TypeString,
	}

	return sshConnectionDeploymentTargetSchema
}

func setSSHConnectionDeploymentTarget(ctx context.Context, d *schema.ResourceData, deploymentTarget *octopusdeploy.DeploymentTarget) error {
	endpointResource, err := octopusdeploy.ToEndpointResource(deploymentTarget.Endpoint)
	if err != nil {
		return err
	}

	d.Set("account_id", endpointResource.AccountID)
	d.Set("dot_net_core_platform", endpointResource.DotNetCorePlatform)
	d.Set("fingerprint", endpointResource.Fingerprint)
	d.Set("host", endpointResource.Host)
	d.Set("port", endpointResource.Port)
	d.Set("proxy_id", endpointResource.ProxyID)

	return setDeploymentTarget(d, deploymentTarget)
}
