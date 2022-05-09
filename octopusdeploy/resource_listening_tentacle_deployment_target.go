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

func resourceListeningTentacleDeploymentTarget() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceListeningTentacleDeploymentTargetCreate,
		DeleteContext: resourceListeningTentacleDeploymentTargetDelete,
		Description:   "This resource manages listening tentacle deployment targets in Octopus Deploy.",
		Importer:      getImporter(),
		ReadContext:   resourceListeningTentacleDeploymentTargetRead,
		Schema:        getListeningTentacleDeploymentTargetResourceSchema(),
		UpdateContext: resourceListeningTentacleDeploymentTargetUpdate,
		SchemaVersion: 1,
		StateUpgraders: []schema.StateUpgrader{
			{
				Type:    resourceListeningTentacleDeploymentTargetSchemaV0().CoreConfigSchema().ImpliedType(),
				Upgrade: resourceListeningTentacleDeploymentTargetStateUpgradeV0,
				Version: 0,
			},
		},
	}
}

func getListeningTentacleDeploymentTargetResourceSchema() map[string]*schema.Schema {
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
		"space_id":                          getSpaceIDSchema(),
		"tenanted_deployment_participation": getTenantedDeploymentSchema(),
		"tenants":                           getTenantsSchema(),
		"tenant_tags":                       getTenantTagsSchema(),
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
		"lock_upgrade": {
			Default:     false,
			Description: "Whether to lock the tentacle version to prevent upgrades.",
			Optional:    true,
			Type:        schema.TypeBool,
		},
	}
}

func resourceListeningTentacleDeploymentTargetSchemaV0() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
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
		},
	}
}

func resourceListeningTentacleDeploymentTargetStateUpgradeV0(_ context.Context, rawState map[string]interface{}, _ interface{}) (map[string]interface{}, error) {
	delete(rawState, "tentacle_version_details")
	delete(rawState, "status_summary")
	delete(rawState, "status")
	delete(rawState, "shell_version")
	delete(rawState, "shell_name")
	delete(rawState, "operating_system")
	delete(rawState, "is_in_process")
	delete(rawState, "health_status")
	delete(rawState, "has_latest_calamari")
	delete(rawState, "uri")
	delete(rawState, "certificate_signature_algorithm")

	return rawState, nil
}

func resourceListeningTentacleDeploymentTargetCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	deploymentTarget := expandListeningTentacleDeploymentTarget(d)

	log.Printf("[INFO] creating listening tentacle deployment target: %#v", deploymentTarget)

	client := m.(*octopusdeploy.Client)
	createdDeploymentTarget, err := client.Machines.Add(deploymentTarget)
	if err != nil {
		return diag.FromErr(err)
	}

	if err := setListeningTentacleDeploymentTarget(ctx, d, createdDeploymentTarget); err != nil {
		return diag.FromErr(err)
	}

	d.SetId(createdDeploymentTarget.GetID())

	log.Printf("[INFO] listening tentacle deployment target created (%s)", d.Id())
	return nil
}

func resourceListeningTentacleDeploymentTargetDelete(_ context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[INFO] deleting listening tentacle deployment target (%s)", d.Id())

	client := m.(*octopusdeploy.Client)
	if err := client.Machines.DeleteByID(d.Id()); err != nil {
		return diag.FromErr(err)
	}

	d.SetId("")

	log.Printf("[INFO] listening tentacle deployment target deleted")
	return nil
}

func resourceListeningTentacleDeploymentTargetRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[INFO] reading listening tentacle deployment target (%s)", d.Id())

	client := m.(*octopusdeploy.Client)
	deploymentTarget, err := client.Machines.GetByID(d.Id())
	if err != nil {
		if apiError, ok := err.(*octopusdeploy.APIError); ok {
			if apiError.StatusCode == 404 {
				log.Printf("[INFO] listening tentacle deployment target (%s) not found; deleting from state", d.Id())
				d.SetId("")
				return nil
			}
		}
		return diag.FromErr(err)
	}

	if err := setListeningTentacleDeploymentTarget(ctx, d, deploymentTarget); err != nil {
		return diag.FromErr(err)
	}

	log.Printf("[INFO] listening tentacle deployment target read (%s)", d.Id())
	return nil
}

func resourceListeningTentacleDeploymentTargetUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[INFO] updating listening tentacle deployment target (%s)", d.Id())

	deploymentTarget := expandListeningTentacleDeploymentTarget(d)
	client := m.(*octopusdeploy.Client)
	updatedDeploymentTarget, err := client.Machines.Update(deploymentTarget)
	if err != nil {
		return diag.FromErr(err)
	}

	if err := setListeningTentacleDeploymentTarget(ctx, d, updatedDeploymentTarget); err != nil {
		return diag.FromErr(err)
	}

	log.Printf("[INFO] listening tentacle deployment target updated (%s)", d.Id())
	return nil
}

func setListeningTentacleDeploymentTarget(ctx context.Context, d *schema.ResourceData, deploymentTarget *octopusdeploy.DeploymentTarget) error {
	endpointResource, err := octopusdeploy.ToEndpointResource(deploymentTarget.Endpoint)
	if err != nil {
		return err
	}

	d.Set("proxy_id", endpointResource.ProxyID)
	d.Set("tentacle_url", endpointResource.URI.String())
	d.Set("lock_upgrade", endpointResource.TentacleVersionDetails.UpgradeLocked)

	return setDeploymentTarget(d, deploymentTarget)
}

func expandListeningTentacleDeploymentTarget(d *schema.ResourceData) *octopusdeploy.DeploymentTarget {
	tentacleURL, _ := url.Parse(d.Get("tentacle_url").(string))
	thumbprint := d.Get("thumbprint").(string)

	endpoint := octopusdeploy.NewListeningTentacleEndpoint(tentacleURL, thumbprint)

	if v, ok := d.GetOk("proxy_id"); ok {
		endpoint.ProxyID = v.(string)
	}

	if v, ok := d.GetOk("lock_upgrade"); ok {
		endpoint.TentacleVersionDetails = &octopusdeploy.TentacleVersionDetails{
			UpgradeLocked: v.(bool),
		}
	}

	deploymentTarget := expandDeploymentTarget(d)
	deploymentTarget.Endpoint = endpoint
	return deploymentTarget
}
