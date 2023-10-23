package octopusdeploy

import (
	"context"
	"fmt"

	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/core"
	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/machines"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func expandDeploymentTarget(d *schema.ResourceData) *machines.DeploymentTarget {
	deploymentMode := core.TenantedDeploymentMode(d.Get("tenanted_deployment_participation").(string))
	endpoint := expandEndpoint(d.Get("endpoint"))
	environments := getSliceFromTerraformTypeList(d.Get("environments"))
	name := d.Get("name").(string)
	roles := getSliceFromTerraformTypeList(d.Get("roles"))
	tenantIDs := getSliceFromTerraformTypeList(d.Get("tenants"))
	tenantTags := getSliceFromTerraformTypeList(d.Get("tenant_tags"))

	deploymentTarget := machines.NewDeploymentTarget(name, endpoint, environments, roles)
	deploymentTarget.ID = d.Id()
	deploymentTarget.TenantedDeploymentMode = deploymentMode
	deploymentTarget.TenantIDs = tenantIDs
	deploymentTarget.TenantTags = tenantTags

	if v, ok := d.GetOk("machine_policy_id"); ok {
		deploymentTarget.MachinePolicyID = v.(string)
	}

	if v, ok := d.GetOk("is_disabled"); ok {
		deploymentTarget.IsDisabled = v.(bool)
	}

	if v, ok := d.GetOk("thumbprint"); ok {
		deploymentTarget.Thumbprint = v.(string)
	}

	if v, ok := d.GetOk("uri"); ok {
		deploymentTarget.URI = v.(string)
	}

	if v, ok := d.GetOk("space_id"); ok {
		deploymentTarget.SpaceID = v.(string)
	}

	return deploymentTarget
}

func flattenDeploymentTarget(deploymentTarget *machines.DeploymentTarget) map[string]interface{} {
	if deploymentTarget == nil {
		return nil
	}

	endpointResource, _ := machines.ToEndpointResource(deploymentTarget.Endpoint)

	return map[string]interface{}{
		"endpoint":                          flattenEndpointResource(endpointResource),
		"environments":                      deploymentTarget.EnvironmentIDs,
		"has_latest_calamari":               deploymentTarget.HasLatestCalamari,
		"health_status":                     deploymentTarget.HealthStatus,
		"id":                                deploymentTarget.GetID(),
		"is_disabled":                       deploymentTarget.IsDisabled,
		"is_in_process":                     deploymentTarget.IsInProcess,
		"machine_policy_id":                 deploymentTarget.MachinePolicyID,
		"name":                              deploymentTarget.Name,
		"operating_system":                  deploymentTarget.OperatingSystem,
		"roles":                             deploymentTarget.Roles,
		"shell_name":                        deploymentTarget.ShellName,
		"shell_version":                     deploymentTarget.ShellVersion,
		"space_id":                          deploymentTarget.SpaceID,
		"status":                            deploymentTarget.Status,
		"status_summary":                    deploymentTarget.StatusSummary,
		"tenanted_deployment_participation": deploymentTarget.TenantedDeploymentMode,
		"tenants":                           deploymentTarget.TenantIDs,
		"tenant_tags":                       deploymentTarget.TenantTags,
		"thumbprint":                        deploymentTarget.Thumbprint,
		"uri":                               deploymentTarget.URI,
	}
}

func getDeploymentTargetDataSchema() map[string]*schema.Schema {
	dataSchema := getDeploymentTargetSchema()
	setDataSchema(&dataSchema)

	return map[string]*schema.Schema{
		"communication_styles": getQueryCommunicationStyles(),
		"deployment_id":        getQueryDeploymentID(),
		"deployment_targets": {
			Computed:    true,
			Description: "A list of deployment targets that match the filter(s).",
			Elem:        &schema.Resource{Schema: dataSchema},
			Optional:    true,
			Type:        schema.TypeList,
		},
		"environments":    getQueryEnvironments(),
		"health_statuses": getQueryHealthStatuses(),
		"ids":             getQueryIDs(),
		"is_disabled":     getQueryIsDisabled(),
		"name":            getQueryName(),
		"partial_name":    getQueryPartialName(),
		"roles":           getQueryRoles(),
		"shell_names":     getQueryShellNames(),
		"skip":            getQuerySkip(),
		"take":            getQueryTake(),
		"tenants":         getQueryTenants(),
		"tenant_tags":     getQueryTenantTags(),
		"thumbprint":      getQueryThumbprint(),
		"space_id":        getSpaceIDSchema(),
	}
}

func getDeploymentTargetSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
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
	}
}

func setDeploymentTarget(ctx context.Context, d *schema.ResourceData, deploymentTarget *machines.DeploymentTarget) error {
	d.Set("has_latest_calamari", deploymentTarget.HasLatestCalamari)
	d.Set("health_status", deploymentTarget.HealthStatus)
	d.Set("is_disabled", deploymentTarget.IsDisabled)
	d.Set("is_in_process", deploymentTarget.IsInProcess)
	d.Set("machine_policy_id", deploymentTarget.MachinePolicyID)
	d.Set("name", deploymentTarget.Name)
	d.Set("operating_system", deploymentTarget.OperatingSystem)
	d.Set("shell_name", deploymentTarget.ShellName)
	d.Set("shell_version", deploymentTarget.ShellVersion)
	d.Set("space_id", deploymentTarget.SpaceID)
	d.Set("status", deploymentTarget.Status)
	d.Set("status_summary", deploymentTarget.StatusSummary)
	d.Set("tenanted_deployment_participation", deploymentTarget.TenantedDeploymentMode)
	d.Set("thumbprint", deploymentTarget.Thumbprint)
	d.Set("uri", deploymentTarget.URI)
	d.Set("space_id", deploymentTarget.SpaceID)

	endpointResource, err := machines.ToEndpointResource(deploymentTarget.Endpoint)
	if err != nil {
		return fmt.Errorf("error setting endpoint: %s", err)
	}

	if d.Get("endpoint") != nil {
		if err := d.Set("endpoint", flattenEndpointResource(endpointResource)); err != nil {
			return fmt.Errorf("error setting endpoint: %s", err)
		}
	}

	if err := d.Set("environments", deploymentTarget.EnvironmentIDs); err != nil {
		return fmt.Errorf("error setting environments: %s", err)
	}

	if err := d.Set("roles", deploymentTarget.Roles); err != nil {
		return fmt.Errorf("error setting roles: %s", err)
	}

	if err := d.Set("tenants", deploymentTarget.TenantIDs); err != nil {
		return fmt.Errorf("error setting tenants: %s", err)
	}

	if err := d.Set("tenant_tags", deploymentTarget.TenantTags); err != nil {
		return fmt.Errorf("error setting tenant_tags: %s", err)
	}

	d.SetId(deploymentTarget.GetID())

	return nil
}
