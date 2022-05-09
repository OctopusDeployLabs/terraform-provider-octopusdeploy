package octopusdeploy

import (
	"github.com/OctopusDeploy/go-octopusdeploy/octopusdeploy"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

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
	dataSchema := map[string]*schema.Schema{
		"certificate_signature_algorithm": {
			Computed: true,
			Optional: true,
			Type:     schema.TypeString,
		},
		"environments": {
			Description: "A list of environment IDs associated with this listening tentacle.",
			Elem:        &schema.Schema{Type: schema.TypeString},
			Required:    true,
			MinItems:    1,
			Type:        schema.TypeList,
		},
		"has_latest_calamari": {
			Computed: true,
			Type:     schema.TypeBool,
		},
		"health_status": getHealthStatusSchema(),
		"id":            getIDSchema(),
		"is_disabled": {
			Computed:    true,
			Description: "Represents the disabled status of this deployment target.",
			Optional:    true,
			Type:        schema.TypeBool,
		},
		"is_in_process": {
			Computed:    true,
			Description: "Represents the in-process status of this deployment target.",
			Optional:    true,
			Type:        schema.TypeBool,
		},
		"machine_policy_id": {
			Computed:    true,
			Description: "The machine policy ID that is associated with this deployment target.",
			Optional:    true,
			Type:        schema.TypeString,
		},
		"name": getNameSchema(true),
		"operating_system": {
			Computed:    true,
			Description: "The operating system that is associated with this deployment target.",
			Optional:    true,
			Type:        schema.TypeString,
		},
		"proxy_id": {
			Computed:    true,
			Description: "The proxy ID that is associated with this deployment target.",
			Optional:    true,
			Type:        schema.TypeString,
		},
		"roles": {
			Description: "A list of role IDs that are associated with this deployment target.",
			Elem:        &schema.Schema{Type: schema.TypeString},
			MinItems:    1,
			Required:    true,
			Type:        schema.TypeList,
		},
		"shell_name": {
			Computed:    true,
			Description: "The shell name associated with this deployment target.",
			Optional:    true,
			Type:        schema.TypeString,
		},
		"shell_version": {
			Computed:    true,
			Description: "The shell version associated with this deployment target.",
			Optional:    true,
			Type:        schema.TypeString,
		},
		"space_id":                          getSpaceIDSchema(),
		"status":                            getStatusSchema(),
		"status_summary":                    getStatusSummarySchema(),
		"tenanted_deployment_participation": getTenantedDeploymentSchema(),
		"tenants":                           getTenantsSchema(),
		"tenant_tags":                       getTenantTagsSchema(),
		"tentacle_version_details": {
			Computed: true,
			Elem:     &schema.Resource{Schema: getTentacleVersionDetailsSchema()},
			Optional: true,
			Type:     schema.TypeList,
		},
		"tentacle_url": {
			Description:      "The tenant URL of this deployment target.",
			Required:         true,
			Type:             schema.TypeString,
			ValidateDiagFunc: validation.ToDiagFunc(validation.IsURLWithHTTPS),
		},
		"thumbprint": {
			Description:      "The thumbprint of this deployment target.",
			Required:         true,
			Type:             schema.TypeString,
			ValidateDiagFunc: validation.ToDiagFunc(validation.StringIsNotEmpty),
		},
		"uri": {
			Computed:    true,
			Description: "The URI of this deployment target.",
			Optional:    true,
			Type:        schema.TypeString,
		},
	}

	setDataSchema(&dataSchema)

	deploymentTargetDataSchema := getDeploymentTargetDataSchema()

	deploymentTargetDataSchema["listening_tentacle_deployment_targets"] = &schema.Schema{
		Computed:    true,
		Description: "A list of listening tentacle deployment targets that match the filter(s).",
		Elem:        &schema.Resource{Schema: dataSchema},
		Optional:    true,
		Type:        schema.TypeList,
	}

	delete(deploymentTargetDataSchema, "communication_styles")
	delete(deploymentTargetDataSchema, "deployment_targets")
	deploymentTargetDataSchema["id"] = getDataSchemaID()

	return deploymentTargetDataSchema
}
