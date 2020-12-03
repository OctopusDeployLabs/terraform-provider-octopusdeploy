package octopusdeploy

import (
	"context"
	"net/url"

	"github.com/OctopusDeploy/go-octopusdeploy/octopusdeploy"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func expandListeningTentacleDeploymentTarget(d *schema.ResourceData) *octopusdeploy.DeploymentTarget {
	tentacleURL, _ := url.Parse(d.Get("tentacle_url").(string))
	thumbprint := d.Get("thumbprint").(string)

	endpoint := octopusdeploy.NewListeningTentacleEndpoint(tentacleURL, thumbprint)

	if v, ok := d.GetOk("certificate_signature_algorithm"); ok {
		endpoint.CertificateSignatureAlgorithm = v.(string)
	}

	if v, ok := d.GetOk("proxy_id"); ok {
		endpoint.ProxyID = v.(string)
	}

	if v, ok := d.GetOk("tentacle_version_details"); ok {
		endpoint.TentacleVersionDetails = expandTentacleVersionDetails(v)
	}

	deploymentTarget := expandDeploymentTarget(d)
	deploymentTarget.Endpoint = endpoint
	return deploymentTarget
}

func flattenListeningTentacleDeploymentTarget(deploymentTarget *octopusdeploy.DeploymentTarget) map[string]interface{} {
	if deploymentTarget == nil {
		return nil
	}

	flattenedDeploymentTarget := flattenDeploymentTarget(deploymentTarget)
	endpointResource, _ := octopusdeploy.ToEndpointResource(deploymentTarget.Endpoint)
	flattenedDeploymentTarget["certificate_signature_algorithm"] = endpointResource.CertificateSignatureAlgorithm
	flattenedDeploymentTarget["proxy_id"] = endpointResource.ProxyID
	flattenedDeploymentTarget["tentacle_version_details"] = flattenTentacleVersionDetails(endpointResource.TentacleVersionDetails)
	flattenedDeploymentTarget["tentacle_url"] = endpointResource.URI.String()
	return flattenedDeploymentTarget
}

func getListeningTentacleDeploymentTargetDataSchema() map[string]*schema.Schema {
	dataSchema := getListeningTentacleDeploymentTargetSchema()
	setDataSchema(&dataSchema)

	deploymentTargetDataSchema := getDeploymentTargetDataSchema()

	deploymentTargetDataSchema["listening_tentacles"] = &schema.Schema{
		Computed:    true,
		Description: "A list of listening tentacle deployment targets that match the filter(s).",
		Elem:        &schema.Resource{Schema: dataSchema},
		Optional:    true,
		Type:        schema.TypeList,
	}

	delete(deploymentTargetDataSchema, "communication_styles")
	delete(deploymentTargetDataSchema, "deployment_targets")
	deploymentTargetDataSchema["id"] = getIDDataSchema()

	return deploymentTargetDataSchema
}

func getListeningTentacleDeploymentTargetSchema() map[string]*schema.Schema {
	listeningTentacleDeploymentTargetSchema := getDeploymentTargetSchema()

	listeningTentacleDeploymentTargetSchema["certificate_signature_algorithm"] = &schema.Schema{
		Computed: true,
		Optional: true,
		Type:     schema.TypeString,
	}

	listeningTentacleDeploymentTargetSchema["proxy_id"] = &schema.Schema{
		Computed: true,
		Optional: true,
		Type:     schema.TypeString,
	}

	listeningTentacleDeploymentTargetSchema["tentacle_version_details"] = &schema.Schema{
		Computed: true,
		Elem:     &schema.Resource{Schema: getTentacleVersionDetailsSchema()},
		Optional: true,
		Type:     schema.TypeList,
	}

	listeningTentacleDeploymentTargetSchema["tentacle_url"] = &schema.Schema{
		Required: true,
		Type:     schema.TypeString,
	}

	return listeningTentacleDeploymentTargetSchema
}

func setListeningTentacleDeploymentTarget(ctx context.Context, d *schema.ResourceData, deploymentTarget *octopusdeploy.DeploymentTarget) {
	if deploymentTarget == nil {
		return
	}

	endpointResource, err := octopusdeploy.ToEndpointResource(deploymentTarget.Endpoint)
	if err != nil {
		return
	}

	d.Set("certificate_signature_algorithm", endpointResource.CertificateSignatureAlgorithm)
	d.Set("proxy_id", endpointResource.ProxyID)
	d.Set("tentacle_version_details", flattenTentacleVersionDetails(endpointResource.TentacleVersionDetails))
	d.Set("tentacle_url", endpointResource.URI.String())

	setDeploymentTarget(ctx, d, deploymentTarget)
}
