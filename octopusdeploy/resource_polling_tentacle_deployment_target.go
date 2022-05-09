package octopusdeploy

import (
	"context"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"log"
	"net/url"

	"github.com/OctopusDeploy/go-octopusdeploy/octopusdeploy"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourcePollingTentacleDeploymentTarget() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourcePollingTentacleDeploymentTargetCreate,
		DeleteContext: resourcePollingTentacleDeploymentTargetDelete,
		Description:   "This resource manages polling tentacle deployment targets in Octopus Deploy.",
		Importer:      getImporter(),
		ReadContext:   resourcePollingTentacleDeploymentTargetRead,
		Schema:        getPollingTentacleDeploymentTargetResourceSchema(),
		UpdateContext: resourcePollingTentacleDeploymentTargetUpdate,
		SchemaVersion: 1,
		StateUpgraders: []schema.StateUpgrader{
			{
				Type:    resourcePollingTentacleDeploymentTargetSchemaV0().CoreConfigSchema().ImpliedType(),
				Upgrade: resourcePollingTentacleDeploymentTargetStateUpgradeV0,
				Version: 0,
			},
		},
	}
}

func resourcePollingTentacleDeploymentTargetSchemaV0() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"certificate_signature_algorithm": {
				Optional: true,
				Type:     schema.TypeString,
			},
			"tentacle_version_details": {
				Computed: true,
				Elem:     &schema.Resource{Schema: getTentacleVersionDetailsSchema()},
				Optional: true,
				Type:     schema.TypeList,
			},
			"tentacle_url": {
				Required: true,
				Type:     schema.TypeString,
			},
			"endpoint": {
				Computed: true,
				Elem:     &schema.Resource{Schema: getEndpointSchema()},
				MinItems: 1,
				Optional: true,
				Type:     schema.TypeList,
			},
			"environments": {
				Description: "A list of environment IDs associated with this resource.",
				Elem:        &schema.Schema{Type: schema.TypeString},
				MinItems:    1,
				Required:    true,
				Type:        schema.TypeList,
			},
			"has_latest_calamari": {
				Computed: true,
				Type:     schema.TypeBool,
			},
			"health_status": getHealthStatusSchema(),
			"id":            getIDSchema(),
			"is_disabled": {
				Computed: true,
				Optional: true,
				Type:     schema.TypeBool,
			},
			"is_in_process": {
				Computed: true,
				Type:     schema.TypeBool,
			},
			"machine_policy_id": {
				Computed: true,
				Optional: true,
				Type:     schema.TypeString,
			},
			"name": getNameSchema(true),
			"operating_system": {
				Computed: true,
				Optional: true,
				Type:     schema.TypeString,
			},
			"roles": {
				Elem:     &schema.Schema{Type: schema.TypeString},
				MinItems: 1,
				Required: true,
				Type:     schema.TypeList,
			},
			"shell_name": {
				Computed: true,
				Optional: true,
				Type:     schema.TypeString,
			},
			"shell_version": {
				Computed: true,
				Optional: true,
				Type:     schema.TypeString,
			},
			"space_id":                          getSpaceIDSchema(),
			"status":                            getStatusSchema(),
			"status_summary":                    getStatusSummarySchema(),
			"tenanted_deployment_participation": getTenantedDeploymentSchema(),
			"tenants":                           getTenantsSchema(),
			"tenant_tags":                       getTenantTagsSchema(),
			"thumbprint": {
				Computed: true,
				Optional: true,
				Type:     schema.TypeString,
			},
			"uri": {
				Computed: true,
				Optional: true,
				Type:     schema.TypeString,
			},
		},
	}
}

func resourcePollingTentacleDeploymentTargetStateUpgradeV0(_ context.Context, rawState map[string]interface{}, _ interface{}) (map[string]interface{}, error) {
	delete(rawState, "tentacle_version_details")
	delete(rawState, "endpoint")
	delete(rawState, "status_summary")
	delete(rawState, "status")
	delete(rawState, "shell_version")
	delete(rawState, "shell_name")
	delete(rawState, "operating_system")
	delete(rawState, "is_in_process")
	delete(rawState, "health_status")
	delete(rawState, "has_latest_calamari")
	delete(rawState, "certificate_signature_algorithm")
	rawState["subscription_id"] = rawState["tentacle_url"]
	delete(rawState, "uri")
	delete(rawState, "tentacle_url")
	return rawState, nil
}

func getPollingTentacleDeploymentTargetResourceSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"environments": {
			Description: "A list of environment IDs associated with this listening tentacle.",
			Elem:        &schema.Schema{Type: schema.TypeString},
			Required:    true,
			MinItems:    1,
			Type:        schema.TypeList,
		},
		"id": getIDSchema(),
		"is_disabled": {
			Computed:    true,
			Description: "Represents the disabled status of this deployment target.",
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
		"roles": {
			Description: "A list of role IDs that are associated with this deployment target.",
			Elem:        &schema.Schema{Type: schema.TypeString},
			MinItems:    1,
			Required:    true,
			Type:        schema.TypeList,
		},
		"space_id":                          getSpaceIDSchema(),
		"tenanted_deployment_participation": getTenantedDeploymentSchema(),
		"tenants":                           getTenantsSchema(),
		"tenant_tags":                       getTenantTagsSchema(),
		"subscription_id": {
			Description:      "The subscription id is a random 20 character id that is used to queue messages from the server to the Polling Tentacle. This should match the value in the Tentacle config file.",
			Required:         true,
			Type:             schema.TypeString,
			ValidateDiagFunc: validation.ToDiagFunc(validation.StringIsNotEmpty),
		},
		"thumbprint": {
			Description:      "The thumbprint of this deployment target.",
			Required:         true,
			Type:             schema.TypeString,
			ValidateDiagFunc: validation.ToDiagFunc(validation.StringIsNotEmpty),
		},
		"lock_upgrade": {
			Default:     false,
			Description: "Whether to lock the tentacle version to prevent upgrades.",
			Optional:    true,
			Type:        schema.TypeBool,
		},
	}
}

func resourcePollingTentacleDeploymentTargetCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	deploymentTarget := expandPollingTentacleDeploymentTarget(d)

	log.Printf("[INFO] creating polling tentacle deployment target: %#v", deploymentTarget)

	client := m.(*octopusdeploy.Client)
	createdDeploymentTarget, err := client.Machines.Add(deploymentTarget)
	if err != nil {
		return diag.FromErr(err)
	}

	if err := setPollingTentacleDeploymentTarget(ctx, d, createdDeploymentTarget); err != nil {
		return diag.FromErr(err)
	}

	d.SetId(createdDeploymentTarget.GetID())

	log.Printf("[INFO] polling tentacle deployment target created (%s)", d.Id())
	return nil
}

func resourcePollingTentacleDeploymentTargetDelete(_ context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[INFO] deleting polling tentacle deployment target (%s)", d.Id())

	client := m.(*octopusdeploy.Client)
	if err := client.Machines.DeleteByID(d.Id()); err != nil {
		return diag.FromErr(err)
	}

	d.SetId("")

	log.Printf("[INFO] polling tentacle deployment target deleted")
	return nil
}

func resourcePollingTentacleDeploymentTargetRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[INFO] reading polling tentacle deployment target (%s)", d.Id())

	client := m.(*octopusdeploy.Client)
	deploymentTarget, err := client.Machines.GetByID(d.Id())
	if err != nil {
		if apiError, ok := err.(*octopusdeploy.APIError); ok {
			if apiError.StatusCode == 404 {
				log.Printf("[INFO] polling tentacle deployment target (%s) not found; deleting from state", d.Id())
				d.SetId("")
				return nil
			}
		}
		return diag.FromErr(err)
	}

	if err := setPollingTentacleDeploymentTarget(ctx, d, deploymentTarget); err != nil {
		return diag.FromErr(err)
	}

	log.Printf("[INFO] polling tentacle deployment target read (%s)", d.Id())
	return nil
}

func resourcePollingTentacleDeploymentTargetUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[INFO] updating polling tentacle deployment target (%s)", d.Id())

	deploymentTarget := expandPollingTentacleDeploymentTarget(d)
	client := m.(*octopusdeploy.Client)
	updatedDeploymentTarget, err := client.Machines.Update(deploymentTarget)
	if err != nil {
		return diag.FromErr(err)
	}

	if err := setPollingTentacleDeploymentTarget(ctx, d, updatedDeploymentTarget); err != nil {
		return diag.FromErr(err)
	}

	log.Printf("[INFO] polling tentacle deployment target updated (%s)", d.Id())
	return nil
}

func expandPollingTentacleDeploymentTarget(d *schema.ResourceData) *octopusdeploy.DeploymentTarget {
	subscriptionId, _ := url.Parse(d.Get("subscription_id").(string))
	thumbprint := d.Get("thumbprint").(string)

	endpoint := octopusdeploy.NewPollingTentacleEndpoint(subscriptionId, thumbprint)

	if v, ok := d.GetOk("lock_upgrade"); ok {
		endpoint.TentacleVersionDetails = &octopusdeploy.TentacleVersionDetails{
			UpgradeLocked: v.(bool),
		}
	}

	deploymentTarget := expandDeploymentTarget(d)
	deploymentTarget.Endpoint = endpoint
	return deploymentTarget
}

func setPollingTentacleDeploymentTarget(ctx context.Context, d *schema.ResourceData, deploymentTarget *octopusdeploy.DeploymentTarget) error {
	endpointResource, err := octopusdeploy.ToEndpointResource(deploymentTarget.Endpoint)
	if err != nil {
		return err
	}

	d.Set("lock_upgrade", endpointResource.TentacleVersionDetails.UpgradeLocked)
	d.Set("subscription_id", endpointResource.URI.String())

	return setDeploymentTarget(d, deploymentTarget)
}
