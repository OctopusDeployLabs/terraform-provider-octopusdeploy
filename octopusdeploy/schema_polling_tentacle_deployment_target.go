package octopusdeploy

import (
	"context"
	"fmt"
	"net/url"

	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/machines"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func expandPollingTentacleDeploymentTarget(d *schema.ResourceData) *machines.DeploymentTarget {
	tentacleURL, _ := url.Parse(d.Get("tentacle_url").(string))
	thumbprint := d.Get("thumbprint").(string)

	endpoint := machines.NewPollingTentacleEndpoint(tentacleURL, thumbprint)

	if v, ok := d.GetOk("certificate_signature_algorithm"); ok {
		endpoint.CertificateSignatureAlgorithm = v.(string)
	}

	if v, ok := d.GetOk("tentacle_version_details"); ok {
		endpoint.TentacleVersionDetails = expandTentacleVersionDetails(v)
	}

	deploymentTarget := expandDeploymentTarget(d)
	deploymentTarget.Endpoint = endpoint
	return deploymentTarget
}

func flattenPollingTentacleDeploymentTarget(deploymentTarget *machines.DeploymentTarget) map[string]interface{} {
	if deploymentTarget == nil {
		return nil
	}

	flattenedDeploymentTarget := flattenDeploymentTarget(deploymentTarget)
	endpointResource, _ := machines.ToEndpointResource(deploymentTarget.Endpoint)
	flattenedDeploymentTarget["certificate_signature_algorithm"] = endpointResource.CertificateSignatureAlgorithm
	flattenedDeploymentTarget["tentacle_version_details"] = flattenTentacleVersionDetails(endpointResource.TentacleVersionDetails)
	flattenedDeploymentTarget["tentacle_url"] = endpointResource.URI.String()
	return flattenedDeploymentTarget
}

func getPollingTentacleDeploymentTargetDataSchema() map[string]*schema.Schema {
	dataSchema := getPollingTentacleDeploymentTargetSchema()
	setDataSchema(&dataSchema)

	deploymentTargetDataSchema := getDeploymentTargetDataSchema()

	deploymentTargetDataSchema["polling_tentacle_deployment_targets"] = &schema.Schema{
		Computed:    true,
		Description: "A list of polling tentacle deployment targets that match the filter(s).",
		Elem:        &schema.Resource{Schema: dataSchema},
		Optional:    true,
		Type:        schema.TypeList,
	}

	delete(deploymentTargetDataSchema, "communication_styles")
	delete(deploymentTargetDataSchema, "deployment_targets")
	deploymentTargetDataSchema["id"] = getDataSchemaID()

	return deploymentTargetDataSchema
}

func getPollingTentacleDeploymentTargetSchema() map[string]*schema.Schema {
	pollingTentacleDeploymentTargetSchema := getDeploymentTargetSchema()

	pollingTentacleDeploymentTargetSchema["certificate_signature_algorithm"] = &schema.Schema{
		Optional: true,
		Type:     schema.TypeString,
	}

	pollingTentacleDeploymentTargetSchema["tentacle_version_details"] = &schema.Schema{
		Computed: true,
		Elem:     &schema.Resource{Schema: getTentacleVersionDetailsSchema()},
		Optional: true,
		Type:     schema.TypeList,
	}

	pollingTentacleDeploymentTargetSchema["tentacle_url"] = &schema.Schema{
		Required: true,
		Type:     schema.TypeString,
	}

	return pollingTentacleDeploymentTargetSchema
}

func setPollingTentacleDeploymentTarget(ctx context.Context, d *schema.ResourceData, deploymentTarget *machines.DeploymentTarget) error {
	endpointResource, err := machines.ToEndpointResource(deploymentTarget.Endpoint)
	if err != nil {
		return err
	}

	d.Set("certificate_signature_algorithm", endpointResource.CertificateSignatureAlgorithm)

	if err := d.Set("tentacle_version_details", flattenTentacleVersionDetails(endpointResource.TentacleVersionDetails)); err != nil {
		return fmt.Errorf("error setting tentacle_version_details: %s", err)
	}

	d.Set("tentacle_url", endpointResource.URI.String())

	return setDeploymentTarget(ctx, d, deploymentTarget)
}
