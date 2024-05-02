package octopusdeploy

import (
	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/machines"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"net/url"
)

// From TF
func expandKubernetesAgentDeploymentTarget(kubernetesAgent *schema.ResourceData) *machines.DeploymentTarget {
	name := kubernetesAgent.Get("name").(string)
	isDisabled := kubernetesAgent.Get("is_disabled").(bool)
	environments := expandArray(kubernetesAgent.Get("environments").([]interface{}))
	roles := expandArray(kubernetesAgent.Get("roles").([]interface{}))

	// Endpoint
	uri, _ := url.Parse(kubernetesAgent.Get("uri").(string))
	thumbprint := kubernetesAgent.Get("thumbprint").(string)

	if kubernetesAgent.Get("uri").(string) == "" {
		uri, _ = url.Parse("poll://aaaaaaaaaaaaaaaaaaaa/")
	}
	if kubernetesAgent.Get("thumbprint").(string) == "" {
		thumbprint = "0000000"
	}

	namespace := kubernetesAgent.Get("namespace").(string)
	communicationsMode := kubernetesAgent.Get("communication_mode").(string)
	upgradeLocked := kubernetesAgent.Get("upgrade_locked").(bool)
	var endpoint machines.IEndpoint = machines.NewKubernetesTentacleEndpoint(uri, thumbprint, upgradeLocked, communicationsMode, namespace)

	deploymentTarget := machines.NewDeploymentTarget(name, endpoint, environments, roles)
	deploymentTarget.Thumbprint = thumbprint

	deploymentTarget.IsDisabled = isDisabled
	// TODO handle tenants

	return deploymentTarget
}

// From API
func flattenKubernetesAgentDeploymentTarget(deploymentTarget *machines.DeploymentTarget) map[string]interface{} {
	if deploymentTarget == nil {
		return nil
	}

	//if deploymentTarget.Endpoint.GetCommunicationStyle() != "KubernetesTentacle" {
	//	return nil
	//}
	endpoint := deploymentTarget.Endpoint.(*machines.KubernetesTentacleEndpoint)

	flattenedDeploymentTarget := map[string]interface{}{}
	flattenedDeploymentTarget["id"] = deploymentTarget.GetID()
	flattenedDeploymentTarget["space_id"] = deploymentTarget.SpaceID
	flattenedDeploymentTarget["name"] = deploymentTarget.Name
	flattenedDeploymentTarget["machine_policy_id"] = deploymentTarget.MachinePolicyID
	flattenedDeploymentTarget["is_disabled"] = deploymentTarget.IsDisabled
	flattenedDeploymentTarget["thumbprint"] = endpoint.TentacleEndpointConfiguration.Thumbprint
	flattenedDeploymentTarget["uri"] = endpoint.TentacleEndpointConfiguration.URI.String()
	flattenedDeploymentTarget["communication_mode"] = endpoint.TentacleEndpointConfiguration.CommunicationMode
	flattenedDeploymentTarget["namespace"] = endpoint.DefaultNamespace

	return flattenedDeploymentTarget
}

// TODO add tenants
func getKubernetesAgentDeploymentTargetSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"space_id": getSpaceIDSchema(),
		"name":     getNameSchema(true),
		"environments": {
			Description: "A list of environment IDs associated with this resource.",
			Elem:        &schema.Schema{Type: schema.TypeString},
			MinItems:    1,
			Required:    true,
			Type:        schema.TypeList,
		},
		"roles": {
			Elem:     &schema.Schema{Type: schema.TypeString},
			MinItems: 1,
			Required: true,
			Type:     schema.TypeList,
		},
		"communication_mode": {
			Optional: true,
			Type:     schema.TypeString,
			Default:  "Polling",
		},
		"machine_policy_id": {
			Computed: true,
			Optional: true,
			Type:     schema.TypeString,
		},
		"thumbprint": {
			Optional: true,
			Computed: true,
			Type:     schema.TypeString,
			//DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
			//	if new == "0000000" {
			//		return true
			//	}
			//
			//	return false
			//},
		},
		"uri": {
			Optional: true,
			Computed: true,
			Type:     schema.TypeString,
			//DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
			//	if new == "poll://aaaaaaaaaaaaaaaaaaaa/" {
			//		return true
			//	}
			//
			//	return false
			//},
		},
		"namespace": {
			Computed: true,
			Optional: true,
			Type:     schema.TypeString,
		},
		"upgrade_locked": {
			Optional: true,
			Type:     schema.TypeBool,
		},
		"is_disabled": {
			Optional: true,
			Type:     schema.TypeBool,
		},
	}
}
