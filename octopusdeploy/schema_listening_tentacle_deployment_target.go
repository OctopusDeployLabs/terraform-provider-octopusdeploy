package octopusdeploy

import (
	"context"
	"fmt"
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

func getListeningTentacleDeploymentTargetSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"certificate_signature_algorithm": {
			Computed: true,
			Optional: true,
			Type:     schema.TypeString,
		},
		"environments": getEnvironmentsSchema(),
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
			Description: "The tenant URL of this deployment target.",
			Required:    true,
			Type:        schema.TypeString,
			// ValidateDiagFunc: validateDiagFunc(validation.IsURLWithHTTPorHTTPS),
		},
		"thumbprint": {
			Computed:    true,
			Description: "The thumbprint of this deployment target.",
			Optional:    true,
			Type:        schema.TypeString,
		},
		"uri": {
			Computed:    true,
			Description: "The URI of this deployment target.",
			Optional:    true,
			Type:        schema.TypeString,
			// ValidateDiagFunc: validateDiagFunc(validation.IsURLWithHTTPorHTTPS),
		},
	}
}

func setListeningTentacleDeploymentTarget(ctx context.Context, d *schema.ResourceData, deploymentTarget *octopusdeploy.DeploymentTarget) error {
	endpointResource, err := octopusdeploy.ToEndpointResource(deploymentTarget.Endpoint)
	if err != nil {
		return err
	}

	d.Set("certificate_signature_algorithm", endpointResource.CertificateSignatureAlgorithm)
	d.Set("proxy_id", endpointResource.ProxyID)
	d.Set("tentacle_url", endpointResource.URI.String())

	if err := d.Set("tentacle_version_details", flattenTentacleVersionDetails(endpointResource.TentacleVersionDetails)); err != nil {
		return fmt.Errorf("error setting tentacle_version_details: %s", err)
	}

	return setDeploymentTarget(ctx, d, deploymentTarget)
}
