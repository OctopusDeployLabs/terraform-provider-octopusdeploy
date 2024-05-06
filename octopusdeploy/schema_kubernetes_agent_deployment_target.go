package octopusdeploy

import (
	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/core"
	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/machines"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"net/url"
)

const PlaceholderKubernetesTentacleUri = "poll://aaaaaaaaaaaaaaaaaaaa/"
const PlaceholderKubernetesTentacleThumbprint = "0000000"

func expandKubernetesAgentDeploymentTarget(kubernetesAgent *schema.ResourceData) *machines.DeploymentTarget {
	uri, _ := url.Parse(kubernetesAgent.Get("uri").(string))
	thumbprint := kubernetesAgent.Get("thumbprint").(string)

	// If the URI or thumbprint is empty, we set it to these placeholder values as we cannot create a KubernetesTentacleEndpoint on server without them
	// Once the Agent with the same target name registers itself with server it will override these values with the correct ones
	// We also want these values to be settable in TF to allow users to manage an existing installation of an Agent
	if kubernetesAgent.Get("uri").(string) == "" {
		uri, _ = url.Parse(PlaceholderKubernetesTentacleUri)
	}
	if kubernetesAgent.Get("thumbprint").(string) == "" {
		thumbprint = PlaceholderKubernetesTentacleThumbprint
	}

	defaultNamespace := kubernetesAgent.Get("default_namespace").(string)
	communicationsMode := kubernetesAgent.Get("communication_mode").(string)
	upgradeLocked := kubernetesAgent.Get("upgrade_locked").(bool)
	var endpoint machines.IEndpoint = machines.NewKubernetesTentacleEndpoint(uri, thumbprint, upgradeLocked, communicationsMode, defaultNamespace)

	name := kubernetesAgent.Get("name").(string)
	environments := expandArray(kubernetesAgent.Get("environments").([]interface{}))
	roles := expandArray(kubernetesAgent.Get("roles").([]interface{}))
	deploymentTarget := machines.NewDeploymentTarget(name, endpoint, environments, roles)

	deploymentTarget.IsDisabled = kubernetesAgent.Get("is_disabled").(bool)
	deploymentTarget.Thumbprint = thumbprint

	deploymentTarget.TenantedDeploymentMode = core.TenantedDeploymentMode(kubernetesAgent.Get("tenanted_deployment_participation").(string))
	deploymentTarget.TenantIDs = expandArray(kubernetesAgent.Get("tenants").([]interface{}))
	deploymentTarget.TenantTags = expandArray(kubernetesAgent.Get("tenant_tags").([]interface{}))

	return deploymentTarget
}

func flattenKubernetesAgentDeploymentTarget(deploymentTarget *machines.DeploymentTarget) map[string]interface{} {
	if deploymentTarget == nil {
		return nil
	}

	if deploymentTarget.Endpoint.GetCommunicationStyle() != "KubernetesTentacle" {
		return nil
	}

	endpoint := deploymentTarget.Endpoint.(*machines.KubernetesTentacleEndpoint)

	flattenedDeploymentTarget := map[string]interface{}{}
	flattenedDeploymentTarget["id"] = deploymentTarget.GetID()
	flattenedDeploymentTarget["space_id"] = deploymentTarget.SpaceID
	flattenedDeploymentTarget["name"] = deploymentTarget.Name
	flattenedDeploymentTarget["machine_policy_id"] = deploymentTarget.MachinePolicyID
	flattenedDeploymentTarget["is_disabled"] = deploymentTarget.IsDisabled
	flattenedDeploymentTarget["tenanted_deployment_participation"] = deploymentTarget.TenantedDeploymentMode
	flattenedDeploymentTarget["tenants"] = deploymentTarget.TenantIDs
	flattenedDeploymentTarget["tenant_tags"] = deploymentTarget.TenantTags

	flattenedDeploymentTarget["thumbprint"] = endpoint.TentacleEndpointConfiguration.Thumbprint
	flattenedDeploymentTarget["uri"] = endpoint.TentacleEndpointConfiguration.URI.String()
	flattenedDeploymentTarget["communication_mode"] = endpoint.TentacleEndpointConfiguration.CommunicationMode
	flattenedDeploymentTarget["default_namespace"] = endpoint.DefaultNamespace

	return flattenedDeploymentTarget
}

func getKubernetesAgentDeploymentTargetSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"space_id": getSpaceIDSchema(),
		"name":     getNameSchema(true),
		"environments": {
			Description: "A list of environment IDs this Kubernetes agent can deploy to.",
			Elem:        &schema.Schema{Type: schema.TypeString},
			MinItems:    1,
			Required:    true,
			Type:        schema.TypeList,
		},
		"roles": {
			Description: "A list of target roles that are associated to this Kubernetes agent.",
			Elem:        &schema.Schema{Type: schema.TypeString},
			MinItems:    1,
			Required:    true,
			Type:        schema.TypeList,
		},
		"tenanted_deployment_participation": getTenantedDeploymentSchema(),
		"tenants":                           getTenantsSchema(),
		"tenant_tags":                       getTenantTagsSchema(),
		"communication_mode": {
			Optional:         true,
			Description:      "The communication mode used by the Kubernetes agent to communicate with Octopus Server. Currently, the only supported value is 'Polling'.",
			Type:             schema.TypeString,
			Default:          "Polling",
			ValidateDiagFunc: validation.ToDiagFunc(validation.StringInSlice([]string{"Polling"}, false)),
		},
		"machine_policy_id": {
			Description: "Optional ID of the machine policy that the Kubernetes agent will use. If not provided the default machine policy will be used.",
			Computed:    true,
			Optional:    true,
			Type:        schema.TypeString,
		},
		"thumbprint": {
			Description: "The thumbprint of the Kubernetes agent's certificate used by server to verify the identify of the agent. This is optional as a new installation of the agent will update this value when registers with the server.",
			Optional:    true,
			Computed:    true,
			Type:        schema.TypeString,
		},
		"uri": {
			Description: "The URI of the Kubernetes agent's used by the server to queue messages. This is optional as a new installation of the agent will update this value when registers with the server.",
			Optional:    true,
			Computed:    true,
			Type:        schema.TypeString,
		},
		"default_namespace": {
			Description: "Optional default namespace that will be used when using Kubernetes deployment steps, can be overrides within step configurations.",
			Computed:    true,
			Optional:    true,
			Type:        schema.TypeString,
		},
		"upgrade_locked": {
			Description: "If enabled the Kubernetes agent will not automatically upgrade and will stay on it's current version, even if the associated machine policy is configured to automatically upgrade.",
			Optional:    true,
			Type:        schema.TypeBool,
		},
		"is_disabled": {
			Description: "Whether the Kubernetes agent is disabled. If the agent is disabled, it will not be included in any deployments.",
			Optional:    true,
			Default:     false,
			Type:        schema.TypeBool,
		},
	}
}
