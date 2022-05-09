package octopusdeploy

import (
	"github.com/OctopusDeploy/go-octopusdeploy/octopusdeploy"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func flattenPollingTentacleDeploymentTarget(deploymentTarget *octopusdeploy.DeploymentTarget) map[string]interface{} {
	if deploymentTarget == nil {
		return nil
	}

	flattenedDeploymentTarget := flattenDeploymentTarget(deploymentTarget)
	endpointResource, _ := octopusdeploy.ToEndpointResource(deploymentTarget.Endpoint)
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
