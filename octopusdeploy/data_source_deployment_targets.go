package octopusdeploy

import (
	"context"
	"time"

	"github.com/OctopusDeploy/go-octopusdeploy/octopusdeploy"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceDeploymentTargets() *schema.Resource {
	return &schema.Resource{
		Description: "Provides information about existing deployment targets.",
		ReadContext: dataSourceDeploymentTargetsRead,
		Schema:      getDeploymentTargetDataSchema(),
	}
}

func dataSourceDeploymentTargetsRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	query := octopusdeploy.MachinesQuery{
		CommunicationStyles: expandArray(d.Get("communication_styles").([]interface{})),
		DeploymentID:        d.Get("deployment_id").(string),
		EnvironmentIDs:      expandArray(d.Get("environments").([]interface{})),
		HealthStatuses:      expandArray(d.Get("health_statuses").([]interface{})),
		IDs:                 expandArray(d.Get("ids").([]interface{})),
		IsDisabled:          d.Get("is_disabled").(bool),
		Name:                d.Get("name").(string),
		PartialName:         d.Get("partial_name").(string),
		Roles:               expandArray(d.Get("roles").([]interface{})),
		ShellNames:          expandArray(d.Get("shell_names").([]interface{})),
		Skip:                d.Get("skip").(int),
		Take:                d.Get("take").(int),
		TenantIDs:           expandArray(d.Get("tenants").([]interface{})),
		TenantTags:          expandArray(d.Get("tenant_tags").([]interface{})),
		Thumbprint:          d.Get("thumbprint").(string),
	}

	client := m.(*octopusdeploy.Client)
	deploymentTargets, err := client.Machines.Get(query)
	if err != nil {
		return diag.FromErr(err)
	}

	flattenedDeploymentTargets := []interface{}{}
	flattenedListeningTentacleDeploymentTargets := []interface{}{}
	flattenedOfflinePackageDropDeploymentTargets := []interface{}{}
	flattenedPollingTentacleDeploymentTargets := []interface{}{}

	for _, deploymentTarget := range deploymentTargets.Items {
		switch deploymentTarget.Endpoint.GetCommunicationStyle() {
		case "OfflineDrop":
			flattenedOfflinePackageDropDeploymentTargets = append(flattenedOfflinePackageDropDeploymentTargets, flattenOfflinePackageDropDeploymentTarget(deploymentTarget))
		case "TentacleActive":
			flattenedPollingTentacleDeploymentTargets = append(flattenedPollingTentacleDeploymentTargets, flattenPollingTentacleDeploymentTarget(deploymentTarget))
		case "TentaclePassive":
			flattenedListeningTentacleDeploymentTargets = append(flattenedListeningTentacleDeploymentTargets, flattenListeningTentacleDeploymentTarget(deploymentTarget))
		default:
			flattenedDeploymentTargets = append(flattenedDeploymentTargets, flattenDeploymentTarget(deploymentTarget))
		}
	}

	d.Set("deployment_targets", flattenedDeploymentTargets)
	d.Set("listening_tentacles", flattenedListeningTentacleDeploymentTargets)
	d.Set("offline_package_drops", flattenedOfflinePackageDropDeploymentTargets)
	d.Set("polling_tentacles", flattenedPollingTentacleDeploymentTargets)

	d.SetId("DeploymentTargets " + time.Now().UTC().String())

	return nil
}
